package sstcptrans

import (
	"github.com/wsrf16/swiss/netkit/layer/transport/sstrans"
	"github.com/wsrf16/swiss/netkit/layer/transport/tcptrans"
	"github.com/wsrf16/swiss/sugar/base/lambda"
	"github.com/wsrf16/swiss/sugar/base/timekit"
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"log"
	"net"
)

// -s 'ss://AEAD_CHACHA20_POLY1305:123456@:8388' -verbose
func TransferFromListenAddress(lAddress string, config *sstrans.ShadowSocksConfig, keepListening bool, output *SSTransferContext) error {
	lAddr, err := tcptrans.NewTCPAddr(lAddress)
	if err != nil {
		return err
	}

	return TransferFromListen(lAddr, config, keepListening, output)
}

func TransferFromListen(lAddr *net.TCPAddr, config *sstrans.ShadowSocksConfig, keepListening bool, output *SSTransferContext) error {
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

	For:
		for {
			src, err := tcptrans.Accept(lListener)
			if err != nil {
				return err
			}
			src.SetDeadline(timekit.Time1Year())

			// go TransferTCPOrHTTP(src, true, config, true)
			go func() {
				_, err := Transfer(src, true, config, true)
				if err != nil {
					log.Println(err)
				}
			}()

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

func Transfer(srcEncrypted net.Conn, closed bool, config *sstrans.ShadowSocksConfig, recovered bool) (chan iokit.Direction, error) {
	if closed {
		defer iokit.Close(srcEncrypted)
	}
	if recovered {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()
	}

	//cipher := *config.Cipher
	//src := cipher.StreamConn(srcEncrypted)
	src, err := toSSConn(srcEncrypted, config)
	if err != nil {
		return nil, err
	}

	//dst, err := parseDest(src)
	//if err != nil {
	//    log.Println(err)
	//    return nil, err
	//}
	//if closed {
	//    defer transport.Close(dst)
	//}
	//
	//return transport.CopyDuplex(src, dst, closed), nil

	// -------------------------
	dst, err := ParseDestConnFrom(src)
	if err != nil {
		return nil, err
	}
	if closed {
		defer iokit.Close(dst)
	}

	return iokit.CopyDuplex(src, dst, closed), nil
}
