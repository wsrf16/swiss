package route

import (
	"fmt"
	"github.com/wsrf16/swiss/sugar/console/shellkit"
)

func Add(destination, mask, gateway string, metric int) (string, error) {
	// args := make([]string, 0)
	// args = append(args, "route")
	// args = append(args, fmt.Sprintf("add %v", destination))
	// args = append(args, fmt.Sprintf("mask %v %v", mask, gateway))
	// args = append(args, fmt.Sprintf("metric %v", metric))
	// if persist {
	//     args = append(args, fmt.Sprintf("-p"))
	// }

	cmd := fmt.Sprintf("route add %v mask %v %v metric %v", destination, mask, gateway, metric)
	code, stdout, stderr, err := shellkit.ExecuteSingleLine(cmd)
	if code < 0 {
		return stdout, shellkit.NewError(stderr, err)
	}

	return stdout, nil
}

func Delete(destination, mask, gateway string) (string, error) {
	// args := make([]string, 0)
	// args = append(args, "route")
	// args = append(args, fmt.Sprintf("delete %v", destination))
	// args = append(args, fmt.Sprintf("mask %v %v", mask, gateway))

	cmd := fmt.Sprintf("route delete %v mask %v %v", destination, mask, gateway)
	code, stdout, stderr, err := shellkit.ExecuteSingleLine(cmd)
	if code < 0 {
		return stdout, shellkit.NewError(stderr, err)
	}
	return stdout, nil
}

func PrintIPV4(destination, mask, gateway string) (string, error) {
	cmd := fmt.Sprintf("route print -4")
	code, stdout, stderr, err := shellkit.ExecuteSingleLine(cmd)
	if code < 0 {
		return stdout, shellkit.NewError(stderr, err)
	}
	return stdout, nil
}

func PrintIPV6(destination, mask, gateway string) (string, error) {
	cmd := fmt.Sprintf("route print -6")
	code, stdout, stderr, err := shellkit.ExecuteSingleLine(cmd)
	if code < 0 {
		return stdout, shellkit.NewError(stderr, err)
	}
	return stdout, nil
}
