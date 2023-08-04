package sockskit

import (
	"errors"
	"github.com/wsrf16/swiss/sugar/base/lambda"
	"github.com/wsrf16/swiss/sugar/base/timekit"
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"github.com/wsrf16/swiss/sugar/netkit/socket/socketkit"
	"github.com/wsrf16/swiss/sugar/netkit/socket/tcpkit"
	"golang.org/x/net/proxy"
	"log"
	"net"
)

func TransferFromListenAddress(lAddress string, config *SocksConfig, keepListening bool, listenerChannel chan *net.TCPListener) error {
	lAddr, err := tcpkit.NewTCPAddr(lAddress)
	if err != nil {
		return err
	}

	return TransferFromListen(lAddr, config, keepListening, listenerChannel)
}

func TransferFromListen(laddr *net.TCPAddr, config *SocksConfig, keepListening bool, listenerChannel chan *net.TCPListener) error {
	return lambda.LoopAlwaysReturn(keepListening, func() error {
		listener, err := tcpkit.Listen(laddr)
		if err != nil {
			return err
		}
		if listenerChannel != nil {
			listenerChannel <- listener
		}

		for {
			src, err := tcpkit.Accept(listener)
			if err != nil {
				return err
			}
			src.SetDeadline(timekit.Time1Year())

			//go Transfer(src, true, config, true)
			go func() {
				_, err := Transfer(src, true, config, true)
				if err != nil {
					log.Println(err)
				}
			}()
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

func Transfer(src net.Conn, closed bool, config *SocksConfig, recoverd bool) ([]int, error) {
	if closed {
		defer socketkit.Close(src)
	}
	if recoverd {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()
	}

	// round 1
	readBytes, err := iokit.ReadAllBytesBlockless(src)
	if err != nil {
		return nil, err
	}
	packet1, err := ResolvePacket1(readBytes)
	if err != nil {
		return nil, err
	}
	if packet1.Version != 5 {
		return nil, errors.New("该协议不是socks5协议")
	} else {
		log.Printf("socks: version-%d\n", packet1.Version)
	}
	if config == nil {
		src.Write([]byte{0x05, 0x00})
	} else {
		src.Write([]byte{0x05, 0x02})

		// round 2
		readBytes, err = iokit.ReadAllBytesBlockless(src)
		if err != nil {
			return nil, err
		}
		packet2 := ResolvePacket2(readBytes)
		credential := config.Credential

		if packet2.UserName == nil || packet2.Password == nil {
			return nil, errors.New("incorrect username or password")
		}

		if string(packet2.UserName) == credential.Username && string(packet2.Password) == credential.Password {
			_, err := src.Write([]byte{0x05, 0x00})
			if err != nil {
				return nil, err
			}
		} else {
			src.Write([]byte{0x05, 0x01})
			return nil, errors.New("authentication failed")
		}
	}

	// round 3
	readBytes, err = iokit.ReadAllBytesBlockless(src)
	if err != nil {
		return nil, err
	}
	packet3 := BuildPacket3(readBytes)
	if packet3.Version != 5 {
		return nil, errors.New("该协议不是socks5协议")
	}
	src.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	address, err := packet3.GetAddress()
	if err != nil {
		return nil, err
	}
	log.Printf("type: socks  version: %d  address: \"%s\"\n", packet3.Version, address)

	dst, err := tcpkit.DialAddress("", address)
	if err != nil {
		return nil, err
	}
	dst.SetDeadline(timekit.Time1Year())

	if closed {
		defer socketkit.Close(dst)
	}

	return socketkit.TransferRoundTripWaitForCompleted(src, dst, closed), nil
}

func Dial(network, address string, auth *proxy.Auth, forward proxy.Dialer) (net.Conn, error) {
	auth = &proxy.Auth{User: auth.User, Password: auth.Password}
	dialer, err := proxy.SOCKS5(network, address, auth, forward)
	if err != nil {
		return nil, err
	}
	conn, err := dialer.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
