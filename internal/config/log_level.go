package config

import "fmt"

// LogLevel Indicates the log level setting that can be specified.
type LogLevel string

const (
	// DebugLevel Indicates the debug log level.
	DebugLevel LogLevel = "debug"
	// InfoLevel Indicates the info log level.
	InfoLevel LogLevel = "info"
	// WarnLevel Indicates the warn log level.
	WarnLevel LogLevel = "warn"
	// ErrorLevel Indicates the error log level.
	ErrorLevel LogLevel = "error"
	// CriticalLevel Indicates the critical log level.
	CriticalLevel LogLevel = "critical"
)

// parseLogLevel Parses the log level setting string.
func parseLogLevel(v string) (any, error) {
	if v == "" {
		return InfoLevel, nil
	}

	switch lv := LogLevel(v); lv {
	case DebugLevel, InfoLevel, WarnLevel, ErrorLevel, CriticalLevel:
		return lv, nil
	default:
		return nil, fmt.Errorf("specify one of the following items: %s, %s, %s, %s, %s",
			DebugLevel, InfoLevel, WarnLevel, ErrorLevel, CriticalLevel)
	}
}
