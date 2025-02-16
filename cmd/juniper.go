/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// juniperCmd represents the juniper command
var juniperCmd = &cobra.Command{
	Use:   "juniper",
	Short: "Juniper subscribers.",
	Long: `
Juniper subscribers. For example:

	macconv juniper <file_path>`,
	Run: juniper_text,
}

func init() {
	rootCmd.AddCommand(juniperCmd)
}

func juniper_text(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("Error: missing arguments.")
		cmd.Help()
		os.Exit(1)
	}

	filePath := args[0]
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error: ", err)
		cmd.Help()
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if idx := strings.Index(line, "IP Address:"); idx != -1 {
			fmt.Print("\n")
			fmt.Print(strings.TrimSpace(line[idx+len("IP Address: "):]))
			fmt.Print("\t")
		}
		if idx := strings.Index(line, "MAC Address: "); idx != -1 {
			fmt.Print(strings.TrimSpace(line[idx+len("MAC Address: "):]))
			fmt.Print("\t")
		}
		if idx := strings.Index(line, "IPv4 Input Filter Name: "); idx != -1 {
			fmt.Print(strings.TrimSpace(line[idx+len("IPv4 Input Filter Name: "):]))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error: ", err)
		return
	}
}