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

const (
	levelDebug = "DEBUG"
	levelInfo  = "INFO"
	levelWarn  = "WARN"
	levelError = "ERROR"
	levelFatal = "FATAL"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

type Logger struct {
	level  LogLevel
	logger *log.Logger
}

var (
	DefaultLogger = NewLogger(WARN, os.Stderr)
)

func NewLogger(level LogLevel, output io.Writer) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(output, "", log.LstdFlags),
	}
}

func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *Logger) logf(level LogLevel, levelStr, format string, v ...interface{}) {
	if l.level <= level {
		l.logger.Printf("[%s] "+format, append([]interface{}{levelStr}, v...)...)
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logf(DEBUG, levelDebug, format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.logf(INFO, levelInfo, format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logf(WARN, levelWarn, format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logf(ERROR, levelError, format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logger.Printf("[%s] "+format, append([]interface{}{levelFatal}, v...)...)
	os.Exit(1)
}

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

func PrintError(err error) {
	if err != nil {
		Errorf("%v", err)
	}
}

func PrintErrorWithMessage(message string, err error) {
	Errorf("%s: %v", message, err)
}

func PrintValidationError(message string) {
	Errorf("Validation error: %s", message)
}
