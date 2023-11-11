package logo

import (
	"github.com/wsrf16/swiss/sugar/runtimekit"
)

func Panic(summary string, data any, message ...any) *Logger {
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
	}).Panic(message...)
	return log
}

func Panicf(summary string, data any, format string, a ...any) *Logger {
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
	}).Panicf(format, a...)
	return log
}

func P(message ...any) *Logger {
	info := runtimekit.GetStackTraceInfos()[1]
	stacks := runtimekit.GetStackTraceLines()[1:]
	log.WithFields(Fields{
		"stackoverflow": stacks,
		"file":          info.File,
		"line":          info.Line,
		"func":          info.Function,
		"pc":            info.ProgramCounter,
	}).Panic(message...)
	return log
}

func Pf(format string, a ...any) *Logger {
	info := runtimekit.GetStackTraceInfos()[1]
	stacks := runtimekit.GetStackTraceLines()[1:]
	log.WithFields(Fields{
		"stackoverflow": stacks,
		"file":          info.File,
		"line":          info.Line,
		"func":          info.Function,
		"pc":            info.ProgramCounter,
	}).Panicf(format, a...)
	return log
}

func Fatal(summary string, data any, message ...any) *Logger {
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
	}).Fatal(message...)
	return log
}

func Fatalf(summary string, data any, format string, a ...any) *Logger {
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
	}).Fatalf(format, a...)
	return log
}

func F(message ...any) *Logger {
	info := runtimekit.GetStackTraceInfos()[1]
	stacks := runtimekit.GetStackTraceLines()[1:]
	log.WithFields(Fields{
		"stackoverflow": stacks,
		"file":          info.File,
		"line":          info.Line,
		"func":          info.Function,
		"pc":            info.ProgramCounter,
	}).Fatal(message...)
	return log
}

func Ff(format string, a ...any) *Logger {
	info := runtimekit.GetStackTraceInfos()[1]
	stacks := runtimekit.GetStackTraceLines()[1:]
	log.WithFields(Fields{
		"stackoverflow": stacks,
		"file":          info.File,
		"line":          info.Line,
		"func":          info.Function,
		"pc":            info.ProgramCounter,
	}).Fatalf(format, a...)
	return log
}

func Error(summary string, data any, message ...any) *Logger {
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
	}).Error(message...)
	return log
}

func Errorf(summary string, data any, format string, a ...any) *Logger {
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
	}).Errorf(format, a...)
	return log
}

func E(message ...any) *Logger {
	info := runtimekit.GetStackTraceInfos()[1]
	stacks := runtimekit.GetStackTraceLines()[1:]
	log.WithFields(Fields{
		"stackoverflow": stacks,
		"file":          info.File,
		"line":          info.Line,
		"func":          info.Function,
		"pc":            info.ProgramCounter,
	}).Error(message...)
	return log
}

func Ef(format string, a ...any) *Logger {
	info := runtimekit.GetStackTraceInfos()[1]
	stacks := runtimekit.GetStackTraceLines()[1:]
	log.WithFields(Fields{
		"stackoverflow": stacks,
		"file":          info.File,
		"line":          info.Line,
		"func":          info.Function,
		"pc":            info.ProgramCounter,
	}).Errorf(format, a...)
	return log
}

func Warn(summary string, data any, message ...any) *Logger {
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
	}).Warn(message...)
	return log
}

func Warnf(summary string, data any, format string, a ...any) *Logger {
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
	}).Warnf(format, a...)
	return log
}

func W(message ...any) *Logger {
	info := runtimekit.GetStackTraceInfos()[1]
	stacks := runtimekit.GetStackTraceLines()[1:]
	log.WithFields(Fields{
		"stackoverflow": stacks,
		"file":          info.File,
		"line":          info.Line,
		"func":          info.Function,
		"pc":            info.ProgramCounter,
	}).Warn(message...)
	return log
}

func Wf(format string, a ...any) *Logger {
	info := runtimekit.GetStackTraceInfos()[1]
	stacks := runtimekit.GetStackTraceLines()[1:]
	log.WithFields(Fields{
		"stackoverflow": stacks,
		"file":          info.File,
		"line":          info.Line,
		"func":          info.Function,
		"pc":            info.ProgramCounter,
	}).Warnf(format, a...)
	return log
}

func Info(summary string, data any, message ...any) *Logger {
	log.WithFields(Fields{
		"summary": summary,
		"data":    data,
	}).Info(message...)
	return log
}

func Infof(summary string, data any, format string, a ...any) *Logger {
	log.WithFields(Fields{
		"summary": summary,
		"data":    data,
	}).Infof(format, a...)
	return log
}

func I(message ...any) *Logger {
	log.WithFields(Fields{}).Info(message...)
	return log
}

func If(format string, a ...any) *Logger {
	log.WithFields(Fields{}).Infof(format, a...)
	return log
}

func Debug(summary string, data any, message ...any) *Logger {
	log.WithFields(Fields{
		"summary": summary,
		"data":    data,
	}).Debug(message...)
	return log
}

func Debugf(summary string, data any, format string, a ...any) *Logger {
	log.WithFields(Fields{
		"summary": summary,
		"data":    data,
	}).Debugf(format, a...)
	return log
}

func D(message ...any) *Logger {
	log.WithFields(Fields{}).Debug(message...)
	return log
}

func Df(format string, a ...any) *Logger {
	log.WithFields(Fields{}).Debugf(format, a...)
	return log
}

func Trace(summary string, data any, message ...any) *Logger {
	log.WithFields(Fields{
		"summary": summary,
		"data":    data,
	}).Trace(message...)
	return log
}

func Tracef(summary string, data any, format string, a ...any) *Logger {
	log.WithFields(Fields{
		"summary": summary,
		"data":    data,
	}).Tracef(format, a...)
	return log
}

func T(summary string, data any, message ...any) *Logger {
	log.WithFields(Fields{}).Trace(message...)
	return log
}

func Tf(summary string, data any, format string, a ...any) *Logger {
	log.WithFields(Fields{}).Tracef(format, a...)
	return log
}
