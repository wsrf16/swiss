package udptrans

import (
	"github.com/wsrf16/swiss/sugar/base/timekit"
	"net"
)

func NewUDPAddr(address string) (*net.UDPAddr, error) {
	if address == "" {
		return nil, nil
	} else {
		return net.ResolveUDPAddr("udp", address)
	}
}

func NewUDP4Addr(address string) (*net.UDPAddr, error) {
	if address == "" {
		return nil, nil
	} else {
		return net.ResolveUDPAddr("udp4", address)
	}
}

func NewUDP6Addr(address string) (*net.UDPAddr, error) {
	if address == "" {
		return nil, nil
	} else {
		return net.ResolveUDPAddr("udp6", address)
	}
}

func NewUDPAddrFromIPPort(ip string, port int, zone string) *net.UDPAddr {
	return &net.UDPAddr{IP: net.ParseIP(ip), Port: port, Zone: zone}
}

func NewCoupleUDPAddr(lAddress string, dAddress string) (*net.UDPAddr, *net.UDPAddr, error) {
	lAddr, err := NewUDPAddr(lAddress)
	if err != nil {
		return nil, nil, err
	}
	dAddr, err := NewUDPAddr(dAddress)
	if err != nil {
		return nil, nil, err
	}
	return lAddr, dAddr, err
}

func DialAddress(lAddress, rAddress string) (net.Conn, error) {
	rAddr, err := NewUDPAddr(rAddress)
	if err != nil {
		return nil, err
	}
	lAddr, err := NewUDPAddr(lAddress)
	if err != nil {
		return nil, err
	}
	return DialAddr(lAddr, rAddr)
}

func DialAddr(lAddr, rAddr *net.UDPAddr) (net.Conn, error) {
	conn, err := net.DialUDP(rAddr.Network(), lAddr, rAddr)
	if err != nil {
		return nil, err
	}
	conn.SetDeadline(timekit.Time3Minutes())
	return conn, nil
}

func ListenAddress(address string) (net.Conn, error) {
	addr, err := NewUDPAddr(address)
	if err != nil {
		return nil, err
	}
	return Listen(addr)
}

func Listen(addr *net.UDPAddr) (net.Conn, error) {
	return net.ListenUDP(addr.Network(), addr)
}
