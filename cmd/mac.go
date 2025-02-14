package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// macCmd represents the mac command
var macCmd = &cobra.Command{
	Use:   "mac",
	Short: "Convert mac address",
	Long:  `Convert mac address`,
	Run:   getMacAddress,
}

func init() {
	rootCmd.AddCommand(macCmd)
}

func isValidMACAddress(mac string) bool {
	if len(mac) != 12 {
		return false
	}
	pattern := `^[0-9a-f]{12}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(mac)
}

func getMacAddress(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("missing arguments")
		return
	}
	origin := args[0]
	// Remove all separators and convert to lowercase
	macAddress := strings.ReplaceAll(origin, "-", "")
	macAddress = strings.ReplaceAll(macAddress, ".", "")
	macAddress = strings.ReplaceAll(macAddress, ":", "")
	macAddress = strings.ToLower(macAddress)

	if !isValidMACAddress(macAddress) {
		fmt.Println("invalid mac address")
		return
	}

	fmt.Println(macAddress)
	fmt.Println(convertMacAddress(macAddress, 2, ":"))
	fmt.Println(convertMacAddress(macAddress, 4, "."))
	fmt.Println(convertMacAddress(macAddress, 4, "-"))
	
	fmt.Println(strings.ToUpper(macAddress))
	fmt.Println(strings.ToUpper(convertMacAddress(macAddress, 2, ":")))
	fmt.Println(strings.ToUpper(convertMacAddress(macAddress, 4, ".")))
	fmt.Println(strings.ToUpper(convertMacAddress(macAddress, 4, "-")))
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
