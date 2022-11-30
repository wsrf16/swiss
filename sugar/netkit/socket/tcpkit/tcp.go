package tcpkit

import (
	"github.com/wsrf16/swiss/sugar/base/control"
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"github.com/wsrf16/swiss/sugar/netkit/socket"
	"net"
	"time"
)

func NewTCPAddr(address string) (*net.TCPAddr, error) {
	return net.ResolveTCPAddr("tcp", address)
}

func NewTCP4Addr(address string) (*net.TCPAddr, error) {
	return net.ResolveTCPAddr("tcp4", address)
}

func NewTCP6Addr(address string) (*net.TCPAddr, error) {
	return net.ResolveTCPAddr("tcp6", address)
}

func NewTCPAddrFromIPPort(ip string, port int, zone string) *net.TCPAddr {
	return &net.TCPAddr{IP: net.ParseIP(ip), Port: port, Zone: zone}
}

func dialAddrNullable(addr *net.TCPAddr) (net.Conn, error) {
	if addr == nil {
		return nil, nil
	} else {
		return dialAddr(addr)
	}
}

func dialAddr(addr *net.TCPAddr) (net.Conn, error) {
	conn, err := net.DialTCP("tcp", nil, addr)
	if conn != nil {
		conn.SetDeadline(time.Now().Add(socket.DefaultDeadLineDuration))
	}
	return conn, err
}

func TransferToServe(lAddress string, dstAddress string) error {
	lAddr, err := NewTCPAddr(lAddress)
	if err != nil {
		return err
	}
	dstAddr, err := NewTCPAddr(dstAddress)
	if err != nil {
		return err
	}

	return ListenThenAcceptThenTransferToTCP(lAddr, dstAddr, true)
}

func TransferToHostServe(lAddress string) error {
	lAddr, err := NewTCPAddr(lAddress)
	if err != nil {
		return err
	}

	return ListenThenAcceptThenTransferToTCP(lAddr, nil, true)
}

func ListenThenAcceptThenTransferToTCP(laddr *net.TCPAddr, daddr *net.TCPAddr, autoReconnect bool) error {
	dstConnFactory := func() (net.Conn, error) {
		return dialAddrNullable(daddr)
	}
	return ListenThenAcceptThenTransferTo(laddr, dstConnFactory, autoReconnect)
}

func ListenThenAcceptThenTransferTo(laddr *net.TCPAddr, dstConnFactory socket.ConnFactoryFunc, autoReconnect bool) error {
	return control.LoopAlwaysReturn(autoReconnect, func() error {
		return listenThenAcceptThenTransferTo(laddr, dstConnFactory)
	})
}

func listenThenAcceptThenTransferTo(laddr *net.TCPAddr, dstConnFactory socket.ConnFactoryFunc) error {
	listener, err := net.ListenTCP(laddr.Network(), laddr)
	if err != nil {
		return err
	}

	for {
		client, err := listener.AcceptTCP()
		if err != nil {
			return err
		}
		client.SetDeadline(time.Now().Add(socket.DefaultDeadLineDuration))

		go TransferDynamic(client, dstConnFactory, true)
	}
}

func TransferDynamic(client net.Conn, serverConnFactory socket.ConnFactoryFunc, closeConn bool) error {
	server, err := serverConnFactory()
	if err != nil {
		return err
	}

	return Transfer(client, server, closeConn)
}

func Transfer(client net.Conn, server net.Conn, closeConn bool) error {
	if closeConn {
		if client != nil {
			defer client.Close()
		}

		if server != nil {
			defer server.Close()
		}
	}

	readBytes, err := iokit.ReadAllBytesBlockless(client)
	if err != nil {
		return err
	}

	packet, err := ResolvePacket(readBytes)
	if err != nil {
		return err
	}

	server, err = control.IfFuncPair(server == nil, packet.ResolveDSTConn, func() (net.Conn, error) {
		return server, nil
	})
	if err != nil {
		return err
	}
	if closeConn {
		if server != nil {
			defer server.Close()
		}
	}

	if packet.IsMethodConnect() {
		_, err := iokit.WriteString(client, ConnectEstablished)
		if err != nil {
			return err
		}
	} else {
		_, err := iokit.Write(server, readBytes)
		if err != nil {
			return err
		}
	}

	iokit.TransferRoundTrip(client, server)

	return nil
}
