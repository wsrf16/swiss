package socketkit

import (
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"io"
	"log"
	"net"
	"time"
)

// DefaultDeadLineDuration keep-alive
const DefaultDeadLineDuration = 30 * time.Second

type ConnFactoryFunc func() (net.Conn, error)

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

//func TransferDynamicRoundTripThenClose(clientConnFactory ConnFactoryFunc, serverConnFactory ConnFactoryFunc, wait int, closeClient bool, closeServer bool) ([]int, error) {
//    client, server, err := NewConnPair(clientConnFactory, serverConnFactory)
//    if err != nil {
//        log.Println(err)
//        return nil, err
//    }
//    return TransferRoundTripThenClose(client, server, wait, closeClient, closeServer), nil
//}

func TransferRoundTripThenClose(client net.Conn, server net.Conn, wait int, closeClient bool, closeServer bool) []int {
	if closeClient && client != nil {
		defer client.Close()
	}
	if closeServer && server != nil {
		defer server.Close()
	}

	return TransferRoundTripWaitForCompleted(client, server, wait)
}

func TransferRoundTripWaitForCompleted(src net.Conn, dst net.Conn, wait int) []int {
	completed := make(chan int, 2)
	// client - server
	go forward(src, dst, completed, 0)
	// server - client
	//reverse(dst, src, completed, 1)
	go reverse(dst, src, completed, 1)
	completedOrder := make([]int, 0)
	if wait == 1 || wait == 2 {
		if wait > 0 {
			completedOrder = append(completedOrder, <-completed)
		}
		if wait > 1 {
			completedOrder = append(completedOrder, <-completed)
		}
	} else {
		log.Println("wait should be 1 or 2")
		completedOrder = append(completedOrder, <-completed)
	}
	return completedOrder
}

func forward(src io.Reader, dst io.Writer, completed chan int, order int) error {
	return Transfer(src, dst, completed, order)
}

func reverse(dst io.Reader, src io.Writer, completed chan int, order int) error {
	return Transfer(dst, src, completed, order)
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
