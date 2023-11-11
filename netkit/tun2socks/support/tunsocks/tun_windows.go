package tunsocks

import (
	"github.com/wsrf16/swiss/netkit/settings/hardware/tun/wintun"
	"github.com/wsrf16/swiss/netkit/settings/windows/iface"
	"github.com/wsrf16/swiss/netkit/tun2socks/support/t2sconfig"
	"io"
)

//func CreateAndOpenTunDevice(config support.TunConfig) (io.ReadWriteCloser, error) {
//    var dev io.ReadWriteCloser
//    var err error
//    if *config.TunMode == support.TunModeWinSyscall {
//        dev, err = syscalltun.OpenTunDevice(*config.TunName, *config.TunAddress, *config.TunMask, *config.TunGateway, config.GetTunDNS())
//        if err != nil {
//            return nil, err
//        }
//    } else {
//        dev, err = createAndOpenTunDevice(*config.TunName, *config.TunAddress, *config.TunMask, *config.TunGateway, *config.TunGWMetric, *config.MTU, config.GetTunDNS())
//        if err != nil {
//            return nil, err
//        }
//    }
//    return dev, err
//}
//
//func createAndOpenTunDevice(name, addr, mask, gw string, gwmetric, mtu int, dns []string) (io.ReadWriteCloser, error) {
//    dev, err := wintun.CreateAndOpenTunDevice(name, mtu)
//    if err != nil {
//        return nil, err
//    }
//
//    _, err = iface.SetAddressAndDNSStatic(name, addr, mask, gw, gwmetric, dns)
//    if err != nil {
//        return nil, err
//    }
//    return dev, err
//}

func CreateAndOpenTunDevice(config t2sconfig.TunConfig) (io.ReadWriteCloser, error) {
	dev, err := wintun.CreateAndOpenTunDevice(*config.TunName, *config.MTU)
	if err != nil {
		return nil, err
	}

	_, err = iface.SetAddressAndDNSStatic(*config.TunName, *config.TunAddress, *config.TunMask, *config.TunGateway, *config.TunGWMetric, config.GetTunDNS())
	if err != nil {
		return nil, err
	}
	return dev, err
}
