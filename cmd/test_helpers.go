//go:build unit

/*
Copyright © 2024-2025 Auska <luodan0709@live.cn>

*/

package cmd

import (
	"bytes"
	"os"
)

func captureOutput(buf *bytes.Buffer) *os.File {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	go func() {
		_, _ = buf.ReadFrom(r)
	}()

	return oldStdout
}

func restoreOutput(oldStdout *os.File) {
	os.Stdout = oldStdout
}
