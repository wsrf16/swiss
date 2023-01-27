package tcpkit

import (
	"github.com/wsrf16/swiss/sugar/base/control"
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
	conn.SetDeadline(time.Now().Add(socketkit.DefaultDeadLineDuration))
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
	conn.SetDeadline(time.Now().Add(socketkit.DefaultDeadLineDuration))

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

func TransferFromListenAddress(lAddress string) error {
	lAddr, err := NewTCPAddr(lAddress)
	if err != nil {
		return err
	}

	return TransferFromListenToDial(lAddr, nil, true)
}

// laddress: local-address, listen-address
// daddress: dial-address, destination-address
func TransferFromListenToDialAddress(lAddress string, dAddress string) error {
	lAddr, err := NewTCPAddr(lAddress)
	if err != nil {
		return err
	}
	dAddr, err := NewTCPAddr(dAddress)
	if err != nil {
		return err
	}

	return TransferFromListenToDial(lAddr, dAddr, true)
}

func TransferFromListenToListenAddress(lAddressFrom string, lAddressTo string) error {
	lAddrFrom, err := NewTCPAddr(lAddressFrom)
	if err != nil {
		return err
	}
	lAddrTo, err := NewTCPAddr(lAddressTo)
	if err != nil {
		return err
	}

	return TransferFromListenToListen(lAddrFrom, lAddrTo, true)
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

func TransferFromListenToDial(lAddr *net.TCPAddr, dAddr *net.TCPAddr, autoReconnect bool) error {
	return control.LoopAlwaysReturn(autoReconnect, func() error {
		listenerFrom, err := Listen(lAddr)
		if err != nil {
			return err
		}
		clientConnFactory := func() (net.Conn, error) {
			return Accept(listenerFrom)
		}
		serverConnFactory := func() (net.Conn, error) {
			if dAddr == nil {
				return nil, nil
			} else {
				return DialAddr(nil, dAddr)
			}
		}

		for {
			client, err := clientConnFactory()
			if err != nil {
				log.Println(err)
				return err
			}
			//server, err := serverConnFactory()
			//if err != nil {
			//    log.Println(err)
			//    return err
			//}
			//go Transfer(client, server, 1, true, true)
			go TransferDynamic(client, serverConnFactory, 1, true, true)
		}
		return nil
	})
}

func TransferFromListenToListen(lAddrFrom *net.TCPAddr, lAddrTo *net.TCPAddr, autoReconnect bool) error {
	return control.LoopAlwaysReturn(autoReconnect, func() error {
		listenerFrom, err := Listen(lAddrFrom)
		if err != nil {
			return err
		}
		listenerTo, err := Listen(lAddrTo)
		if err != nil {
			return err
		}

		serverConnFactory := func() (net.Conn, error) {
			return Accept(listenerTo)
		}
		clientConnFactory := func() (net.Conn, error) {
			return Accept(listenerFrom)
		}

		for {
			client, err := clientConnFactory()
			if err != nil {
				log.Println(err)
				return err
			}
			//server, err := serverConnFactory()
			//if err != nil {
			//    log.Println(err)
			//    time.Sleep(3 * time.Second)
			//    return err
			//}
			go TransferDynamic(client, serverConnFactory, 1, true, true)
		}
		return nil
	})
}

func TransferFromDialToDial(dAddrFrom *net.TCPAddr, dAddrTo *net.TCPAddr, autoReconnect bool) error {
	return control.LoopAlwaysReturn(autoReconnect, func() error {
		clientConnFactory := func() (net.Conn, error) {
			if dAddrFrom == nil {
				return nil, nil
			} else {
				return DialAddr(nil, dAddrFrom)
			}
		}
		serverConnFactory := func() (net.Conn, error) {
			if dAddrTo == nil {
				return nil, nil
			} else {
				return DialAddr(nil, dAddrTo)
			}
		}

		for {
			client, err := clientConnFactory()
			if err != nil {
				log.Println(err)
				time.Sleep(5 * time.Second)
				return err
			}
			TransferDynamic(client, serverConnFactory, 1, true, true)
		}
		return nil
	})
}

func TransferDynamic(client net.Conn, serverConnFactory socketkit.ConnFactoryFunc, wait int, closeClient bool, closeServer bool) ([]int, error) {
	server, err := serverConnFactory()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer server.Close()
	return Transfer(client, server, wait, closeClient, closeServer)
}

func Transfer(client net.Conn, server net.Conn, wait int, closeClient bool, closeServer bool) ([]int, error) {
	if client != nil && closeClient {
		defer client.Close()
	}
	if server != nil && closeServer {
		defer server.Close()
	}

	if server == nil {
		readBytes, err := iokit.ReadAllBytesBlockless(client)
		if err != nil {
			return nil, err
		}

		packet, err := ResolvePacket(readBytes)
		if err != nil {
			return nil, err
		}

		server, err = packet.DialDSTConn()
		if err != nil {
			return nil, err
		}
		if server != nil && closeServer {
			defer server.Close()
		}

		if packet.IsMethodConnect() {
			_, err := iokit.WriteString(client, ConnectEstablished)
			if err != nil {
				return nil, err
			}
		} else {
			_, err := iokit.Write(server, readBytes)
			if err != nil {
				return nil, err
			}
		}
	}

	//iokit.TransferRoundTrip(client, server)
	return socketkit.TransferRoundTripThenClose(client, server, wait, closeClient, closeServer), nil
	//return nil
}

func Telnet(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil || conn == nil {
		return false
	} else {
		return true
	}
}
