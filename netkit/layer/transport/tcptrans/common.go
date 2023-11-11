package tcptrans

import (
	"github.com/wsrf16/swiss/sugar/base/timekit"
	"net"
)

func NewTCPAddr(address string) (*net.TCPAddr, error) {
	if address == "" {
		return nil, nil
	} else {
		return net.ResolveTCPAddr("tcp", address)
	}
}

func NewTCP4Addr(address string) (*net.TCPAddr, error) {
	if address == "" {
		return nil, nil
	} else {
		return net.ResolveTCPAddr("tcp4", address)
	}
}

func NewTCP6Addr(address string) (*net.TCPAddr, error) {
	if address == "" {
		return nil, nil
	} else {
		return net.ResolveTCPAddr("tcp6", address)
	}
}

func NewTCPAddrFromIPPort(ip string, port int, zone string) *net.TCPAddr {
	return &net.TCPAddr{IP: net.ParseIP(ip), Port: port, Zone: zone}
}

func NewCoupleTCPAddr(lAddress string, dAddress string) (*net.TCPAddr, *net.TCPAddr, error) {
	lAddr, err := NewTCPAddr(lAddress)
	if err != nil {
		return nil, nil, err
	}
	dAddr, err := NewTCPAddr(dAddress)
	if err != nil {
		return nil, nil, err
	}
	return lAddr, dAddr, err
}

// laddress: local-address
// raddress: remote-address
func DialAddress(lAddress, rAddress string) (net.Conn, error) {
	lAddr, err := NewTCPAddr(lAddress)
	rAddr, err := NewTCPAddr(rAddress)
	if err != nil {
		return nil, err
	}

	return DialAddr(lAddr, rAddr)
}

func DialAddr(lAddr, rAddr *net.TCPAddr) (net.Conn, error) {
	conn, err := net.DialTCP("tcp", lAddr, rAddr)
	if err != nil {
		return nil, err
	}
	conn.SetDeadline(timekit.Time1Hour())
	return conn, err
}

func ListenAndAcceptAddress(address string) (*net.TCPListener, net.Conn, error) {
	addr, err := NewTCPAddr(address)
	if err != nil {
		return nil, nil, err
	}
	return ListenAndAccept(addr)
}

func ListenAndAccept(addr *net.TCPAddr) (*net.TCPListener, net.Conn, error) {
	listener, err := Listen(addr)
	if err != nil {
		return nil, nil, err
	}
	conn, err := Accept(listener)
	if err != nil {
		return listener, nil, err
	}

	return listener, conn, nil
}

func Listen(addr *net.TCPAddr) (*net.TCPListener, error) {
	listener, err := net.ListenTCP(addr.Network(), addr)
	return listener, err
}

func Accept(listener *net.TCPListener) (net.Conn, error) {
	conn, err := listener.AcceptTCP()
	if err != nil {
		return nil, err
	}
	conn.SetDeadline(timekit.Time1Hour())
	return conn, nil
}
