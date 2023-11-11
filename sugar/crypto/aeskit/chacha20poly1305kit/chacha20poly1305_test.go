package chacha20poly1305kit

import (
	"testing"
)

func TestAES(t *testing.T) {
	encrypt, err := Encrypt([]byte("1234567890"), []byte("woshimima11451410086123456789012"))
	err = err
	decrypt, err := Decrypt(encrypt, []byte("woshimima11451410086123456789012"))
	t.Log(string(decrypt))
}
