package rsakit

import (
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"testing"
)

func TestGenerateRSAKeyPair(t *testing.T) {
	key := GenerateRSAKeyPair(2048)
	if err := iokit.WriteToCurrentFile("id_rsa", key.PrivateKeyBuff.Bytes()); err != nil {
		panic(err)
	}
	if err := iokit.WriteToCurrentFile("id_rsa.pub", key.AuthorizedKeyBuff.Bytes()); err != nil {
		panic(err)
	}
	if err := iokit.WriteToCurrentFile("id_rsa.pub.plain", key.PublicKeyBuff.Bytes()); err != nil {
		panic(err)
	}
	t.Log(string(key.PrivateKeyBuff.Bytes()))
	t.Log(key.AuthorizedKeyBuff.String())
	t.Log(key.PrivateKeyBuff)
}
