package base62kit

import (
	"testing"
)

func TestBase62(t *testing.T) {
	encode := Encode(62)
	t.Log(encode)
	decode := Decode(encode)
	t.Log(decode)
}
