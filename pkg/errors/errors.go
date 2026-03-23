/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

// Package errors provides unified error handling for the macconv application.
// It defines custom error types and utilities for error wrapping and classification.
package errors

import (
	"fmt"
)

type ErrorType int

const (
	ValidationError ErrorType = iota
	NetworkError
	FileSystemError
	ParseError
)

type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(errorType ErrorType, message string) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
	}
}

func Wrap(errorType ErrorType, message string, err error) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
		Err:     err,
	}
}

func isErrorType(err error, errorType ErrorType) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == errorType
	}
	return false
}

func IsValidationError(err error) bool {
	return isErrorType(err, ValidationError)
}

func IsNetworkError(err error) bool {
	return isErrorType(err, NetworkError)
}

func IsFileSystemError(err error) bool {
	return isErrorType(err, FileSystemError)
}

func IsParseError(err error) bool {
	return isErrorType(err, ParseError)
}
