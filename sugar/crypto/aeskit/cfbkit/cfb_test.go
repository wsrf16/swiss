package cfbkit

import (
	"testing"
)

func TestAES(t *testing.T) {
	encrypt, err := Encrypt([]byte("12345678901234567890123456789011"), []byte("woshimima11451410086123456789012"))
	err = err
	decrypt, err := Decrypt(encrypt, []byte("woshimima11451410086123456789012"))
	t.Log(string(decrypt))
}
