package ssudp

import (
	"github.com/wsrf16/swiss/netkit/layer/transport/sstrans"
	"github.com/wsrf16/swiss/sugar/base/lambda"
	"github.com/wsrf16/swiss/sugar/base/timekit"
	"github.com/wsrf16/swiss/sugar/io/iokit"
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

func TransferFromListen(lAddrFrom *net.UDPAddr, config *sstrans.ShadowSocksConfig, keepListening bool) error {
	return lambda.LoopReturn(keepListening, func() error {
		srcConnFactory := func() (net.Conn, error) {
			return Listen(lAddrFrom)
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
			// Transfer(src, dst, true)
			TransferDynamic(src, config, true)
		}
		return nil
	})
}

func TransferDynamic(srcEncrypted net.Conn, config *sstrans.ShadowSocksConfig, closed bool) (chan iokit.Direction, error) {
	cipher := *config.Cipher
	src := cipher.StreamConn(srcEncrypted)

	dst, err := parseDest(src)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if closed {
		defer iokit.Close(dst)
	}

	return Transfer(src, dst, closed)
}

func Transfer(src net.Conn, dst net.Conn, closed bool) (chan iokit.Direction, error) {
	return iokit.CopyDuplex(src, dst, closed), nil
}
