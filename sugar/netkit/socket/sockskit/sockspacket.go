package sockskit

import (
	"encoding/hex"
	"errors"
	"github.com/wsrf16/swiss/sugar/base/convertkit"
	"github.com/wsrf16/swiss/sugar/netkit"
	"strconv"
)

type SocksPacket1 struct {
	Packet   []byte
	Version  byte
	NMethods byte
	Methods  []byte
}

func ResolvePacket1(packet []byte) (*SocksPacket1, error) {
	packet1 := SocksPacket1{}
	packet1.Packet = packet
	if len(packet) < 4 {
		return nil, errors.New("wrong socks packet.")
	}
	packet1.Version = packet[0]
	packet1.NMethods = packet[1]
	packet1.Methods = packet[2:]
	return &packet1, nil
}

type SocksPacket11 struct {
	Packet    []byte
	Version   byte
	NUserName byte
	UserName  []byte
	NPassword byte
	Password  []byte
}

func ResolvePacket2(packet []byte) *SocksPacket11 {
	packet11 := SocksPacket11{}
	packet11.Packet = packet
	cursor := byte(0)
	packet11.Version = packet[cursor]
	cursor += 1
	packet11.NUserName = packet[cursor]
	cursor += 1
	packet11.UserName = packet[cursor : cursor+packet11.NUserName]
	cursor += packet11.NUserName
	packet11.NPassword = packet[cursor]
	cursor += 1
	packet11.Password = packet[cursor : cursor+packet11.NPassword]
	cursor += packet11.NPassword
	return &packet11
}

// func (p SocksPacket1) IsSocks() bool {
//    if p.Version != 4 && p.Version != 5 {
//        return false
//    }
//    return true
// }

type SocksPacket2 struct {
	Packet   []byte
	Version  byte
	CMD      byte
	RSV      byte
	ATYP     byte
	DST_ADDR []byte
	DST_PORT []byte
}

func BuildPacket3(packet []byte) SocksPacket2 {
	packet2 := SocksPacket2{}
	packet2.Packet = packet
	packet2.Version = packet[0]
	packet2.CMD = packet[1]
	packet2.RSV = packet[2]
	packet2.ATYP = packet[3]
	packet2.DST_ADDR, packet2.DST_PORT = packet2.getAddrAndPort()
	return packet2
}

func (p SocksPacket2) getAddrAndPort() ([]byte, []byte) {
	b := p.Packet[3]
	var length int
	switch b {
	case 0x01:
		length = 4
		return p.Packet[4 : 4+length], p.Packet[4+length : 4+length+2]
	case 0x03:
		length = int(p.Packet[4])
		return p.Packet[5 : 5+length], p.Packet[5+length : 5+length+2]
	case 0x04:
		length = 16
		return p.Packet[4 : 4+length], p.Packet[4+length : 4+length+2]
	}
	return nil, nil
}

func (p SocksPacket2) GetPort() int {
	return int(convertkit.BytesToUint16(p.DST_PORT))
}

func (p SocksPacket2) GetAddr() (string, error) {
	switch p.Packet[3] {
	case 0x01:
		addr, err := netkit.BytesToIPv4String(p.DST_ADDR)
		if err != nil {
			return "", err
		}
		return addr, nil
	case 0x03:
		return string(p.DST_ADDR), nil
	case 0x04:
		addr, err := bytesToIPv6(p.DST_ADDR)
		if err != nil {
			return "", err
		}
		return addr, nil
	default:
		return "", errors.New("invalid value")
	}
}

func (p SocksPacket2) GetAddress() (string, error) {
	addr, err := p.GetAddr()
	if err != nil {
		return "", err
	}
	return addr + ":" + strconv.Itoa(p.GetPort()), nil
}

func bytesToIPv6(b []byte) (string, error) {
	if len(b) != 16 {
		return "", errors.New("invalid data")
	}
	return "[" + hex.EncodeToString(b[0:2]) + "::" + hex.EncodeToString(b[8:10]) + "::" + hex.EncodeToString(b[10:12]) + "::" + hex.EncodeToString(b[12:14]) + "::" + hex.EncodeToString(b[14:16]) + "]", nil
}
