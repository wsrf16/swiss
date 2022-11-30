package sockskit

import (
	"testing"
)

func TestTransferToHostServe(t *testing.T) {
	// curl --proxy "socks5://127.0.0.1:1080" https://job.toutiao.com/s/JxLbWby
	err := TransferToHostServe(":1080")
	t.Error(err)
}
