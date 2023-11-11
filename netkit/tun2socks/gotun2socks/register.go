package gotun2socks

import (
	"github.com/wsrf16/swiss/netkit/tun2socks/gotun2socks/core"
	"github.com/wsrf16/swiss/netkit/tun2socks/support/t2sconfig"
	"github.com/wsrf16/swiss/netkit/tun2socks/support/tunsocks"
)

func registerHandler(config *t2sconfig.TunConfig) error {
	if err := registerTun(config); err != nil {
		return err

	}

	return nil
}

func registerTun(config *t2sconfig.TunConfig) error {
	dev, err := tunsocks.CreateAndOpenTunDevice(*config)
	if err != nil {
		return err
	}

	tun := core.New(dev, *config.ProxyServer, config.GetTunDNS(), *config.PublicOnly, *config.EnableDNSCache)
	go tun.Run()
	return nil
}
