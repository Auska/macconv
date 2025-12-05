/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package main

import (
	"flag"
	"fmt"
	"os"

	"macconv/cmd"
)

// These variables are set at build time
var (
	version   = "dev"
	buildDate = "unknown"
)

func main() {
	// Check for version flag
	if showVersion := flag.Bool("version", false, "Show version information"); *showVersion {
		fmt.Printf("macconv version %s\n", version)
		fmt.Printf("Built on: %s\n", buildDate)
		fmt.Printf("Author: LuoDan\n")
		fmt.Printf("Email: luodan0709@live.cn\n")
		os.Exit(0)
	}
	flag.Parse()

	// Set version information in cmd package
	cmd.SetVersionInfo(version, buildDate)
	
	// Execute the command
	cmd.Execute()
}
