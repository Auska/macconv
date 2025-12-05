/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"macconv/pkg/logger"
)

// Version information
var (
	appVersion = "dev"
	appBuildDate = "unknown"
)

// SetVersionInfo sets the version information
func SetVersionInfo(version, buildDate string) {
	appVersion = version
	appBuildDate = buildDate
}

// rootCmd represents the base command when called without any subcommands
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
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// 全局日志级别标志
	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "Set log level (debug, info, warn, error)")
	
	// Cobra 也支持本地标志，这些标志只会在直接调用此操作时运行
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	
	// 设置日志级别
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
		logger.DefaultLogger.SetLevel(logger.INFO)
		logger.Warn("Unknown log level: %s, using info level", logLevel)
	}
	
	logger.Debug("Logger initialized with level: %s", logLevel)
}
