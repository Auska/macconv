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

func ipToHex(ip net.IP) string {
	ip = ip.To4()
	return fmt.Sprintf("%02x%02x%02x%02x", ip[0], ip[1], ip[2], ip[3])
}

func ipToHexBytes(ip net.IP) string {
	ip = ip.To4()
	var sb strings.Builder
	sb.Grow(14)
	sb.WriteString(fmt.Sprintf("0x%02x", ip[0]))
	for i := 1; i < 4; i++ {
		sb.WriteString(fmt.Sprintf(" 0x%02x", ip[i]))
	}
	return sb.String()
}

func toPXEFormat(ips []net.IP) string {
	var sb strings.Builder
	sb.Grow(8 + len(ips)*8)
	sb.WriteString("80")
	sb.WriteString(fmt.Sprintf("%02x", len(ips)*4+3))
	sb.WriteString("0000")
	sb.WriteString(fmt.Sprintf("%02x", len(ips)))
	for _, ip := range ips {
		sb.WriteString(ipToHex(ip))
	}
	return sb.String()
}

func toACSFormat(ips []net.IP) string {
	var sb strings.Builder
	sb.Grow(4 + len(ips)*8)
	sb.WriteString("01")
	sb.WriteString(fmt.Sprintf("%02x", len(ips)*4))
	for _, ip := range ips {
		sb.WriteString(ipToHex(ip))
	}
	return sb.String()
}

func toPXEFormatBytes(ips []net.IP) string {
	var sb strings.Builder
	sb.Grow(20 + len(ips)*15)
	sb.WriteString("0x80")
	sb.WriteString(fmt.Sprintf(" 0x%02x", len(ips)*4+3))
	sb.WriteString(" 0x00 0x00")
	sb.WriteString(fmt.Sprintf(" 0x%02x", len(ips)))
	for _, ip := range ips {
		sb.WriteString(" ")
		sb.WriteString(ipToHexBytes(ip))
	}
	return sb.String()
}

func toACSFormatBytes(ips []net.IP) string {
	var sb strings.Builder
	sb.Grow(8 + len(ips)*15)
	sb.WriteString("0x01")
	sb.WriteString(fmt.Sprintf(" 0x%02x", len(ips)*4))
	for _, ip := range ips {
		sb.WriteString(" ")
		sb.WriteString(ipToHexBytes(ip))
	}
	return sb.String()
}

func dhcp(cmd *cobra.Command, args []string) {
	if len(args) < 1 || len(args) > 2 {
		logger.PrintValidationError("invalid number of arguments, expected 1 or 2 IP addresses")
		if err := cmd.Help(); err != nil {
			logger.PrintErrorWithMessage("failed to show help", err)
		}
		return
	}

	ips, err := parseIPAddresses(args)
	if err != nil {
		logger.PrintErrorWithMessage("failed to parse IP addresses", err)
		if err := cmd.Help(); err != nil {
			logger.PrintErrorWithMessage("failed to show help", err)
		}
		return
	}

	logger.Debugf("Processing %d IP addresses for DHCP option 43", len(ips))

	fmt.Println("PXE Format: ", toPXEFormat(ips))
	fmt.Println("ACS Format: ", toACSFormat(ips))
	fmt.Println("PXE Format (Bytes): ", toPXEFormatBytes(ips))
	fmt.Println("ACS Format (Bytes): ", toACSFormatBytes(ips))

	logger.Infof("Successfully processed DHCP option 43 conversion for %d IP addresses", len(ips))
}

func parseIPAddresses(args []string) ([]net.IP, error) {
	ips := make([]net.IP, 0, len(args))

	for _, arg := range args {
		ip := net.ParseIP(arg)
		if ip == nil {
			return nil, errors.New(errors.ValidationError, fmt.Sprintf("invalid IP address: %s", arg))
		}

		if ip.To4() == nil {
			return nil, errors.New(errors.ValidationError, fmt.Sprintf("IPv6 address not supported: %s", arg))
		}

		ips = append(ips, ip)
	}

	return ips, nil
}
