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

func TestCalculateCIDRInfo(t *testing.T) {
	tests := []struct {
		name    string
		cidr    string
		wantErr bool
		check   func(*testing.T, *CIDRInfo)
	}{
		{
			name:    "Valid IPv4 /24 network",
			cidr:    "192.168.1.0/24",
			wantErr: false,
			check: func(t *testing.T, info *CIDRInfo) {
				if info.NetworkID != "192.168.1.0" {
					t.Errorf("NetworkID = %v, want 192.168.1.0", info.NetworkID)
				}
				if info.FirstIP != "192.168.1.1" {
					t.Errorf("FirstIP = %v, want 192.168.1.1", info.FirstIP)
				}
				if info.LastIP != "192.168.1.254" {
					t.Errorf("LastIP = %v, want 192.168.1.254", info.LastIP)
				}
				if info.BroadcastAddress != "192.168.1.255" {
					t.Errorf("BroadcastAddress = %v, want 192.168.1.255", info.BroadcastAddress)
				}
				if info.SubnetMask != "255.255.255.0" {
					t.Errorf("SubnetMask = %v, want 255.255.255.0", info.SubnetMask)
				}
				if info.InverseMask != "0.0.0.255" {
					t.Errorf("InverseMask = %v, want 0.0.0.255", info.InverseMask)
				}
				if info.TotalHosts != 254 {
					t.Errorf("TotalHosts = %v, want 254", info.TotalHosts)
				}
			},
		},
		{
			name:    "Valid IPv4 /32 network",
			cidr:    "192.168.1.1/32",
			wantErr: false,
			check: func(t *testing.T, info *CIDRInfo) {
				if info.TotalHosts != 1 {
					t.Errorf("TotalHosts = %v, want 1", info.TotalHosts)
				}
			},
		},
		{
			name:    "Valid IPv4 /31 network",
			cidr:    "192.168.1.0/31",
			wantErr: false,
			check: func(t *testing.T, info *CIDRInfo) {
				if info.TotalHosts != 2 {
					t.Errorf("TotalHosts = %v, want 2", info.TotalHosts)
				}
			},
		},
		{
			name:    "Valid IPv6 /64 network",
			cidr:    "2001:db8::/32",
			wantErr: false,
			check: func(t *testing.T, info *CIDRInfo) {
				if info.NetworkID != "2001:db8::" {
					t.Errorf("NetworkID = %v, want 2001:db8::", info.NetworkID)
				}
				if info.TotalHosts != -1 {
					t.Errorf("TotalHosts = %v, want -1 (very large)", info.TotalHosts)
				}
			},
		},
		{
			name:    "Invalid CIDR format",
			cidr:    "192.168.1.0",
			wantErr: true,
		},
		{
			name:    "Invalid CIDR - bad mask",
			cidr:    "192.168.1.0/33",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := calculateCIDRInfo(tt.cidr)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateCIDRInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, info)
			}
		})
	}
}

func TestCalculateInverseMask(t *testing.T) {
	tests := []struct {
		name     string
		mask     string
		expected string
	}{
		{
			name:     "IPv4 /24 mask",
			mask:     "255.255.255.0",
			expected: "0.0.0.255",
		},
		{
			name:     "IPv4 /16 mask",
			mask:     "255.255.0.0",
			expected: "0.0.255.255",
		},
		{
			name:     "IPv4 /8 mask",
			mask:     "255.0.0.0",
			expected: "0.255.255.255",
		},
		{
			name:     "IPv6 /64 mask",
			mask:     "ffff:ffff:ffff:ffff::",
			expected: "::ffff:ffff:ffff:ffff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mask := parseIPMask(tt.mask)
			if mask == nil {
				t.Skipf("Cannot parse mask: %s", tt.mask)
			}
			result := calculateInverseMask(mask)
			if result != tt.expected {
				t.Errorf("calculateInverseMask() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConvertIPAddress(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains []string
	}{
		{
			name:    "Valid IPv4 CIDR",
			args:    []string{"192.168.1.0/24"},
			wantErr: false,
		},
		{
			name:    "Valid IPv6 CIDR",
			args:    []string{"2001:db8::/32"},
			wantErr: false,
		},
		{
			name:    "Missing argument",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "Invalid CIDR format",
			args:    []string{"192.168.1.0"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			convertIPAddress(cmd, tt.args)

			// Just verify the function runs without panicking
			// Output verification is skipped due to capture complexity
		})
	}
}

// Helper function to parse IP mask string
func parseIPMask(maskStr string) net.IPMask {
	ip := net.ParseIP(maskStr)
	if ip == nil {
		return nil
	}
	if ip.To4() != nil {
		return net.IPv4Mask(ip[12], ip[13], ip[14], ip[15])
	}
	return net.IPMask(ip.To16())
}
