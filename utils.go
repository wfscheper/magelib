package magelib

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Say prints the formatted string to stdout.
func Say(format string, args ...interface{}) {
	format = strings.TrimSpace(format)
	fmt.Printf("▶ "+format+"…\n", args...)
}

func execName(p string) string {
	if runtime.GOOS == "windows" {
		return p + ".exe"
	}

	return p
}

func version(ctx context.Context) (string, error) {
	mg.CtxDeps(ctx, VersionDeps...)

	var args []string
	if IgnoreModules {
		args = append(args, "-modules=false")
	}

	output, err := sh.Output(execName(gotaggerPath), args...)
	if err != nil {
		return "", err
	}

	version := strings.Split(output, "\n")[0]
	return version, nil
}
