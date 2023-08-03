package netkit

import (
	"github.com/wsrf16/swiss/sugar/netkit/icmpkit"
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
		t.Log(err)
	} else {
		t.Log(msg)
	}
}
