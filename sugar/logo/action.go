package logo

import (
	"fmt"
	"github.com/wsrf16/swiss/sugar/runtimekit"
)

func Panic(summary string, data any, message string) *Logger {
	info := runtimekit.GetStackTraceInfos()[1]
	stacks := runtimekit.GetStackTraceLines()[1:]
	log.WithFields(Fields{
		"summary":       summary,
		"data":          data,
		"stackoverflow": stacks,
		"file":          info.File,
		"line":          info.Line,
		"func":          info.Function,
		"pc":            info.ProgramCounter,
	}).Panic(message)
	return log
}

func PanicF(summary string, data any, message string, a ...any) *Logger {
	info := runtimekit.GetStackTraceInfos()[1]
	stacks := runtimekit.GetStackTraceLines()[1:]
	log.WithFields(Fields{
		"summary":       summary,
		"data":          data,
		"stackoverflow": stacks,
		"file":          info.File,
		"line":          info.Line,
		"func":          info.Function,
		"pc":            info.ProgramCounter,
	}).Panic(fmt.Sprintf(message, a))
	return log
}

func Fatal(summary string, data any, message string) *Logger {
	info := runtimekit.GetStackTraceInfos()[1]
	stacks := runtimekit.GetStackTraceLines()[1:]
	log.WithFields(Fields{
		"summary":       summary,
		"data":          data,
		"stackoverflow": stacks,
		"file":          info.File,
		"line":          info.Line,
		"func":          info.Function,
		"pc":            info.ProgramCounter,
	}).Fatal(message)
	return log
}

func FatalF(summary string, data any, message string, a ...any) *Logger {
	info := runtimekit.GetStackTraceInfos()[1]
	stacks := runtimekit.GetStackTraceLines()[1:]
	log.WithFields(Fields{
		"summary":       summary,
		"data":          data,
		"stackoverflow": stacks,
		"file":          info.File,
		"line":          info.Line,
		"func":          info.Function,
		"pc":            info.ProgramCounter,
	}).Fatal(fmt.Sprintf(message, a))
	return log
}

func Error(summary string, data any, message string) *Logger {
	info := runtimekit.GetStackTraceInfos()[1]
	stacks := runtimekit.GetStackTraceLines()[1:]
	log.WithFields(Fields{
		"summary":       summary,
		"data":          data,
		"stackoverflow": stacks,
		"file":          info.File,
		"line":          info.Line,
		"func":          info.Function,
		"pc":            info.ProgramCounter,
	}).Error(message)
	return log
}

func ErrorF(summary string, data any, message string, a ...any) *Logger {
	info := runtimekit.GetStackTraceInfos()[1]
	stacks := runtimekit.GetStackTraceLines()[1:]
	log.WithFields(Fields{
		"summary":       summary,
		"data":          data,
		"stackoverflow": stacks,
		"file":          info.File,
		"line":          info.Line,
		"func":          info.Function,
		"pc":            info.ProgramCounter,
	}).Errorf(message, a)
	return log
}

func Warn(summary string, data any, message string) *Logger {
	info := runtimekit.GetStackTraceInfos()[1]
	stacks := runtimekit.GetStackTraceLines()[1:]
	log.WithFields(Fields{
		"summary":       summary,
		"data":          data,
		"stackoverflow": stacks,
		"file":          info.File,
		"line":          info.Line,
		"func":          info.Function,
		"pc":            info.ProgramCounter,
	}).Warn(message)
	return log
}

func WarnF(summary string, data any, message string, a ...any) *Logger {
	info := runtimekit.GetStackTraceInfos()[1]
	stacks := runtimekit.GetStackTraceLines()[1:]
	log.WithFields(Fields{
		"summary":       summary,
		"data":          data,
		"stackoverflow": stacks,
		"file":          info.File,
		"line":          info.Line,
		"func":          info.Function,
		"pc":            info.ProgramCounter,
	}).Warn(fmt.Sprintf(message, a))
	return log
}

func Info(summary string, data any, message string) *Logger {
	log.WithFields(Fields{
		"summary": summary,
		"data":    data,
	}).Info(message)
	return log
}

func InfoF(summary string, data any, message string, a ...any) *Logger {
	log.WithFields(Fields{
		"summary": summary,
		"data":    data,
	}).Info(fmt.Sprintf(message, a))
	return log
}

func Debug(summary string, data any, message string) *Logger {
	log.WithFields(Fields{
		"summary": summary,
		"data":    data,
	}).Debug(message)
	return log
}

func DebugF(summary string, data any, message string, a ...any) *Logger {
	log.WithFields(Fields{
		"summary": summary,
		"data":    data,
	}).Debug(fmt.Sprintf(message, a))
	return log
}

func Trace(summary string, data any, message string) *Logger {
	log.WithFields(Fields{
		"summary": summary,
		"data":    data,
	}).Trace(message)
	return log
}

func TraceF(summary string, data any, message string, a ...any) *Logger {
	log.WithFields(Fields{
		"summary": summary,
		"data":    data,
	}).Trace(fmt.Sprintf(message, a))
	return log
}
