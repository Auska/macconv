/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

// Package validator provides common validation utilities for the macconv application.
// It includes validation functions for MAC addresses, IP addresses, ports, and other network-related inputs.
package validator

import (
	"net"
	"regexp"
	"strings"

	"macconv/pkg/errors"
)

const (
	macAddressLength  = 12
	maxFilePathLength = 4096
	minPort           = 1
	maxPort           = 65535
)

var (
	macAddressPattern = regexp.MustCompile(`^[0-9a-f]{12}$`)
)

func ValidateMACAddress(mac string) error {
	if len(mac) != macAddressLength {
		return errors.New(errors.ValidationError, "MAC address must be 12 characters after normalization")
	}

	if !macAddressPattern.MatchString(mac) {
		return errors.New(errors.ValidationError, "MAC address contains invalid characters")
	}

	return nil
}

func ValidateIPAddress(ip string) error {
	if net.ParseIP(ip) == nil {
		return errors.New(errors.ValidationError, "invalid IP address format")
	}
	return nil
}

func ValidateIPv4Address(ip string) error {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return errors.New(errors.ValidationError, "invalid IP address format")
	}
	if parsedIP.To4() == nil {
		return errors.New(errors.ValidationError, "IPv6 address not supported, expected IPv4")
	}
	return nil
}

func ValidatePort(portStr string) (int, error) {
	port, err := parsePort(portStr)
	if err != nil {
		return 0, errors.New(errors.ValidationError, "port must be a number")
	}
	if port < minPort || port > maxPort {
		return 0, errors.New(errors.ValidationError, "port must be between 1 and 65535")
	}
	return port, nil
}

func parsePort(portStr string) (int, error) {
	var port int
	for _, c := range portStr {
		if c < '0' || c > '9' {
			return 0, errors.New(errors.ValidationError, "port must be a number")
		}
		port = port*10 + int(c-'0')
		if port > maxPort {
			return 0, errors.New(errors.ValidationError, "port must be between 1 and 65535")
		}
	}
	return port, nil
}

func ValidateCIDR(cidr string) error {
	_, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return errors.Wrap(errors.ParseError, "invalid CIDR format", err)
	}
	return nil
}

func ValidateFilePath(filePath string) error {
	if filePath == "" {
		return errors.New(errors.ValidationError, "file path cannot be empty")
	}

	if len(filePath) > maxFilePathLength {
		return errors.New(errors.ValidationError, "file path too long")
	}

	if strings.Contains(filePath, "..") {
		return errors.New(errors.ValidationError, "file path contains invalid traversal sequence")
	}

	return nil
}
