package netkit

import (
	"net"
	"time"
)

func checkSum(msg []byte) uint16 {
	sum := 0
	for n := 0; n < len(msg); n += 2 {
		sum += int(msg[n])<<8 + int(msg[n+1])
	}
	sum = (sum >> 16) + sum&0xffff
	sum += sum >> 16
	return uint16(^sum)
}

func req() []byte {
	msg := make([]byte, 8)
	msg[0] = 8
	msg[1] = 0
	msg[2] = 0
	msg[3] = 0
	msg[4] = 0
	msg[5] = 13
	msg[6] = 0
	msg[7] = 37
	sum := checkSum(msg)
	msg[2] = byte(sum >> 8)
	msg[3] = byte(sum & 0xff)
	return msg
}

const IcmpLen = 8

var REQ = req()

func Ping3s(host string) ([]byte, error) {
	return Ping(host, 3*time.Second)
}

func Ping(host string, timeout time.Duration) ([]byte, error) {
	conn, err := net.Dial("ip4:icmp", host)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write(REQ)
	if err != nil {
		return nil, err
	}
	// 请求超时
	conn.SetReadDeadline(time.Now().Add(timeout))

	msg := make([]byte, 512)
	n, err := conn.Read(msg)
	if err != nil {
		return nil, err
	}

	return msg[:n], nil
}
