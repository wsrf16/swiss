package redirect

import (
	"errors"
	"fmt"
	"github.com/wsrf16/swiss/netkit/tun2socks/ctun2socks/core"
	"github.com/wsrf16/swiss/sugar/logo"
	"net"
	"sync"
	"time"
)

type udpHandler struct {
	sync.Mutex

	timeout        time.Duration
	udpConns       map[core.UDPConn]*net.UDPConn
	udpTargetAddrs map[core.UDPConn]*net.UDPAddr
	target         string
}

func NewUDPHandler(target string, timeout time.Duration) core.UDPConnHandler {
	return &udpHandler{
		timeout:        timeout,
		udpConns:       make(map[core.UDPConn]*net.UDPConn, 8),
		udpTargetAddrs: make(map[core.UDPConn]*net.UDPAddr, 8),
		target:         target,
	}
}

func (h *udpHandler) fetchUDPInput(conn core.UDPConn, pc *net.UDPConn) {
	buf := core.NewBytes(core.BufSize)

	defer func() {
		h.Close(conn)
		core.FreeBytes(buf)
	}()

	for {
		pc.SetDeadline(time.Now().Add(h.timeout))
		n, addr, err := pc.ReadFromUDP(buf)
		if err != nil {
			// log.Printf("failed to read UDP data from remote: %v", err)
			return
		}

		_, err = conn.WriteFrom(buf[:n], addr)
		if err != nil {
			logo.Wf("failed to write UDP data to TUN")
			return
		}
	}
}

func (h *udpHandler) Connect(conn core.UDPConn, target *net.UDPAddr) error {
	bindAddr := &net.UDPAddr{IP: nil, Port: 0}
	pc, err := net.ListenUDP("udp", bindAddr)
	if err != nil {
		logo.Ef("failed to bind udp address")
		return err
	}
	tgtAddr, _ := net.ResolveUDPAddr("udp", h.target)
	h.Lock()
	h.udpTargetAddrs[conn] = tgtAddr
	h.udpConns[conn] = pc
	h.Unlock()
	go h.fetchUDPInput(conn, pc)
	logo.If("new proxy connection for proxyAddress: %s:%s", target.Network(), target.String())
	return nil
}

func (h *udpHandler) ReceiveTo(conn core.UDPConn, data []byte, addr *net.UDPAddr) error {
	h.Lock()
	pc, ok1 := h.udpConns[conn]
	tgtAddr, ok2 := h.udpTargetAddrs[conn]
	h.Unlock()

	if ok1 && ok2 {
		_, err := pc.WriteToUDP(data, tgtAddr)
		if err != nil {
			logo.Wf("failed to write UDP payload to SOCKS5 server: %v", err)
			return errors.New("failed to write UDP data")
		}
		return nil
	} else {
		return errors.New(fmt.Sprintf("proxy connection %v->%v does not exists", conn.LocalAddr(), addr))
	}
}

func (h *udpHandler) Close(conn core.UDPConn) {
	conn.Close()

	h.Lock()
	defer h.Unlock()

	if _, ok := h.udpTargetAddrs[conn]; ok {
		delete(h.udpTargetAddrs, conn)
	}
	if pc, ok := h.udpConns[conn]; ok {
		pc.Close()
		delete(h.udpConns, conn)
	}
}
