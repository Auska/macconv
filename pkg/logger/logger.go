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

// Debugf 记录调试信息
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level <= DEBUG {
		l.logger.Printf("[DEBUG] "+format, v...)
	}
}

// Infof 记录信息
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level <= INFO {
		l.logger.Printf("[INFO] "+format, v...)
	}
}

// Warnf 记录警告
func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.level <= WARN {
		l.logger.Printf("[WARN] "+format, v...)
	}
}

// Errorf 记录错误
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level <= ERROR {
		l.logger.Printf("[ERROR] "+format, v...)
	}
}

// Fatalf 记录致命错误并退出程序
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logger.Printf("[FATAL] "+format, v...)
	os.Exit(1)
}

// 全局便捷函数
func Debugf(format string, v ...interface{}) {
	DefaultLogger.Debugf(format, v...)
}

func Infof(format string, v ...interface{}) {
	DefaultLogger.Infof(format, v...)
}

func Warnf(format string, v ...interface{}) {
	DefaultLogger.Warnf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	DefaultLogger.Errorf(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	DefaultLogger.Fatalf(format, v...)
}

// 保持向后兼容的别名
func Debug(format string, v ...interface{}) {
	Debugf(format, v...)
}

func Info(format string, v ...interface{}) {
	Infof(format, v...)
}

func Warn(format string, v ...interface{}) {
	Warnf(format, v...)
}

func Error(format string, v ...interface{}) {
	Errorf(format, v...)
}

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

// PrintNetworkError 打印网络错误
func PrintNetworkError(message string, err error) {
	Error("Network error: %s: %v", message, err)
}

// PrintFileSystemError 打印文件系统错误
func PrintFileSystemError(message string, err error) {
	Error("File system error: %s: %v", message, err)
}

// PrintParseError 打印解析错误
func PrintParseError(message string, err error) {
	Error("Parse error: %s: %v", message, err)
}
