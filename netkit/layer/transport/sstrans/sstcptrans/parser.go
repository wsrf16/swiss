package sstcptrans

import (
	"github.com/wsrf16/swiss/netkit/layer/transport/tcptrans"
	"github.com/wsrf16/swiss/sugar/base/gokit"
	"github.com/wsrf16/swiss/sugar/base/timekit"
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"log"
	"net"
)

func ParseDestConnFrom(src net.Conn) (net.Conn, error) {
	ssServer := &SSServer{Conn: src}
	packet, err := ssServer.ParseIntoPacket()
	if err != nil {
		return nil, err
	}

	address, err := packet.GetAddress()
	if err != nil {
		return nil, err
	}
	log.Printf("type: shadowsock  address: \"%s\"\n", address)

	dst, err := tcptrans.DialAddress("", address)
	if err != nil {
		return nil, err
	} else {
		log.Printf("socks proxy(go-routine-id: %v): {from: %s <-> %s to: %s <-> %s}\n", gokit.GetGid(), src.LocalAddr(), src.RemoteAddr(), dst.LocalAddr(), address)
	}
	dst.SetDeadline(timekit.Time1Year())

	return dst, err
}

func (s SSServer) ParseIntoPacket() (*SSDetailPacket, error) {
	detailPacket, err := s.receiveAndParseIntoDetailPacket()
	if err != nil {
		return nil, err
	}
	return detailPacket, nil
}

func (s SSServer) receiveAndParseIntoDetailPacket() (*SSDetailPacket, error) {
	bytes, err := iokit.ReadAllBytesNonBlocking(s.Conn)
	if err != nil {
		return nil, err
	}
	detailPacket := parseIntoDetailPacket(bytes)

	return &detailPacket, nil
}
