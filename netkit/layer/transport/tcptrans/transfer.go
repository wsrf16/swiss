package tcptrans

import (
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"net"
	"os"
	"os/signal"
	"syscall"
)

//type TransferContext struct {
//    LAddr         *net.TCPAddr
//    DAddr         *net.TCPAddr
//    KeepListening bool
//    LListener     *net.TCPListener
//    DListener     *net.TCPListener
//    StopChan      *chan os.Signal
//}

type TCPTransferContext struct {
	LAddr         *net.TCPAddr
	DAddr         *net.TCPAddr
	KeepListening bool
	LListener     *net.TCPListener
	DListener     *net.TCPListener
	StopChan      *chan os.Signal
}

func BuildTransfer(lAddress string, dAddress string, keepListening bool) (*TCPTransferContext, error) {
	lAddr, dAddr, err := NewCoupleTCPAddr(lAddress, dAddress)
	if err != nil {
		return nil, err
	}
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	return &TCPTransferContext{LAddr: lAddr, DAddr: dAddr, KeepListening: keepListening, StopChan: &stopChan}, nil
}

func (t TCPTransferContext) Stop() {
	*t.StopChan <- os.Interrupt
	iokit.Close(t.LListener, t.DListener)
}

func (t TCPTransferContext) GetLListener() *net.TCPListener {
	return t.LListener
}

func (t TCPTransferContext) GetDListener() *net.TCPListener {
	return t.DListener
}

func (t TCPTransferContext) TransferFromListenToDial() error {
	lAddr := t.LAddr
	dAddr := t.DAddr
	keepListening := t.KeepListening

	return TransferFromListenToDial(lAddr, dAddr, keepListening, &t)
}

func (t TCPTransferContext) TransferFromListenToListen() error {
	lAddr := t.LAddr
	dAddr := t.DAddr
	keepListening := t.KeepListening

	return TransferFromListenToListen(lAddr, dAddr, keepListening, &t)
}

func (t TCPTransferContext) TransferFromDialToDial() error {
	lAddr := t.LAddr
	dAddr := t.DAddr
	keepListening := t.KeepListening

	return TransferFromDialToDial(lAddr, dAddr, keepListening)
}
