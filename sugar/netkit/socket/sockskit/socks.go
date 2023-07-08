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

func TransferFromListenAddress(lAddress string, config *SocksConfig) error {
	lAddr, err := tcpkit.NewTCPAddr(lAddress)
	if err != nil {
		return err
	}

	return TransferFromListen(lAddr, true, config)
}

func TransferFromListen(laddr *net.TCPAddr, autoReconnect bool, config *SocksConfig) error {
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

			go Transfer(client, true, config)
		}
		return nil
	})
}

type SocksConfig struct {
	Credential Credential
}

func NewSocksConfig(username string, password string) *SocksConfig {
	var config SocksConfig
	config.Credential = Credential{}
	config.Credential.Username = username
	config.Credential.Password = password
	return &config
}

type Credential struct {
	Username string
	Password string
}

func Transfer(client net.Conn, closed bool, config *SocksConfig) ([]int, error) {
	if closed {
		defer socketkit.Close(client)
	}

	// round 1
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
	if config == nil {
		client.Write([]byte{0x05, 0x00})
	} else {
		client.Write([]byte{0x05, 0x02})

		// round 2
		readBytes, err = iokit.ReadAllBytesBlockless(client)
		if err != nil {
			return nil, err
		}
		packet2 := ResolvePacket2(readBytes)
		credential := config.Credential

		if packet2.UserName == nil || packet2.Password == nil {
			return nil, errors.New("incorrect username or password")
		}

		if string(packet2.UserName) == credential.Username && string(packet2.Password) == credential.Password {
			_, err := client.Write([]byte{0x05, 0x00})
			if err != nil {
				return nil, err
			}
		} else {
			client.Write([]byte{0x05, 0x01})
			return nil, errors.New("authentication failed")
		}
	}

	// round 3
	readBytes, err = iokit.ReadAllBytesBlockless(client)
	if err != nil {
		return nil, err
	}
	packet3 := BuildPacket3(readBytes)
	if packet3.Version != 5 {
		return nil, errors.New("该协议不是socks5协议")
	}
	client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	address, err := packet3.GetAddress()
	if err != nil {
		return nil, err
	}
	log.Printf("type: socks  version: %d  address: \"%s\"\n", packet3.Version, address)

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
