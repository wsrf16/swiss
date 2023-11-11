package sstcptrans

import (
	"github.com/wsrf16/swiss/netkit/layer/transport/sstrans"
	"github.com/wsrf16/swiss/netkit/layer/transport/tcptrans"
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type SSTransferContext struct {
	LAddr         *net.TCPAddr
	DAddr         *net.TCPAddr
	KeepListening bool
	LListener     *net.TCPListener
	DListener     *net.TCPListener
	StopChan      *chan os.Signal
}

func BuildTransfer(lAddress string, keepListening bool) (*SSTransferContext, error) {
	lAddr, err := tcptrans.NewTCPAddr(lAddress)
	if err != nil {
		return nil, err
	}
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	return &SSTransferContext{LAddr: lAddr, KeepListening: keepListening, StopChan: &stopChan}, nil
}

func (t SSTransferContext) Stop() {
	*t.StopChan <- os.Interrupt
	iokit.Close(t.LListener, t.DListener)
}

func (t SSTransferContext) GetLListener() *net.TCPListener {
	return t.LListener
}

func (t SSTransferContext) GetDListener() *net.TCPListener {
	return t.DListener
}

func (t SSTransferContext) TransferFromListen(config *sstrans.ShadowSocksConfig) error {
	lAddr := t.LAddr
	keepListening := t.KeepListening

	return TransferFromListen(lAddr, config, keepListening, &t)
}
