package ctun2socks

import (
	"github.com/wsrf16/swiss/netkit/tun2socks/support/loader"
	"github.com/wsrf16/swiss/netkit/tun2socks/support/t2sconfig"
	"github.com/wsrf16/swiss/sugar/logo"
	"os"
	"os/signal"
	"syscall"
)

// Serve
// -tunName 本地连接 -tunAddr 10.255.0.2 -tunGw 10.255.0.1 -tunMask 255.255.255.0 -tunDns 8.8.8.8,8.8.4.4 -proxyType socks -proxyServer 192.168.0.103:1080
// -tunName 本地连接 -proxyType socks -proxyServer 192.168.0.103:1080
func Serve(config *t2sconfig.TunConfig) error {
	config.FillWithDefaultValue()
	defer func() {
		loader.Unload(config)
	}()

	if err := registerHandler(config); err != nil {
		return err
	}

	logo.I("Running tun2socks")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch,
		os.Interrupt,
		os.Kill,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	logo.Info("Program tun2socks exit.", "", "Got signal:", <-ch)
	return nil
}
