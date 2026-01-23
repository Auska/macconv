/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version.",
	Long:  `Used to display information such as version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version %s\n", appVersion)
		fmt.Printf("built on %s\n", appBuildDate)
		fmt.Println("author  LuoDan")
		fmt.Println("E-Mail  luodan0709@live.cn")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
