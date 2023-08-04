package sockskit

import (
	"testing"
)

func TestTransferFromListenAddress(t *testing.T) {
	// curl --proxy "socks5://127.0.0.1:1080" https://job.toutiao.com/s/JxLbWby
	err := TransferFromListenAddress(":1080", nil, true, nil)
	t.Error(err)
}
