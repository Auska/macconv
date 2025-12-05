/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package errors

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	err := New(ValidationError, "test error")
	if err.Type != ValidationError {
		t.Errorf("New() error.Type = %v, want %v", err.Type, ValidationError)
	}
	if err.Message != "test error" {
		t.Errorf("New() error.Message = %v, want %v", err.Message, "test error")
	}
	if err.Err != nil {
		t.Errorf("New() error.Err = %v, want nil", err.Err)
	}
}

func TestWrap(t *testing.T) {
	originalErr := errors.New("original error")
	err := Wrap(NetworkError, "wrapped error", originalErr)
	
	if err.Type != NetworkError {
		t.Errorf("Wrap() error.Type = %v, want %v", err.Type, NetworkError)
	}
	if err.Message != "wrapped error" {
		t.Errorf("Wrap() error.Message = %v, want %v", err.Message, "wrapped error")
	}
	if err.Err != originalErr {
		t.Errorf("Wrap() error.Err = %v, want %v", err.Err, originalErr)
	}
}

func TestError(t *testing.T) {
	// Test without wrapped error
	err1 := New(ValidationError, "test error")
	expected1 := "test error"
	if err1.Error() != expected1 {
		t.Errorf("Error() = %v, want %v", err1.Error(), expected1)
	}
	
	// Test with wrapped error
	originalErr := errors.New("original error")
	err2 := Wrap(NetworkError, "wrapped error", originalErr)
	expected2 := "wrapped error: original error"
	if err2.Error() != expected2 {
		t.Errorf("Error() = %v, want %v", err2.Error(), expected2)
	}
}

func TestUnwrap(t *testing.T) {
	// Test without wrapped error
	err1 := New(ValidationError, "test error")
	if err1.Unwrap() != nil {
		t.Errorf("Unwrap() = %v, want nil", err1.Unwrap())
	}
	
	// Test with wrapped error
	originalErr := errors.New("original error")
	err2 := Wrap(NetworkError, "wrapped error", originalErr)
	if err2.Unwrap() != originalErr {
		t.Errorf("Unwrap() = %v, want %v", err2.Unwrap(), originalErr)
	}
}

func TestIsValidationError(t *testing.T) {
	// Test validation error
	err1 := New(ValidationError, "test error")
	if !IsValidationError(err1) {
		t.Errorf("IsValidationError() = false, want true for validation error")
	}
	
	// Test network error
	err2 := New(NetworkError, "test error")
	if IsValidationError(err2) {
		t.Errorf("IsValidationError() = true, want false for network error")
	}
	
	// Test regular error
	err3 := errors.New("regular error")
	if IsValidationError(err3) {
		t.Errorf("IsValidationError() = true, want false for regular error")
	}
}

func TestIsNetworkError(t *testing.T) {
	// Test network error
	err1 := New(NetworkError, "test error")
	if !IsNetworkError(err1) {
		t.Errorf("IsNetworkError() = false, want true for network error")
	}
	
	// Test validation error
	err2 := New(ValidationError, "test error")
	if IsNetworkError(err2) {
		t.Errorf("IsNetworkError() = true, want false for validation error")
	}
	
	// Test regular error
	err3 := errors.New("regular error")
	if IsNetworkError(err3) {
		t.Errorf("IsNetworkError() = true, want false for regular error")
	}
}

func TestIsFileSystemError(t *testing.T) {
	// Test file system error
	err1 := New(FileSystemError, "test error")
	if !IsFileSystemError(err1) {
		t.Errorf("IsFileSystemError() = false, want true for file system error")
	}
	
	// Test validation error
	err2 := New(ValidationError, "test error")
	if IsFileSystemError(err2) {
		t.Errorf("IsFileSystemError() = true, want false for validation error")
	}
	
	// Test regular error
	err3 := errors.New("regular error")
	if IsFileSystemError(err3) {
		t.Errorf("IsFileSystemError() = true, want false for regular error")
	}
}

func TestIsParseError(t *testing.T) {
	// Test parse error
	err1 := New(ParseError, "test error")
	if !IsParseError(err1) {
		t.Errorf("IsParseError() = false, want true for parse error")
	}
	
	// Test validation error
	err2 := New(ValidationError, "test error")
	if IsParseError(err2) {
		t.Errorf("IsParseError() = true, want false for validation error")
	}
	
	// Test regular error
	err3 := errors.New("regular error")
	if IsParseError(err3) {
		t.Errorf("IsParseError() = true, want false for regular error")
	}
}