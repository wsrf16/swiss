package logo

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
)

type Level = logrus.Level

type Hook = logrus.Hook

var AllLevels = []Level{
	PanicLevel,
	FatalLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
	TraceLevel,
}

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

type Formatter = logrus.Formatter

type JSONFormatter = logrus.JSONFormatter

type TextFormatter = logrus.TextFormatter

type Fields = logrus.Fields

type Entry = logrus.Entry

type Logger = logrus.Logger

//type JSONFormatter struct {
//    logrus.JSONFormatter
//}

//type TextFormatter struct {
//    logrus.TextFormatter
//}

var log *Logger = logrus.New()

var once sync.Once

func init() {
	once.Do(func() {
		//SetReportCaller(true)
		//SetFormatter(&JSONFormatter{})
		SetFormatter(&TextFormatter{
			FullTimestamp: true,
			ForceQuote:    true,
		})
		SetOutput(os.Stdout)
		SetLevel(InfoLevel)
		AddHook(&DefaultFieldsHook{})
	})
}

func SetReportCaller(reportCaller bool) {
	log.SetReportCaller(reportCaller)
}

func SetFormatter(formatter Formatter) {
	log.SetFormatter(formatter)
}

func SetLevel(level Level) {
	log.SetLevel(level)
}

func SetOutput(output io.Writer) {
	log.SetOutput(output)
}

func WithFields(fields Fields) *Entry {
	return log.WithFields(fields)
}

func AddHook(hook Hook) {
	log.AddHook(hook)
}
