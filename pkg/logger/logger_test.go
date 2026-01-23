//go:build unit

/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package logger

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger(WARN, nil)
	if logger == nil {
		t.Error("NewLogger() returned nil")
	}
	if logger.level != WARN {
		t.Errorf("Expected level WARN, got %v", logger.level)
	}
}

func TestSetLevel(t *testing.T) {
	logger := NewLogger(INFO, nil)

	logger.SetLevel(DEBUG)
	if logger.level != DEBUG {
		t.Errorf("Expected level DEBUG, got %v", logger.level)
	}

	logger.SetLevel(ERROR)
	if logger.level != ERROR {
		t.Errorf("Expected level ERROR, got %v", logger.level)
	}
}

func TestLogLevelOutput(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
		logFunc  func(*Logger, string, ...interface{})
		expected bool
	}{
		{
			name:     "DEBUG level - DEBUG message",
			level:    DEBUG,
			logFunc:  (*Logger).Debugf,
			expected: true,
		},
		{
			name:     "INFO level - DEBUG message",
			level:    INFO,
			logFunc:  (*Logger).Debugf,
			expected: false,
		},
		{
			name:     "INFO level - INFO message",
			level:    INFO,
			logFunc:  (*Logger).Infof,
			expected: true,
		},
		{
			name:     "WARN level - INFO message",
			level:    WARN,
			logFunc:  (*Logger).Infof,
			expected: false,
		},
		{
			name:     "WARN level - WARN message",
			level:    WARN,
			logFunc:  (*Logger).Warnf,
			expected: true,
		},
		{
			name:     "ERROR level - WARN message",
			level:    ERROR,
			logFunc:  (*Logger).Warnf,
			expected: false,
		},
		{
			name:     "ERROR level - ERROR message",
			level:    ERROR,
			logFunc:  (*Logger).Errorf,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger(tt.level, &buf)

			tt.logFunc(logger, "test message")

			output := buf.String()
			hasOutput := strings.TrimSpace(output) != ""

			if hasOutput != tt.expected {
				t.Errorf("Expected output: %v, got output: %v", tt.expected, hasOutput)
				if hasOutput {
					t.Logf("Output: %s", output)
				}
			}
		})
	}
}

func TestDefaultLogger(t *testing.T) {
	if DefaultLogger == nil {
		t.Error("DefaultLogger is nil")
	}

	// Test default logger functions
	var buf bytes.Buffer
	originalOutput := DefaultLogger.logger.Writer()
	DefaultLogger.logger.SetOutput(&buf)

	Debug("test debug")
	Info("test info")
	Warn("test warn")
	Error("test error")

	output := buf.String()

	// Default logger is set to INFO level, so debug should not appear
	if strings.Contains(output, "[DEBUG]") {
		t.Error("DEBUG message should not appear at INFO level")
	}

	if !strings.Contains(output, "[INFO]") {
		t.Error("INFO message should appear")
	}

	if !strings.Contains(output, "[WARN]") {
		t.Error("WARN message should appear")
	}

	if !strings.Contains(output, "[ERROR]") {
		t.Error("ERROR message should appear")
	}

	// Restore original output
	DefaultLogger.logger.SetOutput(originalOutput)
}

func TestPrintError(t *testing.T) {
	var buf bytes.Buffer
	originalOutput := DefaultLogger.logger.Writer()
	DefaultLogger.logger.SetOutput(&buf)

	PrintError(nil)
	if buf.String() != "" {
		t.Error("PrintError with nil should not output anything")
	}

	buf.Reset()
	PrintError(errors.New("test error"))
	output := buf.String()
	if !strings.Contains(output, "error") {
		t.Error("PrintError should output error message")
	}

	DefaultLogger.logger.SetOutput(originalOutput)
}

func TestPrintErrorWithMessage(t *testing.T) {
	var buf bytes.Buffer
	originalOutput := DefaultLogger.logger.Writer()
	DefaultLogger.logger.SetOutput(&buf)

	PrintErrorWithMessage("test message", nil)
	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Error("PrintErrorWithMessage should output message")
	}

	buf.Reset()
	PrintErrorWithMessage("test message", errors.New("test error"))
	output = buf.String()
	if !strings.Contains(output, "test message") {
		t.Error("PrintErrorWithMessage should output message")
	}

	DefaultLogger.logger.SetOutput(originalOutput)
}

func TestPrintValidationError(t *testing.T) {
	var buf bytes.Buffer
	originalOutput := DefaultLogger.logger.Writer()
	DefaultLogger.logger.SetOutput(&buf)

	PrintValidationError("test validation")
	output := buf.String()
	if !strings.Contains(output, "Validation error") {
		t.Error("PrintValidationError should output validation error")
	}
	if !strings.Contains(output, "test validation") {
		t.Error("PrintValidationError should output validation message")
	}

	DefaultLogger.logger.SetOutput(originalOutput)
}

func TestLogLevelConstants(t *testing.T) {
	if DEBUG != 0 {
		t.Errorf("DEBUG constant should be 0, got %d", DEBUG)
	}
	if INFO != 1 {
		t.Errorf("INFO constant should be 1, got %d", INFO)
	}
	if WARN != 2 {
		t.Errorf("WARN constant should be 2, got %d", WARN)
	}
	if ERROR != 3 {
		t.Errorf("ERROR constant should be 3, got %d", ERROR)
	}
}

func TestFatalf(t *testing.T) {
	// Fatalf calls os.Exit(1), so we can't test it directly
	// We can only verify that the function exists and has the right signature
	logger := NewLogger(ERROR, nil)

	// This should panic with exit status 1, but we can't test that without subprocess
	// So we just verify the function is callable
	defer func() {
		if r := recover(); r != nil {
			// Expected to panic due to os.Exit
		}
	}()

	// We won't actually call it to avoid exiting the test
	_ = logger.Fatalf
}

func TestLevelComparison(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
		message  LogLevel
		expected bool
	}{
		{"DEBUG <= DEBUG", DEBUG, DEBUG, true},
		{"DEBUG <= INFO", DEBUG, INFO, true},
		{"INFO <= DEBUG", INFO, DEBUG, false},
		{"INFO <= INFO", INFO, INFO, true},
		{"INFO <= WARN", INFO, WARN, true},
		{"WARN <= INFO", WARN, INFO, false},
		{"WARN <= WARN", WARN, WARN, true},
		{"WARN <= ERROR", WARN, ERROR, true},
		{"ERROR <= WARN", ERROR, WARN, false},
		{"ERROR <= ERROR", ERROR, ERROR, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.level <= tt.message
			if result != tt.expected {
				t.Errorf("Level comparison failed: %d <= %d = %v, expected %v",
					tt.level, tt.message, result, tt.expected)
			}
		})
	}
}
