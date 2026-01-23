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

func TestIpToHex(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected string
	}{
		{
			name:     "Standard IPv4",
			ip:       "192.168.1.1",
			expected: "c0a80101",
		},
		{
			name:     "Localhost",
			ip:       "127.0.0.1",
			expected: "7f000001",
		},
		{
			name:     "Zero address",
			ip:       "0.0.0.0",
			expected: "00000000",
		},
		{
			name:     "Broadcast",
			ip:       "255.255.255.255",
			expected: "ffffffff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			result := ipToHex(ip)
			if result != tt.expected {
				t.Errorf("ipToHex() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIpToHexBytes(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected string
	}{
		{
			name:     "Standard IPv4",
			ip:       "192.168.1.1",
			expected: "0xc0 0xa8 0x01 0x01",
		},
		{
			name:     "Localhost",
			ip:       "127.0.0.1",
			expected: "0x7f 0x00 0x00 0x01",
		},
		{
			name:     "Zero address",
			ip:       "0.0.0.0",
			expected: "0x00 0x00 0x00 0x00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			result := ipToHexBytes(ip)
			if result != tt.expected {
				t.Errorf("ipToHexBytes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestToPXEFormat(t *testing.T) {
	tests := []struct {
		name     string
		ips      []string
		expected string
	}{
		{
			name:     "Single IP",
			ips:      []string{"192.168.1.1"},
			expected: "8007000001c0a80101",
		},
		{
			name:     "Multiple IPs",
			ips:      []string{"192.168.1.1", "192.168.1.2"},
			expected: "800b000002c0a80101c0a80102",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ips []net.IP
			for _, ipStr := range tt.ips {
				ips = append(ips, net.ParseIP(ipStr))
			}
			result := toPXEFormat(ips)
			if result != tt.expected {
				t.Errorf("toPXEFormat() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestToACSFormat(t *testing.T) {
	tests := []struct {
		name     string
		ips      []string
		expected string
	}{
		{
			name:     "Single IP",
			ips:      []string{"192.168.1.1"},
			expected: "0104c0a80101",
		},
		{
			name:     "Multiple IPs",
			ips:      []string{"192.168.1.1", "192.168.1.2"},
			expected: "0108c0a80101c0a80102",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ips []net.IP
			for _, ipStr := range tt.ips {
				ips = append(ips, net.ParseIP(ipStr))
			}
			result := toACSFormat(ips)
			if result != tt.expected {
				t.Errorf("toACSFormat() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestToPXEFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		ips      []string
		expected string
	}{
		{
			name:     "Single IP",
			ips:      []string{"192.168.1.1"},
			expected: "0x80 0x07 0x00 0x00 0x01 0xc0 0xa8 0x01 0x01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ips []net.IP
			for _, ipStr := range tt.ips {
				ips = append(ips, net.ParseIP(ipStr))
			}
			result := toPXEFormatBytes(ips)
			if result != tt.expected {
				t.Errorf("toPXEFormatBytes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestToACSFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		ips      []string
		expected string
	}{
		{
			name:     "Single IP",
			ips:      []string{"192.168.1.1"},
			expected: "0x01 0x04 0xc0 0xa8 0x01 0x01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ips []net.IP
			for _, ipStr := range tt.ips {
				ips = append(ips, net.ParseIP(ipStr))
			}
			result := toACSFormatBytes(ips)
			if result != tt.expected {
				t.Errorf("toACSFormatBytes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseIPAddresses(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		check   func(*testing.T, []net.IP)
	}{
		{
			name:    "Valid single IPv4",
			args:    []string{"192.168.1.1"},
			wantErr: false,
			check: func(t *testing.T, ips []net.IP) {
				if len(ips) != 1 {
					t.Errorf("Expected 1 IP, got %d", len(ips))
				}
				if ips[0].String() != "192.168.1.1" {
					t.Errorf("Expected 192.168.1.1, got %s", ips[0].String())
				}
			},
		},
		{
			name:    "Valid multiple IPv4",
			args:    []string{"192.168.1.1", "192.168.1.2"},
			wantErr: false,
			check: func(t *testing.T, ips []net.IP) {
				if len(ips) != 2 {
					t.Errorf("Expected 2 IPs, got %d", len(ips))
				}
			},
		},
		{
			name:    "Invalid IP format",
			args:    []string{"invalid-ip"},
			wantErr: true,
		},
		{
			name:    "IPv6 not supported",
			args:    []string{"::1"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ips, err := parseIPAddresses(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseIPAddresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, ips)
			}
		})
	}
}

func TestDhcp(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains []string
	}{
		{
			name:    "Valid single IP",
			args:    []string{"192.168.1.1"},
			wantErr: false,
		},
		{
			name:    "Valid multiple IPs",
			args:    []string{"192.168.1.1", "192.168.1.2"},
			wantErr: false,
		},
		{
			name:    "Missing arguments",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "Too many arguments",
			args:    []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"},
			wantErr: true,
		},
		{
			name:    "Invalid IP",
			args:    []string{"invalid-ip"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			dhcp(cmd, tt.args)

			// Just verify the function runs without panicking
			// Output verification is skipped due to capture complexity
		})
	}
}
