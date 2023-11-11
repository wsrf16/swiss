package ctun2socks

import (
	"errors"
	"fmt"
	"github.com/wsrf16/swiss/netkit/settings/hardware/tunkit/blocker"
	"github.com/wsrf16/swiss/netkit/settings/hardware/tunkit/dnsfallback"
	"github.com/wsrf16/swiss/netkit/tun2socks/ctun2socks/core"
	"github.com/wsrf16/swiss/netkit/tun2socks/ctun2socks/proxy/redirect"
	"github.com/wsrf16/swiss/netkit/tun2socks/ctun2socks/proxy/socks"
	"github.com/wsrf16/swiss/netkit/tun2socks/support/t2sconfig"
	"github.com/wsrf16/swiss/netkit/tun2socks/support/tunsocks"
	"github.com/wsrf16/swiss/sugar/logo"
	"io"
	"runtime"
)

func registerHandler(config *t2sconfig.TunConfig) error {

	if err := registerBlockOutsideDNS(config); err != nil {
		return err
	}

	registerProxy(config)

	registerDNSFallback(config)

	if err := registerTun(config); err != nil {
		return err

	}

	return nil
}

func registerProxy(config *t2sconfig.TunConfig) {
	// Register TCP and UDP handlers to handle accepted connections.
	if *config.ProxyType == "redirect" {
		core.RegisterTCPConnHandler(redirect.NewTCPHandler(*config.ProxyServer))
		core.RegisterUDPConnHandler(redirect.NewUDPHandler(*config.ProxyServer, *config.UDPTimeout))
	} else {
		//proxyAddr, err := net.ResolveTCPAddr("tcp", *config.ProxyServer)
		//if err != nil {
		//    logo.Ff("invalid proxy server address: %v", err)
		//}
		//proxyHost := proxyAddr.IP.String()
		//proxyPort := uint16(proxyAddr.Port)

		core.RegisterTCPConnHandler(socks.NewTCPHandler(*config.ProxyServer))
		core.RegisterUDPConnHandler(socks.NewUDPHandler(*config.ProxyServer, *config.UDPTimeout))
	}
}

func registerBlockOutsideDNS(config *t2sconfig.TunConfig) error {
	if *config.BlockOutsideDNS && runtime.GOOS == "windows" {
		if err := blocker.BlockOutsideDNS(*config.TunName); err != nil {
			logo.Ff("failed to block outside DNS: %v", err)
			return err
		}
	}
	return nil
}

func registerDNSFallback(config *t2sconfig.TunConfig) {
	if config.DNSFallback != nil && *config.DNSFallback {
		// Override the UDP handler with a DNS-over-TCP (fallback) UDP handler.
		core.RegisterUDPConnHandler(dnsfallback.NewUDPHandler())
	}
}

func registerTun(config *t2sconfig.TunConfig) error {
	dev, err := tunsocks.CreateAndOpenTunDevice(*config)
	if err != nil {
		return err
	}

	registerReaderWriter(dev, *config.MTU)
	return nil
}

func registerReaderWriter(dev io.ReadWriteCloser, mtu int) {
	// Register an output callback to write packets output from lwip stack to tun
	// device, output function should be set before input any packets.
	core.RegisterOutputFn(func(data []byte) (int, error) {
		return dev.Write(data)
	})

	// Setup TCP/IP stack.
	lwipWriter := core.NewLWIPStack().(io.Writer)

	// Copy packets from tun device to lwip stack, it's the main loop.
	go func() {
		_, err := io.CopyBuffer(lwipWriter, dev, make([]byte, mtu))
		if err != nil {
			errors.New(fmt.Sprintf("copying data failed: %v", err))
		}
	}()
}
