package log

import (
	"fmt"
	"io"
	"log"
	"os"
)

// InitStdlog ...
func InitStdlog(level Level, options ...Option) {
	opts := Options{}
	for _, option := range options {
		option.F(&opts)
	}

	formatter := log.LstdFlags | log.Lshortfile

	if level >= DEBUG {
		debugFileOutput, err := opts.NewFileOutput(DEBUG)
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
		infoFileOutput, err := opts.NewFileOutput(INFO)
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

	if level >= WARN {
		warnFileOutput, err := opts.NewFileOutput(WARN)
		if err != nil {
			panic(err)
		}
		output := warnFileOutput
		if opts.IsOutputStd {
			output = io.MultiWriter(output, os.Stdout)
		}
		warnStdLog := log.New(output, "", formatter)
		defaultInfoLog = &StdLogger{l: warnStdLog}
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
		errorStdLog := log.New(output, "", formatter)
		defaultErrorLog = &StdLogger{l: errorStdLog}
	}
}

// StdLogger ...
type StdLogger struct {
	l *log.Logger
}

// Error ...
func (s *StdLogger) Error(args ...interface{}) {
	_ = s.l.Output(3, fmt.Sprintln("ERROR ", args))
}

// Warn ...
func (s *StdLogger) Warn(args ...interface{}) {
	_ = s.l.Output(3, fmt.Sprintln("Warn ", args))
}

// Info ...
func (s *StdLogger) Info(args ...interface{}) {
	_ = s.l.Output(3, fmt.Sprintln("INFO  ", args))
}

// Debug ...
func (s *StdLogger) Debug(args ...interface{}) {
	_ = s.l.Output(3, fmt.Sprintln("DEBUG ", args))
}
