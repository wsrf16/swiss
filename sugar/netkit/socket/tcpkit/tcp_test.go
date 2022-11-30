package tcpkit

import (
	"fmt"
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"net"
	"testing"
)

func TestSample(t *testing.T) {
	addr, err := NewTCP4Addr("0.0.0.0:8082")
	println(err)
	listen, err := net.ListenTCP("tcp", addr)

	for {
		client, err := listen.AcceptTCP()
		if err != nil {
			fmt.Println("accept error: ", err.Error())
			continue
		}
		//go Channal(client, dstaddr, dsthost)

		go func() {
			//addr, err := tcpkit.NewTCP4Addr("192.168.0.133:8080")
			tcpAddr, _ := net.ResolveTCPAddr("tcp", "192.168.0.133:8080")
			conn, err := net.DialTCP("tcp", nil, tcpAddr)
			println(err)
			iokit.TransferRoundTrip(client, conn)
			client.Close()
			conn.Close()
		}()
	}
}

func TestTransferToServe(t *testing.T) {
	err := TransferToServe(":8082", "mecs.com:8080")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestTransferToHostServe(t *testing.T) {
	err := TransferToHostServe(":8082")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestClassic(t *testing.T) {
	//socket.TransferToHostServe("tcp", "0.0.0.0:8081")
	//tcpkit.TransferToHostServe(":8081")
	//udpkit.TransferToServe("udp", ":8081", "localhost:8084)
	//socket.TransferToServe("tcp", ":8082", "192.168.0.133:22")
	//tcpkit.TransferToServe("tcp", "0.0.0.0:8082", "192.168.0.133:8080")
}
