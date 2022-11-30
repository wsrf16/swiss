package bufferkit

import "bytes"

func DeclareBuffer() bytes.Buffer {
	var b bytes.Buffer
	return b
}

func InitialBuffer() bytes.Buffer {
	b := new(bytes.Buffer)
	return *b
}

func NewBufferByte(buf []byte) bytes.Buffer {
	b := bytes.NewBuffer(buf)
	return *b
}

func NewBufferString(s string) bytes.Buffer {
	b := bytes.NewBufferString(s)
	return *b
}

func Write(buffer bytes.Buffer, p []byte) {
	buffer.Write(p)
}

func WriteString(buffer bytes.Buffer, s string) {
	buffer.WriteString(s)
}

func WriteByte(buffer bytes.Buffer, c byte) {
	buffer.WriteByte(c)
}

func WriteRune(buffer bytes.Buffer, r rune) {
	buffer.WriteRune(r)
}

func Read(buffer bytes.Buffer, c []byte) (int, error) {
	return buffer.Read(c)
}

func ReadByte(buffer bytes.Buffer) (byte, error) {
	return buffer.ReadByte()
}

func ReadRune(buffer bytes.Buffer) (r rune, size int, err error) {
	return buffer.ReadRune()
}
