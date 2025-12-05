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

// ipCmd represents the ip command
var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "CIDR mask conversion",
	Long: `
CIDR mask conversion. For example:

	macconv ip 192.168.1.1/24`,
	Run: convertIPAddress,
}

func init() {
	rootCmd.AddCommand(ipCmd)
}

func convertIPAddress(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		logger.PrintValidationError("missing CIDR address argument")
		if err := cmd.Help(); err != nil {
			logger.PrintErrorWithMessage("failed to show help", err)
		}
		return
	}

	cidr := args[0]
	logger.Debug("Processing CIDR address: %s", cidr)

	info, err := calculateCIDRInfo(cidr)
	if err != nil {
		logger.PrintErrorWithMessage("failed to parse CIDR address", err)
		if err := cmd.Help(); err != nil {
			logger.PrintErrorWithMessage("failed to show help", err)
		}
		return
	}

	fmt.Println("CIDR Address Range:", info.FirstIP, "-", info.LastIP)
	fmt.Println("Subnet Mask:", info.SubnetMask)
	fmt.Println("Inverse Mask:", info.InverseMask)
	fmt.Println("Network ID:", info.NetworkID)

	// IPv6 没有广播地址概念
	if strings.Contains(info.NetworkID, ":") {
		fmt.Println("Network Type: IPv6")
	} else {
		fmt.Println("Broadcast Address:", info.BroadcastAddress)
	}

	// 显示主机数
	if info.TotalHosts == -1 {
		fmt.Println("Total Hosts: Very large number (>2^63)")
	} else {
		fmt.Println("Total Hosts:", info.TotalHosts)
	}

	logger.Info("Successfully processed CIDR address: %s", cidr)
}

// CIDRInfo 包含 CIDR 网络信息
type CIDRInfo struct {
	NetworkID        string
	FirstIP          string
	LastIP           string
	BroadcastAddress string
	SubnetMask       string
	InverseMask      string
	TotalHosts       int
}

func calculateCIDRInfo(cidr string) (*CIDRInfo, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, errors.Wrap(errors.ParseError, "invalid CIDR format", err)
	}

	// 计算网络号
	network := ip.Mask(ipnet.Mask)

	// 计算第一个可用IP（网络地址+1）
	firstIP := net.IP(make([]byte, len(network)))
	copy(firstIP, network)
	if len(firstIP) == 4 { // IPv4
		firstIP[3]++
	}

	// 计算广播地址
	broadcast := net.IP(make([]byte, len(network)))
	copy(broadcast, network)
	for i := range broadcast {
		broadcast[i] |= ^ipnet.Mask[i]
	}

	// 计算最后一个可用IP（广播地址-1）
	lastIP := net.IP(make([]byte, len(broadcast)))
	copy(lastIP, broadcast)
	if len(lastIP) == 4 { // IPv4
		lastIP[3]--
	}

	// 计算总主机数
	ones, bits := ipnet.Mask.Size()
	var totalHosts int

	if bits == 32 { // IPv4
		if ones == 32 {
			// /32 网络只有1个主机
			totalHosts = 1
		} else if ones == 31 {
			// /31 网络有2个主机（通常用于点对点连接）
			totalHosts = 2
		} else {
			// 其他情况计算可用主机数
			totalHosts = 1 << (bits - ones)
			totalHosts -= 2 // 减去网络地址和广播地址
		}
	} else {
		// IPv6 - 计算主机数
		hostBits := bits - ones
		if hostBits > 63 {
			// 如果主机位数超过63位，使用科学计数法表示
			totalHosts = -1 // 表示非常大的数
		} else {
			totalHosts = 1 << hostBits
		}
	}

	// 计算反掩码
	inverseMask := calculateInverseMask(ipnet.Mask)

	return &CIDRInfo{
		NetworkID:        network.String(),
		FirstIP:          firstIP.String(),
		LastIP:           lastIP.String(),
		BroadcastAddress: broadcast.String(),
		SubnetMask:       net.IP(ipnet.Mask).String(),
		InverseMask:      inverseMask,
		TotalHosts:       totalHosts,
	}, nil
}

// calculateInverseMask 计算反掩码（通配符掩码）
func calculateInverseMask(mask net.IPMask) string {
	if len(mask) != 4 && len(mask) != 16 {
		return ""
	}

	inverse := make(net.IPMask, len(mask))
	for i := 0; i < len(mask); i++ {
		inverse[i] = ^mask[i]
	}

	return net.IP(inverse).String()
}
