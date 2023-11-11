package ipsetting

import (
	"github.com/wsrf16/swiss/sugar/console/shellkit"
	"strconv"
)

// ip route
func RouteAddIPWithDev(ip, name string) (int, string, string, error) {
	return shellkit.ExecuteSingleLine("/sbin/ip route add " + ip + " dev " + name)
}

func RouteAddIPWithDevAndMetric(ip, name string, metric int) (int, string, string, error) {
	return shellkit.ExecuteSingleLine("/sbin/ip route add " + ip + " dev " + name + " metric " + strconv.Itoa(metric))
}

func RouteAddIPWithGateway(ip, gw string) (int, string, string, error) {
	return shellkit.ExecuteSingleLine("/sbin/ip route add " + ip + " via " + gw)
}

func RouteAddIPsWithGateway(ip []string, gw string) (code int, stdout string, stderr string, err error) {
	for _, item := range ip {
		code, stdout, stderr, err := RouteAddIPWithGateway(item, gw)
		if err != nil || code != 0 {
			return code, stdout, stderr, err
		}
	}
	return code, stdout, stderr, err
}

func RouteAddIPWithGatewayAndMetric(ip, gw string, metric int) (int, string, string, error) {
	return shellkit.ExecuteSingleLine("/sbin/ip route add " + ip + " via " + gw + " metric" + strconv.Itoa(metric))
}

func RouteDelIP(ip string) (int, string, string, error) {
	return shellkit.ExecuteSingleLine("/sbin/ip route delete " + ip)
}

func RouteDelIPs(ip []string) (code int, stdout string, stderr string, err error) {
	for _, item := range ip {
		code, stdout, stderr, err := RouteDelIP(item)
		if err != nil || code != 0 {
			return code, stdout, stderr, err
		}
	}
	return code, stdout, stderr, err
}

func RouteDelIPWithGateway(ip, gw string) (int, string, string, error) {
	return shellkit.ExecuteSingleLine("/sbin/ip route delete " + ip + " via " + gw)
}

func RouteDelIPsWithGateway(ip []string, gw string) (code int, stdout string, stderr string, err error) {
	for _, item := range ip {
		code, stdout, stderr, err := RouteDelIPWithGateway(item, gw)
		if err != nil || code != 0 {
			return code, stdout, stderr, err
		}
	}
	return code, stdout, stderr, err
}

func RouteDelIPWithDev(ip, name string) (int, string, string, error) {
	return shellkit.ExecuteSingleLine("/sbin/ip route delete " + ip + " dev " + name)
}

func RouteDelIPsWithDev(ip []string, name string) (code int, stdout string, stderr string, err error) {
	for _, item := range ip {
		code, stdout, stderr, err := RouteDelIPWithDev(item, name)
		if err != nil || code != 0 {
			return code, stdout, stderr, err
		}
	}
	return code, stdout, stderr, err
}

// ip addr
func AddrAddIPAndCIDR(name, cidr string) (int, string, string, error) {
	return shellkit.ExecuteSingleLine("/sbin/ip addr add dev " + name + " " + cidr)
	// Both
	//return shellkit.ExecuteSingleLine("/sbin/ip addr add dev ens33 192.168.0.1/24")
	//return shellkit.ExecuteSingleLine("/sbin/ip addr 192.168.0.1/24 add dev ens33")
}

func AddrDelIPAndCIDR(name, cidr string) (int, string, string, error) {
	return shellkit.ExecuteSingleLine("/sbin/ip addr del dev " + name + " " + cidr)
}

func AddrShow() (int, string, string, error) {
	return shellkit.ExecuteSingleLine("/sbin/ip addr show")
}

func AddrList() (int, string, string, error) {
	return shellkit.ExecuteSingleLine("/sbin/ip addr list")
}

// ip link
func LinkSetMTU(name string, mtu int) (int, string, string, error) {
	return shellkit.ExecuteSingleLine("/sbin/ip link set dev " + name + " mtu " + strconv.Itoa(mtu))
}

func LinkRename(name, rename string) (int, string, string, error) {
	return shellkit.ExecuteSingleLine("/sbin/ip link set dev " + name + " name " + rename)
}

func LinkSetDevUp(name string) (int, string, string, error) {
	return shellkit.ExecuteSingleLine("/sbin/ip link set dev " + name + " up")
}

func LinkSetDevDown(name string) (int, string, string, error) {
	return shellkit.ExecuteSingleLine("/sbin/ip link set dev " + name + " down")
}

// ip tuntap
func TunTapCreateTAP(name string) (int, string, string, error) {
	return shellkit.ExecuteSingleLine("/sbin/ip tuntap add dev " + name + " mode tap")
}

func TunTapCreateTUN(name string) (int, string, string, error) {
	return shellkit.ExecuteSingleLine("/sbin/ip tuntap add dev " + name + " mode tun")
}

func TunTapDelTAP(name string) (int, string, string, error) {
	return shellkit.ExecuteSingleLine("/sbin/ip tuntap del dev " + name + " mode tap")

}

func TunTapDelTUN(name string) (int, string, string, error) {
	return shellkit.ExecuteSingleLine("/sbin/ip tuntap del dev " + name + " mode tun")
}

// func TunTapCreateTAP(name string, mode string) (int, string, string, error) {
//    return shellkit.ExecuteSingleLine("/sbin/ip add dev " + name + " mod tap")
// }
//
// func TunTapCreateTUN(name string, mode string) (int, string, string, error) {
//    return shellkit.ExecuteSingleLine("/sbin/ip add dev " + name + " mod tun")
// }
//
// func TunTapDelTAP(name string, mode string) (int, string, string, error) {
//    return shellkit.ExecuteSingleLine("/sbin/ip add dev " + name + " mod tap")
// }
//
// func TunTapDelTUN(name string, mode string) (int, string, string, error) {
//    return shellkit.ExecuteSingleLine("/sbin/ip add dev " + name + " mod tun")
// }
