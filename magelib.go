package magelib

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	// ChangelogDeps is a slice of targets that are dependencies for the changelog target
	ChangelogDeps []interface{}

	// tools
	toolsBinDir = filepath.Join(toolsDir, "bin")
	stentorPath = filepath.Join(toolsBinDir, "stentor")
)

// Changelog reports the changelog update for the current release
func Changelog(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, setup)
	mg.CtxDeps(ctx, ChangelogDeps...)

	exe := stentorPath
	if runtime.GOOS == "windows" {
		exe += ".exe"
	}

	next := os.Getenv("MAGELIB_NEXT_VERSION")
	if next == "" {
		output, err := sh.Output("git", "describe", "--dirty", "--always", "--long")
		if err != nil {
			return err
		}
		next = strings.TrimSpace(output)
	}

	output, err := sh.Output("git", "tag", "--list", "--sort", "version:refname", "v*")
	if err != nil {
		return err
	}

	tags := strings.Split(output, "\n")

	if os.Getenv("MAGELIB_DRY_RUN") == "false" {
		Say("updating changelog")
		return sh.Run(exe, "-release", next, tags[len(tags)-1])
	}

	Say("printing changelog update")
	return sh.RunV(exe, next, tags[len(tags)-1])
}
