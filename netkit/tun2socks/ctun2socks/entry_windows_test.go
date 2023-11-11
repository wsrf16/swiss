package ctun2socks

import (
	"flag"
	"github.com/wsrf16/swiss/netkit/tun2socks/support/t2sconfig"
	"testing"
)

func TestConvert(t *testing.T) {
	var config = new(t2sconfig.TunConfig)
	config.TunName = flag.String("tunName", "MyNIC", "TUN interface name")
	config.TunAddress = flag.String("tunAddr", "10.255.0.2", "TUN interface address")
	config.TunGateway = flag.String("tunGw", "10.255.0.1", "TUN interface gateway")
	config.TunMask = flag.String("tunMask", "255.255.255.0", "TUN interface netmask, it should be a prefixlen (a number) for IPv6 address")
	config.TunDNS = flag.String("tunDNS", "8.8.8.8,8.8.4.4", "DNS resolvers for TUN interface (only need on Windows)")
	config.TunPersist = flag.Bool("tunPersist", false, "Persist TUN interface after the program exits or the last open file descriptor is closed (Linux only)")
	config.BlockOutsideDNS = flag.Bool("blockOutsideDNS", false, "Prevent DNS leaks by blocking plaintext DNS queries going out through non-TUN interface (may require admin privileges) (Windows only) ")
	config.ProxyType = flag.String("proxyType", "socks", "Proxy handler type: socks/redirect")
	config.ProxyServer = flag.String("proxyServer", "127.0.0.1:1080", "proxyServer")
	config.DNSFallback = flag.Bool("dnsFallback", false, "Enable DNS fallback over TCP (overrides the UDP proxy handler).")
	flag.Parse()
	{
		// err := gotun2socks.Serve(config)
		err := Serve(config)
		if err != nil {
			t.Error(err)
		}
	}
}
