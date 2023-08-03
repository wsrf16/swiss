package tcpkit

import (
	"github.com/wsrf16/swiss/sugar/base/lambda"
	"github.com/wsrf16/swiss/sugar/base/timekit"
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"github.com/wsrf16/swiss/sugar/netkit/socket/socketkit"
	"log"
	"net"
	"time"
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

func Listen(addr *net.TCPAddr) (*net.TCPListener, error) {
	listener, err := net.ListenTCP(addr.Network(), addr)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

func Accept(listener *net.TCPListener) (net.Conn, error) {
	conn, err := listener.AcceptTCP()
	if err != nil {
		return nil, err
	}
	conn.SetDeadline(timekit.Time1Hour())
	return conn, nil
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

func TransferFromListenAddress(lAddress string, keepListening bool, listenerChannel chan *net.TCPListener) error {
	lAddr, err := NewTCPAddr(lAddress)
	if err != nil {
		return err
	}

	return TransferFromListenToDial(lAddr, nil, keepListening, listenerChannel)
}

// laddress: local-address, listen-address
// daddress: dial-address, destination-address
func TransferFromListenToDialAddress(lAddress string, dAddress string, keepListening bool, listenerChannel chan *net.TCPListener) error {
	lAddr, err := NewTCPAddr(lAddress)
	if err != nil {
		return err
	}
	dAddr, err := NewTCPAddr(dAddress)
	if err != nil {
		return err
	}

	return TransferFromListenToDial(lAddr, dAddr, keepListening, listenerChannel)
}

func TransferFromListenToListenAddress(lAddressFrom string, lAddressTo string, keepListening bool, srcListenerChannel chan *net.TCPListener, dstListenerChannel chan *net.TCPListener) error {
	lAddrFrom, err := NewTCPAddr(lAddressFrom)
	if err != nil {
		return err
	}
	lAddrTo, err := NewTCPAddr(lAddressTo)
	if err != nil {
		return err
	}

	return TransferFromListenToListen(lAddrFrom, lAddrTo, keepListening, srcListenerChannel, dstListenerChannel)
}
func TransferFromDialToDialAddress(dAddressFrom string, dAddressTo string) error {
	dAddrFrom, err := NewTCPAddr(dAddressFrom)
	if err != nil {
		return err
	}
	dAddrTo, err := NewTCPAddr(dAddressTo)
	if err != nil {
		return err
	}

	return TransferFromDialToDial(dAddrFrom, dAddrTo, true)
}

func TransferFromListenToDial(lAddr *net.TCPAddr, dAddr *net.TCPAddr, keepListening bool, listenerChannel chan *net.TCPListener) error {
	return lambda.LoopAlwaysReturn(keepListening, func() error {
		listenerFrom, err := Listen(lAddr)
		if err != nil {
			return err
		}
		if listenerChannel != nil {
			listenerChannel <- listenerFrom
		}

		srcConnFactory := func() (net.Conn, error) {
			return Accept(listenerFrom)
		}
		dstConnFactory := func() (net.Conn, error) {
			if dAddr == nil {
				return nil, nil
			} else {
				return DialAddr(nil, dAddr)
			}
		}
		for {
			src, err := srcConnFactory()
			if err != nil {
				log.Println(err)
				return err
			}
			// dst, err := dstConnFactory()
			// if err != nil {
			//    log.Println(err)
			//    return err
			// }
			// go Transfer(src, dst, 1, true, true)
			go TransferDynamic(src, dstConnFactory, true)
		}
		return nil
	})
}

func TransferFromListenToListen(lAddrFrom *net.TCPAddr, lAddrTo *net.TCPAddr, keepListening bool, srcListenerChannel chan *net.TCPListener, dstListenerChannel chan *net.TCPListener) error {
	return lambda.LoopAlwaysReturn(keepListening, func() error {
		listenerFrom, err := Listen(lAddrFrom)
		if err != nil {
			return err
		}
		if srcListenerChannel != nil {
			srcListenerChannel <- listenerFrom
		}
		listenerTo, err := Listen(lAddrTo)
		if err != nil {
			return err
		}
		if dstListenerChannel != nil {
			dstListenerChannel <- listenerTo
		}

		dstConnFactory := func() (net.Conn, error) {
			return Accept(listenerTo)
		}
		srcConnFactory := func() (net.Conn, error) {
			return Accept(listenerFrom)
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
			closedOrder, err := Transfer(src, dst, true)
			if err != nil {
				log.Println(err)
				closedOrder = closedOrder
				return err
			}
			// go TransferDynamic(src, dstConnFactory, true)

		}
		return nil
	})
}

func TransferFromDialToDial(dAddrFrom *net.TCPAddr, dAddrTo *net.TCPAddr, keepListening bool) error {
	return lambda.LoopAlwaysReturn(keepListening, func() error {
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
			closedOrder, err := Transfer(src, dst, true)
			if err != nil {
				log.Println(err)
				closedOrder = closedOrder
				break
			}
			// go TransferDynamic(src, dstConnFactory, true)
		}
		return nil
	})
}

func TransferDynamic(src net.Conn, dstConnFactory socketkit.ConnFactoryFunc, closed bool) ([]int, error) {
	dst, err := dstConnFactory()
	// dst.SetDeadline(timekit.Time3Minutes())
	// if dst != nil {
	//    dst.Close()
	// }
	// src.Close()
	// return nil, nil

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return Transfer(src, dst, closed)
}

func Transfer(src net.Conn, dst net.Conn, closed bool) ([]int, error) {
	if closed {
		defer socketkit.Close(src, dst)
	}

	if dst == nil {
		readBytes, err := iokit.ReadAllBytesBlockless(src)
		if err != nil {
			return nil, err
		}

		packet, err := ResolvePacket(readBytes)
		if err != nil {
			return nil, err
		}

		dst, err = packet.DialDSTConn()
		if err != nil {
			return nil, err
		}
		if dst != nil && closed {
			defer socketkit.Close(dst)
		}

		if packet.IsMethodConnect() {
			_, err := iokit.WriteString(src, ConnectEstablished)
			if err != nil {
				return nil, err
			}
		} else {
			_, err := iokit.Write(dst, readBytes)
			if err != nil {
				return nil, err
			}
		}
	}

	// iokit.TransferRoundTrip(src, dst)
	return socketkit.TransferRoundTripWaitForCompleted(src, dst, closed), nil
	// return nil
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
