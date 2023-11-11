package sh

import (
	"fmt"
	"github.com/wsrf16/swiss/sugar/console/shellkit"
	"github.com/wsrf16/swiss/sugar/io/iokit"
	"github.com/wsrf16/swiss/sugar/logo"
)

func Initial(host string) (int, string, string, error) {
	var json = "{\"insecure-registries\":[\"" + host + "\"]}"
	iokit.WriteToFile([]byte(json), "/etc/docker/daemon.json")
	fmt.Println(json)

	code, stdout, stderr, err := shellkit.ExecuteSingleLine("systemctl daemon-reload")
	if code != 0 {
		logo.E(stdout, stderr)
		return code, stdout, stderr, err
	}
	code, stdout, stderr, err = shellkit.ExecuteSingleLine("systemctl restart docker")
	if code != 0 {
		logo.E(stdout, stderr)
		return code, stdout, stderr, err
	}
	return 0, "", "", nil
}

func Login(host, username, password string) (int, string, string, error) {
	return shellkit.ExecuteSingleLine(fmt.Sprintf("docker login --username=%s -p=\"%s\" %s", username, password, host))
}

func Pull(image string) (int, string, string, error) {
	return shellkit.ExecuteSingleLine(fmt.Sprintf("docker pull %s", image))
}

func Save(image, filename string) (int, string, string, error) {
	return shellkit.ExecuteSingleLine(fmt.Sprintf("docker save %s -o %s", image, filename))
}
