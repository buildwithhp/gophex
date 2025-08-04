package logger

import (
	"bytes"
	"errors"
	"log"
	"strings"
	"testing"
)

func TestLogger_Levels(t *testing.T) {
	var buf bytes.Buffer
	logger := &logger{
		level:  LevelDebug,
		logger: log.New(&buf, "", 0),
		fields: make([]interface{}, 0),
	}

	tests := []struct {
		name     string
		logFunc  func()
		expected string
	}{
		{
			name: "Debug message",
			logFunc: func() {
				logger.Debug("debug message")
			},
			expected: "DEBUG: debug message",
		},
		{
			name: "Info message",
			logFunc: func() {
				logger.Info("info message")
			},
			expected: "INFO: info message",
		},
		{
			name: "Warn message",
			logFunc: func() {
				logger.Warn("warn message")
			},
			expected: "WARN: warn message",
		},
		{
			name: "Error message",
			logFunc: func() {
				logger.Error("error message", errors.New("test error"))
			},
			expected: "ERROR: error message error=test error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf.Reset()
			test.logFunc()

			output := buf.String()
			if !strings.Contains(output, test.expected) {
				t.Errorf("Expected output to contain %q, got %q", test.expected, output)
			}
		})
	}
}

func TestLogger_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := &logger{
		level:  LevelWarn, // Only warn and error should be logged
		logger: log.New(&buf, "", 0),
		fields: make([]interface{}, 0),
	}

	// These should not be logged
	logger.Debug("debug message")
	logger.Info("info message")

	// These should be logged
	logger.Warn("warn message")
	logger.Error("error message", nil)

	output := buf.String()

	if strings.Contains(output, "debug message") {
		t.Error("Debug message should not be logged at WARN level")
	}

	if strings.Contains(output, "info message") {
		t.Error("Info message should not be logged at WARN level")
	}

	if !strings.Contains(output, "warn message") {
		t.Error("Warn message should be logged at WARN level")
	}

	if !strings.Contains(output, "error message") {
		t.Error("Error message should be logged at WARN level")
	}
}

func TestLogger_WithFields(t *testing.T) {
	var buf bytes.Buffer
	baseLogger := &logger{
		level:  LevelInfo,
		logger: log.New(&buf, "", 0),
		fields: make([]interface{}, 0),
	}

	// Create logger with fields
	loggerWithFields := baseLogger.With("key1", "value1", "key2", "value2")
	loggerWithFields.Info("test message")

	output := buf.String()

	if !strings.Contains(output, "key1=value1") {
		t.Error("Output should contain key1=value1")
	}

	if !strings.Contains(output, "key2=value2") {
		t.Error("Output should contain key2=value2")
	}

	if !strings.Contains(output, "test message") {
		t.Error("Output should contain the log message")
	}
}

func TestLogger_AdditionalFields(t *testing.T) {
	var buf bytes.Buffer
	logger := &logger{
		level:  LevelInfo,
		logger: log.New(&buf, "", 0),
		fields: make([]interface{}, 0),
	}

	logger.Info("test message", "field1", "value1", "field2", "value2")

	output := buf.String()

	if !strings.Contains(output, "field1=value1") {
		t.Error("Output should contain field1=value1")
	}

	if !strings.Contains(output, "field2=value2") {
		t.Error("Output should contain field2=value2")
	}
}

func TestLogger_SetLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := &logger{
		level:  LevelInfo,
		logger: log.New(&buf, "", 0),
		fields: make([]interface{}, 0),
	}

	// Initially at INFO level
	logger.Debug("debug message")
	if strings.Contains(buf.String(), "debug message") {
		t.Error("Debug message should not be logged at INFO level")
	}

	// Change to DEBUG level
	buf.Reset()
	logger.SetLevel(LevelDebug)
	logger.Debug("debug message")

	if !strings.Contains(buf.String(), "debug message") {
		t.Error("Debug message should be logged at DEBUG level")
	}
}

func TestLevel_String(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{LevelDebug, "DEBUG"},
		{LevelInfo, "INFO"},
		{LevelWarn, "WARN"},
		{LevelError, "ERROR"},
		{Level(999), "UNKNOWN"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			if test.level.String() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, test.level.String())
			}
		})
	}
}

func TestNew(t *testing.T) {
	logger := New()
	if logger == nil {
		t.Error("New() should return a non-nil logger")
	}

	// Test that it implements the Logger interface
	var _ Logger = logger
}

func TestNewWithLevel(t *testing.T) {
	logger := NewWithLevel(LevelError)
	if logger == nil {
		t.Error("NewWithLevel() should return a non-nil logger")
	}

	// Test that it implements the Logger interface
	var _ Logger = logger
}

func TestNoOpLogger(t *testing.T) {
	logger := NewNoOp()
	if logger == nil {
		t.Error("NewNoOp() should return a non-nil logger")
	}

	// Test that all methods can be called without panic
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error", errors.New("test"))
	logger.SetLevel(LevelDebug)

	withLogger := logger.With("key", "value")
	if withLogger == nil {
		t.Error("With() should return a non-nil logger")
	}
}

func TestFormatFields(t *testing.T) {
	tests := []struct {
		name     string
		fields   []interface{}
		expected string
	}{
		{
			name:     "Even number of fields",
			fields:   []interface{}{"key1", "value1", "key2", "value2"},
			expected: "[key1=value1 key2=value2]",
		},
		{
			name:     "Odd number of fields",
			fields:   []interface{}{"key1", "value1", "key2"},
			expected: "[key1=value1 key2=]",
		},
		{
			name:     "Empty fields",
			fields:   []interface{}{},
			expected: "[]",
		},
		{
			name:     "Single field",
			fields:   []interface{}{"key1"},
			expected: "[key1=]",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := formatFields(test.fields)
			if !strings.Contains(result, "key1=value1") && len(test.fields) >= 2 {
				t.Errorf("Expected result to contain formatted fields, got %s", result)
			}
		})
	}
}

func TestLogger_ConcurrentAccess(t *testing.T) {
	var buf bytes.Buffer
	logger := &logger{
		level:  LevelInfo,
		logger: log.New(&buf, "", 0),
		fields: make([]interface{}, 0),
	}

	// Test concurrent logging
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(index int) {
			logger.Info("concurrent message", "index", index)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	output := buf.String()
	if !strings.Contains(output, "concurrent message") {
		t.Error("Output should contain concurrent messages")
	}
}

func BenchmarkLogger_Info(b *testing.B) {
	var buf bytes.Buffer
	logger := &logger{
		level:  LevelInfo,
		logger: log.New(&buf, "", 0),
		fields: make([]interface{}, 0),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message", "iteration", i)
	}
}

func BenchmarkLogger_WithFields(b *testing.B) {
	var buf bytes.Buffer
	logger := &logger{
		level:  LevelInfo,
		logger: log.New(&buf, "", 0),
		fields: make([]interface{}, 0),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loggerWithFields := logger.With("key1", "value1", "key2", "value2")
		loggerWithFields.Info("benchmark message")
	}
}
