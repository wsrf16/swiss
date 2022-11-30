package regexpkit

import "regexp"

func ReplaceAll(pattern string, src, repl string) (string, error) {
	if compile, err := regexp.Compile(pattern); err != nil {
		return "", err
	} else {
		return string(compile.ReplaceAll([]byte(src), []byte(repl))), nil
	}
}

func FindStringSubmatch(pattern string, src string) ([]string, error) {
	if compile, err := regexp.Compile(pattern); err != nil {
		return nil, err
	} else {
		return compile.FindStringSubmatch(src), nil
	}
}
