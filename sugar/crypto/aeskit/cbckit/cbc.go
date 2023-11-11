package cbckit

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func Encrypt(plain []byte, key []byte) ([]byte, error) {
	return EncryptWithNonceSize(plain, key, aes.BlockSize)
}

// 32 bit key for AES-256
// 24 bit key for AES-192
// 16 bit key for AES-128
func EncryptWithNonceSize(plain []byte, key []byte, nonceSize int) ([]byte, error) {
	// ^ 使密文成为大小为 BlockSize + 消息长度的字节切片,这样传值后修改不会更改底层数组
	// ^ nonce 是初始化向量 (16字节)
	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	return EncryptWithNonce(plain, key, nonce)
}

func EncryptWithNonce(plain []byte, key []byte, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plain = pkcs5Padding(plain, block.BlockSize())
	cbc := cipher.NewCBCEncrypter(block, nonce)

	// ^ 使密文成为大小为 BlockSize + 消息长度的字节切片,这样传值后修改不会更改底层数组
	plainCipher := make([]byte, len(plain))
	cbc.CryptBlocks(plainCipher, plain)

	return plainCipher, nil
}

func DecryptWithNonceSize(cipherBytes []byte, key []byte, nonceSize int) ([]byte, error) {
	newCipher, err := aes.NewCipher(key)
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
	cipherBytes = cipherBytes[nonceSize:]

	stream := cipher.NewCFBDecrypter(newCipher, nonce)
	stream.XORKeyStream(cipherBytes, cipherBytes)

	return cipherBytes, err
}

func Decrypt(cipherText []byte, key []byte) ([]byte, error) {
	return DecryptWithNonceSize(cipherText, key, aes.BlockSize)
}
