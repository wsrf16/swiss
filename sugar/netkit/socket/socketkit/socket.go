package socketkit

import (
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"io"
	"log"
	"net"
	"time"
)

// DefaultDeadLineDuration keep-alive
//const DefaultDeadLineDuration = 30 * time.Second
//const LongDeadLineDuration = 365 * 24 * time.Hour

type ConnFactoryFunc func() (net.Conn, error)

func SetDeadLine(conn net.Conn, d time.Duration) {
	conn.SetDeadline(time.Now().Add(d))
}

func NewConnPair(clientConnFactory ConnFactoryFunc, serverConnFactory ConnFactoryFunc) (net.Conn, net.Conn, error) {
	server, err := serverConnFactory()
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	client, err := clientConnFactory()
	if err != nil {
		log.Println(err)
		return nil, server, err
	}

	return client, server, nil
}

func NewConn(connFactory ConnFactoryFunc) (net.Conn, error) {
	conn, err := connFactory()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}

func Close(conn ...net.Conn) {
	if conn == nil {
		return
	}
	for _, c := range conn {
		if c != nil {
			c.Close()
		}
	}
}

func ForceClose(conn ...net.Conn) {
	if conn == nil {
		return
	}
	for _, c := range conn {
		if c != nil {
			c.Close()
			t, ok := c.(*net.TCPConn)
			if ok {
				t.SetLinger(0)
			}
		}
	}
}

//func TransferDynamicRoundTripThenClose(clientConnFactory ConnFactoryFunc, serverConnFactory ConnFactoryFunc, wait int, closeClient bool, closeServer bool) ([]int, error) {
//    client, server, err := NewConnPair(clientConnFactory, serverConnFactory)
//    if err != nil {
//        log.Println(err)
//        return nil, err
//    }
//    return TransferRoundTripThenClose(client, server, wait, closeClient, closeServer), nil
//}

func TransferRoundTripWaitForCompleted(src net.Conn, dst net.Conn, closed bool) []int {
	//if closeClient && client != nil {
	//    defer Close(client)
	//}
	//if closeServer && server != nil {
	//    defer Close(server)
	//}

	//return TransferRoundTripWaitForCompleted(client, server, wait)
	//}

	//func TransferRoundTripWaitForCompleted(src net.Conn, dst net.Conn, wait int) []int {
	completed := make(chan int, 2)
	// client - server
	go forward(src, dst, completed, 0)
	// server - client
	//reverse(src, dst, completed, 1)
	go reverse(src, dst, completed, 1)

	completedOrder := make([]int, 0)
	completedOrder = append(completedOrder, <-completed)
	if closed {
		if completedOrder[0] == 0 {
			Close(dst)
		}
		if completedOrder[0] == 1 {
			Close(src)
		}
	}
	completedOrder = append(completedOrder, <-completed)

	return completedOrder
}

func forward(src io.Reader, dst io.Writer, completed chan int, order int) error {
	err := Transfer(src, dst, completed, order)
	return err
}

func reverse(dst io.Writer, src io.Reader, completed chan int, order int) error {
	err := Transfer(src, dst, completed, order)
	return err
}

func Transfer(src io.Reader, dst io.Writer, completed chan int, order int) error {
	defer complete(completed, order)
	_, err := iokit.CopyBufferBlock(dst, src, false)
	if err != nil {
		log.Printf("order: %v, %v", order, err)
		return err
	}
	return err
}

func complete(completed chan int, order int) {
	completed <- order
}
