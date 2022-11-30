package rsakit

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/wsrf16/swiss/sugar/base/regexpkit"
	"github.com/wsrf16/swiss/sugar/crypto/hashkit"
	"golang.org/x/crypto/ssh"
)

type RSAKey struct {
	PrivateKey       *rsa.PrivateKey
	PrivateKeyBuff   *bytes.Buffer
	PrivateKeyText   string
	PrivateKeyBase64 string

	PublicKey       *rsa.PublicKey
	PublicKeyBuff   *bytes.Buffer
	PublicKeyText   string
	PublicKeyBase64 string

	AuthorizedKey       *ssh.PublicKey
	AuthorizedKeyBuff   *bytes.Buffer
	AuthorizedKeyText   string
	AuthorizedKeyBase64 string
}

func FormatKey(key string) (string, error) {
	replaced := key
	replaced, err := regexpkit.ReplaceAll("(-----BEGIN [a-zA-Z]* PRIVATE KEY-----)", replaced, "$1\n")
	if err != nil {
		return "", err
	}
	replaced, err = regexpkit.ReplaceAll("(-----END [a-zA-Z]* PRIVATE KEY-----)", replaced, "\n$1")
	if err != nil {
		return "", err
	}
	replaced, err = regexpkit.ReplaceAll("(-----BEGIN [a-zA-Z]* PUBLIC KEY-----)", replaced, "$1\n")
	if err != nil {
		return "", err
	}
	replaced, err = regexpkit.ReplaceAll("(-----END [a-zA-Z]* PUBLIC KEY-----)", replaced, "\n$1")
	if err != nil {
		return "", err
	}
	return replaced, nil
}

func ParsePrivateKey(privateKey string) *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(privateKey))
	if prk, err := x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		return nil
	} else {
		return prk
	}
}

func ParsePublicKey(publicKey string) *rsa.PublicKey {
	block, _ := pem.Decode([]byte(publicKey))
	if puk, err := x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		return nil
	} else {
		return puk.(*rsa.PublicKey)
	}
}

func GenerateRSAKeyPair(bits int) RSAKey {
	// 1. 生成私钥
	// GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}

	privateKeyBuff := ParseIntoPrivateKeyBuf(privateKey)
	publicKeyBuff := ParseIntoPublicKeyBuf(&privateKey.PublicKey)
	authorizedKey, authorizedKeyBuff := ParseIntoAuthorizedKeyBuf(&privateKey.PublicKey)

	key := RSAKey{
		PrivateKey:       privateKey,
		PrivateKeyBuff:   privateKeyBuff,
		PrivateKeyText:   privateKeyBuff.String(),
		PrivateKeyBase64: base64.StdEncoding.EncodeToString(privateKeyBuff.Bytes()),

		PublicKey:       &privateKey.PublicKey,
		PublicKeyBuff:   publicKeyBuff,
		PublicKeyText:   publicKeyBuff.String(),
		PublicKeyBase64: base64.StdEncoding.EncodeToString(publicKeyBuff.Bytes()),

		AuthorizedKey:       authorizedKey,
		AuthorizedKeyBuff:   authorizedKeyBuff,
		AuthorizedKeyText:   authorizedKeyBuff.String(),
		AuthorizedKeyBase64: base64.StdEncoding.EncodeToString(authorizedKeyBuff.Bytes()),
	}
	return key
}

func createPrivateKeyBlock(privateKey *rsa.PrivateKey) *pem.Block {
	// 2. MarshalPKCS1PrivateKey将rsa私钥序列化为ASN.1 PKCS#1 DER编码
	derPrivateStream := x509.MarshalPKCS1PrivateKey(privateKey)

	// 3. Block代表PEM编码的结构, 对其进行设置
	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derPrivateStream,
	}
	return &block
}

func ParseIntoPrivateKeyBuf(privateKey *rsa.PrivateKey) *bytes.Buffer {
	block := createPrivateKeyBlock(privateKey)
	// 4. 使用pem编码, 写入byte
	privateBuff := new(bytes.Buffer)
	if pem.Encode(privateBuff, block) != nil {
		return nil
	}
	return privateBuff
}

func createPublicKeyBlock(publicKey *rsa.PublicKey) *pem.Block {
	derPublicStream, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil
	}

	block := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPublicStream,
	}
	// block, _ = pem.Decox509.MarshalPKIXPublicKey(publicKey)
	return &block
}

func ParseIntoPublicKeyBuf(publicKey *rsa.PublicKey) *bytes.Buffer {
	block := createPublicKeyBlock(publicKey)
	// 2. 编码公钥, 写入byte
	publicBuff := new(bytes.Buffer)
	if pem.Encode(publicBuff, block) != nil {
		return nil
	}
	return publicBuff
}

func createAuthorizedKey(publicKey *rsa.PublicKey) ssh.PublicKey {
	sshPublicKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		return nil
	}
	return sshPublicKey
}

func ParseIntoAuthorizedKeyBuf(publicKey *rsa.PublicKey) (*ssh.PublicKey, *bytes.Buffer) {
	sshPublicKey := createAuthorizedKey(publicKey)
	buf := ssh.MarshalAuthorizedKey(sshPublicKey)
	return &sshPublicKey, bytes.NewBuffer(buf)
}

func EncryptIntoBase64(publicKey rsa.PublicKey, plain string) (string, error) {
	if b, err := EncryptOAEP(publicKey, []byte(plain)); err != nil {
		return "", err
	} else {
		return base64.StdEncoding.EncodeToString(b), nil
	}
}

func DecryptFromBase64(privateKey rsa.PrivateKey, cipher string) (string, error) {
	if bytes, err := base64.StdEncoding.DecodeString(cipher); err != nil {
		return "", err
	} else {
		decrypt, err := DecryptOAEP(privateKey, bytes)
		return string(decrypt), err
	}
}

func EncryptOAEP(publicKey rsa.PublicKey, plain []byte) ([]byte, error) {
	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		&publicKey,
		plain,
		nil)
	return encryptedBytes, err
}

func DecryptOAEP(privateKey rsa.PrivateKey, cipher []byte) ([]byte, error) {
	decryptedBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, &privateKey, cipher, nil)
	// decryptedBytes, err := privateKey.Decrypt(nil, cipher, &rsa.OAEPOptions{Hash: crypto.SHA256})
	return decryptedBytes, err
}

func EncryptPKCS1v15(publicKey rsa.PublicKey, plain []byte) ([]byte, error) {
	encryptedBytes, err := rsa.EncryptPKCS1v15(
		rand.Reader,
		&publicKey,
		plain)
	return encryptedBytes, err
}

func DecryptPKCS1v15(privateKey rsa.PrivateKey, cipher []byte) ([]byte, error) {
	decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, &privateKey, cipher)
	return decryptedBytes, err
}

func SignPSS(privateKey *rsa.PrivateKey, target []byte) ([]byte, error) {
	// Before signing, we need to hash our message
	// The hash is what we actually sign
	sum := hashkit.SHA256(target)

	// In order to generate the signature, we provide a random number generator,
	// our private key, the hashing algorithm that we used, and the hash sum
	// of our message
	if signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, sum, nil); err != nil {
		return nil, err
	} else {
		return signature, nil
	}
}

func VerifyPSS(publicKey *rsa.PublicKey, signature []byte) error {
	// Before signing, we need to hash our message
	// The hash is what we actually sign
	sum := hashkit.SHA256(signature)

	// To verify the signature, we provide the public key, the hashing algorithm
	// the hash sum of our message and the signature we generated previously
	// there is an optional "options" parameter which can omit for now
	if err := rsa.VerifyPSS(publicKey, crypto.SHA256, sum, signature, nil); err != nil {
		return err
	} else {
		return nil
	}
}

func SignPKCS1v15(privateKey *rsa.PrivateKey, target []byte) ([]byte, error) {
	// Before signing, we need to hash our message
	// The hash is what we actually sign
	sum := hashkit.SHA256(target)

	// In order to generate the signature, we provide a random number generator,
	// our private key, the hashing algorithm that we used, and the hash sum
	// of our message
	if signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, sum); err != nil {
		return nil, err
	} else {
		return signature, nil
	}
}

func VerifyPKCS1v15(publicKey *rsa.PublicKey, target []byte, signature []byte) error {
	// Before signing, we need to hash our message
	// The hash is what we actually sign
	sum := hashkit.SHA256(target)

	// To verify the signature, we provide the public key, the hashing algorithm
	// the hash sum of our message and the signature we generated previously
	// there is an optional "options" parameter which can omit for now
	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, sum, signature); err != nil {
		return err
	} else {
		return nil
	}
}

func SignText(privateKey *rsa.PrivateKey, plain string) string {
	if signature, err := SignPKCS1v15(privateKey, []byte(plain)); err != nil {
		return ""
	} else {
		return base64.StdEncoding.EncodeToString(signature)
	}
}

func VerifyText(publicKey *rsa.PublicKey, content string, signature string) bool {
	if decode, err := base64.StdEncoding.DecodeString(signature); err != nil {
		return false
	} else {
		return VerifyPKCS1v15(publicKey, []byte(content), decode) == nil
	}
}
