package rsakit

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/wsrf16/swiss/sugar/base/regexpkit"
	"golang.org/x/crypto/ssh"
)

type RSAKey struct {
	PrivateKey       *ecdsa.PrivateKey
	PrivateKeyBuff   *bytes.Buffer
	PrivateKeyText   string
	PrivateKeyBase64 string

	PublicKey       *ecdsa.PublicKey
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

func ParsePrivateKey(privateKey string) *ecdsa.PrivateKey {
	block, _ := pem.Decode([]byte(privateKey))
	if prk, err := x509.ParseECPrivateKey(block.Bytes); err != nil {
		return nil
	} else {
		return prk
	}
}

// ///////////////////
func ParsePublicKey(publicKey string) *ecdsa.PublicKey {
	block, _ := pem.Decode([]byte(publicKey))
	if puk, err := x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		return nil
	} else {
		return puk.(*ecdsa.PublicKey)
	}
}

func curveFromLength(l int) elliptic.Curve {
	switch l {
	case 224:
		return elliptic.P224()
	case 256:
		return elliptic.P256()
	case 348:
		return elliptic.P384()
	case 521:
		return elliptic.P521()
	}
	return elliptic.P384()
}

func GenerateDSAKeyPair(bits int) RSAKey {
	// 1. 生成私钥
	// GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	privateKey, err := ecdsa.GenerateKey(curveFromLength(bits), rand.Reader)
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

func createPrivateKeyBlock(privateKey *ecdsa.PrivateKey) *pem.Block {
	// 2. MarshalPKCS1PrivateKey将rsa私钥序列化为ASN.1 PKCS#1 DER编码
	derPrivateStream, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil
	}

	// 3. Block代表PEM编码的结构, 对其进行设置
	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derPrivateStream,
	}
	return &block
}

func createECDSAPrivateKeyBlock(privateKey *ecdsa.PrivateKey) *pem.Block {
	derPrivateStream, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil
	}

	// 3. Block代表PEM编码的结构, 对其进行设置
	block := pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: derPrivateStream,
	}
	return &block
}

func ParseIntoPrivateKeyBuf(privateKey *ecdsa.PrivateKey) *bytes.Buffer {
	block := createPrivateKeyBlock(privateKey)
	// 4. 使用pem编码, 写入byte
	privateBuff := new(bytes.Buffer)
	if pem.Encode(privateBuff, block) != nil {
		return nil
	}
	return privateBuff
}

func createPublicKeyBlock(publicKey *ecdsa.PublicKey) *pem.Block {
	derPublicStream, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil
	}

	block := pem.Block{
		Type:  "EC PUBLIC KEY",
		Bytes: derPublicStream,
	}
	// block, _ = pem.Decox509.MarshalPKIXPublicKey(publicKey)
	return &block
}

func ParseIntoPublicKeyBuf(publicKey *ecdsa.PublicKey) *bytes.Buffer {
	block := createPublicKeyBlock(publicKey)
	// 2. 编码公钥, 写入byte
	publicBuff := new(bytes.Buffer)
	if pem.Encode(publicBuff, block) != nil {
		return nil
	}
	return publicBuff
}

func createAuthorizedKey(publicKey *ecdsa.PublicKey) ssh.PublicKey {
	sshPublicKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		return nil
	}
	return sshPublicKey
}

func ParseIntoAuthorizedKeyBuf(publicKey *ecdsa.PublicKey) (*ssh.PublicKey, *bytes.Buffer) {
	sshPublicKey := createAuthorizedKey(publicKey)
	buf := ssh.MarshalAuthorizedKey(sshPublicKey)
	return &sshPublicKey, bytes.NewBuffer(buf)
}
