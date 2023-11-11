package gosocks

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

type ClientAuthenticator interface {
	ClientAuthenticate(conn *SocksConn) error
}

type SocksDialer struct {
	Timeout time.Duration
	Auth    ClientAuthenticator
}

type AnonymousClientAuthenticator struct{}

func (a *AnonymousClientAuthenticator) ClientAuthenticate(conn *SocksConn) (err error) {
	conn.SetWriteDeadline(time.Now().Add(conn.Timeout))
	var req [3]byte
	req[0] = SocksVersion
	req[1] = 1
	req[2] = SocksNoAuthentication
	_, err = conn.Write(req[:])
	if err != nil {
		return
	}

	conn.SetReadDeadline(time.Now().Add(conn.Timeout))
	var resp [2]byte
	r := bufio.NewReader(conn)
	_, err = io.ReadFull(r, resp[:2])
	if err != nil {
		return
	}
	if resp[0] != SocksVersion || resp[1] != SocksNoAuthentication {
		err = fmt.Errorf("Fail to pass anonymous authentication: (0x%02x, 0x%02x)", resp[0], resp[1])
		return
	}
	return
}

func (d *SocksDialer) Dial(address string) (socksConn *SocksConn, err error) {
	conn, err := net.DialTimeout("tcp", address, d.Timeout)
	if err != nil {
		return
	}
	socksConn = &SocksConn{conn.(*net.TCPConn), d.Timeout}
	err = d.Auth.ClientAuthenticate(socksConn)
	if err != nil {
		socksConn.Close()
		return
	}
	return

	//proxy.SOCKS5("tcp", address, &proxy.Auth{User: user, Password: password}, proxy.Direct)
	//dial, err := proxy.SOCKS5("tcp", address, nil, proxy.Direct)
	//if err != nil {
	//    return nil, err
	//}
	//tcpConn, err := dial.Dial("tcp", address)
	//if err != nil {
	//    return nil, err
	//}
	//socksConn = &SocksConn{tcpConn.(*net.TCPConn), d.Timeout}
	//return socksConn, nil

}

func ClientAuthAnonymous(conn *SocksConn) (err error) {
	conn.SetWriteDeadline(time.Now().Add(conn.Timeout))
	var req [3]byte
	req[0] = SocksVersion
	req[1] = 1
	req[2] = SocksNoAuthentication
	_, err = conn.Write(req[:])
	if err != nil {
		return
	}

	conn.SetReadDeadline(time.Now().Add(conn.Timeout))
	var resp [2]byte
	r := bufio.NewReader(conn)
	_, err = io.ReadFull(r, resp[:2])
	if err != nil {
		return
	}
	if resp[0] != SocksVersion || resp[1] != SocksNoAuthentication {
		err = fmt.Errorf("Fail to pass anonymous authentication: (0x%02x, 0x%02x)", resp[0], resp[1])
		return
	}
	return
}

func ClientRequest(conn *SocksConn, req *SocksRequest) (reply *SocksReply, err error) {
	conn.SetWriteDeadline(time.Now().Add(conn.Timeout))
	_, err = WriteSocksRequest(conn, req)
	if err != nil {
		return
	}
	conn.SetReadDeadline(time.Now().Add(conn.Timeout))
	reply, err = ReadSocksReply(conn)
	return
}
