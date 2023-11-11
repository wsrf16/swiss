package udptrans

import (
	"github.com/wsrf16/swiss/sugar/base/lambda"
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"log"
	"net"
)

func TransferFromListenToDialAddress(lAddress string, dAddress string, keepListening bool, output *UDPTransferContext) error {
	lAddr, err := NewUDPAddr(lAddress)
	if err != nil {
		return err
	}
	dAddr, err := NewUDPAddr(dAddress)
	if err != nil {
		return err
	}

	return TransferFromListenToDial(lAddr, dAddr, true, output)
}

func TransferFromListenToListenAddress(lAddressFrom string, lAddressTo string, output *UDPTransferContext) error {
	lAddrFrom, err := NewUDPAddr(lAddressFrom)
	if err != nil {
		return err
	}
	lAddrTo, err := NewUDPAddr(lAddressTo)
	if err != nil {
		return err
	}

	return TransferFromListenToListen(lAddrFrom, lAddrTo, true, output)
}
func TransferFromDialToDialAddress(dAddressFrom string, dAddressTo string, output *UDPTransferContext) error {
	dAddrFrom, err := NewUDPAddr(dAddressFrom)
	if err != nil {
		return err
	}
	dAddrTo, err := NewUDPAddr(dAddressTo)
	if err != nil {
		return err
	}

	return TransferFromDialToDial(dAddrFrom, dAddrTo, true, output)
}

func TransferFromListenToDial(lAddr *net.UDPAddr, dAddr *net.UDPAddr, keepListening bool, output *UDPTransferContext) error {
	return lambda.LoopReturn(keepListening, func() error {
		if output != nil {
			output.LAddr = lAddr
			output.DAddr = dAddr
			output.KeepListening = keepListening
		}

		srcConnFactory := func() (net.Conn, error) {
			return Listen(lAddr)
		}
		dstConnFactory := func() (net.Conn, error) {
			if dAddr == nil {
				return nil, nil
			} else {
				return DialAddr(nil, dAddr)
			}
		}

	For:
		for {
			src, err := srcConnFactory()
			if err != nil {
				log.Println(err)
				return err
			}

			// dst, err := dstConnFactory()
			// if err != nil {
			//    log.Println(err)
			//    return err
			// }
			// Transfer(src, dst, 1, true, true)
			TransferDynamic(src, dstConnFactory, true)

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

func TransferFromListenToListen(lAddr *net.UDPAddr, dAddr *net.UDPAddr, keepListening bool, output *UDPTransferContext) error {
	return lambda.LoopReturn(keepListening, func() error {
		if output != nil {
			output.LAddr = lAddr
			output.DAddr = dAddr
			output.KeepListening = keepListening
		}

		srcConnFactory := func() (net.Conn, error) {
			return Listen(lAddr)
		}
		dstConnFactory := func() (net.Conn, error) {
			return Listen(dAddr)
		}

	For:
		for {
			src, err := srcConnFactory()
			if err != nil {
				log.Println(err)
				return err
			}
			// dst, err := dstConnFactory()
			// if err != nil {
			//    log.Println(err)
			//    return err
			// }
			// Transfer(src, dst, true)
			TransferDynamic(src, dstConnFactory, true)

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

func TransferFromDialToDial(lAddr *net.UDPAddr, dAddr *net.UDPAddr, keepListening bool, output *UDPTransferContext) error {
	return lambda.LoopReturn(keepListening, func() error {
		if output != nil {
			output.LAddr = lAddr
			output.DAddr = dAddr
			output.KeepListening = keepListening
		}

		srcConnFactory := func() (net.Conn, error) {
			return DialAddr(nil, lAddr)
		}
		dstConnFactory := func() (net.Conn, error) {
			return DialAddr(nil, dAddr)
		}

	For:
		for {
			src, err := srcConnFactory()
			if err != nil {
				log.Println(err)
				return err
			}
			// dst, err := dstConnFactory()
			// if err != nil {
			//    log.Println(err)
			//    return err
			// }
			// Transfer(src, dst, 1, true, true)
			TransferDynamic(src, dstConnFactory, true)

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

func TransferDynamic(src net.Conn, dstConnFactory iokit.ConnFactoryFunc, closed bool) (chan iokit.Direction, error) {
	dst, err := dstConnFactory()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return Transfer(src, dst, closed)
}

func Transfer(src net.Conn, dst net.Conn, closed bool) (chan iokit.Direction, error) {
	return iokit.CopyDuplex(src, dst, closed), nil
}
