package bytekit

import (
	"bytes"
	"encoding/binary"
)

func Unmarshal(bin []byte, data interface{}) {
	reader := bytes.NewReader(bin)
	binary.Read(reader, binary.LittleEndian, data)
}

func UnmarshalToT[T any](bin []byte) T {
	t := new(T)
	reader := bytes.NewReader(bin)
	binary.Read(reader, binary.LittleEndian, *t)
	return *t
}

func Marshal(data interface{}) []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, data)
	return buffer.Bytes()
}

func SubString(b []byte, c byte) []byte {
	return b[:bytes.IndexByte(b, c)]
}
