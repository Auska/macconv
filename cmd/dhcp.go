/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

// dhcpCmd represents the dhcp command
var dhcpCmd = &cobra.Command{
	Use:   "dhcp",
	Short: "DHCP option 43 conversion",
	Long: `Commonly used for configuring DHCP servers, converting IP addresses to hexadecimal strings, and transforming them into PXE and ACS formats.`,
	Run: dhcp,
}

func init() {
	rootCmd.AddCommand(dhcpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dhcpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dhcpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// 将单个IPv4地址转换为十六进制字符串
func ipToHex(ip net.IP) string {
	ip = ip.To4()
	return fmt.Sprintf("%02x%02x%02x%02x", ip[0], ip[1], ip[2], ip[3])
}

// 转换为PXE格式
func toPXEFormat(ips []net.IP) string {
	var hexString string
	hexString += "80" // PXE format identifier
	hexString += fmt.Sprintf("%02x", len(ips)*4+3) // Length of the following data
	hexString += "0000"                             // Fixed value
	hexString += fmt.Sprintf("%02x", len(ips))      // Number of IPs
	for _, ip := range ips {
		ipHex := ipToHex(ip)
		// fmt.Println(ipHex)
		// fmt.Println(ip)
		hexString += ipHex
	}
	return hexString
}

// 转换为ACS格式
func toACSFormat(ips []net.IP) string {
	var hexString string
	hexString += "01" // ACS format identifier
	hexString += fmt.Sprintf("%02x", len(ips)*4) // Length of the following data (each IP is 4 bytes)
	for _, ip := range ips {
		ipHex := ipToHex(ip)
		// fmt.Println(ipHex)
		// fmt.Println(ip)
		hexString += ipHex
	}
	return hexString
}

// 主函数用于测试
func dhcp(cmd *cobra.Command, args []string) {
	var ips []net.IP
	if len(args) != 1 && len(args) != 2 {
		fmt.Println("Error: Invalid number of arguments. Expected 1 or 2 IP addresses.")
		return
	}

	for _, arg := range args {
		ip := net.ParseIP(arg)
		if ip == nil {
			fmt.Printf("Error: Invalid IP address: %s\n", arg)
			return
		}
		ips = append(ips, ip)
	}

	fmt.Println("PXE Format: ", toPXEFormat(ips))
	fmt.Println("ACS Format: ", toACSFormat(ips))
}