package netkit

import "testing"

func TestConvert(t *testing.T) {
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
	beInSegment1 := BeInSegment("1.1.1.1", "1.1.1.0-1.1.1.3")
	beInSegment2 := BeInSegment("1.1.1.1", "1.1.1.2-1.1.1.3")
	t.Log(beInSegment1)
	t.Log(beInSegment2)
}
