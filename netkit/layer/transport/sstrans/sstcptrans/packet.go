package sstcptrans

import (
	"errors"
	"github.com/wsrf16/swiss/sugar/base/convertkit"
	"github.com/wsrf16/swiss/sugar/base/ipkit"
	"net"
	"strconv"
)

type SSServer struct {
	Conn net.Conn
}

/*
*

	+----+-----+-------+------+----------+----------+
	|VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	+----+-----+-------+------+----------+----------+
	| 1  |  1  | X'00' |  1   | Variable |    2     |
	+----+-----+-------+------+----------+----------+
*/
type SSDetailPacket struct {
	Packet   []byte
	ATYP     byte
	DST_ADDR []byte
	DST_PORT []byte
}

func parseIntoDetailPacket(packet []byte) SSDetailPacket {
	packet2 := SSDetailPacket{}
	packet2.Packet = packet
	packet2.ATYP = packet[0]
	packet2.DST_ADDR, packet2.DST_PORT = getAddrAndPort(packet)
	return packet2
}

const (
	ipv4   = 0x01
	domain = 0x03
	ipv6   = 0x04
)

const MaxAddrLen = 1 + 1 + 255 + 2

func getAddrAndPort(packet []byte) ([]byte, []byte) {
	var length int
	atyp := packet[0]
	switch atyp {
	case ipv4:
		length = net.IPv4len
		return packet[1 : 1+length], packet[1+length : 1+length+2]
	case domain:
		length = int(packet[1])
		return packet[2 : 2+length], packet[2+length : 2+length+2]
	case ipv6:
		length = net.IPv6len
		return packet[1 : 1+length], packet[1+length : 1+length+2]
	}
	return nil, nil
}

func (p SSDetailPacket) GetPort() int {
	return int(convertkit.BytesToUint16(p.DST_PORT))
}

func (p SSDetailPacket) GetAddr() (string, error) {
	switch p.ATYP {
	case ipv4:
		addr, err := ipkit.BytesToIPv4String(p.DST_ADDR)
		if err != nil {
			return "", err
		}
		return addr, nil
	case domain:
		return string(p.DST_ADDR), nil
	case ipv6:
		addr, err := ipkit.BytesToIPv6String(p.DST_ADDR)
		if err != nil {
			return "", err
		}
		return addr, nil
	default:
		return "", errors.New("invalid value")
	}
}

func (p SSDetailPacket) GetAddress() (string, error) {
	addr, err := p.GetAddr()
	if err != nil {
		return "", err
	}
	return addr + ":" + strconv.Itoa(p.GetPort()), nil
}
