package httptrans

import (
	"github.com/wsrf16/swiss/netkit/layer/transport/tcptrans"
	"github.com/wsrf16/swiss/sugar/base/regexpkit"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const ConnectEstablished = "HTTP/1.1 200 Connection established\n\n"

type HttpPacket struct {
	Data     []byte
	Protocol string
	Method   string
	Host     string
	Port     int
}

func (p HttpPacket) IsMethodConnect() bool {
	return p.Method == http.MethodConnect
	// if p.Port == 443 {
	//	return true
	// } else {
	//	return false
	// }
}

func (p HttpPacket) GetAddress() string {
	return net.JoinHostPort(p.Host, strconv.Itoa(p.Port))
}

func (p HttpPacket) DialDSTConn() (net.Conn, error) {
	return tcptrans.DialAddress("", p.GetAddress())
}

type ParsePacketToConnFunc func(*HttpPacket) (net.Conn, error)

const RegExpMethods = http.MethodConnect + "|" +
	http.MethodGet + "|" +
	http.MethodHead + "|" +
	http.MethodPost + "|" +
	http.MethodPut + "|" +
	http.MethodPatch + "|" +
	http.MethodDelete + "|" +
	http.MethodOptions + "|" +
	http.MethodTrace

func ResolvePacket(readLine []byte) (*HttpPacket, error) {
	packet := new(HttpPacket)
	packet.Data = readLine

	pattern := "^(" + RegExpMethods + ")" + " (.*) HTTP/.*"
	isHTTPOrHTTPS, err := regexp.Match(pattern, readLine)
	if err != nil {
		return nil, err
	}

	if isHTTPOrHTTPS {
		matches, err := regexpkit.FindStringSubmatch(pattern, string(readLine))
		if err != nil {
			return nil, err
		}
		var headerMethod, headerHost = matches[1], matches[2]
		// fmt.Sscanf(string(readLine), "%s%s", &headerMethod, &headerHost)

		packet.Method = headerMethod
		if packet.IsMethodConnect() {
			packet.Protocol = "HTTPS"
		} else {
			packet.Protocol = "HTTP"
		}

		uri, err := url.Parse(headerHost)
		if err != nil {
			return nil, err
		}

		if uri.Opaque == "443" {
			packet.Host = uri.Scheme
			packet.Port = 443
		} else {
			if strings.Index(uri.Host, ":") == -1 {
				packet.Host = uri.Host
				packet.Port = 80
			} else {
				host, port, err := net.SplitHostPort(uri.Host)
				if err != nil {
					return nil, err
				}
				atoi, err := strconv.Atoi(port)
				if err != nil {
					return nil, err
				}
				packet.Host = host
				packet.Port = atoi
			}
		}
	}

	return packet, nil
}
