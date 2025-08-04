package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Level represents the logging level
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// String returns the string representation of the level
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger defines the logging interface
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, err error, fields ...interface{})
	With(fields ...interface{}) Logger
	SetLevel(level Level)
}

// logger implements the Logger interface
type logger struct {
	level  Level
	logger *log.Logger
	fields []interface{}
}

// New creates a new logger instance
func New() Logger {
	return &logger{
		level:  LevelInfo,
		logger: log.New(os.Stdout, "", 0),
		fields: make([]interface{}, 0),
	}
}

// NewWithLevel creates a new logger with a specific level
func NewWithLevel(level Level) Logger {
	return &logger{
		level:  level,
		logger: log.New(os.Stdout, "", 0),
		fields: make([]interface{}, 0),
	}
}

// Debug logs a debug message
func (l *logger) Debug(msg string, fields ...interface{}) {
	if l.level <= LevelDebug {
		l.log(LevelDebug, msg, nil, fields...)
	}
}

// Info logs an info message
func (l *logger) Info(msg string, fields ...interface{}) {
	if l.level <= LevelInfo {
		l.log(LevelInfo, msg, nil, fields...)
	}
}

// Warn logs a warning message
func (l *logger) Warn(msg string, fields ...interface{}) {
	if l.level <= LevelWarn {
		l.log(LevelWarn, msg, nil, fields...)
	}
}

// Error logs an error message
func (l *logger) Error(msg string, err error, fields ...interface{}) {
	if l.level <= LevelError {
		l.log(LevelError, msg, err, fields...)
	}
}

// With creates a new logger with additional fields
func (l *logger) With(fields ...interface{}) Logger {
	newFields := make([]interface{}, len(l.fields)+len(fields))
	copy(newFields, l.fields)
	copy(newFields[len(l.fields):], fields)

	return &logger{
		level:  l.level,
		logger: l.logger,
		fields: newFields,
	}
}

// SetLevel sets the logging level
func (l *logger) SetLevel(level Level) {
	l.level = level
}

// log performs the actual logging
func (l *logger) log(level Level, msg string, err error, fields ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Build the log message
	logMsg := fmt.Sprintf("[%s] %s: %s", timestamp, level.String(), msg)

	// Add error if present
	if err != nil {
		logMsg += fmt.Sprintf(" error=%v", err)
	}

	// Add persistent fields
	if len(l.fields) > 0 {
		logMsg += " " + formatFields(l.fields)
	}

	// Add additional fields
	if len(fields) > 0 {
		logMsg += " " + formatFields(fields)
	}

	l.logger.Println(logMsg)
}

// formatFields formats key-value pairs for logging
func formatFields(fields []interface{}) string {
	if len(fields)%2 != 0 {
		// If odd number of fields, add the last one as a value with empty key
		fields = append(fields, "")
	}

	var parts []string
	for i := 0; i < len(fields); i += 2 {
		key := fmt.Sprintf("%v", fields[i])
		value := fmt.Sprintf("%v", fields[i+1])
		parts = append(parts, fmt.Sprintf("%s=%s", key, value))
	}

	return fmt.Sprintf("[%s]", fmt.Sprintf("%s", parts))
}

// NoOpLogger is a logger that does nothing (useful for testing)
type NoOpLogger struct{}

// NewNoOp creates a new no-op logger
func NewNoOp() Logger {
	return &NoOpLogger{}
}

func (n *NoOpLogger) Debug(msg string, fields ...interface{})            {}
func (n *NoOpLogger) Info(msg string, fields ...interface{})             {}
func (n *NoOpLogger) Warn(msg string, fields ...interface{})             {}
func (n *NoOpLogger) Error(msg string, err error, fields ...interface{}) {}
func (n *NoOpLogger) With(fields ...interface{}) Logger                  { return n }
func (n *NoOpLogger) SetLevel(level Level)                               {}
