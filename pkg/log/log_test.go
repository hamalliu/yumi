package log

import (
	"testing"
)

func TestLogrus(t *testing.T) {
	InitLogrus(INFO, SetStorageDir("logfile"), SetIsOutputStd(true), SetFileName("yumi"))
	Error("Logrus boom ....")
}

func TestStdlog(t *testing.T) {
	InitStdlog(INFO, SetStorageDir("logfile"), SetIsOutputStd(true), SetFileName("yumi"))
	Error("Stdlog boom ....", "xxxxxxxxx")
}
