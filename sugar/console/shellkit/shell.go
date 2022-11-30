package shellkit

import (
	"bytes"
	"github.com/wsrf16/swiss/sugar/encoding/encodekit"
	"os/exec"
	"runtime"
)

type Result struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Err    error  `json:"err,omitempty"`
}
type ResultTotal struct {
	Results []*Result `json:"results,omitempty"`
	Err     string    `json:"err,omitempty"`
}

func (t *ResultTotal) Append(result Result) {
	t.Results = append(t.Results, &result)
}

func ExecuteBatch(commands []string) (ResultTotal, error) {
	results := make([]*Result, 0, 128)

	for _, v := range commands {
		stdout, stderr, err := Execute(v)
		result := &Result{Stdout: stdout, Stderr: stderr, Err: err}
		results = append(results, result)
		if err != nil {
			return ResultTotal{Results: results, Err: err.Error()}, err
		}
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

func Execute(command string) (string, string, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("/bin/bash", "-c", command)
	case "windows":
		cmd = exec.Command("cmd", "/C", command)
	}

	stdoutBuffer := new(bytes.Buffer)
	stderrBuffer := new(bytes.Buffer)
	cmd.Stdout, cmd.Stderr = stdoutBuffer, stderrBuffer

	var stdout, stderr string
	cmderr := cmd.Run()

	// if s, err := decode(stdoutBuffer.Bytes()); err != nil {
	//     return &s, nil, err
	// } else {
	//     stdout = s
	// }
	//
	// if s, err := decode(stderrBuffer.Bytes()); err != nil {
	//     return &s, nil, err
	// } else {
	//     stderr = s
	// }

	stdout, _ = decode(stdoutBuffer.Bytes())
	stderr, _ = decode(stderrBuffer.Bytes())

	return stdout, stderr, cmderr
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
