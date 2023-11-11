package sshkit

import (
	"github.com/wsrf16/swiss/sugar/crypto/rsakit"
	"golang.org/x/crypto/ssh"
	"time"
)

type SSHProperty struct {
	Id         string `yaml:"id"`
	Addr       string `yaml:"addr"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	PublicKey  string `yaml:"publicKey"`
	PrivateKey string `yaml:"privateKey"`
}

func (p SSHProperty) GetPrivateKey() (string, error) {
	key := p.PrivateKey
	replaced, err := rsakit.FormatKey(key)
	return replaced, err
}

func (p SSHProperty) GenSSHClientConfig() (*ssh.ClientConfig, error) {
	privateKey, err := p.GetPrivateKey()
	if err != nil {
		return nil, err
	}
	authMethods := make([]ssh.AuthMethod, 0, 2)
	if len(p.Password) > 0 {
		authMethods = append(authMethods, ssh.Password(p.Password))
	}

	if len(privateKey) > 0 {
		if signer, err := ssh.ParsePrivateKey([]byte(privateKey)); err != nil {
			return nil, err
		} else {
			authMethods = append(authMethods, ssh.PublicKeys(signer))
		}
	}

	clientConfig := &ssh.ClientConfig{
		Timeout: time.Second,
		User:    p.User,
		Auth:    authMethods,
		// Auth: []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以, 但是不够安全
		//HostKeyCallback: ssh.FixedHostKey(signera.PublicKey()),
		//HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		//    // key: server ecdsa-sha2-nistp256
		//    f := signer.PublicKey()
		//
		//    if f == nil {
		//        return fmt.Errorf("ssh: required host key was nil")
		//    }
		//
		//    if !bytes.Equal(key.Marshal(), f.Marshal()) {
		//        return fmt.Errorf("ssh: host key mismatch")
		//    }
		//    return nil
		//},
	}
	return clientConfig, nil

}

func NewSSHConfigSpecifyId(id string, addr string) *SSHProperty {
	return &SSHProperty{Id: id, Addr: addr}
}

func NewSSHConfigDefaultId(addr string) *SSHProperty {
	id := addr
	return &SSHProperty{Id: id, Addr: addr}
}
