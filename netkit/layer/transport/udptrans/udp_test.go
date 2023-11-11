package udptrans

import (
	"testing"
)

func TestTransferFromListenToDialAddress(t *testing.T) {
	err := TransferFromListenToDialAddress(":8082", "mecs.com:8080", false, nil)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestTransferFromDialToDialAddress(t *testing.T) {
	err := TransferFromDialToDialAddress(":9091", "mecs.com:8080", nil)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestTransferFromListenToListenAddress(t *testing.T) {
	err := TransferFromListenToListenAddress(":9090", ":9091", nil)
	if err != nil {
		t.Error(err.Error())
	}
}

//func TestTransferToHostServe(t *testing.T) {
//    err := TransferToHostServe(":8082")
//    if err != nil {
//        t.Error(err.Error())
//    }
//}
