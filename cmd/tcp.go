/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/seancfoley/ipaddress-go/ipaddr"
	"github.com/spf13/cobra"
)

// tcpCmd represents the tcp command
var tcpCmd = &cobra.Command{
	Use:   "tcp",
	Short: "check host port",
	Long: `Check if the host port is open. For example:

macconv tcp 192.168.1.1 22`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Printf("missing arguments")
			return
		}
		checkPort(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(tcpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tcpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tcpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func checkPort(ipStr, portStr string) {

	addrString := ipaddr.NewIPAddressString(ipStr)
	addr := addrString.GetAddress()
	if addr == nil {
		fmt.Printf("invalid IP address")
		return
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Printf("invalid port")
		return
	}

	version := addrString.GetIPVersion()
	target := fmt.Sprintf("%s:%d", ipStr, port) // 默认使用IPv4格式
	if version == ipaddr.IPv6 {
		target = fmt.Sprintf("[%s]:%d", addr, port) // 如果是IPv6，则调整格式
	}
	count := 0
	for count < 5 {
		// 尝试连接到目标主机的指定端口
		conn, err := net.DialTimeout("tcp", target, 2*time.Second)
		if err != nil {
			fmt.Printf("Port %d on %s is close\n", port, ipStr)
		} else {
			fmt.Printf("Port %d on %s is open\n", port, ipStr)
			count++
			conn.Close()
		}
		time.Sleep(time.Second)

	}
}
