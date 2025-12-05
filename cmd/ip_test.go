//go:build unit

/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package cmd

import (
	"net"
	"testing"
)

func TestCalculateCIDRInfo(t *testing.T) {
	tests := []struct {
		name     string
		cidr     string
		wantErr  bool
		expected *CIDRInfo
	}{
		{
			name:    "Valid /24 network",
			cidr:    "192.168.1.0/24",
			wantErr: false,
			expected: &CIDRInfo{
				NetworkID:        "192.168.1.0",
				FirstIP:          "192.168.1.1",
				LastIP:           "192.168.1.254",
				BroadcastAddress: "192.168.1.255",
				SubnetMask:       "255.255.255.0",
				InverseMask:      "0.0.0.255",
				TotalHosts:       254,
			},
		},
		{
			name:    "Valid /32 network",
			cidr:    "192.168.1.1/32",
			wantErr: false,
			expected: &CIDRInfo{
				NetworkID:        "192.168.1.1",
				FirstIP:          "192.168.1.2", // Note: This will wrap around, but it's expected behavior
				LastIP:           "192.168.1.0", // Note: This will wrap around
				BroadcastAddress: "192.168.1.1",
				SubnetMask:       "255.255.255.255",
				InverseMask:      "0.0.0.0",
				TotalHosts:       1,
			},
		},
		{
			name:    "Valid /8 network",
			cidr:    "10.0.0.0/8",
			wantErr: false,
			expected: &CIDRInfo{
				NetworkID:        "10.0.0.0",
				FirstIP:          "10.0.0.1",
				LastIP:           "10.255.255.254",
				BroadcastAddress: "10.255.255.255",
				SubnetMask:       "255.0.0.0",
				InverseMask:      "0.255.255.255",
				TotalHosts:       16777214,
			},
		},
		{
			name:     "Invalid CIDR - no mask",
			cidr:     "192.168.1.0",
			wantErr:  true,
			expected: nil,
		},
		{
			name:     "Invalid CIDR - invalid mask",
			cidr:     "192.168.1.0/33",
			wantErr:  true,
			expected: nil,
		},
		{
			name:     "Invalid CIDR - invalid IP",
			cidr:     "256.256.256.256/24",
			wantErr:  true,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateCIDRInfo(tt.cidr)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateCIDRInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.NetworkID != tt.expected.NetworkID {
					t.Errorf("calculateCIDRInfo() NetworkID = %v, want %v", got.NetworkID, tt.expected.NetworkID)
				}
				if got.SubnetMask != tt.expected.SubnetMask {
					t.Errorf("calculateCIDRInfo() SubnetMask = %v, want %v", got.SubnetMask, tt.expected.SubnetMask)
				}
				if got.TotalHosts != tt.expected.TotalHosts {
					t.Errorf("calculateCIDRInfo() TotalHosts = %v, want %v", got.TotalHosts, tt.expected.TotalHosts)
				}
			}
		})
	}
}

func TestBuildTargetAddress(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		port     int
		expected string
	}{
		{
			name:     "IPv4 address",
			ip:       "192.168.1.1",
			port:     80,
			expected: "192.168.1.1:80",
		},
		{
			name:     "IPv6 address",
			ip:       "::1",
			port:     80,
			expected: "[::1]:80",
		},
		{
			name:     "IPv6 full address",
			ip:       "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			port:     443,
			expected: "[2001:db8:85a3::8a2e:370:7334]:443",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			result := buildTargetAddress(ip, tt.port)
			if result != tt.expected {
				t.Errorf("buildTargetAddress() = %v, want %v", result, tt.expected)
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
			name:     "/24 mask",
			mask:     "255.255.255.0",
			expected: "0.0.0.255",
		},
		{
			name:     "/16 mask",
			mask:     "255.255.0.0",
			expected: "0.0.255.255",
		},
		{
			name:     "/8 mask",
			mask:     "255.0.0.0",
			expected: "0.255.255.255",
		},
		{
			name:     "/32 mask",
			mask:     "255.255.255.255",
			expected: "0.0.0.0",
		},
		{
			name:     "/0 mask",
			mask:     "0.0.0.0",
			expected: "255.255.255.255",
		},
		{
			name:     "/26 mask",
			mask:     "255.255.255.192",
			expected: "0.0.0.63",
		},
		{
			name:     "/30 mask",
			mask:     "255.255.255.252",
			expected: "0.0.0.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipMask := net.IPMask(net.ParseIP(tt.mask).To4())
			result := calculateInverseMask(ipMask)
			if result != tt.expected {
				t.Errorf("calculateInverseMask() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCalculateCIDRInfoIPv6(t *testing.T) {
	tests := []struct {
		name     string
		cidr     string
		wantErr  bool
		expected *CIDRInfo
	}{
		{
			name:    "Valid IPv6 /126 network",
			cidr:    "2001:db8::/126",
			wantErr: false,
			expected: &CIDRInfo{
				NetworkID:        "2001:db8::",
				FirstIP:          "2001:db8::",
				LastIP:           "2001:db8::3",
				BroadcastAddress: "", // IPv6 doesn't have broadcast
				SubnetMask:       "ffff:ffff:ffff:ffff:ffff:ffff:ffff:fffc",
				InverseMask:      "::3",
				TotalHosts:       4,
			},
		},
		{
			name:    "Valid IPv6 /128 network",
			cidr:    "2001:db8::1/128",
			wantErr: false,
			expected: &CIDRInfo{
				NetworkID:        "2001:db8::1",
				FirstIP:          "2001:db8::1",
				LastIP:           "2001:db8::1",
				BroadcastAddress: "", // IPv6 doesn't have broadcast
				SubnetMask:       "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
				InverseMask:      "::",
				TotalHosts:       1,
			},
		},
		{
			name:    "Valid IPv6 /32 network (large)",
			cidr:    "2001:db8::/32",
			wantErr: false,
			expected: &CIDRInfo{
				NetworkID:        "2001:db8::",
				FirstIP:          "2001:db8::",
				LastIP:           "2001:db8:ffff:ffff:ffff:ffff:ffff:ffff",
				BroadcastAddress: "", // IPv6 doesn't have broadcast
				SubnetMask:       "ffff:ffff::",
				InverseMask:      "::ffff:ffff:ffff:ffff:ffff:ffff",
				TotalHosts:       -1, // Very large number
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateCIDRInfo(tt.cidr)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateCIDRInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.NetworkID != tt.expected.NetworkID {
					t.Errorf("calculateCIDRInfo() NetworkID = %v, want %v", got.NetworkID, tt.expected.NetworkID)
				}
				if got.SubnetMask != tt.expected.SubnetMask {
					t.Errorf("calculateCIDRInfo() SubnetMask = %v, want %v", got.SubnetMask, tt.expected.SubnetMask)
				}
				if got.InverseMask != tt.expected.InverseMask {
					t.Errorf("calculateCIDRInfo() InverseMask = %v, want %v", got.InverseMask, tt.expected.InverseMask)
				}
				if got.TotalHosts != tt.expected.TotalHosts {
					t.Errorf("calculateCIDRInfo() TotalHosts = %v, want %v", got.TotalHosts, tt.expected.TotalHosts)
				}
			}
		})
	}
}
