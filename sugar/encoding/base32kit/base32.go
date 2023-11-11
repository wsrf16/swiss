package base32kit

import "encoding/base32"

func Encode(src []byte) []byte {
	enc := base32.StdEncoding
	dst := make([]byte, enc.EncodedLen(len(src)))
	enc.Encode(dst, src)
	return dst
}

func EncodeToString(src []byte) string {
	return base32.StdEncoding.EncodeToString(src)
}

func EncodeStringToString(s string) string {
	return base32.StdEncoding.EncodeToString([]byte(s))
}

func EncodeString(s string) []byte {
	return Encode([]byte(s))
}

func Decode(src []byte) ([]byte, error) {
	enc := base32.StdEncoding
	dst := make([]byte, enc.DecodedLen(len(src)))
	n, err := enc.Decode(dst, src)
	return dst[:n], err
}

func DecodeToString(src []byte) (string, error) {
	decrypted, err := Decode(src)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

func DecodeStringToString(s string) (string, error) {
	return DecodeToString([]byte(s))
}

func DecodeString(src string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(src)
}
