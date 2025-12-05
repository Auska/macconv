/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package cmd

import (
	"fmt"
	"net"

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
		cmd.Help()
		return
	}
	
	cidr := args[0]
	logger.Debug("Processing CIDR address: %s", cidr)
	
	info, err := calculateCIDRInfo(cidr)
	if err != nil {
		logger.PrintErrorWithMessage("failed to parse CIDR address", err)
		cmd.Help()
		return
	}

	fmt.Println("CIDR Address Range:", info.FirstIP, "-", info.LastIP)
	fmt.Println("Subnet Mask:", info.SubnetMask)
	fmt.Println("Network ID:", info.NetworkID)
	fmt.Println("Broadcast Address:", info.BroadcastAddress)
	fmt.Println("Total Hosts:", info.TotalHosts)
	
	logger.Info("Successfully processed CIDR address: %s", cidr)
}
// CIDRInfo 包含 CIDR 网络信息
type CIDRInfo struct {
	NetworkID       string
	FirstIP         string
	LastIP          string
	BroadcastAddress string
	SubnetMask      string
	TotalHosts      int
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
	totalHosts := 1 << (bits - ones)
	if bits == 32 { // IPv4
		totalHosts -= 2 // 减去网络地址和广播地址
	}

	return &CIDRInfo{
		NetworkID:       network.String(),
		FirstIP:         firstIP.String(),
		LastIP:          lastIP.String(),
		BroadcastAddress: broadcast.String(),
		SubnetMask:      net.IP(ipnet.Mask).String(),
		TotalHosts:      totalHosts,
	}, nil
}
