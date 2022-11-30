package sshkit

import (
	"github.com/wsrf16/swiss/sugar/netkit/socket"
	//"github.com/wsrf16/swiss/sugar/netkit/socket/tcpkit"
	"github.com/wsrf16/swiss/sugar/netkit/socket/tcpkit"
	"net"
)

func TunnelServe(lAddress string, mediumProperty *SSHProperty, dstAddress string) error {
	lAddr, err := net.ResolveTCPAddr("tcp", lAddress)
	if err != nil {
		return err
	}

	connFunc, err := connFactoryFunc(mediumProperty, dstAddress)
	if err != nil {
		return err
	}

	return tcpkit.ListenThenAcceptThenTransferTo(lAddr, connFunc, true)
}

func connFactoryFunc(mediumProperty *SSHProperty, dstAddress string) (socket.ConnFactoryFunc, error) {
	serverSSHClient, err := BuildClientByProperty(mediumProperty)
	if err != nil {
		return nil, err
	}

	addr, err := tcpkit.NewTCP4Addr(dstAddress)
	if err != nil {
		return nil, err
	}

	connFunc := func() (net.Conn, error) {
		return serverSSHClient.DialTCP("tcp", nil, addr)
	}
	return connFunc, nil
}
