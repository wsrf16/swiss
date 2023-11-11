package logo

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	log := new()
	log.SetReportCaller(true)
	log.SetFormatter(&JSONFormatter{})
	log.SetFormatter(&TextFormatter{
		FullTimestamp: true,
		ForceQuote:    true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(InfoLevel)
	log.AddHook(&DefaultFieldsHook{})

	log.Debug("abcdefg")
	log.Debugf("ab%vefg", "cd")
}

func TestGlobal(t *testing.T) {
	SetReportCaller(true)
	SetFormatter(&JSONFormatter{})
	SetFormatter(&TextFormatter{
		FullTimestamp: true,
		ForceQuote:    true,
	})
	SetOutput(os.Stdout)
	SetLevel(InfoLevel)
	AddHook(&DefaultFieldsHook{})

	Debug("summary", "abcdefg")
	Debugf("summary", nil, "ab%vefg", "cd")
}
