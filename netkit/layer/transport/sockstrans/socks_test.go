package sockstrans

import (
	"golang.org/x/net/proxy"
	"testing"
	"time"
)

func TestTransferFromListenAddress(t *testing.T) {
	// curl --proxy "socks5://127.0.0.1:1080" https://job.toutiao.com/s/JxLbWby
	auth := &proxy.Auth{User: "root", Password: "1"}
	err := TransferFromListenAddress(":1080", auth, true, nil)
	t.Error(err)
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

	err = transfer.TransferFromListen(nil)
	if err != nil {
		t.Error(err)
	}
}
