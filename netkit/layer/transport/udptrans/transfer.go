package udptrans

import (
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

type UDPTransferContext struct {
	LAddr         *net.UDPAddr
	DAddr         *net.UDPAddr
	KeepListening bool
	//LListener     *net.TCPListener
	//DListener     *net.TCPListener
	StopChan *chan os.Signal
}

func BuildTransfer(lAddress string, dAddress string, keepListening bool) (*UDPTransferContext, error) {
	lAddr, dAddr, err := NewCoupleUDPAddr(lAddress, dAddress)
	if err != nil {
		return nil, err
	}
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	return &UDPTransferContext{LAddr: lAddr, DAddr: dAddr, KeepListening: keepListening, StopChan: &stopChan}, nil
}

func (t UDPTransferContext) Stop() {
	*t.StopChan <- os.Interrupt
	//iokit.Close(t.LListener, t.DListener)
}

//func (t UDPTransferContext) GetLListener() *net.TCPListener {
//    return t.LListener
//}
//
//func (t UDPTransferContext) GetDListener() *net.TCPListener {
//    return t.DListener
//}

func (t UDPTransferContext) TransferFromListenToDial() error {
	lAddr := t.LAddr
	dAddr := t.DAddr
	keepListening := t.KeepListening

	return TransferFromListenToDial(lAddr, dAddr, keepListening, &t)
}

func (t UDPTransferContext) TransferFromListenToListen() error {
	lAddr := t.LAddr
	dAddr := t.DAddr
	keepListening := t.KeepListening

	return TransferFromListenToListen(lAddr, dAddr, keepListening, &t)
}

func (t UDPTransferContext) TransferFromDialToDial() error {
	lAddr := t.LAddr
	dAddr := t.DAddr
	keepListening := t.KeepListening

	return TransferFromDialToDial(lAddr, dAddr, keepListening, &t)
}
