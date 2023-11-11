package ipkit

import (
	"errors"
	"fmt"
	"github.com/wsrf16/swiss/sugar/base/stringkit"
	"net"
	"regexp"
	"strconv"
	"strings"
)

func GetHostIp() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.To4(), nil
			}
		}
	}

	return nil, errors.New("no valid ipv4 address founded")
}

func BytesToIPv4(b []byte) (net.IP, error) {
	//return net.IPv4(b[0], b[1], b[2], b[3])
	if len(b) != net.IPv4len {
		return nil, errors.New("invalid data")
	}
	return net.IP(b), nil
}

func BytesToIPv4String(b []byte) (string, error) {
	//return strconv.Itoa(int(b[0])) + "." + strconv.Itoa(int(b[1])) + "." + strconv.Itoa(int(b[2])) + "." + strconv.Itoa(int(b[3])), nil

	//ipv4 := fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3])
	//return ipv4, nil
	//return net.IPv4(b[0], b[1], b[2], b[3]).String(), nil
	ipv4, err := BytesToIPv4(b)
	if err != nil {
		return "", err
	} else {
		return ipv4.String(), nil
	}
}

func BytesToIPv6(b []byte) (net.IP, error) {
	//return net.IPv4(b[0], b[1], b[2], b[3])
	if len(b) != net.IPv6len {
		return nil, errors.New("invalid data")
	}
	return net.IP(b), nil
}

func BytesToIPv6String(b []byte) (string, error) {
	// [255 255 :: 255 255 :: 255 255 :: 255 255 :: 255 255]: 65535
	// [FEDC:BA48:7654:3210:FEDC:BA98:7654:3211]: 65535
	//return hex.EncodeToString(b[0:2]) + "::" + hex.EncodeToString(b[2:4]) + "::" + hex.EncodeToString(b[4:6]) + "::" + hex.EncodeToString(b[6:8]) + "::" + hex.EncodeToString(b[8:10])
	//return net.IP(b), nil
	ipv4, err := BytesToIPv6(b)
	if err != nil {
		return "", err
	} else {
		return ipv4.String(), nil
	}
}

func IPv4ToInt(ip string) uint32 {
	if len(ip) == 0 {
		return 0
	}
	bits := strings.Split(ip, ".")
	if len(bits) < 4 {
		return 0
	}
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum uint32
	sum += uint32(b0) << 24
	sum += uint32(b1) << 16
	sum += uint32(b2) << 8
	sum += uint32(b3)

	return sum
}

// 2130706433=127.0.0.1
func IntToIPv4(v uint32) string {
	b0 := byte(v >> 24)
	b1 := byte(v >> 16)
	b2 := byte(v >> 8)
	b3 := byte(v)
	return fmt.Sprintf("%d.%d.%d.%d", b0, b1, b2, b3)
}

func ResolveIP(ip string) (net.IP, error) {
	if len(ip) == 0 {
		return nil, errors.New("wrong format ip")
	}
	bits := strings.Split(ip, ".")
	if len(bits) < 4 {
		return nil, errors.New("wrong format ip")
	}
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])
	ipv4 := net.IPv4(byte(b0), byte(b1), byte(b2), byte(b3))
	return ipv4, nil
}

func Between(ip string, from string, to string) bool {
	return IPv4ToInt(ip) >= IPv4ToInt(from) && IPv4ToInt(ip) <= IPv4ToInt(to)
}

func InSegment(ip string, segment string) bool {
	split := stringkit.SplitPath(segment, "-")
	from, to := split[0], split[1]
	return IPv4ToInt(ip) >= IPv4ToInt(from) && IPv4ToInt(ip) <= IPv4ToInt(to)
}

func IsPort(port string) bool {
	PortNum, err := strconv.Atoi(port)
	if err != nil {
		return false
		// log.Fatalln("[x]", "port should be a number")
	}
	if PortNum < 1 || PortNum > 65535 {
		return false
		// log.Fatalln("[x]", "port should be a number and the range is [1,65536)")
	}
	return true
}

func IsIPV4(ip string) bool {
	pattern := `^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$`
	ok, err := regexp.MatchString(pattern, ip)
	if err != nil || !ok {
		return false
		// log.Fatalln("[x]", "ip error. ")
	}
	return true
}

func GetIPv4Gateway(ip string) (string, error) {
	if len(ip) < 1 {
		return ip, errors.New("invalid ip")
	}
	bits := strings.Split(ip, ".")
	if len(bits) < 4 {
		return ip, errors.New("invalid ip")
	}
	return fmt.Sprintf("%v.%v.%v.%v", bits[0], bits[1], bits[2], 1), nil
}
