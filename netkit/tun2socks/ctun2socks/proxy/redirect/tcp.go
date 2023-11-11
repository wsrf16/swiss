package redirect

import (
	"github.com/wsrf16/swiss/netkit/tun2socks/ctun2socks/core"
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"github.com/wsrf16/swiss/sugar/logo"
	"io"
	"net"
)

// To do a benchmark using iperf3 locally, you may follow these steps:
//
//  1. Setup and configure the TUN device and start tun2socks with the
//     redirect handler using the following command:
//     tun2socks -proxyType redirect -proxyServer 127.0.0.1:1234
//     Tun2socks will redirect all traffic to 127.0.0.1:1234.
//
//  2. Route traffic targeting 1.2.3.4 to the TUN interface (240.0.0.1):
//     route add 1.2.3.4/32 240.0.0.1
//
//  3. Run iperf3 server locally and listening on 1234 port:
//     iperf3 -s -p 1234
//
//  4. Run iperf3 client locally and connect to 1.2.3.4:1234:
//     iperf3 -c 1.2.3.4 -p 1234
//
// It works this way:
// iperf3 client -> 1.2.3.4:1234 -> routing table -> TUN (240.0.0.1) -> tun2socks -> tun2socks redirect anything to 127.0.0.1:1234 -> iperf3 server
type tcpHandler struct {
	proxyAddress string
}

type duplexConn interface {
	net.Conn
	CloseWrite() error
	CloseRead() error
}

func NewTCPHandler(proxyAddress string) core.TCPConnHandler {
	return &tcpHandler{proxyAddress: proxyAddress}
}

func (h *tcpHandler) handleInput(conn net.Conn, input io.ReadCloser) {
	defer func() {
		if tcpConn, ok := conn.(core.TCPConn); ok {
			tcpConn.CloseWrite()
		} else {
			conn.Close()
		}
		if tcpInput, ok := input.(duplexConn); ok {
			tcpInput.CloseRead()
		} else {
			input.Close()
		}
	}()

	io.Copy(conn, input)
}

func (h *tcpHandler) handleOutput(conn net.Conn, output io.WriteCloser) {
	defer func() {
		if tcpConn, ok := conn.(core.TCPConn); ok {
			tcpConn.CloseRead()
		} else {
			conn.Close()
		}
		if tcpOutput, ok := output.(duplexConn); ok {
			tcpOutput.CloseWrite()
		} else {
			output.Close()
		}
	}()

	io.Copy(output, conn)
}

func (h *tcpHandler) HandleOld(conn net.Conn, target *net.TCPAddr) error {
	c, err := net.Dial("tcp", h.proxyAddress)
	if err != nil {
		return err
	}
	go h.handleInput(conn, c)
	go h.handleOutput(conn, c)
	logo.If("new proxy connection for proxyAddress: %s:%s", target.Network(), target.String())
	return nil
}

func (h *tcpHandler) Handle(conn net.Conn, targetAddress string) error {
	c, err := net.Dial("tcp", h.proxyAddress)
	if err != nil {
		return err
	}
	// ？？？？？？？？？？？？？？？？
	go h.handleInput(conn, c)
	go h.handleOutput(conn, c)
	go iokit.CopyDuplex(c, conn, true)
	logo.If("new proxy connection for proxyAddress: %s:%s", "tcp", targetAddress)
	return nil
}
