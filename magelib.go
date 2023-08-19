package magelib

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	// ChangelogDeps is a slice of targets that are dependencies for the changelog target
	ChangelogDeps []interface{}

	// VersionDeps is a slice of targets that are dependencies for the version target
	VersionDeps []interface{}

	toolsBinDir = "bin"

	// tools
	gotaggerPath = filepath.Join(toolsBinDir, "gotagger")
	stentorPath  = filepath.Join(toolsBinDir, "stentor")
)

// Changelog reports the changelog update for the current release
func Changelog(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, setup)
	mg.CtxDeps(ctx, ChangelogDeps...)

	output, err := sh.Output("git", "tag", "--list", "--sort", "version:refname", "v*")
	if err != nil {
		return err
	}

	tags := strings.Split(output, "\n")

	exe := execName(stentorPath)
	if DryRun {
		Say("printing changelog update")
		return sh.RunV(exe, ProjectVersion, tags[len(tags)-1])
	}

	Say("updating changelog")
	return sh.Run(exe, "-release", ProjectVersion, tags[len(tags)-1])
}

// Version reports the current project version
func Version(ctx context.Context) {
	mg.SerialCtxDeps(ctx, setup)
	fmt.Println(ProjectVersion)
}
