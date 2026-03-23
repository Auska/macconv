/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

// Package cmd implements the command-line interface for the macconv application.
// It provides commands for MAC address conversion, CIDR operations, port checking, and DHCP configuration.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"macconv/pkg/logger"
)

var (
	appVersion   = "dev"
	appBuildDate = "unknown"
)

func SetVersionInfo(version, buildDate string) {
	appVersion = version
	appBuildDate = buildDate
}

var rootCmd = &cobra.Command{
	Use:   "macconv",
	Short: "Parse mac address",
	Long: `Used to convert mac addresses between different devices. 
For example:
	macconv mac 00:11:22:33:44:55
	macconv ip 192.168.1.1/24
	macconv tcp 192.168.1.1 22
	macconv dhcp 192.168.1.1
`,
	Run: func(cmd *cobra.Command, args []string) {
		if version, _ := cmd.PersistentFlags().GetBool("version"); version {
			printVersion()
			os.Exit(0)
		}
		if err := cmd.Help(); err != nil {
			logger.PrintErrorWithMessage("failed to show help", err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Show version information")
	rootCmd.PersistentFlags().StringP("log-level", "l", "warn", "Set log level (debug, info, warn, error)")
	cobra.OnInitialize(initLogger)
}

func initLogger() {
	logLevel, _ := rootCmd.PersistentFlags().GetString("log-level")

	switch logLevel {
	case "debug":
		logger.DefaultLogger.SetLevel(logger.DEBUG)
	case "info":
		logger.DefaultLogger.SetLevel(logger.INFO)
	case "warn":
		logger.DefaultLogger.SetLevel(logger.WARN)
	case "error":
		logger.DefaultLogger.SetLevel(logger.ERROR)
	default:
		logger.DefaultLogger.SetLevel(logger.WARN)
		logger.Warnf("Unknown log level: %s, using warn level", logLevel)
	}

	logger.Debugf("Logger initialized with level: %s", logLevel)
}

func printVersion() {
	fmt.Printf("macconv version %s\n", appVersion)
	fmt.Printf("Built on: %s\n", appBuildDate)
	fmt.Printf("Author: LuoDan\n")
	fmt.Printf("Email: luodan0709@live.cn\n")
}
