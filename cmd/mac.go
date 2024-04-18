/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// macCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// macCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func isValidMACAddress(mac string) bool {
	strLength := len("001122334455")
	compareStrLength := len(mac)

	if strLength != compareStrLength {
		return false
	}
	// 正则表达式模式用于验证 MAC 地址
	pattern := `[0-9A-Fa-f]{12}`

	// 编译正则表达式
	re := regexp.MustCompile(pattern)

	// 使用正则表达式进行匹配
	return re.MatchString(mac)
}
func getMacAddress(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Printf("missing arguments")
		return
	}
	origin := args[0]
	mac_adress := strings.Replace(strings.Replace(strings.Replace(origin, "-", "", -1), ".", "", -1), ":", "", -1)

	//fmt.Println("mac called", origin)
	//fmt.Println("mac called", mac_adress)

	if isValidMACAddress(mac_adress) {
		//fmt.Println("valid mac address")
		fmt.Println(convertMacAddress(mac_adress, 2, ":"))
		fmt.Println(convertMacAddress(mac_adress, 4, "."))
		fmt.Println(convertMacAddress(mac_adress, 4, "-"))
	} else {
		fmt.Println("invalid mac address")
	}
}

func convertMacAddress(mac_adress string, mac_step_length int, mac_step_str string) string {
	var mac_address_str string
	for i := 0; i < len(mac_adress); i += mac_step_length {
		mac_address_str += mac_adress[i:i+mac_step_length] + mac_step_str
	}
	return mac_address_str[:len(mac_address_str)-1]
}
