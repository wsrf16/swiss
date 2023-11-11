package sshkit

import "testing"

func TestPty(t *testing.T) {
	property := &SSHProperty{Id: "1", Addr: "mecs.com:22", User: "root", Password: "1"}
	//config, err := property.GenSSHClientConfig()
	client, err := BuildClientByProperty(property)
	if err != nil {
		t.Error(err.Error())
	}

	err = SimplePty(client)
	if err != nil {
		t.Error(err.Error())
	}

}
