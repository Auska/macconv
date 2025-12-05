/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

// Package errors provides unified error handling for the macconv application.
// It defines custom error types and utilities for error wrapping and classification.
package errors

import (
	"fmt"
)

// ErrorType 定义错误类型
type ErrorType int

const (
	// ValidationError 验证错误
	ValidationError ErrorType = iota
	// NetworkError 网络错误
	NetworkError
	// FileSystemError 文件系统错误
	FileSystemError
	// ParseError 解析错误
	ParseError
)

// AppError 应用程序错误
type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 支持errors.Unwrap
func (e *AppError) Unwrap() error {
	return e.Err
}

// New 创建新的应用程序错误
func New(errorType ErrorType, message string) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
	}
}

// Wrap 包装已有错误
func Wrap(errorType ErrorType, message string, err error) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
		Err:     err,
	}
}

// IsValidationError 检查是否为验证错误
func IsValidationError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == ValidationError
	}
	return false
}

// IsNetworkError 检查是否为网络错误
func IsNetworkError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == NetworkError
	}
	return false
}

// IsFileSystemError 检查是否为文件系统错误
func IsFileSystemError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == FileSystemError
	}
	return false
}

// IsParseError 检查是否为解析错误
func IsParseError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == ParseError
	}
	return false
}
