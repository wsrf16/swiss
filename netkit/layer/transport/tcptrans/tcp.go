package tcptrans

import (
	"github.com/wsrf16/swiss/sugar/base/lambda"
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"log"
	"net"
	"time"
)

// laddress: local-address, listen-address
// daddress: dial-address, destination-address
func TransferFromListenToDialAddress(lAddress string, dAddress string, keepListening bool, output *TCPTransferContext) error {
	lAddr, dAddr, err := NewCoupleTCPAddr(lAddress, dAddress)
	if err != nil {
		return err
	}

	return TransferFromListenToDial(lAddr, dAddr, keepListening, output)
}

func TransferFromListenToListenAddress(lAddress string, dAddress string, keepListening bool, output *TCPTransferContext) error {
	lAddr, dAddr, err := NewCoupleTCPAddr(lAddress, dAddress)
	if err != nil {
		return err
	}

	return TransferFromListenToListen(lAddr, dAddr, keepListening, output)
}

func TransferFromDialToDialAddress(lAddress string, dAddress string, keepListening bool) error {
	lAddr, dAddr, err := NewCoupleTCPAddr(lAddress, dAddress)
	if err != nil {
		return err
	}

	return TransferFromDialToDial(lAddr, dAddr, keepListening)
}

func TransferFromListenToDial(lAddr *net.TCPAddr, dAddr *net.TCPAddr, keepListening bool, output *TCPTransferContext) error {
	return lambda.LoopReturn(keepListening, func() error {
		lListener, err := Listen(lAddr)
		if err != nil {
			return err
		}
		if output != nil {
			output.LAddr = lAddr
			output.DAddr = dAddr
			output.KeepListening = keepListening
			output.LListener = lListener
		}

		srcConnFactory := func() (net.Conn, error) {
			return Accept(lListener)
		}
		dstConnFactory := func() (net.Conn, error) {
			return DialAddr(nil, dAddr)
		}
	For:
		for {
			src, err := srcConnFactory()
			if err != nil {
				log.Println(err)
				return err
			}
			dst, err := dstConnFactory()
			if err != nil {
				log.Println(err)
				return err
			}

			go TransferTCP(src, dst, true)

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

func TransferFromListenToListen(lAddr *net.TCPAddr, dAddr *net.TCPAddr, keepListening bool, output *TCPTransferContext) error {
	return lambda.LoopReturn(keepListening, func() error {
		lListener, err := Listen(lAddr)
		if err != nil {
			return err
		}
		if output != nil {
			output.LAddr = lAddr
			output.DAddr = dAddr
			output.KeepListening = keepListening
			output.LListener = lListener
		}
		dListener, err := Listen(dAddr)
		if err != nil {
			return err
		}
		if output != nil {
			output.DListener = dListener
		}

		dstConnFactory := func() (net.Conn, error) {
			return Accept(dListener)
		}
		srcConnFactory := func() (net.Conn, error) {
			return Accept(lListener)
		}
		for {
			dst, err := dstConnFactory()
			if err != nil {
				log.Println(err)
				time.Sleep(3 * time.Second)
				return err
			}
			src, err := srcConnFactory()
			if err != nil {
				log.Println(err)
				return err
			}
			// fmt.Println("-->" + src.RemoteAddr().String())
			closedOrder, err := TransferTCP(src, dst, true)
			if err != nil {
				log.Println(err)
				closedOrder = closedOrder
				return err
			}
			// go TransferTCPOrHTTPDynamic(src, dstConnFactory, true)

		}
		return nil
	})
}

func TransferFromDialToDial(dAddrFrom *net.TCPAddr, dAddrTo *net.TCPAddr, keepListening bool) error {
	return lambda.LoopReturn(keepListening, func() error {
		srcConnFactory := func() (net.Conn, error) {
			if dAddrFrom == nil {
				return nil, nil
			} else {
				return DialAddr(nil, dAddrFrom)
			}
		}
		dstConnFactory := func() (net.Conn, error) {
			if dAddrTo == nil {
				return nil, nil
			} else {
				return DialAddr(nil, dAddrTo)
			}
		}

		for {
			src, err := srcConnFactory()
			if err != nil {
				log.Println(err)
				time.Sleep(5 * time.Second)
				return err
			}
			dst, err := dstConnFactory()
			if err != nil {
				log.Println(err)
				time.Sleep(5 * time.Second)
				return err
			}
			// fmt.Println("-->" + dst.LocalAddr().String())
			closedOrder, err := TransferTCP(src, dst, true)
			if err != nil {
				log.Println(err)
				closedOrder = closedOrder
				break
			}
			// go TransferTCPOrHTTPDynamic(src, dstConnFactory, true)
		}
		return nil
	})
}

func TransferTCPDynamic(srcConnFactory iokit.ConnFactoryFunc, dstConnFactory iokit.ConnFactoryFunc, closed bool) (chan iokit.Direction, error) {
	src, err := srcConnFactory()
	if err != nil {
		return nil, err
	}

	dst, err := dstConnFactory()
	if err != nil {
		return nil, err
	}
	return iokit.CopyDuplex(src, dst, closed), nil
}

func TransferTCP(src net.Conn, dst net.Conn, closed bool) (chan iokit.Direction, error) {
	return iokit.CopyDuplex(src, dst, closed), nil
}

func Telnet(addr string, timeout time.Duration) (bool, error) {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil || conn == nil {
		return false, err
	} else {
		defer conn.Close()
		return true, err
	}
}

func TelnetN(addr string, times int, timeout time.Duration) (bool, error) {
	for n := 0; n <= times; n++ {
		ok, err := Telnet(addr, timeout)
		if !ok && n < times {
			continue
		}
		return ok, err
	}
	return true, nil
}
