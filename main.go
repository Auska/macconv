/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package main

import (
	"macconv/cmd"
)

// These variables are set at build time
var (
	version   = "dev"
	buildDate = "unknown"
)

func main() {
	// Set version information in cmd package
	cmd.SetVersionInfo(version, buildDate)

	// Execute the command
	cmd.Execute()
}
