package sshkit

import (
	"github.com/wsrf16/swiss/sugar/netkit/socket/socketkit"
	"github.com/wsrf16/swiss/sugar/netkit/socket/tcpkit"
	"net"
)

func TunnelServe(lAddress string, mediumProperty *SSHProperty, dAddress string) error {
	_, client, err := tcpkit.ListenAndAcceptAddress(lAddress)
	if err != nil {
		return err
	}

	//connFunc, err := connFactoryFunc(mediumProperty, dAddress)
	//if err != nil {
	//    return err
	//}
	//
	//server, err := connFunc()
	//if err != nil {
	//    return err
	//}

	server, err := dialAddress(mediumProperty, dAddress)
	if err != nil {
		return err
	}
	_, err = tcpkit.Transfer(client, server, true)
	return err
}

func connFactoryFunc(lSSHProperty *SSHProperty, dAddress string) (socketkit.ConnFactoryFunc, error) {
	serverSSHClient, err := BuildClientByProperty(lSSHProperty)
	if err != nil {
		return nil, err
	}

	addr, err := tcpkit.NewTCP4Addr(dAddress)
	if err != nil {
		return nil, err
	}

	connFunc := func() (net.Conn, error) {
		return serverSSHClient.DialTCP("tcp", nil, addr)
	}
	return connFunc, nil
}

func dialAddress(lSSHProperty *SSHProperty, dAddress string) (net.Conn, error) {
	serverSSHClient, err := BuildClientByProperty(lSSHProperty)
	if err != nil {
		return nil, err
	}

	addr, err := tcpkit.NewTCP4Addr(dAddress)
	if err != nil {
		return nil, err
	}

	//connFunc := func() (net.Conn, error) {
	//    return serverSSHClient.DialTCP("tcp", nil, addr)
	//}
	conn, err := serverSSHClient.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
