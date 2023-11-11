package httptrans

import (
	"github.com/wsrf16/swiss/netkit/layer/transport/tcptrans"
	"github.com/wsrf16/swiss/sugar/base/gokit"
	"github.com/wsrf16/swiss/sugar/base/lambda"
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"log"
	"net"
)

func TransferFromListenAddress(lAddress string, keepListening bool, output *HttpTransferContext) error {
	lAddr, err := tcptrans.NewTCPAddr(lAddress)
	if err != nil {
		return err
	}

	return TransferFromListen(lAddr, keepListening, output)
}

func TransferFromListen(lAddr *net.TCPAddr, keepListening bool, output *HttpTransferContext) error {
	return lambda.LoopReturn(keepListening, func() error {
		lListener, err := tcptrans.Listen(lAddr)
		if err != nil {
			return err
		}
		if output != nil {
			output.LAddr = lAddr
			output.KeepListening = keepListening
			output.LListener = lListener
		}

		srcConnFactory := func() (net.Conn, error) {
			return tcptrans.Accept(lListener)
		}

	For:
		for {
			src, err := srcConnFactory()
			if err != nil {
				log.Println(err)
				return err
			}

			go TransferHTTP(src, true)

			if output != nil && output.StopChan != nil {
				select {
				case <-*output.StopChan:
					break For
				default:
				}
			}
		}
		return nil
	})
}

func TransferHTTP(src net.Conn, closed bool) (chan iokit.Direction, error) {
	var dst net.Conn

	if closed {
		defer iokit.Close(src, dst)
	}

	readBytes, err := iokit.ReadAllBytesNonBlocking(src)
	if err != nil {
		return nil, err
	}

	packet, err := ResolvePacket(readBytes)
	if err != nil {
		return nil, err
	}

	dst, err = packet.DialDSTConn()
	if err != nil {
		return nil, err
	} else {
		log.Printf("tcp proxy(go-routine-id: %v): {from: %s <-> %s to: %s <-> %s}\n", gokit.GetGid(), src.LocalAddr(), src.RemoteAddr(), dst.LocalAddr(), packet.GetAddress())
	}

	if dst != nil && closed {
		defer iokit.Close(dst)
	}

	if packet.IsMethodConnect() {
		_, err := iokit.WriteString(src, ConnectEstablished)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := iokit.Write(dst, readBytes)
		if err != nil {
			return nil, err
		}
	}

	return iokit.CopyDuplex(src, dst, closed), nil
}

func TransferTCPOrHTTP(src net.Conn, dst net.Conn, closed bool) (chan iokit.Direction, error) {
	if closed {
		defer iokit.Close(src, dst)
	}
	if dst == nil {
		return TransferHTTP(src, closed)
	} else {
		return tcptrans.TransferTCP(src, dst, closed)
	}
}

func TransferTCPOrHTTPDynamic(src net.Conn, dstConnFactory iokit.ConnFactoryFunc, closed bool) (chan iokit.Direction, error) {
	dst, err := dstConnFactory()
	if err != nil {
		return nil, err
	}
	return TransferTCPOrHTTP(src, dst, closed)
}
