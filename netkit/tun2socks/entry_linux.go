package tun2socks

import (
	"github.com/wsrf16/swiss/netkit/tun2socks/gotun2socks"
	"github.com/wsrf16/swiss/netkit/tun2socks/support/t2sconfig"
)

func Serve(config *t2sconfig.TunConfig) error {
	return gotun2socks.Serve(config)
}
