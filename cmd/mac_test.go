//go:build unit

/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package cmd

import (
	"testing"
)

func TestNormalizeMACAddress(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"00:11:22:33:44:55", "001122334455"},
		{"00-11-22-33-44-55", "001122334455"},
		{"0011.2233.4455", "001122334455"},
		{"001122334455", "001122334455"},
		{"AA:BB:CC:DD:EE:FF", "aabbccddeeff"},
		{"aa-bb-cc-dd-ee-ff", "aabbccddeeff"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeMACAddress(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeMACAddress() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValidateMACAddress(t *testing.T) {
	tests := []struct {
		name    string
		mac     string
		wantErr bool
	}{
		{"Valid MAC", "001122334455", false},
		{"Valid MAC lowercase", "aabbccddeeff", false},
		{"Invalid MAC - too short", "00112233445", true},
		{"Invalid MAC - too long", "00112233445566", true},
		{"Invalid MAC - invalid chars", "00112233445g", true},
		{"Empty MAC", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMACAddress(tt.mac)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMACAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConvertMacAddress(t *testing.T) {
	tests := []struct {
		input    string
		step     int
		sep      string
		expected string
	}{
		{"001122334455", 2, ":", "00:11:22:33:44:55"},
		{"001122334455", 4, ".", "0011.2233.4455"},
		{"001122334455", 4, "-", "0011-2233-4455"},
		{"aabbccddeeff", 2, ":", "aa:bb:cc:dd:ee:ff"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := convertMacAddress(tt.input, tt.step, tt.sep)
			if result != tt.expected {
				t.Errorf("convertMacAddress() = %v, want %v", result, tt.expected)
			}
		})
	}
}
