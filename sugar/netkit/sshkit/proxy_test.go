package sshkit

import "testing"

func TestTunnelServe(t *testing.T) {
	property := &SSHProperty{Id: "1", Addr: "mecs.com:22", User: "root", Password: "1"}
	err := TunnelServe("0.0.0.0:8081", property, "192.168.0.133:22")
	if err != nil {
		t.Error(err)
	}
}
