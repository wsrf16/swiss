package sstrans

import (
	"github.com/wsrf16/swiss/netkit/layer/transport/sstrans/sstcptrans"
	"testing"
	"time"
)

func TestTransferFromListenAddress(t *testing.T) {
	algorithm := "aes-256-gcm"
	//algorithm := "aes-128-gcm"
	//algorithm := "CHACHA20-IETF-POLY1305"
	//algorithm := "plain"
	config, err := BuildConfig(algorithm, nil, "123456")
	if err != nil {
		t.Error(err)
	}
	err = sstcptrans.TransferFromListenAddress(":8388", config, true, nil)
}

func TestTransferFromListen(t *testing.T) {
	transfer, err := sstcptrans.BuildTransfer(":8082", false)
	if err != nil {
		t.Error(err)
	}

	go func() {
		time.Sleep(10 * time.Second)
		transfer.Stop()
	}()

	algorithm := "aes-256-gcm"
	//algorithm := "aes-128-gcm"
	//algorithm := "CHACHA20-IETF-POLY1305"
	//algorithm := "plain"
	config, err := BuildConfig(algorithm, nil, "123456")
	err = transfer.TransferFromListen(config)
	if err != nil {
		t.Error(err)
	}
}
