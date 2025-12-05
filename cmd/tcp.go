/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package cmd

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"macconv/pkg/logger"
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
		logger.PrintValidationError("missing arguments: IP address and port required")
		cmd.Help()
		return
	}

	ipStr := args[0]
	portStr := args[1]
	
	ip := net.ParseIP(ipStr)
	if ip == nil {
		logger.PrintValidationError(fmt.Sprintf("invalid IP address: %s", ipStr))
		cmd.Help()
		return
	}

	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		logger.PrintValidationError(fmt.Sprintf("invalid port number: %s", portStr))
		cmd.Help()
		return
	}

	target := buildTargetAddress(ip, port)
	logger.Info("Starting port check for %s:%d", ipStr, port)
	
	successCount := 0
	maxAttempts := 5
	
	for i := 0; i < maxAttempts; i++ {
		isOpen := checkSingleConnection(target, ipStr, port)
		if isOpen {
			successCount++
		}
		
		// 在最后一次检查后不需要等待
		if i < maxAttempts-1 {
			time.Sleep(time.Second)
		}
	}
	
	logger.Info("Port check completed: %d/%d successful connections", successCount, maxAttempts)
}

func buildTargetAddress(ip net.IP, port int) string {
	if ip.To4() == nil {
		// IPv6 地址需要用方括号包围
		return fmt.Sprintf("[%s]:%d", ip.String(), port)
	}
	// IPv4 地址
	return fmt.Sprintf("%s:%d", ip.String(), port)
}

func checkSingleConnection(target, ipStr string, port int) bool {
	now := time.Now()
	conn, err := net.DialTimeout("tcp", target, 2*time.Second)
	
	if err != nil {
		logger.Debug("Connection failed to %s:%d: %v", ipStr, port, err)
		fmt.Printf("%v Port %d on %s is closed\n", now.Format(time.RFC3339), port, ipStr)
		return false
	}
	
	defer conn.Close()
	logger.Debug("Connection successful to %s:%d", ipStr, port)
	fmt.Printf("%v Port %d on %s is open\n", now.Format(time.RFC3339), port, ipStr)
	return true
}
