package sockstrans

import (
	"github.com/wsrf16/swiss/netkit/layer/transport/tcptrans"
	"github.com/wsrf16/swiss/sugar/base/lambda"
	"github.com/wsrf16/swiss/sugar/base/timekit"
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"golang.org/x/net/proxy"
	"log"
	"net"
)

func TransferFromListenAddress(lAddress string, auth *proxy.Auth, keepListening bool, output *SocksTransferContext) error {
	lAddr, err := tcptrans.NewTCPAddr(lAddress)
	if err != nil {
		return err
	}

	return TransferFromListen(lAddr, auth, keepListening, output)
}

func TransferFromListen(lAddr *net.TCPAddr, auth *proxy.Auth, keepListening bool, output *SocksTransferContext) error {
	return lambda.LoopReturn(keepListening, func() error {
		lListener, err := tcptrans.Listen(lAddr)
		if err != nil {
			return err
		}
		if output != nil {
			output.LAddr = lAddr
			output.KeepListening = keepListening
			output.LListener = lListener
		}

	For:
		for {
			src, err := tcptrans.Accept(lListener)
			if err != nil {
				return err
			}
			src.SetDeadline(timekit.Time1Year())

			// go TransferTCPOrHTTP(src, true, config, true)
			go func() {
				_, err := Transfer(src, true, auth, true)
				if err != nil {
					log.Println(err)
				}
			}()

			if output != nil && output.StopChan != nil {
				select {
				case <-*output.StopChan:
					break For
				default:
				}
			}
		}
		return nil
	})
}

type SocksConfig struct {
	Authentication *proxy.Auth
}

//func NewSocksConfig(username string, password string) *SocksConfig {
//    var config = new(SocksConfig)
//    if len(username) > 0 {
//        config.Authentication = new(proxy.Auth)
//        config.Authentication.User = username
//        config.Authentication.Password = password
//    }
//    return config
//}

type Credential struct {
	Username string
	Password string
}

func Transfer(src net.Conn, closed bool, auth *proxy.Auth, recovered bool) (chan iokit.Direction, error) {
	if closed {
		defer iokit.Close(src)
	}
	if recovered {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()
	}

	dst, err := ParseDestConnFrom(src, auth)
	if err != nil {
		return nil, err
	}
	if closed {
		defer iokit.Close(dst)
	}

	return iokit.CopyDuplex(src, dst, closed), nil
}

func Dial(network, targetAddress, proxyAddress string, auth *proxy.Auth, forward proxy.Dialer) (net.Conn, error) {
	dialer, err := proxy.SOCKS5(network, proxyAddress, auth, forward)
	if err != nil {
		return nil, err
	}
	conn, err := dialer.Dial(network, targetAddress)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
