/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package cmd

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/spf13/cobra"
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

// 主函数用于测试
func dhcp(cmd *cobra.Command, args []string) {
	var ips []net.IP
	if len(args) != 1 && len(args) != 2 {
		fmt.Println("Error: Invalid number of arguments. Expected 1 or 2 IP addresses.")
		cmd.Help()
		os.Exit(1)
	}

	for _, arg := range args {
		ip := net.ParseIP(arg)
		if ip == nil {
			fmt.Printf("Error: Invalid IP address: %s\n", arg)
			cmd.Help()
			os.Exit(1)
		}
		ips = append(ips, ip)
	}

	fmt.Println("PXE Format: ", toPXEFormat(ips))
	fmt.Println("ACS Format: ", toACSFormat(ips))
	fmt.Println("PXE Format (Bytes): ", toPXEFormatBytes(ips))
	fmt.Println("ACS Format (Bytes): ", toACSFormatBytes(ips))
}
