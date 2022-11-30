package encodekit

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
)

func TransformToUTF8FromGBK(s []byte) ([]byte, error) {
	return TransformToUTF8From(s, simplifiedchinese.GBK.NewDecoder())
}

func TransformToUTF8TextFromGBK(s []byte) (string, error) {
	utf8, e := TransformToUTF8FromGBK(s)
	return string(utf8), e
}

func TransformToUTF8From(s []byte, t transform.Transformer) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), t)
	return io.ReadAll(reader)
}
