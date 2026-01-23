/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

// Package logger provides structured logging functionality for the macconv application.
// It supports multiple log levels and formatted output.
package logger

import (
	"io"
	"log"
	"os"
)

// LogLevel 日志级别
type LogLevel int

const (
	// DEBUG 调试级别
	DEBUG LogLevel = iota
	// INFO 信息级别
	INFO
	// WARN 警告级别
	WARN
	// ERROR 错误级别
	ERROR
)

// Logger 日志记录器
type Logger struct {
	level  LogLevel
	logger *log.Logger
}

var (
	// DefaultLogger 默认日志记录器
	DefaultLogger = NewLogger(INFO, os.Stderr)
)

// NewLogger 创建新的日志记录器
func NewLogger(level LogLevel, output io.Writer) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(output, "", log.LstdFlags),
	}
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// Debugf records debug information with formatting
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level <= DEBUG {
		l.logger.Printf("[DEBUG] "+format, v...)
	}
}

// Infof records information with formatting
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level <= INFO {
		l.logger.Printf("[INFO] "+format, v...)
	}
}

// Warnf records warning information with formatting
func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.level <= WARN {
		l.logger.Printf("[WARN] "+format, v...)
	}
}

// Errorf records error information with formatting
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level <= ERROR {
		l.logger.Printf("[ERROR] "+format, v...)
	}
}

// Fatalf records fatal error with formatting and exits the program
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logger.Printf("[FATAL] "+format, v...)
	os.Exit(1)
}

// Debugf records debug information using the default logger
func Debugf(format string, v ...interface{}) {
	DefaultLogger.Debugf(format, v...)
}

// Infof records information using the default logger
func Infof(format string, v ...interface{}) {
	DefaultLogger.Infof(format, v...)
}

// Warnf records warning information using the default logger
func Warnf(format string, v ...interface{}) {
	DefaultLogger.Warnf(format, v...)
}

// Errorf records error information using the default logger
func Errorf(format string, v ...interface{}) {
	DefaultLogger.Errorf(format, v...)
}

// Fatalf records fatal error using the default logger and exits
func Fatalf(format string, v ...interface{}) {
	DefaultLogger.Fatalf(format, v...)
}

// Debug records debug information using the default logger
func Debug(format string, v ...interface{}) {
	Debugf(format, v...)
}

// Info records information using the default logger
func Info(format string, v ...interface{}) {
	Infof(format, v...)
}

// Warn records warning information using the default logger
func Warn(format string, v ...interface{}) {
	Warnf(format, v...)
}

// Error records error information using the default logger
func Error(format string, v ...interface{}) {
	Errorf(format, v...)
}

// Fatal records fatal error using the default logger and exits
func Fatal(format string, v ...interface{}) {
	Fatalf(format, v...)
}

// PrintError 打印错误信息
func PrintError(err error) {
	if err != nil {
		Error("%v", err)
	}
}

// PrintErrorWithMessage 打印带消息的错误信息
func PrintErrorWithMessage(message string, err error) {
	Error("%s: %v", message, err)
}

// PrintValidationError 打印验证错误
func PrintValidationError(message string) {
	Error("Validation error: %s", message)
}
