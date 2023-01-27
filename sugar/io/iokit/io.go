package iokit

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/wsrf16/swiss/sugar/console/colorkit"
	"github.com/wsrf16/swiss/sugar/io/pathkit"
	"io"
	"log"
	"net"
	"os"
)

//var (
//    Monitor = false
//)

func WriteToFile(bytes []byte, path string) error {
	return os.WriteFile(path, bytes, 0666)
}

func WriteToCurrentFile(relative string, bytes []byte) error {
	path := pathkit.GetPWD(relative)
	return WriteToFile(bytes, path)
}

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func ReadCurrentFile(relative string) ([]byte, error) {
	path := pathkit.GetPWD(relative)
	bytes, err := os.ReadFile(path)
	return bytes, err
}

func ReadFileToBuf(filePath string, bufSize int, receive func([]byte)) error {
	f, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer func() {
		if err = f.Close(); err != nil {
			panic(err)
		}
	}()

	reader := bufio.NewReader(f)
	buf := make([]byte, bufSize)

	for {
		n, err := reader.Read(buf)
		if err != nil {
			return err
		}
		if n <= 0 {
			break
		}
		receive(buf)
	}
	return nil
}

func ScanLine(filePath string, handle func(string)) error {
	f, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer func() {
		if err = f.Close(); err != nil {
			panic(err)
		}
	}()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		handle(scanner.Text())
	}
	return nil
}

func DirectCopy(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}

func CopyBufferBlock(dst io.Writer, src io.Reader, monitor bool) (int64, error) {
	return CopyBuffer(dst, src, 1024, true, monitor)
}

func p(b []byte) {
	format := colorkit.SpellColorString("%s", colorkit.GreenBg, colorkit.Yellow)
	log.Printf(format, b)
}

func CopyBuffer(dst io.Writer, src io.Reader, bufLength int, block bool, monitor bool) (int64, error) {
	if block && monitor == false {
		return io.Copy(dst, src)
	}

	var back CopyBack
	if monitor {
		back = p
	} else {
		back = nil
	}
	buffer, err := CopyBufferCallBack(dst, src, bufLength, block, back)
	return int64(len(buffer)), err
}

type CopyBack func([]byte)

func CopyBufferCallBack(dst io.Writer, src io.Reader, bufLength int, block bool, back CopyBack) (total []byte, err error) {
	//buf不会自动扩容
	//buf := make([]byte, 0, bufLength)
	buf := make([]byte, bufLength)
	total = make([]byte, 0)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			if back != nil {
				back(buf[0:nr])
			}
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errors.New("invalid write result")
				}
			}
			// written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}

			total = append(total, buf[0:nr]...)
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}

		if !block && nr != cap(buf) {
			return total, err
		}
	}
	return total, err
}

func Write(wr io.Writer, b []byte) (int, error) {
	return wr.Write(b)
}

func Read(rd io.Reader, b []byte) (int, error) {
	return rd.Read(b)
}

func WriteString(wr io.Writer, s string) (int, error) {
	return io.WriteString(wr, s)
}

func ReadAllBytesBlockless(rd io.Reader) ([]byte, error) {
	return ReadAllBytesBuffer(rd, 256, false)
}

func ReadAllBytes(rd io.Reader) ([]byte, error) {
	return ReadAllBytesBuffer(rd, 256, true)
}

func ReadAllBytesBuffer(rd io.Reader, bufLength int, block bool) ([]byte, error) {
	buf := make([]byte, bufLength)
	total := make([]byte, 0, 1024)
	for {
		n, err := rd.Read(buf)
		if err != nil {
			if err == io.EOF {
				return total, nil
			} else {
				return total, err
			}
		}

		total = append(total, buf[:n]...)

		if !block && n != cap(buf) {
			return total, err
		}
	}
}

func ReadAllString(rd io.Reader) (string, error) {
	slice, err := ReadAllBytesBlockless(rd)
	return string(slice), err
}

func ReadToByte(rd io.Reader, delim byte) ([]byte, error) {
	return ReadToBytes(rd, []byte{delim})
}

func ReadToBytes(rd io.Reader, delim []byte) ([]byte, error) {
	b := make([]byte, 1, 1)
	total := make([]byte, 0, 1024)
	length := 0
	for {
		n, err := rd.Read(b)
		length += n
		if err != nil {
			return total[0:length], err
		} else {
			total = append(total, b[0])
			if bytes.HasSuffix(total, delim) {
				return total[0:length], err
			}
		}
	}
}

func ReadLine(rd io.Reader) (line []byte, err error) {
	return ReadToBytes(rd, []byte{'\n'})
}

func ReadLines(rd io.Reader) ([][]byte, error) {
	lines := make([][]byte, 0, 512)
	for {
		line, err := ReadLine(rd)
		if err != nil {
			if err == io.EOF {
				return lines, nil
			} else {
				return lines, err
			}
		} else {
			lines = append(lines, line)
		}
	}
}

//func ReadBytes(rd io.Reader, delim byte) (line []byte, err error) {
//    return bufio.NewReader(rd).ReadBytes(delim)
//}
//
//func ReadString(rd io.Reader, delim byte) (string, error) {
//    return bufio.NewReader(rd).ReadString(delim)
//}
//
//func ReadSlice(rd io.Reader, delim byte) (line []byte, err error) {
//    return bufio.NewReader(rd).ReadSlice(delim)
//}

func CloseAll(connSlice ...net.Conn) {
	if connSlice != nil {
		for _, conn := range connSlice {
			if conn != nil {
				conn.Close()
			}
		}
	}
}

func IsConnect(conn net.Conn) bool {
	n, err := conn.Write([]byte("~~~~~~"))
	log.Println(n)
	return err == nil
}

func ReadFirstLine(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}

	firstLine, err := ReadLine(f)
	if err != nil {
		return "", err
	}

	return string(firstLine), nil
}
