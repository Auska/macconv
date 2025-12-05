/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package cmd

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"macconv/pkg/logger"
)

// tcpCmd represents the tcp command
var tcpCmd = &cobra.Command{
	Use:   "tcp",
	Short: "Check host port",
	Long: `
Check if the host port is open. Supports both IP addresses and domain names. For example:

	macconv tcp 192.168.1.1 22
	macconv tcp google.com 80

The command will continuously check the port until it becomes open.
Use Ctrl+C to stop the check.`,
	Run: checkPort,
}

func init() {
	rootCmd.AddCommand(tcpCmd)
}
func checkPort(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		logger.PrintValidationError("missing arguments: IP address or hostname and port required")
		if err := cmd.Help(); err != nil {
			logger.PrintErrorWithMessage("failed to show help", err)
		}
		return
	}

	host := args[0]
	portStr := args[1]

	// 解析主机名（可以是 IP 地址或域名）
	ips, lookupErr := net.LookupIP(host)
	if lookupErr != nil {
		logger.PrintValidationError(fmt.Sprintf("failed to resolve hostname: %s - %v", host, lookupErr))
		if helpErr := cmd.Help(); helpErr != nil {
			logger.PrintErrorWithMessage("failed to show help", helpErr)
		}
		return
	}

	if len(ips) == 0 {
		logger.PrintValidationError(fmt.Sprintf("no IP addresses found for hostname: %s", host))
		if err := cmd.Help(); err != nil {
			logger.PrintErrorWithMessage("failed to show help", err)
		}
		return
	}

	// 使用第一个解析的 IP 地址
	ip := ips[0]
	logger.Debug("Resolved hostname %s to IP %s", host, ip.String())

	port, portErr := strconv.Atoi(portStr)
	if portErr != nil || port < 1 || port > 65535 {
		logger.PrintValidationError(fmt.Sprintf("invalid port number: %s", portStr))
		if helpErr := cmd.Help(); helpErr != nil {
			logger.PrintErrorWithMessage("failed to show help", helpErr)
		}
		return
	}

	target := buildTargetAddress(ip, port)
	logger.Info("Starting continuous port check for %s (%s):%d (use Ctrl+C to stop)", host, ip.String(), port)

	attempt := 0
	consecutiveSuccess := 0
	requiredSuccess := 5

	for {
		attempt++
		isOpen := checkSingleConnection(target, host, port, attempt)

		if isOpen {
			consecutiveSuccess++
			if consecutiveSuccess < requiredSuccess {
				fmt.Printf(" (%d/%d consecutive checks)\n", consecutiveSuccess, requiredSuccess)
			} else {
				fmt.Printf(" (%d/%d consecutive checks) - CONFIRMED\n", consecutiveSuccess, requiredSuccess)
				logger.Info("Port %d on %s confirmed open after %d consecutive successful checks", port, host, requiredSuccess)
				break
			}
		} else {
			consecutiveSuccess = 0 // 重置连续成功计数
		}

		// 等待1秒后再次检查
		time.Sleep(time.Second)
	}
}

func buildTargetAddress(ip net.IP, port int) string {
	if ip.To4() == nil {
		// IPv6 地址需要用方括号包围
		return fmt.Sprintf("[%s]:%d", ip.String(), port)
	}
	// IPv4 地址
	return fmt.Sprintf("%s:%d", ip.String(), port)
}

func checkSingleConnection(target string, host string, port int, attempt int) bool {
	now := time.Now()
	conn, err := net.DialTimeout("tcp", target, 2*time.Second)

	if err != nil {
		logger.Debug("Connection failed to %s:%d (attempt %d): %v", host, port, attempt, err)
		if isHostname(host) {
			ip := extractIPFromTarget(target)
			fmt.Printf("%v [%d] Port %d on %s (%s) is closed\n", now.Format(time.RFC3339), attempt, port, host, ip)
		} else {
			fmt.Printf("%v [%d] Port %d on %s is closed\n", now.Format(time.RFC3339), attempt, port, host)
		}
		return false
	}

	defer func() {
		if err := conn.Close(); err != nil {
			logger.Debug("Error closing connection: %v", err)
		}
	}()
	logger.Debug("Connection successful to %s:%d (attempt %d)", host, port, attempt)
	if isHostname(host) {
		ip := extractIPFromTarget(target)
		fmt.Printf("%v [%d] Port %d on %s (%s) is OPEN ✓", now.Format(time.RFC3339), attempt, port, host, ip)
	} else {
		fmt.Printf("%v [%d] Port %d on %s is OPEN ✓", now.Format(time.RFC3339), attempt, port, host)
	}
	return true
}

// isHostname 判断输入是否是主机名（非IP地址）
func isHostname(host string) bool {
	return net.ParseIP(host) == nil
}

// extractIPFromTarget 从目标地址中提取IP部分
func extractIPFromTarget(target string) string {
	// 处理IPv6地址 [::1]:80 格式
	if strings.HasPrefix(target, "[") {
		if idx := strings.Index(target, "]"); idx != -1 {
			return target[:idx+1]
		}
	}
	// 处理IPv4地址 192.168.1.1:80 格式
	if idx := strings.LastIndex(target, ":"); idx != -1 {
		return target[:idx]
	}
	return target
}
