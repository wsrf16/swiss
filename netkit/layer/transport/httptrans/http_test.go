package httptrans

import (
	"testing"
	"time"
)

// forward proxy
func TestTransferFromListenAddress(t *testing.T) {
	err := TransferFromListenAddress(":8082", true, nil)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestTransferFromListen(t *testing.T) {
	transfer, err := BuildTransfer(":8082", false)
	if err != nil {
		t.Error(err)
	}

	go func() {
		time.Sleep(10 * time.Second)
		transfer.Stop()
	}()

	err = transfer.TransferFromListen()
	if err != nil {
		t.Error(err)
	}
}
