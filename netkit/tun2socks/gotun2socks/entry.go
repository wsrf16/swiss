package gotun2socks

import (
	"github.com/wsrf16/swiss/netkit/tun2socks/support/loader"
	"github.com/wsrf16/swiss/netkit/tun2socks/support/t2sconfig"
	"github.com/wsrf16/swiss/sugar/logo"
	"os"
	"os/signal"
	"syscall"
)

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
