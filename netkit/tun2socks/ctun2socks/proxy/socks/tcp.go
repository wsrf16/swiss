package socks

import (
	"github.com/wsrf16/swiss/netkit/layer/transport/sockstrans"
	"github.com/wsrf16/swiss/netkit/tun2socks/ctun2socks/core"
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"github.com/wsrf16/swiss/sugar/logo"
	"net"
	"sync"
)

type tcpHandler struct {
	sync.Mutex

	//proxyHost string
	//proxyPort uint16

	proxyAddress string
}

//	func NewTCPHandler(proxyHost string, proxyPort uint16) core.TCPConnHandler {
//	   return &tcpHandler{
//	       proxyHost: proxyHost,
//	       proxyPort: proxyPort,
//	   }
//	}
func NewTCPHandler(proxyAddress string) core.TCPConnHandler {
	return &tcpHandler{proxyAddress: proxyAddress}
}

//type direction byte
//
//const (
//    dirUplink direction = iota
//    dirDownlink
//)
//
//type duplexConn interface {
//    net.Conn
//    CloseRead() error
//    CloseWrite() error
//}

//func (h *tcpHandler) relay(lhs, rhs net.Conn) {
//    upCh := make(chan struct{})
//
//    cls := func(dir direction, interrupt bool) {
//        lhsDConn, lhsOk := lhs.(duplexConn)
//        rhsDConn, rhsOk := rhs.(duplexConn)
//        if !interrupt && lhsOk && rhsOk {
//            switch dir {
//            case dirUplink:
//                lhsDConn.CloseRead()
//                rhsDConn.CloseWrite()
//            case dirDownlink:
//                lhsDConn.CloseWrite()
//                rhsDConn.CloseRead()
//            default:
//                panic("unexpected direction")
//            }
//        } else {
//            lhs.Close()
//            rhs.Close()
//        }
//    }
//
//    // Uplink
//    go func() {
//        var err error
//        _, err = io.Copy(rhs, lhs)
//        if err != nil {
//            cls(dirUplink, true) // interrupt the conn if the error is not nil (not EOF)
//        } else {
//            cls(dirUplink, false) // half close uplink direction of the TCP conn if possible
//        }
//        upCh <- struct{}{}
//    }()
//
//    // Downlink
//    var err error
//    _, err = io.Copy(lhs, rhs)
//    if err != nil {
//        cls(dirDownlink, true)
//    } else {
//        cls(dirDownlink, false)
//    }
//
//    <-upCh // Wait for uplink done.
//}

//func (h *tcpHandler) HandleOld(conn net.Conn, target *net.TCPAddr) error {
//    dialer, err := proxy.SOCKS5("tcp", h.proxyAddress, nil, nil)
//    if err != nil {
//        return err
//    }
//
//    c, err := dialer.Dial(target.Network(), target.String())
//    if err != nil {
//        return err
//    }
//
//    //go h.relay(conn, c)
//    go transport.CopyDuplex(conn, c, true)
//
//    logo.If("new proxy connection to %v", target)
//
//    return nil
//}
//
//func (h *tcpHandler) HandleNew(conn net.Conn, target *net.TCPAddr) error {
//    c, err := sockskit.Dial("tcp", target.String(), h.proxyAddress, nil, nil)
//    if err != nil {
//        return err
//    }
//    go transport.CopyDuplex(conn, c, true)
//
//    return nil
//}

func (h *tcpHandler) Handle(conn net.Conn, targetAddress string) error {
	c, err := sockstrans.Dial("tcp", targetAddress, h.proxyAddress, nil, nil)
	if err != nil {
		return err
	}
	go iokit.CopyDuplex(conn, c, true)

	logo.If("new proxy connection to %v", targetAddress)

	return nil
}
