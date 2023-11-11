package wintun

import (
	"golang.zx2c4.com/wireguard/tun"
	"io"
)

type WinTunDevice struct {
	dev tun.Device
}

func (d WinTunDevice) Read(p []byte) (n int, err error) {
	return d.dev.Read(p, 0)
}

func (d WinTunDevice) Write(p []byte) (n int, err error) {
	return d.dev.Write(p, 0)
}

func (d WinTunDevice) Close() (err error) {
	return d.dev.Close()
}

func CreateAndOpenTunDevice(name string, mtu int) (io.ReadWriteCloser, error) {
	dev, err := tun.CreateTUN(name, mtu)
	if err != nil {
		return nil, err
	}
	wdev := &WinTunDevice{dev: dev}
	return wdev, nil
}
