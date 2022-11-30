package sockskit

import (
	"errors"
	"fmt"
	"github.com/wsrf16/swiss/sugar/base/control"
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"github.com/wsrf16/swiss/sugar/netkit/socket"
	"github.com/wsrf16/swiss/sugar/netkit/socket/tcpkit"
	"net"
	"time"
)

func TransferToHostServe(lAddress string) error {
	lAddr, err := tcpkit.NewTCPAddr(lAddress)
	if err != nil {
		return err
	}

	return ListenThenAcceptThenTransferToTCP(lAddr, true)
}

func ListenThenAcceptThenTransferToTCP(laddr *net.TCPAddr, autoReconnect bool) error {
	return ListenThenAcceptThenTransferTo(laddr, autoReconnect)
}

func ListenThenAcceptThenTransferTo(laddr *net.TCPAddr, autoReconnect bool) error {
	return control.LoopAlwaysReturn(autoReconnect, func() error {
		return listenThenAcceptThenTransferTo(laddr)
	})
}

func listenThenAcceptThenTransferTo(laddr *net.TCPAddr) error {
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

		go Transfer(client, true)
	}
}

func Transfer(client net.Conn, closeConn bool) error {
	if closeConn {
		if client != nil {
			defer client.Close()
		}
	}

	readBytes, err := iokit.ReadAllBytesBlockless(client)
	if err != nil {
		return err
	}
	packet1 := BuildPacket1(readBytes)
	if packet1.Version != 5 {
		return errors.New("该协议不是socks5协议")
	} else {
		fmt.Printf("socks: version-%d\n", packet1.Version)
	}
	client.Write([]byte{0x05, 0x00})

	readBytes, err = iokit.ReadAllBytesBlockless(client)
	if err != nil {
		return err
	}
	packet2 := BuildPacket2(readBytes)
	if packet2.Version != 5 {
		return errors.New("该协议不是socks5协议")
	}
	client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	address := packet2.GetAddress()
	fmt.Printf("socks: version-%d address-\"%s\"\n", packet2.Version, address)
	addr, err := tcpkit.NewTCPAddr(address)
	if err != nil {
		return err
	}
	server, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}
	if closeConn {
		if server != nil {
			defer server.Close()
		}
	}

	iokit.TransferRoundTrip(client, server)

	return nil
}
