package magelib

import (
	"fmt"
	"runtime"
	"strings"
)

// Say prints the formatted string to stdout.
func Say(format string, args ...interface{}) {
	format = strings.TrimSpace(format)
	fmt.Printf("▶ "+format+"…\n", args...)
}

func execPath(p string) string {
	if runtime.GOOS == "windows" {
		return p + ".exe"
	}

	return p
}
