package convertkit

import (
	"bytes"
	"encoding/binary"
)

func Uint16ToBytes(x uint16) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func BytesToUint16(b []byte) uint16 {
	bytesBuffer := bytes.NewBuffer(b)
	var x uint16
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x

	//return binary.LittleEndian.Uint16(b)
}

func Uint32ToBytes(x uint32) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func BytesToUint32(b []byte) uint32 {
	bytesBuffer := bytes.NewBuffer(b)
	var x uint32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x

	//return binary.LittleEndian.Uint32(b)
}
