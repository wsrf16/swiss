package iokit

import (
	"github.com/wsrf16/swiss/sugar/base/gokit"
	"io"
	"log"
	"net"
	"time"
)

// DefaultDeadLineDuration keep-alive
// const DefaultDeadLineDuration = 30 * time.Second
// const LongDeadLineDuration = 365 * 24 * time.Hour

type ConnFactoryFunc func() (net.Conn, error)

func SetDeadLine(conn net.Conn, d time.Duration) {
	conn.SetDeadline(time.Now().Add(d))
}

func NewConn(connFactory ConnFactoryFunc) (net.Conn, error) {
	conn, err := connFactory()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func Close(conn ...io.Closer) {
	if conn == nil {
		return
	}
	for _, c := range conn {
		if c != nil {
			c.Close()
		}
	}
}

func ForceClose(conn ...io.Closer) {
	if conn == nil {
		return
	}
	for _, c := range conn {
		if c != nil {
			c.Close()

			if t, ok := c.(*net.TCPConn); ok {
				t.SetLinger(0)
			}
		}
	}
}

func CopyDuplex(src io.ReadWriteCloser, dst io.ReadWriteCloser, thenClose bool) chan Direction {
	if thenClose {
		defer func() {
			Close(src, dst)
		}()
	}

	times := 0
	finishChan := make(chan Direction, 2)

	trans := func(direction Direction) {
		var err error
		switch direction {
		case Upload:
			_, err = copy(dst, src)
		case Download:
			_, err = copy(src, dst)
		}
		finishChan <- direction
		times++
		if err != nil {
			log.Printf("go-routine-id: %v, direction: %v, error: %v", gokit.GetGid(), direction, err)
			if thenClose && times > 1 {
				Close(src, dst)
			}
		}

	}
	go trans(Upload)
	trans(Download)

	return finishChan
}

var Monitor = false

//func Transfer(src io.Reader, dst io.Writer, completedChan chan Direction, direct Direction) error {
//    defer func() { completedChan <- direct }()
//    _, err := copy(dst, src)
//    if err != nil {
//        // srcConn, srcOk := src.(net.Conn)
//        // dstConn, dstOk := dst.(net.Conn)
//        // if srcOk && dstOk {
//        //    log.Printf("go-routine-id: %v, order: srcLocal(%v)-srcRemote(%v)->dstLocal(%v)dstRemote(%v)", gokit.GetGid(), srcConn.LocalAddr(), srcConn.RemoteAddr(), dstConn.LocalAddr(), dstConn.RemoteAddr(), err)
//        // }
//        log.Printf("go-routine-id: %v, %v", gokit.GetGid(), err)
//        return err
//    }
//    return err
//}

func copy(dst io.Writer, src io.Reader) (int64, error) {
	n, err := CopyBufferBlock(dst, src, Monitor)
	return n, err
}
