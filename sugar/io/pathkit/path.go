package pathkit

import (
	"github.com/wsrf16/swiss/sugar/base/stringkit"
	"os"
	"path/filepath"
	"strings"
)

func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func GetWorkDirectory() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

func GetRecursiveFileList(basePath string) ([]string, error) {
	pathLists := make([]string, 16, 64)
	return pathLists, filepath.Walk(basePath,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if !f.IsDir() {
				pathLists = append(pathLists, path)
			}
			return nil
		})
}

func GetRecursiveDirectoryList(basePath string) ([]string, error) {
	pathLists := make([]string, 16, 64)
	return pathLists, filepath.Walk(basePath,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				pathLists = append(pathLists, path)
			}
			return nil
		})
}

func GetPWD(relatives ...string) string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	relative := Join(relatives...)
	return Join(pwd, relative)
}

func JoinBy(delim string, parts ...string) string {
	combine := strings.Builder{}
	for i, part := range parts {
		if i == 0 {
			combine.WriteString(stringkit.TrimAllSuffixes(part, "/", "\\"))
		} else {
			combine.WriteString(delim + stringkit.TrimAll(part, "/", "\\"))
		}

	}
	return combine.String()
}

func Join(parts ...string) string {
	delim := string(os.PathSeparator)
	return JoinBy(delim, parts...)
}

func GetSelfDirAndFileName() (string, string) {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}

	dir, file := filepath.Split(path)
	if err != nil {
		panic(err)
	}

	return dir, file
}
