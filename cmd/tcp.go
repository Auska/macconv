/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package cmd

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// tcpCmd represents the tcp command
var tcpCmd = &cobra.Command{
	Use:   "tcp",
	Short: "Check host port",
	Long: `
Check if the host port is open. For example:

	macconv tcp 192.168.1.1 22`,
	Run: checkPort,
}

func init() {
	rootCmd.AddCommand(tcpCmd)
}
func checkPort(cmd *cobra.Command, args []string) {

	if len(args) != 2 {
		fmt.Printf("Error: missing arguments")
		cmd.Help()
		os.Exit(1)
	}

	ipStr := args[0]
	portStr := args[1]

	if net.ParseIP(ipStr) == nil {
		fmt.Printf("Error: invalid IP address.")
		cmd.Help()
		os.Exit(1)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Printf("Error: invalid port.")
		cmd.Help()
		os.Exit(1)
	}

	var target string
	if net.ParseIP(ipStr).To4() == nil {
		target = fmt.Sprintf("%s:%d", ipStr, port) // 默认使用IPv4格式
	} else {
		target = fmt.Sprintf("[%s]:%d", ipStr, port) // 如果是IPv6，则调整格式
	}
	count := 0
	for count < 5 {
		// 尝试连接到目标主机的指定端口
		now := time.Now()
		conn, err := net.DialTimeout("tcp", target, 2*time.Second)
		if err != nil {
			fmt.Printf("%v Port %d on %s is close\n", now.Format(time.RFC3339), port, ipStr)
		} else {
			fmt.Printf("%v Port %d on %s is open\n", now.Format(time.RFC3339), port, ipStr)
			count++
			conn.Close()
		}
		time.Sleep(time.Second)

	}
}
