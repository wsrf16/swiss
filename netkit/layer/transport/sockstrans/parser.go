package sockstrans

import (
	"errors"
	"github.com/wsrf16/swiss/netkit/layer/transport/tcptrans"
	"github.com/wsrf16/swiss/sugar/base/gokit"
	"github.com/wsrf16/swiss/sugar/base/timekit"
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"golang.org/x/net/proxy"
	"log"
	"net"
)

func ParseDestConnFrom(src net.Conn, auth *proxy.Auth) (net.Conn, error) {
	// // round 1
	// readBytes1, err := iokit.ReadAllBytesNonBlocking(src)
	// if err != nil {
	//    return nil, err
	// }
	// packet1, err := parseToBasicPacket(readBytes1)
	// if err != nil {
	//    return nil, err
	// }
	// switch packet1.Version {
	// case 4:
	//    return nil, errors.New("not support socks4")
	// case 5:
	//    {
	//        if config == nil || config.Authentication == nil {
	//            src.Write([]byte{0x05, 0x00})
	//        } else {
	//            src.Write([]byte{0x05, 0x02})
	//
	//            err := verify(src, config)
	//            if err != nil {
	//                return nil, err
	//            }
	//        }
	//    }
	// default:
	//    return nil, errors.New("only support socks protocol")
	// }
	// log.Printf("socks: version-%d\n", packet1.Version)
	//
	// // round 2
	// readBytes2, err := iokit.ReadAllBytesNonBlocking(src)
	// if err != nil {
	//    return nil, err
	// }
	// packet2 := parseToDetailPacket(readBytes2)
	// if packet2.Version != 5 {
	//    return nil, errors.New("该协议不是socks5协议")
	// }
	// src.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	socksServer := &SocksServer{Conn: src, Auth: auth}
	detailPacket, err := socksServer.ParseIntoPacket()
	if err != nil {
		return nil, err
	}

	address, err := detailPacket.GetAddress()
	if err != nil {
		return nil, err
	}

	dst, err := tcptrans.DialAddress("", address)
	if err != nil {
		return nil, err
	} else {
		log.Printf("socks  version: %d  proxy(go-routine-id: %v): {from: %s <-> %s to: %s <-> %s}\n", detailPacket.Version, gokit.GetGid(), src.LocalAddr(), src.RemoteAddr(), dst.LocalAddr(), address)
	}
	dst.SetDeadline(timekit.Time1Year())

	return dst, err
}

func (s SocksServer) ParseIntoPacket() (*SocksDetailPacket, error) {
	basicPacket, err := s.receiveAndParseIntoBasicPacket()
	if err != nil {
		return nil, err
	}

	err = s.sendAuthenticateInvite(basicPacket)
	if err != nil {
		return nil, err
	}

	detailPacket, err := s.receiveAndParseIntoDetailPacket(basicPacket)
	return detailPacket, err
}

func (s SocksServer) sendAuthenticateInvite(basicPacket *SocksBasicPacket) error {
	version := basicPacket.Version
	switch version {
	case 4:
		return errors.New("not support socks4")
	case 5:
		{
			err := s.inviteAuthenticateIfNeed()
			return err
		}
	default:
		log.Printf("socks: version-%d\n", version)
		return errors.New("only support socks protocol")
	}
}

func (s SocksServer) receiveAndParseIntoDetailPacket(basicPacket *SocksBasicPacket) (*SocksDetailPacket, error) {
	version := basicPacket.Version
	switch version {
	case 4:
		return nil, errors.New("not support socks4")
	case 5:
		//{
		//    err := s.inviteAuthenticateIfNeed()
		//    if err != nil {
		//        return nil, err
		//    }
		//}
		conn := s.Conn
		bytes, err := iokit.ReadAllBytesNonBlocking(conn)
		if err != nil {
			return nil, err
		}
		detailPacket := parseToDetailPacket(bytes)
		conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

		return &detailPacket, nil
	default:
		log.Printf("socks: version-%d\n", version)
		return nil, errors.New("only support socks protocol")
	}

}

func parseToBasicPacket(bytes []byte) (*SocksBasicPacket, error) {
	packet := SocksBasicPacket{}
	packet.Raw = bytes
	if len(bytes) < 3 {
		return nil, errors.New("invalid socks bytes")
	}
	packet.Version = bytes[0]
	packet.NMethods = bytes[1]
	packet.Methods = bytes[2:]
	return &packet, nil
}

func (s SocksServer) receiveAndParseIntoBasicPacket() (*SocksBasicPacket, error) {
	buffer, err := iokit.ReadAllBytesNonBlocking(s.Conn)
	if err != nil {
		return nil, err
	}

	basicInfoPacket, err := parseToBasicPacket(buffer)
	if err != nil {
		return nil, err
	}

	return basicInfoPacket, nil
}

const (
	VerifySuccess = 0x00
	VerifyFailure = 0x01
)

func (s SocksServer) verify() error {
	// round 2
	src := s.Conn
	auth := s.Auth
	buffer, err := iokit.ReadAllBytesNonBlocking(src)
	if err != nil {
		return err
	}
	authPacket := parseToAuthenticationPacket(buffer)
	checked := checkPacket(authPacket, auth)
	if checked {
		_, err := src.Write([]byte{0x05, VerifySuccess})
		if err != nil {
			return err
		}
	} else {
		_, err := src.Write([]byte{0x05, VerifyFailure})
		if err != nil {
			return err
		} else {
			return errors.New("authentication failed: incorrect username or password")
		}
	}
	return nil
}

func (s SocksServer) inviteAuthenticateIfNeed() error {
	if s.Auth == nil {
		s.replyWhetherAuthMethod(false)
	} else {
		s.replyWhetherAuthMethod(true)

		err := s.verify()
		if err != nil {
			return err
		}
	}
	return nil
}
