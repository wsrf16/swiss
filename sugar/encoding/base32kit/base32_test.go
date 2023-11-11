package base32kit

import (
	"testing"
)

func TestBase32(t *testing.T) {
	encode := EncodeStringToString("62")
	t.Log(encode)
	if decode, err := DecodeStringToString(encode); err != nil {
		t.Error(err)
	} else {
		t.Log(decode)
	}
}
