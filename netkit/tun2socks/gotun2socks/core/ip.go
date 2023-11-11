package core

import (
	"github.com/wsrf16/swiss/netkit/tun2socks/gotun2socks/internal/packet"
	"log"
	"net"
)

type ipPacket struct {
	ip     *packet.IPv4
	mtuBuf []byte
	wire   []byte
}

var (
	frags = make(map[uint16]*ipPacket)
)

func procFragment(ip *packet.IPv4, raw []byte) (bool, *packet.IPv4, []byte) {
	exist, ok := frags[ip.Id]
	if !ok {
		if ip.Flags&0x1 == 0 {
			return false, nil, nil
		}
		// first
		log.Printf("first fragment of IPID %d", ip.Id)
		dup := make([]byte, len(raw))
		copy(dup, raw)
		clone := &packet.IPv4{}
		packet.ParseIPv4(dup, clone)
		frags[ip.Id] = &ipPacket{
			ip:   clone,
			wire: dup,
		}
		return false, clone, dup
	} else {
		exist.wire = append(exist.wire, ip.Payload...)
		packet.ParseIPv4(exist.wire, exist.ip)

		last := false
		if ip.Flags&0x1 == 0 {
			log.Printf("last fragment of IPID %d", ip.Id)
			last = true
		} else {
			log.Printf("continue fragment of IPID %d", ip.Id)
		}

		return last, exist.ip, exist.wire
	}
}

func genFragments(first *packet.IPv4, offset uint16, data []byte) []*ipPacket {
	var ret []*ipPacket
	for {
		frag := packet.GetIPv4()

		frag.Version = 4
		frag.Id = first.Id
		frag.SrcIP = make(net.IP, len(first.SrcIP))
		copy(frag.SrcIP, first.SrcIP)
		frag.DstIP = make(net.IP, len(first.DstIP))
		copy(frag.DstIP, first.DstIP)
		frag.TTL = first.TTL
		frag.Protocol = first.Protocol
		frag.FragOffset = offset
		if len(data) <= MTU-20 {
			frag.Payload = data
		} else {
			frag.Flags = 1
			offset += (MTU - 20) / 8
			frag.Payload = data[:MTU-20]
			data = data[MTU-20:]
		}

		pkt := &ipPacket{ip: frag}
		pkt.mtuBuf = GetMTU()

		payloadL := len(frag.Payload)
		payloadStart := MTU - payloadL
		if payloadL != 0 {
			copy(pkt.mtuBuf[payloadStart:], frag.Payload)
		}
		ipHL := frag.HeaderLength()
		ipStart := payloadStart - ipHL
		frag.Serialize(pkt.mtuBuf[ipStart:payloadStart], payloadL)
		pkt.wire = pkt.mtuBuf[ipStart:]
		ret = append(ret, pkt)

		if frag.Flags == 0 {
			return ret
		}
	}
}

func putIPPacket(pkt *ipPacket) {
	packet.PutIPv4(pkt.ip)
	if pkt.mtuBuf != nil {
		PutMTU(pkt.mtuBuf)
	}
	pkt.mtuBuf = nil
	pkt.wire = nil
}
