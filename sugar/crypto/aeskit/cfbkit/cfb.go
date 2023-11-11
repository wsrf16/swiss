package cfbkit

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

//func Encrypt(data, key []byte) ([]byte, error) {
//    //md5sum := md5.Sum(key)
//    block, err := aes.NewCipher(key)
//    if err != nil {
//        return nil, err
//    }
//    iv := make([]byte, block.BlockSize())
//    _, err = io.ReadFull(rand.Reader, iv)
//    if err != nil {
//        return nil, err
//    }
//
//    stream := cipher.NewCFBEncrypter(block, iv)
//    dst := make([]byte, block.BlockSize())
//    stream.XORKeyStream(dst, data)
//    return dst, nil
//}
//
//func Decrypt(data, key []byte) ([]byte, error) {
//   //md5sum := md5.Sum(key)
//   block, err := aes.NewCipher(key)
//   if err != nil {
//       return nil, err
//   }
//   blockSize := block.BlockSize()
//   BlockMode := cipher.NewCFBDecrypter(block, key[:blockSize])
//   dst := make([]byte, len(data))
//   BlockMode.CryptBlocks(dst, data)
//   return dst, nil
//}

//func Decrypt(crypted, key []byte) ([]byte, error) {
//    block, err := aes.NewCipher(key)
//    if err != nil {
//        return nil, err
//    }
//    blockSize := block.BlockSize()
//    blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
//    origData := make([]byte, len(crypted))
//    blockMode.CryptBlocks(origData, crypted)
//    return origData, nil
//}

///////////////////////
//func Encrypt(message string, key []byte) (encoded string, err error) {
//    //从输入字符串创建字节切片
//    plainText := []byte(message)
//
//    //使用密钥创建新的 AES 密码
//    block, err := aes.NewCipher(key)
//
//    //如果 NewCipher 失败，退出：
//    if err != nil {
//        return
//    }
//
//    // ^ 使密文成为大小为 BlockSize + 消息长度的字节切片,这样传值后修改不会更改底层数组
//    cipherText := make([]byte, aes.BlockSize+len(plainText))
//
//    // ^ iv 是初始化向量 (16字节)
//    iv := cipherText[:aes.BlockSize]
//    if _, err = io.ReadFull(rand.Reader, iv); err != nil {
//        return
//    }
//
//    // ^ 加密数据,给定加密算法用的密钥,以及初始化向量
//    stream := cipher.NewCFBEncrypter(block, iv)
//    stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
//
//    //返回以base64编码的字符串
//    return base64.RawStdEncoding.EncodeToString(cipherText), err
//}
//
//func Decrypt(secure string, key []byte) (decoded string, err error) {
//    //删除 base64 编码：
//    cipherText, err := base64.RawStdEncoding.DecodeString(secure)
//
//    //如果解码字符串失败，退出：
//    if err != nil {
//        return
//    }
//
//    //使用密钥和加密消息创建新的 AES 密码
//    block, err := aes.NewCipher(key)
//
//    //如果 NewCipher 失败，退出：
//    if err != nil {
//        return
//    }
//
//    //如果密文的长度小于 16 字节
//    //if len(cipherText) < aes.BlockSize {
//    //    err = errors.New("密文分组长度太小")
//    //    return
//    //}
//
//    // ^ iv 是初始化向量 (16字节)
//    iv := cipherText[:aes.BlockSize]
//    cipherText = cipherText[aes.BlockSize:]
//
//    //解密消息
//    stream := cipher.NewCFBDecrypter(block, iv)
//    stream.XORKeyStream(cipherText, cipherText)
//
//    return string(cipherText), err
//}

func Encrypt(plainBytes []byte, key []byte) ([]byte, error) {
	return EncryptWithNonceSize(plainBytes, key, aes.BlockSize)
}

// 32 bit key for AES-256
// 24 bit key for AES-192
// 16 bit key for AES-128
func EncryptWithNonceSize(plainBytes []byte, key []byte, nonceSize int) ([]byte, error) {
	// ^ 使密文成为大小为 BlockSize + 消息长度的字节切片,这样传值后修改不会更改底层数组
	// ^ nonce 是初始化向量 (16字节)
	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	return EncryptWithNonce(plainBytes, key, nonce)
}

func EncryptWithNonce(plainBytes []byte, key []byte, nonceBytes []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, nonceBytes)

	// ^ 使密文成为大小为 BlockSize + 消息长度的字节切片,这样传值后修改不会更改底层数组
	cipherBytes := make([]byte, len(plainBytes))
	stream.XORKeyStream(cipherBytes, plainBytes)

	bytes := append(nonceBytes, cipherBytes...)
	return bytes, nil
}

func NewStream(key []byte, nonceBytes []byte) (cipher.Stream, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, nonceBytes)

	return stream, nil
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

	// ^ nonceBytes 是初始化向量 (16字节)
	nonceBytes := cipherBytes[:nonceSize]
	//cipherBytes = cipherBytes[nonceSize:]
	plainBytes := make([]byte, len(cipherBytes))

	stream := cipher.NewCFBDecrypter(block, nonceBytes)
	stream.XORKeyStream(plainBytes, cipherBytes[nonceSize:])

	return plainBytes, err
}

func Decrypt(cipherText []byte, key []byte) ([]byte, error) {
	return DecryptWithNonceSize(cipherText, key, aes.BlockSize)
}
