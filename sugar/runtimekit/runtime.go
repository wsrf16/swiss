package runtimekit

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func GetCurrentPackage() string {
	_, file, _, _ := runtime.Caller(1)
	pkg := strings.TrimSuffix(file, filepath.Base(file))
	pkg = strings.TrimSuffix(pkg, "/")
	pkg = filepath.Base(pkg)
	return pkg
}

func GetCurrentFile() string {
	_, file, _, _ := runtime.Caller(1)
	return file
}

func GetStackTraceLines(s ...any) []string {
	lines := make([]string, 0, 16)

	for _, v := range s {
		lines = append(lines, fmt.Sprintf("\"%v\"", v))
	}
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		funcName := runtime.FuncForPC(pc).Name()
		lines = append(lines, fmt.Sprintf("    at %s:%s:%d (0x%x)", file, funcName, line, pc))
	}
	return lines
}

type StackTraceInfo struct {
	ProgramCounter uintptr
	File           string
	Function       string
	Line           int
}

func GetStackTraceInfos() []StackTraceInfo {
	infos := make([]StackTraceInfo, 0, 16)

	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		function := runtime.FuncForPC(pc).Name()
		info := StackTraceInfo{ProgramCounter: pc, File: file, Line: line, Function: function}

		infos = append(infos, info)
	}
	return infos
}

func GetGoRoutineID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
