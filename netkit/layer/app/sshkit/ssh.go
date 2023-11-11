package sshkit

import (
	"github.com/wsrf16/swiss/sugar/console/shellkit"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"time"
)

func BuildClientByProperty(property *SSHProperty) (*ssh.Client, error) {
	config, err := property.GenSSHClientConfig()
	if err != nil {
		return nil, err
	}

	client, err := ssh.Dial("tcp", property.Addr, config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func BuildClientConfig(user string, password string, privateKeys []string) (*ssh.ClientConfig, error) {
	// ClientConfig -> Client -> Conn
	authMethods := make([]ssh.AuthMethod, 0, 2)
	if len(password) > 0 {
		authMethods = append(authMethods, ssh.Password(password))
	}

	for _, key := range privateKeys {
		if signer, err := ssh.ParsePrivateKey([]byte(key)); err != nil {
			log.Println(err)
		} else {
			authMethods = append(authMethods, ssh.PublicKeys(signer))
		}
	}

	clientConfig := &ssh.ClientConfig{
		Timeout: time.Second,
		User:    user,
		Auth:    authMethods,
		// Auth: []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以, 但是不够安全
	}
	return clientConfig, nil
}

func Run(client *ssh.Client, cmd string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	return string(output), err
}

func RunBatch(client *ssh.Client, cmds []string) (*shellkit.ResultTotal, error) {
	total := new(shellkit.ResultTotal)
	for _, cmd := range cmds {
		if stdout, err := Run(client, cmd); err != nil {
			total.Append(shellkit.Result{Stdout: stdout, Stderr: err.Error()})
		} else {
			total.Append(shellkit.Result{Stdout: stdout, Stderr: ""})
		}
	}
	return total, nil
}

func SimplePty(client *ssh.Client) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	return SimplePtySession(session)
}

func SimplePtyFlat(addr string, user string, pwd string, prik string) error {
	property := &SSHProperty{User: user, Password: pwd, Addr: addr, PrivateKey: prik}
	client, err := BuildClientByProperty(property)
	if err != nil {
		panic(err)
	}
	return SimplePty(client)
}

func SimplePtyByProperty(property *SSHProperty) error {
	client, err := BuildClientByProperty(property)
	if err != nil {
		return err
	}
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	return SimplePtySession(session)
}

func PtyClient(client *ssh.Client, term string, h, w int, modes ssh.TerminalModes) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	return PtySession(session, term, h, w, modes)
}

func SimplePtySession(session *ssh.Session) error {
	// 设置Terminal Mode
	term := "linux"
	h, w := 32, 160
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // 关闭回显
		ssh.TTY_OP_ISPEED: 14400, // 设置传输速率
		ssh.TTY_OP_OSPEED: 14400,
	}

	return PtySession(session, term, h, w, modes)
}

func PtySession(session *ssh.Session, term string, h, w int, modes ssh.TerminalModes) error {
	// 请求伪终端
	// err = session.Pty("xterm-256color", termWidth, termHeight, modes)
	// err := session.Pty("linux", 32, 160, modes)

	if err := session.RequestPty(term, h, w, modes); err != nil {
		return err
	}

	if session.Stdin == nil {
		session.Stdin = os.Stdin
	}
	if session.Stdout == nil {
		session.Stdout = os.Stdout
	}
	if session.Stderr == nil {
		session.Stderr = os.Stderr
	}
	session.Shell()
	session.Wait()

	return nil
}
