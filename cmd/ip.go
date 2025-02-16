/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
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
		fmt.Printf("Error: missing arguments.")
		cmd.Help()
		os.Exit(1)
	}
	first, last, mask, err := calculateCIDRInfo(args[0])
	if err != nil {
		fmt.Println("Error:", err)
		cmd.Help()
		os.Exit(1)
	}

	fmt.Println("CIDR Address Range:", first, "-", last)
	fmt.Println("Subnet Mask:", mask)
	fmt.Println("Network ID:", first)
}
func calculateCIDRInfo(cidr string) (string, string, string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", "", "", err
	}

	// 计算网络号
	network := ip.Mask(ipnet.Mask)

	// 计算地址段
	first := network
	last := net.IP(make([]byte, len(network)))
	copy(last, network)
	for i := range last {
		last[i] |= ^ipnet.Mask[i]
	}

	// 转换掩码为字符串
	mask := net.IP(ipnet.Mask).String()

	return first.String(), last.String(), mask, nil
}
