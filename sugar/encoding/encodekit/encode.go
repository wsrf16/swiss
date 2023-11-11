package encodekit

import (
	"bytes"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
)

func TransformToUTF8FromGBK(b []byte) ([]byte, error) {
	return TransformToUTF8From(b, simplifiedchinese.GBK.NewDecoder())
}

func TransformToUTF8TextFromGBK(b []byte) (string, error) {
	utf8, e := TransformToUTF8FromGBK(b)
	return string(utf8), e
}

func TransformToUTF8FromGBKText(str string) ([]byte, error) {
	return TransformToUTF8From([]byte(str), simplifiedchinese.GBK.NewDecoder())
}

func TransformToUTF8From(b []byte, t transform.Transformer) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(b), t)
	return io.ReadAll(reader)
}

// simplifiedchinese.GBK.NewEncoder()
func TransformTextToUTF8From(str string, enc *encoding.Encoder) (string, error) {
	return enc.String(str)
}

func TransformTextToUTF8FromGBKText(str string) (string, error) {
	return simplifiedchinese.GBK.NewDecoder().String(str)
}
