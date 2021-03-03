package magelib

import (
	"context"
	"os"
	"strconv"
	"strings"
)

var (
	// DryRun makes co-operating targets only report what they would do.
	// Defaults to true.
	// Can be set via the MAGELIB_DRY_RUN environment variable.
	DryRun bool

	// IgnoreModules tells the Version target to ignore the presence of go.mod files.
	// Set this to true
	// if your project isn't written in go.
	// Can be set in your magefile
	// or via the MAGELIB_IGNORE_MODULES environment variable.
	IgnoreModules bool

	// ProjectVersion is the version of the project.
	// This can be determined from your commit history via gotagger,
	// set in your magefile,
	// or set via the MAGELIB_VERSION environment variable.
	ProjectVersion string
)

const (
	envPrefix = "MAGELIB_"
)

func setup(ctx context.Context) error {
	DryRun = envBool("dry_run", true)
	IgnoreModules = envBool("ignore_modules", IgnoreModules)

	projectVersion := envString("version", ProjectVersion)
	if projectVersion == "" {
		var err error
		projectVersion, err = version(ctx)
		if err != nil {
			return err
		}
	}
	ProjectVersion = projectVersion

	return nil
}

func envBool(name string, def bool) bool {
	if b, err := strconv.ParseBool(envString(name, "")); err != nil {
		return b
	}

	return def
}

func envString(name string, def string) string {
	if v, ok := os.LookupEnv(envPrefix + strings.ToUpper(name)); ok {
		return v
	}

	return def
}
