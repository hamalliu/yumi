package log

import (
	"fmt"
	"log"
)

func InitStdlog(options ...LogOption) {
	opts := LogOptions{}
	for _, option := range options {
		option.F(&opts)
	}

	formatter := log.LstdFlags | log.Lshortfile

	debugOutput, err := opts.NewOutput(DEBUG)
	if err != nil {
		panic(err)
	}
	debugStdLog := log.New(debugOutput, "", formatter)
	defaultDebugLog = &StdLogger{l: debugStdLog}

	infoOutput, err := opts.NewOutput(INFO)
	if err != nil {
		panic(err)
	}
	infoStdLog := log.New(infoOutput, "", formatter)
	defaultInfoLog = &StdLogger{l: infoStdLog}

	errorOutput, err := opts.NewOutput(ERROR)
	if err != nil {
		panic(err)
	}
	errorStdLog := log.New(errorOutput, "", formatter)
	defaultErrorLog = &StdLogger{l: errorStdLog}
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
