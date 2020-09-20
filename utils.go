package magelib

import (
	"fmt"
	"strings"
)

func Say(format string, args ...interface{}) {
	format = strings.TrimSpace(format)
	fmt.Printf("▶ "+format+"…\n", args...)
}
