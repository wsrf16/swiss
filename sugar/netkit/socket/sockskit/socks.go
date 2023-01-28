package sockskit

import (
	"errors"
	"github.com/wsrf16/swiss/sugar/base/lambda"
	"github.com/wsrf16/swiss/sugar/base/timekit"
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"github.com/wsrf16/swiss/sugar/netkit/socket/socketkit"
	"github.com/wsrf16/swiss/sugar/netkit/socket/tcpkit"
	"log"
	"net"
)

func TransferFromListenAddress(lAddress string) error {
	lAddr, err := tcpkit.NewTCPAddr(lAddress)
	if err != nil {
		return err
	}

	return TransferFromListen(lAddr, true)
}

func TransferFromListen(laddr *net.TCPAddr, autoReconnect bool) error {
	return lambda.LoopAlwaysReturn(autoReconnect, func() error {
		listener, err := tcpkit.Listen(laddr)
		if err != nil {
			return err
		}

		for {
			client, err := tcpkit.Accept(listener)
			if err != nil {
				return err
			}
			client.SetDeadline(timekit.Time1Year())

			go Transfer(client, true)
		}
		return nil
	})
}

func Transfer(client net.Conn, closed bool) ([]int, error) {
	if closed {
		defer socketkit.Close(client)
	}

	readBytes, err := iokit.ReadAllBytesBlockless(client)
	if err != nil {
		return nil, err
	}
	packet1 := ResolvePacket1(readBytes)
	if packet1.Version != 5 {
		return nil, errors.New("该协议不是socks5协议")
	} else {
		log.Printf("socks: version-%d\n", packet1.Version)
	}
	client.Write([]byte{0x05, 0x00})

	readBytes, err = iokit.ReadAllBytesBlockless(client)
	if err != nil {
		return nil, err
	}
	packet2 := BuildPacket2(readBytes)
	if packet2.Version != 5 {
		return nil, errors.New("该协议不是socks5协议")
	}
	client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	address, err := packet2.GetAddress()
	if err != nil {
		return nil, err
	}
	log.Printf("type: socks  version: %d  address: \"%s\"\n", packet2.Version, address)

	server, err := tcpkit.DialAddress("", address)
	if err != nil {
		return nil, err
	}
	server.SetDeadline(timekit.Time1Year())

	if closed {
		defer socketkit.Close(server)
	}

	return socketkit.TransferRoundTripWaitForCompleted(client, server, closed), nil
}
