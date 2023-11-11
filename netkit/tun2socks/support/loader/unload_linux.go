package loader

import (
	"github.com/wsrf16/swiss/netkit/settings/linux/ipsetting"
	"github.com/wsrf16/swiss/netkit/tun2socks/support/t2sconfig"
	"github.com/wsrf16/swiss/sugar/logo"
)

func Unload(config *t2sconfig.TunConfig) error {
	//"ip route add default via 192.168.0.1 dev ens33 proto static metric 100"

	code, stdout, stderr, err := ipsetting.RouteDelIPWithDev("default", *config.TunName)
	if code != 0 || err != nil {
		logo.W("failed to delete route \"default\"'", err, stderr)
	} else {
		logo.I(stdout)
	}

	code, stdout, stderr, err = ipsetting.RouteDelIPWithGateway(config.GetProxyHost(), *config.DefaultGateway)
	if code != 0 || err != nil {
		logo.W("failed to delete route \""+config.GetProxyHost()+"\"'", err, stderr)
	} else {
		logo.I(stdout)
	}

	code, stdout, stderr, err = ipsetting.RouteDelIPsWithGateway(config.GetTunDNS(), *config.DefaultGateway)
	if code != 0 || err != nil {
		logo.W("failed to delete route \""+*config.TunDNS+"\"'", err, stderr)
	} else {
		logo.I(stdout)
	}

	code, stdout, stderr, err = ipsetting.LinkSetDevDown(*config.TunName)
	if code != 0 || err != nil {
		logo.W("failed to down tun device", err, stderr)
	} else {
		logo.I(stdout)
	}

	ipsetting.TunTapDelTUN(*config.TunName)
	//code, stdout, stderr, err = ipsetting.TunTapDelTUN(*config.TunName)
	//if code != 0 || err != nil {
	//    logo.W("failed to delete tun device", err, stderr)
	//} else {
	//    logo.I(stdout)
	//}

	return nil
}
