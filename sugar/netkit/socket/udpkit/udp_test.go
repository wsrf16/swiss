package udpkit

import (
	"testing"
)

func TestTransferToServe(t *testing.T) {
	err := TransferToServe(":8082", "mecs.com:8080")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestTransferToHostServe(t *testing.T) {
	err := TransferToHostServe(":8082")
	if err != nil {
		t.Error(err.Error())
	}
}
