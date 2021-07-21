package log

type Logger interface {
	Error(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
}

var defaultDebugLog Logger
var defaultInfoLog Logger
var defaultErrorLog Logger

//Error ...
func Error(args ...interface{}) {
	if defaultErrorLog != nil {
		defaultErrorLog.Error(args...)
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
