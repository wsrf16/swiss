package sockskit

import (
	"encoding/hex"
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

func BuildPacket1(packet []byte) SocksPacket1 {
	packet1 := SocksPacket1{}
	packet1.Packet = packet
	packet1.Version = packet[0]
	packet1.NMethods = packet[1]
	packet1.Methods = packet[2:]
	return packet1
}

type SocksPacket2 struct {
	Packet   []byte
	Version  byte
	CMD      byte
	RSV      byte
	ATYP     byte
	DST_ADDR []byte
	DST_PORT []byte
}

func BuildPacket2(packet []byte) SocksPacket2 {
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

func (p SocksPacket2) GetAddr() string {
	switch p.Packet[3] {
	case 0x01:
		return netkit.BytesToIPv4(p.DST_ADDR)
	case 0x03:
		return string(p.DST_ADDR)
	case 0x04:
		return bytesToIPv6(p.DST_ADDR)
	}
	return ""
}

func (p SocksPacket2) GetAddress() string {
	return p.GetAddr() + ":" + strconv.Itoa(p.GetPort())
}

func bytesToIPv6(b []byte) string {
	return "[" + hex.EncodeToString(b[0:2]) + "::" + hex.EncodeToString(b[8:10]) + "::" + hex.EncodeToString(b[10:12]) + "::" + hex.EncodeToString(b[12:14]) + "::" + hex.EncodeToString(b[14:16]) + "]"
}
