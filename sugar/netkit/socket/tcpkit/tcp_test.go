package tcpkit

import (
	"fmt"
	"github.com/wsrf16/swiss/sugar/netkit/socket/socketkit"
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

		go func(client *net.TCPConn) {
			//addr, err := tcpkit.NewTCP4Addr("192.168.0.133:8080")
			tcpAddr, _ := net.ResolveTCPAddr("tcp", "192.168.0.133:8080")
			conn, err := net.DialTCP("tcp", nil, tcpAddr)
			println(err)
			socketkit.TransferRoundTripWaitForCompleted(client, conn, 1)
			client.Close()
			conn.Close()
		}(client)
	}
}

// forward proxy
func TestTransferFromListenAddress(t *testing.T) {
	err := TransferFromListenAddress(":8082")
	if err != nil {
		t.Error(err.Error())
	}
}

// reverse proxy
func TestTransferFromListenToDialAddress(t *testing.T) {
	err := TransferFromListenToDialAddress(":8082", "mecs.com:8080")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestNAT(t *testing.T) {
	listenAddressFrom := ":9090"
	listenAddressTo := ":9091"
	go TransferFromListenToListenAddress(listenAddressFrom, listenAddressTo)

	dialAddressFrom := "127.0.0.1:9091"
	dialAddressTo := "mecs.com:7777"
	go TransferFromDialToDialAddress(dialAddressFrom, dialAddressTo)

	select {}
}

func TestClassic(t *testing.T) {
	//socket.TransferToHostServe("tcp", "0.0.0.0:8081")
	//tcpkit.TransferToHostServe(":8081")
	//udpkit.TransferToServe("udp", ":8081", "localhost:8084)
	//socket.TransferToServe("tcp", ":8082", "192.168.0.133:22")
	//tcpkit.TransferToServe("tcp", "0.0.0.0:8082", "192.168.0.133:8080")
}
