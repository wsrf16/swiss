package tcptrans

import (
	"fmt"
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"net"
	"testing"
	"time"
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
		// go Channal(client, dstaddr, dsthost)

		go func(client *net.TCPConn) {
			// addr, err := tcptrans.NewTCP4Addr("192.168.0.133:8080")
			tcpAddr, _ := net.ResolveTCPAddr("tcp", "192.168.0.133:8080")
			conn, err := net.DialTCP("tcp", nil, tcpAddr)
			println(err)
			iokit.CopyDuplex(client, conn, true)
			client.Close()
			conn.Close()
		}(client)
	}
}

// reverse proxy
func TestTransferFromListenToDialAddress(t *testing.T) {
	err := TransferFromListenToDialAddress(":8082", "mecs.com:8080", true, nil)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestTransferFromListenToDial(t *testing.T) {
	transfer, err := BuildTransfer(":8082", "micro.com:3001", false)
	if err != nil {
		t.Error(err)
	}

	go func() {
		time.Sleep(10 * time.Second)
		transfer.Stop()
	}()

	err = transfer.TransferFromListenToDial()
	if err != nil {
		t.Error(err)
	}
}

func TestNAT(t *testing.T) {
	listenAddressFrom := ":9090"
	listenAddressTo := ":9091"

	go TransferFromListenToListenAddress(listenAddressFrom, listenAddressTo, true, nil)

	lAddress := "127.0.0.1:9091"
	dAddress := "mecs.com:7777"
	go TransferFromDialToDialAddress(lAddress, dAddress, false)

}

func TestClassic(t *testing.T) {
	// socket.TransferToHostServe("tcp", "0.0.0.0:8081")
	// tcptrans.TransferToHostServe(":8081")
	// udpkit.TransferToServe("udp", ":8081", "localhost:8084)
	// socket.TransferToServe("tcp", ":8082", "192.168.0.133:22")
	// tcptrans.TransferToServe("tcp", "0.0.0.0:8082", "192.168.0.133:8080")
}

func TestPing(t *testing.T) {
	telnet, err := Telnet("163.com:801", time.Second*2)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(telnet)
	}
}
