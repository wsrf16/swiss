package udpkit

import (
	"github.com/wsrf16/swiss/sugar/base/control"
	"github.com/wsrf16/swiss/sugar/netkit/socket"
	"net"
)

func NewUDPAddr(address string) (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", address)
}

func NewUDP4Addr(address string) (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp4", address)
}

func NewUDP6Addr(address string) (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp6", address)
}

func NewUDPAddrFromIPPort(ip string, port int, zone string) *net.UDPAddr {
	return &net.UDPAddr{IP: net.ParseIP(ip), Port: port, Zone: zone}
}

func TransferToServe(lAddress string, dstAddress string) error {
	lAddr, err := NewUDPAddr(lAddress)
	if err != nil {
		return err
	}
	dstAddr, err := NewUDPAddr(dstAddress)
	if err != nil {
		return err
	}

	return ListenThenTransferToUDP(lAddr, dstAddr, true)
}

func ListenThenTransferToUDP(laddr *net.UDPAddr, daddr *net.UDPAddr, autoReconnect bool) error {
	dstConnFactory := func() (net.Conn, error) {
		return dialAddr(daddr)
	}
	return ListenThenTransferTo(laddr, dstConnFactory, autoReconnect)
}

func ListenThenTransferTo(laddr *net.UDPAddr, dstConnFactory socket.ConnFactoryFunc, autoReconnect bool) error {
	return control.LoopAlwaysReturn(autoReconnect, func() error {
		return listenThenTransferTo(laddr, dstConnFactory)
	})
}

func listenThenTransferTo(laddr *net.UDPAddr, dstConnFactory socket.ConnFactoryFunc) error {
	listener, err := net.ListenUDP(laddr.Network(), laddr)
	if err != nil {
		return err
	}

	go socket.TransferDynamic(listener, dstConnFactory, true)
	return nil
}

func dialAddr(addr *net.UDPAddr) (net.Conn, error) {
	return net.DialUDP(addr.Network(), nil, addr)
}
