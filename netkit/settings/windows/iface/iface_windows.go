package iface

import (
	"fmt"
	"github.com/wsrf16/swiss/sugar/base/stringkit"
	"github.com/wsrf16/swiss/sugar/console/shellkit"
)

func SetAddressDHCP(devName string) (string, error) {
	code, stdout, stderr, err := shellkit.Execute("netsh", "interface", "ip", "set", "address", devName, "dhcp")
	if code < 0 {
		return stdout, shellkit.NewError(stderr, err)
	}

	return stdout, nil
}

func SetDNSDHCP(devName string) (string, error) {
	code, stdout, stderr, err := shellkit.Execute("netsh", "interface", "ip", "set", "dns", devName, "dhcp")
	if code < 0 {
		return stdout, shellkit.NewError(stderr, err)
	}

	return stdout, nil
}

func SetAddressAndDNSDHCP(devName string) (string, error) {
	cmd1, err := SetAddressDHCP(devName)
	if err != nil {
		return cmd1, err
	}

	cmd2, err := SetDNSDHCP(devName)
	if err != nil {
		return cmd2, err
	}

	return cmd1 + stringkit.Newline + cmd2, err
}

func SetAddressStatic(devName, ip, mask, gateway string, gwmetric int) (string, error) {
	// parseCIDR, ipNet, err := net.ParseCIDR(cidr)
	// if err != nil {
	//     return err
	// }
	// "netsh interface ip set address name=\"MyTUN\" source=static addr=ip mask=mask gateway=gateway gwmetric=0"
	args := make([]string, 0)
	args = append(args, "netsh", "interface", "ip", "set", "address")

	args = append(args, fmt.Sprintf("source=static"))
	args = append(args, fmt.Sprintf("name=\"%v\"", devName))

	if len(ip) > 0 {
		args = append(args, fmt.Sprintf("addr=%v", ip))
	}
	if len(mask) > 0 {
		args = append(args, fmt.Sprintf("mask=%v", mask))
	}
	if len(gateway) > 0 {
		args = append(args, fmt.Sprintf("gateway=%v", gateway))
	}
	if gwmetric > 0 {
		args = append(args, fmt.Sprintf("gwmetric=%v", gwmetric))
	}

	// stdout, stderr, err := shellkit.Execute("netsh", "interface", "ip", "set", "address", builder.String())
	code, stdout, stderr, err := shellkit.Execute(args...)
	if code < 0 {
		return stdout, shellkit.NewError(stderr, err)
	}
	return stdout, nil
}

func SetDNSStatic(devName string, dns []string) (string, error) {
	// parseCIDR, ipNet, err := net.ParseCIDR(cidr)
	// if err != nil {
	//     return err
	// }
	// "netsh interface dns set dns name=\"MyTUN\" source=static addr=dns register=PRIMARY"
	args1 := make([]string, 0)
	args1 = append(args1, "netsh", "interface", "ip", "set", "dnsservers")
	args1 = append(args1, fmt.Sprintf("name=\"%v\"", devName))
	args1 = append(args1, fmt.Sprintf("source=static"))
	args1 = append(args1, fmt.Sprintf("addr=%v", dns[0]))
	args1 = append(args1, fmt.Sprintf("register=%v", "PRIMARY"))

	code, stdout, stderr, err := shellkit.Execute(args1...)
	if code < 0 {
		return stdout, shellkit.NewError(stderr, err)
	}

	if len(dns) > 1 {
		args2 := make([]string, 0)
		args2 = append(args2, "netsh", "interface", "ip", "add", "dnsservers")
		args2 = append(args2, fmt.Sprintf("name=\"%v\"", devName))
		// args2 = append(args2, fmt.Sprintf("source=static"))
		args2 = append(args2, fmt.Sprintf("addr=%v", dns[1]))
		args2 = append(args2, fmt.Sprintf("index=%v", 2))

		code, stdout, stderr, err = shellkit.Execute(args2...)
		if code < 0 {
			return stdout, shellkit.NewError(stderr, err)
		}
	}
	return stdout, nil
}

func SetAddressAndDNSStatic(devName, ip, mask, gateway string, gwmetric int, dns []string) (string, error) {
	cmd1, err := SetDNSStatic(devName, dns)
	if err != nil {
		return cmd1, err
	}

	cmd2, err := SetAddressStatic(devName, ip, mask, gateway, gwmetric)
	if err != nil {
		return cmd2, err
	}

	return cmd1 + stringkit.Newline + cmd2, err
}

func Enable(devName string) (string, error) {
	code, stdout, stderr, err := shellkit.ExecuteSingleLine("netsh interface set interface " + devName + " admin=enable")
	if code < 0 {
		return stdout, shellkit.NewError(stderr, err)
	}

	return stdout, nil
}

func Disable(devName string) (string, error) {
	code, stdout, stderr, err := shellkit.ExecuteSingleLine("netsh interface set interface " + devName + " admin=disable")
	if code < 0 {
		return stdout, shellkit.NewError(stderr, err)
	}

	return stdout, nil
}
