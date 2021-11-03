package log

// Level log level
type Level int8

const (
	// ERROR is error level
	ERROR Level = iota
	// WARN is error level
	WARN
	// INFO is error level
	INFO
	// DEBUG is error level
	DEBUG
)

// ToString ...
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
