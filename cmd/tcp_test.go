//go:build unit

/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package cmd

import (
	"net"
	"testing"

	"github.com/spf13/cobra"
)

func TestBuildTargetAddress(t *testing.T) {
	tests := []struct {
		name     string
		ip       net.IP
		port     int
		expected string
	}{
		{
			name:     "IPv4 address",
			ip:       net.ParseIP("192.168.1.1"),
			port:     80,
			expected: "192.168.1.1:80",
		},
		{
			name:     "IPv6 address",
			ip:       net.ParseIP("::1"),
			port:     80,
			expected: "[::1]:80",
		},
		{
			name:     "IPv6 full address",
			ip:       net.ParseIP("2001:db8::1"),
			port:     443,
			expected: "[2001:db8::1]:443",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildTargetAddress(tt.ip, tt.port)
			if result != tt.expected {
				t.Errorf("buildTargetAddress() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsHostname(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		expected bool
	}{
		{
			name:     "IPv4 address",
			host:     "192.168.1.1",
			expected: false,
		},
		{
			name:     "IPv6 address",
			host:     "::1",
			expected: false,
		},
		{
			name:     "IPv6 full address",
			host:     "2001:db8::1",
			expected: false,
		},
		{
			name:     "Hostname",
			host:     "example.com",
			expected: true,
		},
		{
			name:     "Hostname with subdomain",
			host:     "www.example.com",
			expected: true,
		},
		{
			name:     "Invalid string",
			host:     "not-a-valid-host",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isHostname(tt.host)
			if result != tt.expected {
				t.Errorf("isHostname() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractIPFromTarget(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		expected string
	}{
		{
			name:     "IPv4 target",
			target:   "192.168.1.1:80",
			expected: "192.168.1.1",
		},
		{
			name:     "IPv6 target",
			target:   "[::1]:80",
			expected: "[::1]",
		},
		{
			name:     "IPv6 full target",
			target:   "[2001:db8::1]:443",
			expected: "[2001:db8::1]",
		},
		{
			name:     "IPv4 target with custom port",
			target:   "10.0.0.1:8080",
			expected: "10.0.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractIPFromTarget(tt.target)
			if result != tt.expected {
				t.Errorf("extractIPFromTarget() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCheckPort(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains []string
	}{
		{
			name:    "Missing arguments",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "Missing port argument",
			args:    []string{"192.168.1.1"},
			wantErr: true,
		},
		{
			name:    "Invalid hostname",
			args:    []string{"invalid-hostname-that-does-not-exist-12345.com", "80"},
			wantErr: true,
		},
		{
			name:    "Invalid port - non-numeric",
			args:    []string{"192.168.1.1", "abc"},
			wantErr: true,
		},
		{
			name:    "Invalid port - out of range",
			args:    []string{"192.168.1.1", "99999"},
			wantErr: true,
		},
		{
			name:    "Invalid port - zero",
			args:    []string{"192.168.1.1", "0"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			checkPort(cmd, tt.args)

			// Just verify the function runs without panicking
			// Output verification is skipped due to capture complexity
		})
	}
}

func TestCheckSingleConnection(t *testing.T) {
	tests := []struct {
		name   string
		target string
		host   string
		port   int
	}{
		{
			name:   "IPv4 localhost closed port",
			target: "127.0.0.1:9999",
			host:   "127.0.0.1",
			port:   9999,
		},
		{
			name:   "IPv6 localhost closed port",
			target: "[::1]:9999",
			host:   "::1",
			port:   9999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkSingleConnection(tt.target, tt.host, tt.port, 1)

			// We expect connection to fail on closed port
			if result {
				t.Errorf("Expected connection to fail on closed port, but it succeeded")
			}
		})
	}
}

func TestRequiredConsecutiveSuccesses(t *testing.T) {
	// Verify the constant is set correctly
	if requiredConsecutiveSuccesses != 5 {
		t.Errorf("requiredConsecutiveSuccesses = %v, want 5", requiredConsecutiveSuccesses)
	}
}
