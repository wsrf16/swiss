package sockstrans

import (
	"errors"
	"github.com/wsrf16/swiss/sugar/base/convertkit"
	"github.com/wsrf16/swiss/sugar/base/ipkit"
	"golang.org/x/net/proxy"
	"net"
	"strconv"
)

type SocksServer struct {
	Conn net.Conn
	Auth *proxy.Auth
}

const (
	UserPassword = 0x02
	None         = 0x00
)

func (s SocksServer) replyWhetherAuthMethod(authentication bool) {
	if authentication {
		s.Conn.Write([]byte{0x05, UserPassword})
	} else {
		s.Conn.Write([]byte{0x05, None})
	}
}

func checkPacket(authPacket *SocksAuthenticationPacket, auth *proxy.Auth) bool {
	if auth == nil || len(auth.User) < 1 || len(auth.Password) < 1 {
		return true
	}

	if authPacket == nil || len(authPacket.UserName) < 1 || len(authPacket.Password) < 1 {
		return false
	}

	if string(authPacket.UserName) == auth.User && string(authPacket.Password) == auth.Password {
		return true
	}

	return false
}

/**
   The localConn connects to the dstServer, and sends a ver
   identifier/method selection message:
	          +----+----------+----------+
	          |VER | NMETHODS | METHODS  |
	          +----+----------+----------+
	          | 1  |    1     | 1 to 255 |
	          +----+----------+----------+
   The VER field is set to X'05' for this ver of the protocol.  The
   NMETHODS field contains the number of method identifier octets that
   appear in the METHODS field.
*/
// 第一个字段VER代表Socks的版本，Socks5默认为0x05，其固定长度为1个字节
type SocksBasicPacket struct {
	Raw      []byte
	Version  byte
	NMethods byte
	Methods  []byte
}

type SocksAuthenticationPacket struct {
	Packet    []byte
	Version   byte
	NUserName byte
	UserName  []byte
	NPassword byte
	Password  []byte
}

func parseToAuthenticationPacket(bytes []byte) *SocksAuthenticationPacket {
	packet := SocksAuthenticationPacket{}
	packet.Packet = bytes
	cursor := byte(0)
	packet.Version = bytes[cursor]
	cursor += 1
	packet.NUserName = bytes[cursor]
	cursor += 1
	packet.UserName = bytes[cursor : cursor+packet.NUserName]
	cursor += packet.NUserName
	packet.NPassword = bytes[cursor]
	cursor += 1
	packet.Password = bytes[cursor : cursor+packet.NPassword]
	cursor += packet.NPassword
	return &packet
}

// func (p SocksBasicPacket) IsSocks() bool {
//    if p.Version != 4 && p.Version != 5 {
//        return false
//    }
//    return true
// }

/*
*

	+----+-----+-------+------+----------+----------+
	|VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	+----+-----+-------+------+----------+----------+
	| 1  |  1  | X'00' |  1   | Variable |    2     |
	+----+-----+-------+------+----------+----------+
*/
type SocksDetailPacket struct {
	Packet   []byte
	Version  byte
	CMD      byte
	RSV      byte
	ATYP     byte
	DST_ADDR []byte
	DST_PORT []byte
}

const (
	ipv4   = 0x01
	domain = 0x03
	ipv6   = 0x04
)

func parseToDetailPacket(packet []byte) SocksDetailPacket {
	packet2 := SocksDetailPacket{}
	packet2.Packet = packet
	packet2.Version = packet[0]
	packet2.CMD = packet[1]
	packet2.RSV = packet[2]
	packet2.ATYP = packet[3]
	packet2.DST_ADDR, packet2.DST_PORT = getAddrAndPort(packet)
	return packet2
}

func getAddrAndPort(packet []byte) ([]byte, []byte) {
	var length int
	atyp := packet[3]
	switch atyp {
	case ipv4:
		length = net.IPv4len
		return packet[4 : 4+length], packet[4+length : 4+length+2]
	case domain:
		length = int(packet[4])
		return packet[5 : 5+length], packet[5+length : 5+length+2]
	case ipv6:
		length = net.IPv6len
		return packet[4 : 4+length], packet[4+length : 4+length+2]
	}
	return nil, nil
}

func (p SocksDetailPacket) GetPort() int {
	return int(convertkit.BytesToUint16(p.DST_PORT))
}

func (p SocksDetailPacket) GetAddr() (string, error) {
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

func (p SocksDetailPacket) GetAddress() (string, error) {
	addr, err := p.GetAddr()
	if err != nil {
		return "", err
	}
	return addr + ":" + strconv.Itoa(p.GetPort()), nil
}
