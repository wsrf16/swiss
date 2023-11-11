package shellkit

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/wsrf16/swiss/sugar/encoding/encodekit"
	"github.com/wsrf16/swiss/sugar/logo"
	"os/exec"
	"runtime"
	"strings"
)

type Result struct {
	Code   int    `json:"code"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Error  error  `json:"err,omitempty"`
}
type ResultTotal struct {
	Results []*Result `json:"results,omitempty"`
}

func (t *ResultTotal) Append(result Result) {
	t.Results = append(t.Results, &result)
}

func (t *ResultTotal) OK() bool {
	for _, result := range t.Results {
		if result.Code != 0 {
			return false
		}
	}
	return true
}

func ExecuteBatch(commands []string) (ResultTotal, error) {
	results := make([]*Result, 0, 128)

	for _, v := range commands {

		code, stdout, stderr, err := Execute(v)
		results = append(results, &Result{Code: code, Stdout: stdout, Stderr: stderr, Error: err})

		// if err != nil {
		//    return ResultTotal{Results: results, Err: err.Error()}, err
		// }
	}
	return ResultTotal{Results: results}, nil
}

// func ExecuteBatch(commands []string) ([]*string, []*string, error) {
// 	stdouts := make([]*string, 0, 128)
// 	stderrs := make([]*string, 0, 128)
//
// 	for _, v := range commands {
// 		stdout, stderr, err := Execute(v)
// 		stdouts = append(stdouts, stdout)
// 		stderrs = append(stderrs, stderr)
// 		if err != nil {
// 			return stdouts, stderrs, err
// 		}
// 	}
// 	return stdouts, stderrs, nil
// }

func NewError(stderr string, err error) error {
	var exerr string
	if err == nil {
		exerr = ""
	} else {
		exerr = err.Error()
	}
	return errors.New(fmt.Sprintf("%v|%v", stderr, exerr))
}

func ExecuteSingleLine(cmd string) (code int, stdout string, stderr string, err error) {
	cmds := strings.Split(cmd, " ")
	return Execute(cmds...)
}

func Execute(commands ...string) (code int, stdout string, stderr string, err error) {
	join := strings.Join(commands, " ")
	logo.Debug("", "", join)
	var cmd *exec.Cmd
	// switch runtime.GOOS {
	// case "linux":
	//    commands := append([]string{"-c"}, commands...)
	//    cmd = exec.Command("/bin/bash", commands...)
	// case "windows":
	//    commands := append([]string{"/C"}, commands...)
	//    cmd = exec.Command("cmd", commands...)
	//    // commands := command
	// }
	cmd = exec.Command(commands[0], commands[1:]...)

	stdoutBuffer := new(bytes.Buffer)
	stderrBuffer := new(bytes.Buffer)
	cmd.Stdout, cmd.Stderr = stdoutBuffer, stderrBuffer

	err = cmd.Run()
	// stdout, err := cmd.Output()
	stdout, _ = decode(stdoutBuffer.Bytes())
	stderr, _ = decode(stderrBuffer.Bytes())
	code = cmd.ProcessState.ExitCode()

	return code, stdout, stderr, err
}

func decode(b []byte) (string, error) {
	switch runtime.GOOS {
	case "windows":
		s, err := encodekit.TransformToUTF8TextFromGBK(b)
		return s, err
	default:
		return string(b), nil
	}
}
