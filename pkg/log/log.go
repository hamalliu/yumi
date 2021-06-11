package log

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

var infoLog = logrus.New()
var errorLog = logrus.New()

//Init ...
func Init(options ...LogOption) error {
	logs := LogOptions{}
	for _, option := range options {
		option.F(&logs)
	}
	formatter := &logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		PadLevelText:    true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			_, fileName := filepath.Split(frame.File)
			file = fmt.Sprintf("%s:%d", fileName, frame.Line)
			return
		},
		EnvironmentOverrideColors: true,
	}

	infoLog.SetLevel(logrus.InfoLevel)
	infoLog.SetReportCaller(true)
	infoLog.SetFormatter(formatter)
	infoWriter, err := logs.NewFileWriter("INFO")
	if err != nil {
		return err
	}
	infoLog.SetOutput(infoWriter)

	errorLog.SetLevel(logrus.InfoLevel)
	errorLog.SetReportCaller(true)
	errorLog.SetFormatter(formatter)
	errorWriter, err := logs.NewFileWriter("INFO")
	if err != nil {
		return err
	}
	errorLog.SetOutput(errorWriter)

	return nil
}

//Error ...
func Error(args ...interface{}) {
	errorLog.SetReportCaller() = 1
	errorLog.Error(args)
}

//Info ...
func Info(args ...interface{}) {
	infolog.ExtraCalldepth = 1
	infolog.Info(args)
}

//Debug ...
func Debug(args ...interface{}) {
	log.ExtraCalldepth = 1
	log.Debug(args)
}

