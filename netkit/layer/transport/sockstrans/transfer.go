package sockstrans

import (
	"github.com/wsrf16/swiss/netkit/layer/transport/tcptrans"
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"golang.org/x/net/proxy"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type SocksTransferContext struct {
	LAddr         *net.TCPAddr
	DAddr         *net.TCPAddr
	KeepListening bool
	LListener     *net.TCPListener
	DListener     *net.TCPListener
	StopChan      *chan os.Signal
}

func BuildTransfer(lAddress string, keepListening bool) (*SocksTransferContext, error) {
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

	return &SocksTransferContext{LAddr: lAddr, KeepListening: keepListening, StopChan: &stopChan}, nil
}

func (t SocksTransferContext) Stop() {
	*t.StopChan <- os.Interrupt
	iokit.Close(t.LListener, t.DListener)
}

func (t SocksTransferContext) GetLListener() *net.TCPListener {
	return t.LListener
}

func (t SocksTransferContext) GetDListener() *net.TCPListener {
	return t.DListener
}

func (t SocksTransferContext) TransferFromListen(auth *proxy.Auth) error {
	lAddr := t.LAddr
	keepListening := t.KeepListening

	return TransferFromListen(lAddr, auth, keepListening, &t)
}
