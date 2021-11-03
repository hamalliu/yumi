package log

// Logger 抽象接口
type Logger interface {
	Error(args ...interface{})
	Warn(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
}

var defaultDebugLog Logger
var defaultWarnLog Logger
var defaultInfoLog Logger
var defaultErrorLog Logger

//Error ...
func Error(args ...interface{}) {
	if defaultErrorLog != nil {
		defaultErrorLog.Error(args...)
	}
}

// Warn ...
func Warn(args ...interface{}) {
	if defaultWarnLog != nil {
		defaultWarnLog.Warn(args...)
	}
}

//Info ...
func Info(args ...interface{}) {
	if defaultInfoLog != nil {
		defaultInfoLog.Info(args...)
	}
}

//Debug ...
func Debug(args ...interface{}) {
	if defaultDebugLog != nil {
		defaultDebugLog.Debug(args...)
	}
}
