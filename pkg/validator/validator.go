/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

// Package validator provides common validation utilities for the macconv application.
// It includes validation functions for MAC addresses, IP addresses, ports, and other network-related inputs.
package validator

import (
	"net"
	"regexp"
	"strconv"

	"macconv/pkg/errors"
)

// ValidateMACAddress 验证 MAC 地址格式
func ValidateMACAddress(mac string) error {
	if len(mac) != 12 {
		return errors.New(errors.ValidationError, "MAC address must be 12 characters after normalization")
	}

	pattern := `^[0-9a-f]{12}$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(mac) {
		return errors.New(errors.ValidationError, "MAC address contains invalid characters")
	}

	return nil
}

// ValidateIPAddress 验证 IP 地址格式
func ValidateIPAddress(ip string) error {
	if net.ParseIP(ip) == nil {
		return errors.New(errors.ValidationError, "invalid IP address format")
	}
	return nil
}

// ValidateIPv4Address 验证 IPv4 地址格式
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

// ValidatePort 验证端口号
func ValidatePort(portStr string) (int, error) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, errors.New(errors.ValidationError, "port must be a number")
	}
	if port < 1 || port > 65535 {
		return 0, errors.New(errors.ValidationError, "port must be between 1 and 65535")
	}
	return port, nil
}

// ValidateCIDR 验证 CIDR 格式
func ValidateCIDR(cidr string) error {
	_, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return errors.Wrap(errors.ParseError, "invalid CIDR format", err)
	}
	return nil
}

// ValidateFilePath 验证文件路径
func ValidateFilePath(filePath string) error {
	if filePath == "" {
		return errors.New(errors.ValidationError, "file path cannot be empty")
	}

	// 检查路径长度
	if len(filePath) > 4096 {
		return errors.New(errors.ValidationError, "file path too long")
	}

	return nil
}
