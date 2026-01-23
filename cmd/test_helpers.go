//go:build unit

/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package cmd

import (
	"os"
)

// captureOutput is a no-op placeholder for tests
func captureOutput(buf interface{}) *os.File {
	return os.Stdout
}

// restoreOutput is a no-op placeholder for tests
func restoreOutput(oldStdout *os.File) {
	// No-op
}
