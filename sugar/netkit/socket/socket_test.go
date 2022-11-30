package socket

import (
	"github.com/wsrf16/swiss/sugar/netkit/socket/tcpkit"
	"testing"
)

func TestListenThenAcceptThenTransferToSocket(t *testing.T) {
	lAddr, err := tcpkit.NewTCP4Addr("0.0.0.0:8082")
	if err != nil {
		t.Error(err.Error())
	}
	dAddr, err := tcpkit.NewTCP4Addr("mecs.com:8080")
	if err != nil {
		t.Error(err.Error())
	}

	tcpkit.ListenThenAcceptThenTransferToTCP(lAddr, dAddr, true)
}
