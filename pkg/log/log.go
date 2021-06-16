package log

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

type Level int8

const (
	ERROR Level = iota
	INFO
	DEBUG
)

func (l Level) ToString() string {
	switch l {
	case ERROR:
		return "ERROR"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	default:
		return ""
	}
}

var debugLog *logrus.Logger
var infoLog *logrus.Logger
var errorLog *logrus.Logger

//Init ...
func Init(options ...LogOption) error {
	opts := LogOptions{}
	for _, option := range options {
		option.F(&opts)
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

	debugOutput, err := opts.NewOutput(DEBUG)
	if err != nil {
		return err
	}
	if debugOutput != nil {
		debugLog = logrus.New()
		debugLog.SetLevel(logrus.DebugLevel)
		debugLog.SetFormatter(formatter)
		debugLog.SetOutput(debugOutput)	
	}

	infoOutput, err := opts.NewOutput(INFO)
	if err != nil {
		return err
	}
	if infoOutput != nil {
		infoLog = logrus.New()
		infoLog.SetLevel(logrus.InfoLevel)
		infoLog.SetFormatter(formatter)
		infoLog.SetOutput(infoOutput)
	}

	errorOutput, err := opts.NewOutput(ERROR)
	if err != nil {
		return err
	}
	if errorOutput != nil {
		errorLog = logrus.New()
		errorLog.SetLevel(logrus.InfoLevel)
		errorLog.SetFormatter(formatter)
		errorLog.SetOutput(errorOutput)
	}

	return nil
}

//Error ...
func Error(args ...interface{}) {
	if errorLog != nil {
		errorLog.Error(args)
	}
}

//Info ...
func Info(args ...interface{}) {
	if infoLog != nil {
		infoLog.Info(args)
	}
}

//Debug ...
func Debug(args ...interface{}) {
	if debugLog != nil {
		debugLog.Debug(args)
	}
}
