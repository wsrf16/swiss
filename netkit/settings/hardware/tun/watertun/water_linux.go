package watertun

import (
	"io"

	"github.com/songgao/water"
)

func CreateAndOpenTunDevice(name string, persist bool) (io.ReadWriteCloser, error) {
	cfg := water.Config{
		DeviceType: water.TUN,
	}
	cfg.Name = name
	cfg.Persist = persist
	tunDev, err := water.New(cfg)
	if err != nil {
		return nil, err
	}

	return tunDev, nil
}

func CreateAndOpenTapDevice(name string, persist bool) (io.ReadWriteCloser, error) {
	cfg := water.Config{
		DeviceType: water.TAP,
	}
	cfg.Name = name
	cfg.Persist = persist
	tunDev, err := water.New(cfg)
	if err != nil {
		return nil, err
	}

	return tunDev, nil
}
