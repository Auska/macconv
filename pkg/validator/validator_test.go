/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package validator

import (
	"testing"
)

func TestValidateMACAddress(t *testing.T) {
	tests := []struct {
		name    string
		mac     string
		wantErr bool
	}{
		{"Valid MAC", "001122334455", false},
		{"Valid MAC with hex", "aabbccddeeff", false},
		{"Invalid MAC - too short", "00112233445", true},
		{"Invalid MAC - too long", "00112233445566", true},
		{"Invalid MAC - invalid chars", "00112233445g", true},
		{"Empty MAC", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMACAddress(tt.mac)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMACAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateIPAddress(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		wantErr bool
	}{
		{"Valid IPv4", "192.168.1.1", false},
		{"Valid IPv6", "::1", false},
		{"Valid IPv6 full", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", false},
		{"Invalid IP", "not.an.ip", true},
		{"Empty IP", "", true},
		{"Invalid IPv4", "256.256.256.256", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIPAddress(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIPAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateIPv4Address(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		wantErr bool
	}{
		{"Valid IPv4", "192.168.1.1", false},
		{"Invalid IPv6", "::1", true},
		{"Invalid IP", "not.an.ip", true},
		{"Empty IP", "", true},
		{"Invalid IPv4", "256.256.256.256", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIPv4Address(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIPv4Address() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		name     string
		portStr  string
		wantPort int
		wantErr  bool
	}{
		{"Valid port", "80", 80, false},
		{"Valid port max", "65535", 65535, false},
		{"Valid port min", "1", 1, false},
		{"Invalid port - too high", "65536", 0, true},
		{"Invalid port - zero", "0", 0, true},
		{"Invalid port - negative", "-1", 0, true},
		{"Invalid port - not number", "abc", 0, true},
		{"Empty port", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPort, err := ValidatePort(tt.portStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPort != tt.wantPort {
				t.Errorf("ValidatePort() = %v, want %v", gotPort, tt.wantPort)
			}
		})
	}
}

func TestValidateCIDR(t *testing.T) {
	tests := []struct {
		name    string
		cidr    string
		wantErr bool
	}{
		{"Valid CIDR", "192.168.1.0/24", false},
		{"Valid CIDR /32", "192.168.1.1/32", false},
		{"Valid CIDR /8", "10.0.0.0/8", false},
		{"Invalid CIDR - no mask", "192.168.1.0", true},
		{"Invalid CIDR - invalid mask", "192.168.1.0/33", true},
		{"Invalid CIDR - invalid IP", "256.256.256.256/24", true},
		{"Empty CIDR", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCIDR(tt.cidr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCIDR() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFilePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"Valid path", "/tmp/test.txt", false},
		{"Valid relative path", "./test.txt", false},
		{"Empty path", "", true},
		{"Too long path", string(make([]byte, 4097)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFilePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFilePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
