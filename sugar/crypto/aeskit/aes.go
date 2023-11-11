package aeskit

import (
	"bytes"
	"errors"
	"github.com/wsrf16/swiss/sugar/crypto/aeskit/cbckit"
	"github.com/wsrf16/swiss/sugar/crypto/aeskit/cfbkit"
	"github.com/wsrf16/swiss/sugar/crypto/aeskit/gcmkit"
	"github.com/wsrf16/swiss/sugar/crypto/aeskit/plainkit"
	"github.com/wsrf16/swiss/sugar/encoding/base64kit"
	"strings"
)

//type Encrypter interface {
//    Encrypt([]byte, []byte) []byte
//}
//
//type Decrypter interface {
//    Decrypt([]byte, []byte) []byte
//}

type Encrypter interface {
	Encrypt([]byte) []byte
}

type Decrypter interface {
	Decrypt([]byte) []byte
}

type AESTemplate struct {
	Key       []byte
	NonceSize int
	Type      Type
	Encrypter
	Decrypter
}

type Type = string

const (
	PlainType = "PLAIN"
	CBCType   = "CBC"
	CFBType   = "CFB"
	GCMType   = "GCM"
)

func (t AESTemplate) EncryptStringToBase64(plain string) (string, error) {
	encrypt, err := t.Encrypt([]byte(plain))
	if err != nil {
		return "", err
	}

	base64 := base64kit.EncodeToString(encrypt)
	return base64, nil
}

func (t AESTemplate) DecryptBase64ToString(base64 string) (string, error) {
	encrypt, err := base64kit.DecodeString(base64)
	if err != nil {
		return "", err
	}

	decrypted, err := t.Decrypt(encrypt)
	return string(decrypted), nil
}

func (t AESTemplate) Encrypt(plain []byte) ([]byte, error) {
	switch strings.ToUpper(t.Type) {
	case "PLAIN":
		return plainkit.Encrypt(plain)
	case "CBC":
		return cbckit.EncryptWithNonceSize(plain, t.Key, t.NonceSize)
	case "CFB":
		return cfbkit.EncryptWithNonceSize(plain, t.Key, t.NonceSize)
	case "GCM":
		return gcmkit.EncryptWithNonceSize(plain, t.Key, t.NonceSize)
	default:
		return plainkit.Encrypt(plain)
	}
}

func (t AESTemplate) Decrypt(encrypt []byte) ([]byte, error) {
	switch strings.ToUpper(t.Type) {
	case "PLAIN":
		return plainkit.Decrypt(encrypt)
	case "CBC":
		return cbckit.DecryptWithNonceSize(encrypt, t.Key, t.NonceSize)
	case "CFB":
		return cfbkit.DecryptWithNonceSize(encrypt, t.Key, t.NonceSize)
	case "GCM":
		return gcmkit.DecryptWithNonceSize(encrypt, t.Key, t.NonceSize)
	default:
		return plainkit.Decrypt(encrypt)
	}
}

// PKCS5Padding 对数据进行PKCS5填充
func PKCS5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS5Unpadding 去除PKCS5填充
func PKCS5Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("PKCS5 unpadding error: data is empty")
	}
	unpadding := int(data[length-1])
	if unpadding > length {
		return nil, errors.New("PKCS5 unpadding error: invalid padding size")
	}
	return data[:length-unpadding], nil
}

// ZeroPadding 使用ZeroPadding填充数据
func ZeroPadding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{0}, padding)
	return append(data, padText...)
}

// ZeroUnpadding 去除ZeroPadding填充数据
func ZeroUnpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("ZeroUnpadding error: data is empty")
	}
	unpadding := 0
	for i := length - 1; i >= 0; i-- {
		if data[i] == 0 {
			unpadding++
		} else {
			break
		}
	}
	if unpadding == 0 {
		return nil, errors.New("ZeroUnpadding error: no padding bytes found")
	}
	return data[:length-unpadding], nil
}
