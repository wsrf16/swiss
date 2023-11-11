package ssudp

import (
	"errors"
	"github.com/wsrf16/swiss/netkit/layer/transport/udptrans"
	"github.com/wsrf16/swiss/sugar/base/gokit"
	"github.com/wsrf16/swiss/sugar/base/timekit"
	"io"
	"log"
	"net"
	"strconv"
)

const MaxAddrLen = 1 + 1 + 255 + 2

const (
	ATYP_IPv4        = 1
	ATYP_DOMAIN_NAME = 3
	ATYP_IPV6        = 4
)

func toString(a []byte) string {
	var host, port string

	switch a[0] { // address type
	case ATYP_DOMAIN_NAME:
		host = string(a[2 : 2+int(a[1])])
		port = strconv.Itoa((int(a[2+int(a[1])]) << 8) | int(a[2+int(a[1])+1]))
	case ATYP_IPv4:
		host = net.IP(a[1 : 1+net.IPv4len]).String()
		port = strconv.Itoa((int(a[1+net.IPv4len]) << 8) | int(a[1+net.IPv4len+1]))
	case ATYP_IPV6:
		host = net.IP(a[1 : 1+net.IPv6len]).String()
		port = strconv.Itoa((int(a[1+net.IPv6len]) << 8) | int(a[1+net.IPv6len+1]))
	}

	return net.JoinHostPort(host, port)
}

func parseDestAddress(src net.Conn) (string, error) {
	b := make([]byte, MaxAddrLen)
	if len(b) < MaxAddrLen {
		return "", io.ErrShortBuffer
	}
	_, err := io.ReadFull(src, b[:1]) // read 1st byte for address type
	if err != nil {
		return "", err
	}

	switch b[0] {
	case ATYP_DOMAIN_NAME:
		_, err = io.ReadFull(src, b[1:2]) // read 2nd byte for domain length
		if err != nil {
			return "", err
		}
		_, err = io.ReadFull(src, b[2:2+int(b[1])+2])
		return toString(b[:1+1+int(b[1])+2]), err
	case ATYP_IPv4:
		_, err = io.ReadFull(src, b[1:1+net.IPv4len+2])
		return toString(b[:1+net.IPv4len+2]), err
	case ATYP_IPV6:
		_, err = io.ReadFull(src, b[1:1+net.IPv6len+2])
		return toString(b[:1+net.IPv6len+2]), err
	default:
		return "", errors.New("address not support")
	}
}

func parseDest(src net.Conn) (net.Conn, error) {
	destAddress, err := parseDestAddress(src)
	if err != nil {
		return nil, err
	}

	dst, err := udptrans.DialAddress("", destAddress)
	if err != nil {
		return nil, err
	} else {
		log.Printf("shadowsocks udp proxy(go-routine-id: %v): {from: %s <-> %s to: %s <-> %s}\n", gokit.GetGid(), src.LocalAddr(), src.RemoteAddr(), dst.LocalAddr(), destAddress)
	}
	dst.SetDeadline(timekit.Time1Year())

	return dst, err
}
