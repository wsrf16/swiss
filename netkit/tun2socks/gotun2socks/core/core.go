package core

import (
	"github.com/wsrf16/swiss/netkit/tun2socks/gotun2socks/gosocks"
	"github.com/wsrf16/swiss/netkit/tun2socks/gotun2socks/internal/packet"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const (
	MTU = 1500
)

var (
	socksDialer *gosocks.SocksDialer = &gosocks.SocksDialer{
		Auth:    &gosocks.AnonymousClientAuthenticator{},
		Timeout: 1 * time.Second,
	}

	_, ip1, _ = net.ParseCIDR("10.0.0.0/8")
	_, ip2, _ = net.ParseCIDR("172.16.0.0/12")
	_, ip3, _ = net.ParseCIDR("192.168.0.0/24")
)

type Tun2Socks struct {
	dev          io.ReadWriteCloser
	socksAddress string
	publicOnly   bool

	writerStopCh chan bool
	writeCh      chan interface{}

	tcpConnTrackLock sync.Mutex
	tcpConnTrackMap  map[string]*tcpConnTrack

	udpConnTrackLock sync.Mutex
	udpConnTrackMap  map[string]*udpConnTrack

	dnsAddresses []string
	cache        *dnsCache

	wg sync.WaitGroup
}

func isPrivate(ip net.IP) bool {
	return ip1.Contains(ip) || ip2.Contains(ip) || ip3.Contains(ip)
}

func dialSocksServer(socksAddress string) (*gosocks.SocksConn, error) {
	return socksDialer.Dial(socksAddress)
}

func New(dev io.ReadWriteCloser, socksAddress string, dnsAddresses []string, publicOnly bool, enableDNSCache bool) *Tun2Socks {
	t2s := &Tun2Socks{
		dev:             dev,
		socksAddress:    socksAddress,
		publicOnly:      publicOnly,
		writerStopCh:    make(chan bool, 10),
		writeCh:         make(chan interface{}, 10000),
		tcpConnTrackMap: make(map[string]*tcpConnTrack),
		udpConnTrackMap: make(map[string]*udpConnTrack),
		dnsAddresses:    dnsAddresses,
	}
	if enableDNSCache {
		t2s.cache = &dnsCache{
			storage: make(map[string]*dnsCacheEntry),
		}
	}
	return t2s
}

func (t2s *Tun2Socks) Stop() {
	t2s.writerStopCh <- true
	//b := t2s.dev == nil
	//fmt.Println("2222222" + strconv.FormatBool(b))

	t2s.dev.Close()

	t2s.tcpConnTrackLock.Lock()
	defer t2s.tcpConnTrackLock.Unlock()
	for _, tcpTrack := range t2s.tcpConnTrackMap {
		close(tcpTrack.quitByOther)
	}

	t2s.udpConnTrackLock.Lock()
	defer t2s.udpConnTrackLock.Unlock()
	for _, udpTrack := range t2s.udpConnTrackMap {
		close(udpTrack.quitByOther)
	}
	t2s.wg.Wait()
}

func (t2s *Tun2Socks) Run() {
	defer t2s.Stop()

	// writer
	go func() {
		t2s.wg.Add(1)
		defer t2s.wg.Done()
		for {
			select {
			case pkt := <-t2s.writeCh:
				switch pkt.(type) {
				case *tcpPacket:
					tcp := pkt.(*tcpPacket)
					t2s.dev.Write(tcp.wire)
					putTCPPacket(tcp)
				case *udpPacket:
					udp := pkt.(*udpPacket)
					t2s.dev.Write(udp.wire)
					putUDPPacket(udp)
				case *ipPacket:
					ip := pkt.(*ipPacket)
					t2s.dev.Write(ip.wire)
					putIPPacket(ip)
				}
			case <-t2s.writerStopCh:
				log.Printf("quit tun2socks writer")
				return
			}
		}
	}()

	// reader
	var buf [MTU]byte
	var ip packet.IPv4
	var tcp packet.TCP
	var udp packet.UDP

	t2s.wg.Add(1)
	defer t2s.wg.Done()
	for {
		n, err := t2s.dev.Read(buf[:])
		if err != nil {
			// TODO: stop at critical error
			log.Printf("read packet error: %s", err)
			return
		}
		data := buf[:n]
		err = packet.ParseIPv4(data, &ip)
		if err != nil {
			log.Printf("error to parse IPv4: %s", err)
			continue
		}
		if t2s.publicOnly {
			if !ip.DstIP.IsGlobalUnicast() {
				continue
			}
			if isPrivate(ip.DstIP) {
				continue
			}
		}

		if ip.Flags&0x1 != 0 || ip.FragOffset != 0 {
			last, pkt, raw := procFragment(&ip, data)
			if last {
				ip = *pkt
				data = raw
			} else {
				continue
			}
		}

		switch ip.Protocol {
		case packet.IPProtocolTCP:
			err = packet.ParseTCP(ip.Payload, &tcp)
			if err != nil {
				log.Printf("error to parse TCP: %s", err)
				continue
			}
			t2s.tcp(data, &ip, &tcp)

		case packet.IPProtocolUDP:
			err = packet.ParseUDP(ip.Payload, &udp)
			if err != nil {
				log.Printf("error to parse UDP: %s", err)
				continue
			}
			t2s.udp(data, &ip, &udp)

		default:
			// Unsupported packets
			log.Printf("Unsupported packet: protocol %d", ip.Protocol)
		}
	}
}
