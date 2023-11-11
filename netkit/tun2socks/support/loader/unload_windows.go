package loader

import (
	"github.com/wsrf16/swiss/netkit/settings/windows/route"
	"github.com/wsrf16/swiss/netkit/tun2socks/support/t2sconfig"
	"github.com/wsrf16/swiss/sugar/logo"
)

func Unload(config *t2sconfig.TunConfig) error {
	execute, err := route.Delete("0.0.0.0", "0.0.0.0", *config.TunGateway)
	if err != nil {
		logo.W("failed to delete route \""+*config.TunGateway+"\"'", err)
	} else {
		logo.I(execute)
	}

	//if code != 0 || err != nil {
	//    logo.W("failed to delete route \"" + *config.TunDNS + "\"'", err, stderr)
	//} else {
	//    logo.I(stdout)
	//}

	return nil
}
