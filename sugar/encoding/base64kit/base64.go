package base64kit

import "encoding/base64"

func Encode(s string) string {
	bytes := []byte(s)
	encoded := base64.StdEncoding.EncodeToString(bytes)
	return encoded
}

func Decode(s string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	result := string(bytes)
	return result, nil
}
