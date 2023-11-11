package totpkit

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"time"
)

// otpauth://totp/Google%3Ayourname@gmail.com?secret=xxxx&issuer=Google
func GenerateCodeByTime(secret string, period int64, digits int) (string, error) {
	timeStep := time.Now().Unix() / period
	return GenerateCodeByPeriod(secret, timeStep, digits)
}

func GenerateCodeByPeriod(secret string, second int64, digits int) (string, error) {
	// 1. 共享密钥解码：将Base32编码的密钥解码为字节数组
	decodedKey, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}

	// 2. HMAC-SHA-1：计算哈希值
	h := hmac.New(sha1.New, decodedKey)
	binary.Write(h, binary.BigEndian, second)

	// 3. Truncate：对哈希值进行截断，只保留最后 4 个字节
	offset := h.Sum(nil)[19] & 0x0f
	truncatedHash := h.Sum(nil)[offset : offset+4]
	truncatedHash[0] = truncatedHash[0] & 0x7f
	otp := (int32(truncatedHash[0])<<24 | int32(truncatedHash[1])<<16 | int32(truncatedHash[2])<<8 | int32(truncatedHash[3])) % int32(math.Pow(10, float64(digits)))

	// 4. 格式化：将结果转换为十进制6位的字符串
	return fmt.Sprintf("%0"+strconv.Itoa(digits)+"d", otp), nil
}

func CheckCode(secret string, period int64, digits int, targetCode string) bool {
	if code, err := GenerateCodeByPeriod(secret, period, digits); err != nil {
		return false
	} else {
		return code == targetCode
	}

	if code, err := GenerateCodeByPeriod(secret, period-30, digits); err != nil {
		return false
	} else {
		return code == targetCode
	}

	if code, err := GenerateCodeByPeriod(secret, period+30, digits); err != nil {
		return false
	} else {
		return code == targetCode
	}

	return false
}
