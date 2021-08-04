package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

//InitLogrus ...
func InitLogrus(level Level, options ...LogOption) {
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

	if level >= DEBUG {
		debugFileOutput, err := opts.NewFileOutput(DEBUG)
		if err != nil {
			panic(err)
		}
		output := debugFileOutput
		if opts.IsOutputStd {
			output = io.MultiWriter(output, os.Stdout)
		}
		if output != nil {
			debugLog := logrus.New()
			debugLog.SetLevel(logrus.DebugLevel)
			debugLog.SetFormatter(formatter)
			debugLog.SetOutput(output)
			defaultDebugLog = debugLog
		}
	}

	if level >= INFO {
		infoFileOutput, err := opts.NewFileOutput(INFO)
		if err != nil {
			panic(err)
		}
		output := infoFileOutput
		if opts.IsOutputStd {
			output = io.MultiWriter(output, os.Stdout)
		}
		if output != nil {
			infoLog := logrus.New()
			infoLog.SetLevel(logrus.InfoLevel)
			infoLog.SetFormatter(formatter)
			infoLog.SetOutput(output)
			defaultInfoLog = infoLog
		}
	}

	if level >= WARN {
		warnFileOutput, err := opts.NewFileOutput(WARN)
		if err != nil {
			panic(err)
		}
		output := warnFileOutput
		if opts.IsOutputStd {
			output = io.MultiWriter(output, os.Stdout)
		}
		if output != nil {
			warnLog := logrus.New()
			warnLog.SetLevel(logrus.WarnLevel)
			warnLog.SetFormatter(formatter)
			warnLog.SetOutput(output)
			defaultWarnLog = warnLog
		}
	}

	if level >= ERROR {
		errorFileOutput, err := opts.NewFileOutput(ERROR)
		if err != nil {
			panic(err)
		}
		output := errorFileOutput
		if opts.IsOutputStd {
			output = io.MultiWriter(output, os.Stdout)
		}
		if output != nil {
			errorLog := logrus.New()
			errorLog.SetLevel(logrus.InfoLevel)
			errorLog.SetFormatter(formatter)
			errorLog.SetOutput(output)
			defaultErrorLog = errorLog
		}
	}
}
