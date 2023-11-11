//go:build linux
// +build linux

package tunsocks

import (
	"github.com/wsrf16/swiss/netkit/settings/hardware/tun/syscalltun"
	"github.com/wsrf16/swiss/netkit/settings/linux/ipsetting"
	"github.com/wsrf16/swiss/netkit/tun2socks/support/t2sconfig"
	"github.com/wsrf16/swiss/sugar/logo"
	"io"
)

func CreateAndOpenTunDevice(config t2sconfig.TunConfig) (io.ReadWriteCloser, error) {
	//return tundevice.CreateAndOpenTunDevice(*config.TunName, *config.TunAddress, *config.TunMask, *config.TunGateway, *config.TunGWMetric, *config.MTU, config.GetTunDNS())

	dev, err := syscalltun.CreateAndOpenTunDevice(*config.TunName)
	//dev, err := watertun.CreateAndOpenTunDevice(*config.TunName, false)
	//dev, err := wintun.CreateAndOpenTunDevice(*config.TunName, *config.MTU)
	if err != nil {
		return nil, err
	}

	if code, stdout, stderr, err := ipsetting.AddrAddIPAndCIDR(*config.TunName, *config.TunAddress); err != nil || code != 0 {
		logo.E(stdout, stderr)
		return nil, err
	}

	if code, stdout, stderr, err := ipsetting.LinkSetDevUp(*config.TunName); err != nil || code != 0 {
		logo.E(stdout, stderr)
		return nil, err
	}

	if code, stdout, stderr, err := ipsetting.RouteAddIPWithDev("default", *config.TunName); err != nil || code != 0 {
		logo.E(stdout, stderr)
		return nil, err
	}

	if code, stdout, stderr, err := ipsetting.RouteAddIPWithGateway(config.GetProxyHost(), *config.DefaultGateway); err != nil || code != 0 {
		logo.E(stdout, stderr)
		return nil, err
	}

	if code, stdout, stderr, err := ipsetting.RouteAddIPsWithGateway(config.GetTunDNS(), *config.DefaultGateway); err != nil || code != 0 {
		logo.E(stdout, stderr)
		return nil, err
	}

	return dev, err
}
