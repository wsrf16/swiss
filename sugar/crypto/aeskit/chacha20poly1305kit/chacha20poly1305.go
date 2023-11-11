package chacha20poly1305kit

import (
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"golang.org/x/crypto/chacha20poly1305"
	"strconv"
)

func Encrypt(plain []byte, key []byte) ([]byte, error) {
	return EncryptWithNonceSize(plain, key, chacha20poly1305.NonceSize)
}

// 32 bit key for AES-256
// 24 bit key for AES-192
// 16 bit key for AES-128
func EncryptWithNonceSize(plainBytes []byte, key []byte, nonceSize int) ([]byte, error) {
	nonceBytes := make([]byte, nonceSize)
	if _, err := rand.Read(nonceBytes); err != nil {
		return nil, err
	}

	return EncryptWithNonce(plainBytes, key, nonceBytes)
}

func EncryptWithNonce(plainBytes []byte, key []byte, nonceBytes []byte) ([]byte, error) {
	block, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, err
	}

	// ^ 使密文成为大小为 BlockSize + 消息长度的字节切片,这样传值后修改不会更改底层数组
	//cipherBytes := make([]byte, len(plainBytes))
	cipherBytes := block.Seal(nil, nonceBytes, plainBytes, nil)
	bytes := append(nonceBytes, cipherBytes...)

	return bytes, nil
}

func NewAEAD(key []byte) (cipher.AEAD, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, errors.New("keysize is not " + strconv.Itoa(chacha20poly1305.KeySize))
	}

	block, err := chacha20poly1305.New(key)
	return block, err
}

func Decrypt(cipherBytes []byte, key []byte) ([]byte, error) {
	return DecryptWithNonceSize(cipherBytes, key, chacha20poly1305.NonceSize)
}

func DecryptWithNonceSize(cipherBytes []byte, key []byte, nonceSize int) ([]byte, error) {
	block, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, err
	}

	//如果密文的长度小于 16 字节
	//if len(cipherBytes) < nonceSize {
	//    err = errors.New("密文分组长度太小")
	//    return nil, err
	//}

	// ^ nonce 是初始化向量 (16字节)
	nonce := cipherBytes[:nonceSize]
	//cipherBytes = cipherBytes[nonceSize:]
	//plainBytes := make([]byte, len(cipherBytes))

	plainBytes, err := block.Open(nil, nonce, cipherBytes[nonceSize:], nil)
	return plainBytes, err
}
