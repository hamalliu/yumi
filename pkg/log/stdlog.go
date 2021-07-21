package log

import (
	"fmt"
	"io"
	"log"
	"os"
)

func InitStdlog(level Level, options ...LogOption) {
	opts := LogOptions{}
	for _, option := range options {
		option.F(&opts)
	}

	formatter := log.LstdFlags | log.Lshortfile

	var (
		debugFileOutput, infoFileOutput, errorFileOutput io.Writer
		err                                  error
	)
	if level >= DEBUG {
		debugFileOutput, err = opts.NewFileOutput(nil, DEBUG)
		if err != nil {
			panic(err)
		}
		output := debugFileOutput
		if opts.IsOutputStd {
			output = io.MultiWriter(output, os.Stdout)
		}
		debugStdLog := log.New(output, "", formatter)
		defaultDebugLog = &StdLogger{l: debugStdLog}
	}

	if level >= INFO {
		infoFileOutput, err = opts.NewFileOutput(debugFileOutput, INFO)
		if err != nil {
			panic(err)
		}
		output := infoFileOutput
		if opts.IsOutputStd {
			output = io.MultiWriter(output, os.Stdout)
		}
		infoStdLog := log.New(output, "", formatter)
		defaultInfoLog = &StdLogger{l: infoStdLog}
	}

	if level >= ERROR {
		errorFileOutput, err = opts.NewFileOutput(infoFileOutput, ERROR)
		if err != nil {
			panic(err)
		}
		output := errorFileOutput
		if opts.IsOutputStd {
			output = io.MultiWriter(output, os.Stdout)
		}
		errorStdLog := log.New(output, "", formatter)
		defaultErrorLog = &StdLogger{l: errorStdLog}
	}
}

type StdLogger struct {
	l *log.Logger
}

func (s *StdLogger) Error(args ...interface{}) {
	s.l.Output(3, fmt.Sprintln("ERROR ", args))
}

func (s *StdLogger) Info(args ...interface{}) {
	s.l.Output(3, fmt.Sprintln("INFO  ", args))
}

func (s *StdLogger) Debug(args ...interface{}) {
	s.l.Output(3, fmt.Sprintln("DEBUG ", args))
}
