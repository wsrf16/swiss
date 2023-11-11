//go:build !windows
// +build !windows

package blocker

import (
	"errors"
)

func BlockOutsideDNS(tunName string) error {
	return errors.New("not implemented")
}
