package log

import (
	"testing"
)

func TestLogrus(t *testing.T) {
	InitLogrus(SetStorageDir("logfile"), SetIsOutputStd(true), SetFileName("yumi"))
	Error("Logrus boom ....")
}

func TestStdlog(t *testing.T) {
	InitStdlog(SetStorageDir("logfile"), SetIsOutputStd(true), SetFileName("yumi"))
	Error("Stdlog boom ....", "xxxxxxxxx")
}
