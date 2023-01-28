package udpkit

import (
	"github.com/wsrf16/swiss/sugar/base/lambda"
	"github.com/wsrf16/swiss/sugar/base/timekit"
	"github.com/wsrf16/swiss/sugar/netkit/socket/socketkit"
	"log"
	"net"
)

func NewUDPAddr(address string) (*net.UDPAddr, error) {
	if address == "" {
		return nil, nil
	} else {
		return net.ResolveUDPAddr("udp", address)
	}
}

func NewUDP4Addr(address string) (*net.UDPAddr, error) {
	if address == "" {
		return nil, nil
	} else {
		return net.ResolveUDPAddr("udp4", address)
	}
}

func NewUDP6Addr(address string) (*net.UDPAddr, error) {
	if address == "" {
		return nil, nil
	} else {
		return net.ResolveUDPAddr("udp6", address)
	}
}

func NewUDPAddrFromIPPort(ip string, port int, zone string) *net.UDPAddr {
	return &net.UDPAddr{IP: net.ParseIP(ip), Port: port, Zone: zone}
}

func DialAddress(lAddress, rAddress string) (net.Conn, error) {
	rAddr, err := NewUDPAddr(rAddress)
	if err != nil {
		return nil, err
	}
	lAddr, err := NewUDPAddr(lAddress)
	if err != nil {
		return nil, err
	}
	return DialAddr(lAddr, rAddr)
}

func DialAddr(lAddr, rAddr *net.UDPAddr) (net.Conn, error) {
	conn, err := net.DialUDP(rAddr.Network(), lAddr, rAddr)
	if err != nil {
		return nil, err
	}
	conn.SetDeadline(timekit.Time3Minutes())
	return conn, nil
}

func ListenAddress(address string) (net.Conn, error) {
	addr, err := NewUDPAddr(address)
	if err != nil {
		return nil, err
	}
	return Listen(addr)
}

func Listen(addr *net.UDPAddr) (net.Conn, error) {
	return net.ListenUDP(addr.Network(), addr)
}

func TransferFromListenToDialAddress(lAddress string, dAddress string) error {
	lAddr, err := NewUDPAddr(lAddress)
	if err != nil {
		return err
	}
	dAddr, err := NewUDPAddr(dAddress)
	if err != nil {
		return err
	}

	return TransferFromListenToDial(lAddr, dAddr, true)
}

func TransferFromListenToListenAddress(lAddressFrom string, lAddressTo string) error {
	lAddrFrom, err := NewUDPAddr(lAddressFrom)
	if err != nil {
		return err
	}
	lAddrTo, err := NewUDPAddr(lAddressTo)
	if err != nil {
		return err
	}

	return TransferFromListenToListen(lAddrFrom, lAddrTo, true)
}
func TransferFromDialToDialAddress(dAddressFrom string, dAddressTo string) error {
	dAddrFrom, err := NewUDPAddr(dAddressFrom)
	if err != nil {
		return err
	}
	dAddrTo, err := NewUDPAddr(dAddressTo)
	if err != nil {
		return err
	}

	return TransferFromDialToDial(dAddrFrom, dAddrTo, true)
}

func TransferFromListenToDial(lAddr *net.UDPAddr, dAddr *net.UDPAddr, autoReconnect bool) error {
	return lambda.LoopAlwaysReturn(autoReconnect, func() error {
		clientConnFactory := func() (net.Conn, error) {
			return Listen(lAddr)
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
			//Transfer(client, server, 1, true, true)
			TransferDynamic(client, serverConnFactory, true)
		}
		return nil
	})
}

func TransferFromListenToListen(lAddrFrom *net.UDPAddr, lAddrTo *net.UDPAddr, autoReconnect bool) error {
	return lambda.LoopAlwaysReturn(autoReconnect, func() error {
		clientConnFactory := func() (net.Conn, error) {
			return Listen(lAddrFrom)
		}
		serverConnFactory := func() (net.Conn, error) {
			return Listen(lAddrTo)
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
			//Transfer(client, server, true)
			TransferDynamic(client, serverConnFactory, true)
		}
		return nil
	})
}

func TransferFromDialToDial(dAddrFrom *net.UDPAddr, dAddrTo *net.UDPAddr, autoReconnect bool) error {
	return lambda.LoopAlwaysReturn(autoReconnect, func() error {
		clientConnFactory := func() (net.Conn, error) {
			return DialAddr(nil, dAddrFrom)
		}
		serverConnFactory := func() (net.Conn, error) {
			return DialAddr(nil, dAddrTo)
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
			//Transfer(client, server, 1, true, true)
			TransferDynamic(client, serverConnFactory, true)
		}
		return nil
	})
}

func TransferDynamic(client net.Conn, serverConnFactory socketkit.ConnFactoryFunc, closed bool) ([]int, error) {
	server, err := serverConnFactory()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return Transfer(client, server, closed)
}

func Transfer(client net.Conn, server net.Conn, closed bool) ([]int, error) {
	return socketkit.TransferRoundTripWaitForCompleted(client, server, closed), nil
}
