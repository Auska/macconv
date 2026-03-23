/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package cmd

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"macconv/pkg/logger"
)

const (
	requiredConsecutiveSuccesses = 5
	dnsTimeout                   = 5 * time.Second
	connectionTimeout            = 2 * time.Second
	retryInterval                = 1 * time.Second
	maxCheckDuration             = 10 * time.Minute
)

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

	ctx, cancel := context.WithTimeout(context.Background(), maxCheckDuration)
	defer cancel()

	resolver := net.Resolver{}
	ips, lookupErr := resolver.LookupIPAddr(ctx, host)
	if lookupErr != nil {
		logger.PrintValidationError(fmt.Sprintf("failed to resolve hostname: %s - %v", host, lookupErr))
		if helpErr := cmd.Help(); helpErr != nil {
			logger.PrintErrorWithMessage("failed to show help", helpErr)
		}
		return
	}

	ipAddrs := make([]net.IP, len(ips))
	for i, addr := range ips {
		ipAddrs[i] = addr.IP
	}

	if len(ipAddrs) == 0 {
		logger.PrintValidationError(fmt.Sprintf("no IP addresses found for hostname: %s", host))
		if err := cmd.Help(); err != nil {
			logger.PrintErrorWithMessage("failed to show help", err)
		}
		return
	}

	ip := ipAddrs[0]
	logger.Debugf("Resolved hostname %s to IP %s", host, ip.String())

	port, portErr := parsePort(portStr)
	if portErr != nil || port < 1 || port > 65535 {
		logger.PrintValidationError(fmt.Sprintf("invalid port number: %s", portStr))
		if helpErr := cmd.Help(); helpErr != nil {
			logger.PrintErrorWithMessage("failed to show help", helpErr)
		}
		return
	}

	target := buildTargetAddress(ip, port)
	logger.Infof("Starting continuous port check for %s (%s):%d (use Ctrl+C to stop)", host, ip.String(), port)

	attempt := 0
	consecutiveSuccess := 0

	for {
		select {
		case <-ctx.Done():
			logger.Errorf("Port check timed out after %v", maxCheckDuration)
			return
		default:
			attempt++
			isOpen := checkSingleConnection(target, host, port, attempt)

			if isOpen {
				consecutiveSuccess++
				if consecutiveSuccess < requiredConsecutiveSuccesses {
					fmt.Printf(" (%d/%d consecutive checks)\n", consecutiveSuccess, requiredConsecutiveSuccesses)
				} else {
					fmt.Printf(" (%d/%d consecutive checks) - CONFIRMED\n", consecutiveSuccess, requiredConsecutiveSuccesses)
					logger.Infof("Port %d on %s confirmed open after %d consecutive successful checks", port, host, requiredConsecutiveSuccesses)
					return
				}
			} else {
				consecutiveSuccess = 0
			}

			time.Sleep(retryInterval)
		}
	}
}

func parsePort(portStr string) (int, error) {
	var port int
	for _, c := range portStr {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid character in port")
		}
		port = port*10 + int(c-'0')
		if port > 65535 {
			return 0, fmt.Errorf("port exceeds maximum value")
		}
	}
	return port, nil
}

func buildTargetAddress(ip net.IP, port int) string {
	if ip.To4() == nil {
		return fmt.Sprintf("[%s]:%d", ip.String(), port)
	}
	return fmt.Sprintf("%s:%d", ip.String(), port)
}

func checkSingleConnection(target string, host string, port, attempt int) bool {
	now := time.Now()
	conn, err := net.DialTimeout("tcp", target, connectionTimeout)

	if err != nil {
		logger.Debugf("Connection failed to %s:%d (attempt %d): %v", host, port, attempt, err)
		if isHostname(host) {
			ip := extractIPFromTarget(target)
			fmt.Printf("%v [%d] Port %d on %s (%s) is closed\n", now.Format(time.RFC3339), attempt, port, host, ip)
		} else {
			fmt.Printf("%v [%d] Port %d on %s is closed\n", now.Format(time.RFC3339), attempt, port, host)
		}
		return false
	}

	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			logger.Debugf("Error closing connection: %v", closeErr)
		}
	}()

	logger.Debugf("Connection successful to %s:%d (attempt %d)", host, port, attempt)
	if isHostname(host) {
		ip := extractIPFromTarget(target)
		fmt.Printf("%v [%d] Port %d on %s (%s) is OPEN ✓", now.Format(time.RFC3339), attempt, port, host, ip)
	} else {
		fmt.Printf("%v [%d] Port %d on %s is OPEN ✓", now.Format(time.RFC3339), attempt, port, host)
	}
	return true
}

func isHostname(host string) bool {
	return net.ParseIP(host) == nil
}

func extractIPFromTarget(target string) string {
	if strings.HasPrefix(target, "[") {
		if idx := strings.Index(target, "]"); idx != -1 {
			return target[:idx+1]
		}
	}
	if idx := strings.LastIndex(target, ":"); idx != -1 {
		return target[:idx]
	}
	return target
}
