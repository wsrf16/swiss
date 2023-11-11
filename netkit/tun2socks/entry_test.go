package tun2socks

import (
	"github.com/wsrf16/swiss/netkit/tun2socks/support/t2sconfig"
	"github.com/wsrf16/swiss/sugar/base/pointkit"
	"github.com/wsrf16/swiss/sugar/logo"
	"testing"
)

func TestServe(t *testing.T) {
	// -tunName 本地连接1 -tunAddr 10.255.0.8 -proxyType socks -proxyServer micro.com:1080
	var config = new(t2sconfig.TunConfig)
	config.TunAddress = pointkit.ToPoint("10.255.0.4/24")
	config.TunName = pointkit.ToPoint("MYNIC")
	//config.ProxyServer = pointkit.ToPoint("micro.com:1080")
	config.ProxyServer = pointkit.ToPoint("192.168.0.103:1080")
	{
		//err := gotun2socks.Serve(config)
		//err := ctun2socks.Serve(config)
		err := Serve(config)
		if err != nil {
			logo.Error("", err)
		}
	}
}
