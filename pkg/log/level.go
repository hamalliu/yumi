package log

type Level int8

const (
	ERROR Level = iota
	WARN
	INFO
	DEBUG
)

func (l Level) ToString() string {
	switch l {
	case ERROR:
		return "ERROR"
	case WARN:
		return "WARN"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	default:
		return ""
	}
}
