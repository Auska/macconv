/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package cmd

import (
	"fmt"
	"net"
	"strings"

	"github.com/spf13/cobra"
	"macconv/pkg/errors"
	"macconv/pkg/logger"
)

// dhcpCmd represents the dhcp command
var dhcpCmd = &cobra.Command{
	Use:   "dhcp",
	Short: "DHCP option 43 conversion",
	Long: `
Commonly used for configuring DHCP servers, converting IP addresses to hexadecimal strings, and transforming them into PXE and ACS formats.
For example:

	macconv dhcp 192.168.1.1`,
	Run: dhcp,
}

func init() {
	rootCmd.AddCommand(dhcpCmd)
}

// 将单个IPv4地址转换为十六进制字符串
func ipToHex(ip net.IP) string {
	ip = ip.To4()
	return fmt.Sprintf("%02x%02x%02x%02x", ip[0], ip[1], ip[2], ip[3])
}

// 将单个IPv4地址转换为每个字节的十六进制表示
func ipToHexBytes(ip net.IP) string {
	ip = ip.To4()
	var hexBytes []string
	for _, b := range ip {
		hexBytes = append(hexBytes, fmt.Sprintf("0x%02x", b))
	}
	return strings.Join(hexBytes, " ")
}

// 转换为PXE格式
func toPXEFormat(ips []net.IP) string {
	var hexString string
	hexString += "80"                              // PXE format identifier
	hexString += fmt.Sprintf("%02x", len(ips)*4+3) // Length of the following data
	hexString += "0000"                            // Fixed value
	hexString += fmt.Sprintf("%02x", len(ips))     // Number of IPs
	for _, ip := range ips {
		ipHex := ipToHex(ip)
		hexString += ipHex
	}
	return hexString
}

// 转换为ACS格式
func toACSFormat(ips []net.IP) string {
	var hexString string
	hexString += "01"                            // ACS format identifier
	hexString += fmt.Sprintf("%02x", len(ips)*4) // Length of the following data (each IP is 4 bytes)
	for _, ip := range ips {
		ipHex := ipToHex(ip)
		hexString += ipHex
	}
	return hexString
}

// 转换为PXE格式（每个字节的十六进制表示）
func toPXEFormatBytes(ips []net.IP) string {
	var hexString string
	hexString += "0x80"                               // PXE format identifier
	hexString += fmt.Sprintf(" 0x%02x", len(ips)*4+3) // Length of the following data
	hexString += " 0x00 0x00"                         // Fixed value
	hexString += fmt.Sprintf(" 0x%02x", len(ips))     // Number of IPs
	for _, ip := range ips {
		ipHex := ipToHexBytes(ip)
		hexString += " " + ipHex
	}
	return hexString
}

// 转换为ACS格式（每个字节的十六进制表示）
func toACSFormatBytes(ips []net.IP) string {
	var hexString string
	hexString += "0x01"                             // ACS format identifier
	hexString += fmt.Sprintf(" 0x%02x", len(ips)*4) // Length of the following data (each IP is 4 bytes)
	for _, ip := range ips {
		ipHex := ipToHexBytes(ip)
		hexString += " " + ipHex
	}
	return hexString
}

// dhcp 处理 DHCP 选项 43 转换命令
func dhcp(cmd *cobra.Command, args []string) {
	if len(args) < 1 || len(args) > 2 {
		logger.PrintValidationError("invalid number of arguments, expected 1 or 2 IP addresses")
		cmd.Help()
		return
	}

	ips, err := parseIPAddresses(args)
	if err != nil {
		logger.PrintErrorWithMessage("failed to parse IP addresses", err)
		cmd.Help()
		return
	}

	logger.Debug("Processing %d IP addresses for DHCP option 43", len(ips))

	fmt.Println("PXE Format: ", toPXEFormat(ips))
	fmt.Println("ACS Format: ", toACSFormat(ips))
	fmt.Println("PXE Format (Bytes): ", toPXEFormatBytes(ips))
	fmt.Println("ACS Format (Bytes): ", toACSFormatBytes(ips))
	
	logger.Info("Successfully processed DHCP option 43 conversion for %d IP addresses", len(ips))
}

// parseIPAddresses 解析并验证 IP 地址列表
func parseIPAddresses(args []string) ([]net.IP, error) {
	var ips []net.IP
	
	for _, arg := range args {
		ip := net.ParseIP(arg)
		if ip == nil {
			return nil, errors.New(errors.ValidationError, fmt.Sprintf("invalid IP address: %s", arg))
		}
		
		// 只支持 IPv4 地址
		if ip.To4() == nil {
			return nil, errors.New(errors.ValidationError, fmt.Sprintf("IPv6 address not supported: %s", arg))
		}
		
		ips = append(ips, ip)
	}
	
	return ips, nil
}
