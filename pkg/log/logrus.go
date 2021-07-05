package log

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)


//InitLogrus ...
func InitLogrus(options ...LogOption) {
	opts := LogOptions{}
	for _, option := range options {
		option.F(&opts)
	}
	formatter := &logrus.TextFormatter{
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

	debugOutput, err := opts.NewOutput(DEBUG)
	if err != nil {
		panic(err)
	}
	if debugOutput != nil {
		debugLog := logrus.New()
		debugLog.SetLevel(logrus.DebugLevel)
		debugLog.SetFormatter(formatter)
		debugLog.SetOutput(debugOutput)	
		defaultDebugLog = debugLog
	}

	infoOutput, err := opts.NewOutput(INFO)
	if err != nil {
		panic(err)
	}
	if infoOutput != nil {
		infoLog := logrus.New()
		infoLog.SetLevel(logrus.InfoLevel)
		infoLog.SetFormatter(formatter)
		infoLog.SetOutput(infoOutput)
		defaultInfoLog = infoLog
	}

	errorOutput, err := opts.NewOutput(ERROR)
	if err != nil {
		panic(err)
	}
	if errorOutput != nil {
		errorLog := logrus.New()
		errorLog.SetLevel(logrus.InfoLevel)
		errorLog.SetFormatter(formatter)
		errorLog.SetOutput(errorOutput)
		defaultErrorLog = errorLog
	}
}

