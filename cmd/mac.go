/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"macconv/pkg/logger"
	"macconv/pkg/validator"
)

var macCmd = &cobra.Command{
	Use:   "mac",
	Short: "Convert mac address",
	Long: `
Convert mac address to different formats. For example:

	macconv mac 001122334455`,
	Run: getMacAddress,
}

func init() {
	rootCmd.AddCommand(macCmd)
}

func normalizeMACAddress(mac string) string {
	mac = strings.ReplaceAll(mac, "-", "")
	mac = strings.ReplaceAll(mac, ".", "")
	mac = strings.ReplaceAll(mac, ":", "")
	return strings.ToLower(mac)
}

func getMacAddress(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		logger.PrintValidationError("missing MAC address argument")
		if err := cmd.Help(); err != nil {
			logger.PrintErrorWithMessage("failed to show help", err)
		}
		return
	}

	origin := args[0]
	logger.Debugf("Processing MAC address: %s", origin)

	macAddress := normalizeMACAddress(origin)

	if err := validator.ValidateMACAddress(macAddress); err != nil {
		logger.PrintErrorWithMessage("invalid MAC address", err)
		if err := cmd.Help(); err != nil {
			logger.PrintErrorWithMessage("failed to show help", err)
		}
		return
	}

	colonFormat := convertMacAddress(macAddress, 2, ":")
	dotFormat := convertMacAddress(macAddress, 4, ".")
	dashFormat := convertMacAddress(macAddress, 4, "-")

	formats := []string{
		macAddress,
		colonFormat,
		dotFormat,
		dashFormat,
		strings.ToUpper(macAddress),
		strings.ToUpper(colonFormat),
		strings.ToUpper(dotFormat),
		strings.ToUpper(dashFormat),
	}

	for _, format := range formats {
		fmt.Println(format)
	}

	logger.Infof("Successfully processed MAC address: %s", origin)
}

func convertMacAddress(mac string, step int, sep string) string {
	var result strings.Builder
	result.Grow(len(mac) + (len(mac)/step - 1))
	for i := 0; i < len(mac); i += step {
		if i > 0 {
			result.WriteString(sep)
		}
		result.WriteString(mac[i : i+step])
	}
	return result.String()
}
