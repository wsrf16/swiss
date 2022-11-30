package hashkit

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

func SHA256(signature []byte) []byte {
	h := sha256.New()
	if _, err := h.Write(signature); err != nil {
		panic(err)
	}
	sum := h.Sum(nil)
	return sum
}

func MD5(content string) string {
	md5hash := md5.New()
	md5hash.Write([]byte(content))
	cipherStr := md5hash.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
