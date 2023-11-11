package tunsocks

import (
	"github.com/wsrf16/swiss/netkit/settings/hardware/tun/watertun"
	"github.com/wsrf16/swiss/netkit/tun2socks/support/t2sconfig"
	"io"
)

func CreateAndOpenTunDevice(config t2sconfig.TunConfig) (io.ReadWriteCloser, error) {
	return watertun.CreateAndOpenTunDevice(*config.TunName, *config.TunAddress, *config.TunMask, *config.TunGateway, config.GetTunDNS())
}
