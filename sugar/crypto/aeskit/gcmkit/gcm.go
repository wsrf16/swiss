package gcmkit

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

//func encrypt(plain []byte, key []byte) (cipher.AEAD, error) {
//    block, err := aes.NewCipher(key)
//    if err != nil {
//        return nil, err
//    }
//    return cipher.NewGCM(block)
//}

func Encrypt(plainBytes []byte, key []byte) ([]byte, error) {
	return EncryptWithNonceSize(plainBytes, key, aes.BlockSize)
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
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCMWithNonceSize(block, len(nonceBytes))
	if err != nil {
		return nil, err
	}

	// ^ 使密文成为大小为 BlockSize + 消息长度的字节切片,这样传值后修改不会更改底层数组
	//cipherBytes := make([]byte, len(plainBytes))
	cipherBytes := gcm.Seal(nil, nonceBytes, plainBytes, nil)
	bytes := append(nonceBytes, cipherBytes...)

	return bytes, nil
}

func NewAEAD(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(block)
}

func Decrypt(cipherBytes []byte, key []byte) ([]byte, error) {
	return DecryptWithNonceSize(cipherBytes, key, aes.BlockSize)
}

func DecryptWithNonceSize(cipherBytes []byte, key []byte, nonceSize int) ([]byte, error) {
	block, err := aes.NewCipher(key)
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

	aead, err := cipher.NewGCMWithNonceSize(block, len(nonce))
	if err != nil {
		return nil, err
	}

	plainBytes, err := aead.Open(nil, nonce, cipherBytes[nonceSize:], nil)
	return plainBytes, err
}
