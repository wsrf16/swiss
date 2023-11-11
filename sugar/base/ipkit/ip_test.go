package ipkit

import (
	"github.com/wsrf16/swiss/netkit/layer/network/icmpkit"
	"net"
	"testing"
	"time"
)

func TestConvert(t *testing.T) {
	// 2130706433 = 127.0.0.1
	ipInt := IPv4ToInt("1.1.1.1")
	t.Log(ipInt)
	ipString := IntToIPv4(ipInt)
	t.Log(ipString)
}

func TestBetween(t *testing.T) {
	between1 := Between("1.1.1.1", "1.1.1.0", "1.1.1.3")
	between2 := Between("1.1.1.1", "1.1.1.2", "1.1.1.3")
	t.Log(between1)
	t.Log(between2)
}

func TestBeInSegment(t *testing.T) {
	beInSegment1 := InSegment("1.1.1.1", "1.1.1.0-1.1.1.3")
	beInSegment2 := InSegment("1.1.1.1", "1.1.1.2-1.1.1.3")
	t.Log(beInSegment1)
	t.Log(beInSegment2)
}

func TestPing(t *testing.T) {
	msg, err := icmpkit.Ping("163.com", time.Second*2)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(msg)
	}
}

func TestIPv4Ipv6(t *testing.T) {
	b4 := []byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}
	ipv4, err := BytesToIPv6String(b4)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(ipv4)
	}

	b6 := []byte{1, 2, 3, 4}
	ipv6, err := BytesToIPv4String(b6)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(ipv6)
	}
}

func TestParseCIDR(t *testing.T) {
	// -------------------------
	ipv4Addr, ipv4Net, err := net.ParseCIDR("192.0.2.1/24")
	ipv4Net = ipv4Net
	ipv4Addr = ipv4Addr
	if err != nil {
		t.Error(err)
	}
}
