package base64kit

import (
	"testing"
)

func TestBase64(t *testing.T) {
	encode := EncodeStringToString("62")
	t.Log(encode)
	if decode, err := DecodeStringToString(encode); err != nil {
		t.Error(err)
	} else {
		t.Log(decode)
	}
}
