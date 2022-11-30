package socket

import (
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"net"
	"time"
)

const DefaultDeadLineDuration = 180 * time.Second

func TransferThenClose(client net.Conn, server net.Conn) error {
	return Transfer(client, server, true)
}

type ConnFactoryFunc func() (net.Conn, error)

func TransferDynamic(client net.Conn, serverConnFactory ConnFactoryFunc, closeConn bool) error {
	server, err := serverConnFactory()
	if err != nil {
		return err
	}

	return Transfer(client, server, closeConn)
}

func Transfer(client net.Conn, server net.Conn, closeConn bool) error {
	if closeConn {
		if client != nil {
			defer client.Close()
		}

		if server != nil {
			defer server.Close()
		}
	}

	iokit.TransferRoundTrip(client, server)

	return nil
}

//type Packet interface {
//    ResolveDSTConn() (net.Conn, error)
//    PreHandle(client net.Conn, server net.Conn) error
//}
