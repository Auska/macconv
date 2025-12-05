/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"macconv/pkg/errors"
	"macconv/pkg/logger"
)

// macCmd represents the mac command
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
	// Remove all separators and convert to lowercase
	mac = strings.ReplaceAll(mac, "-", "")
	mac = strings.ReplaceAll(mac, ".", "")
	mac = strings.ReplaceAll(mac, ":", "")
	return strings.ToLower(mac)
}

func validateMACAddress(mac string) error {
	if len(mac) != 12 {
		return errors.New(errors.ValidationError, "MAC address must be 12 characters after normalization")
	}

	pattern := `^[0-9a-f]{12}$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(mac) {
		return errors.New(errors.ValidationError, "MAC address contains invalid characters")
	}

	return nil
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
	logger.Debug("Processing MAC address: %s", origin)

	// Remove all separators and convert to lowercase
	macAddress := normalizeMACAddress(origin)

	if err := validateMACAddress(macAddress); err != nil {
		logger.PrintErrorWithMessage("invalid MAC address", err)
		if err := cmd.Help(); err != nil {
			logger.PrintErrorWithMessage("failed to show help", err)
		}
		return
	}

	// Output all MAC address formats
	formats := []string{
		macAddress,
		convertMacAddress(macAddress, 2, ":"),
		convertMacAddress(macAddress, 4, "."),
		convertMacAddress(macAddress, 4, "-"),
		strings.ToUpper(macAddress),
		strings.ToUpper(convertMacAddress(macAddress, 2, ":")),
		strings.ToUpper(convertMacAddress(macAddress, 4, ".")),
		strings.ToUpper(convertMacAddress(macAddress, 4, "-")),
	}

	for _, format := range formats {
		fmt.Println(format)
	}

	logger.Info("Successfully processed MAC address: %s", origin)
}

func convertMacAddress(mac string, step int, sep string) string {
	var result strings.Builder
	for i := 0; i < len(mac); i += step {
		if i > 0 {
			result.WriteString(sep)
		}
		result.WriteString(mac[i : i+step])
	}
	return result.String()
}
