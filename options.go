package magelib

import (
	"os"
	"strconv"
	"strings"
)

var (
	// DryRun makes co-operating targets only report what they would do.
	// Defaults to true.
	// Can be set via the MAGELIB_DRY_RUN environment variable.
	DryRun bool
)

const (
	envPrefix = "MAGELIB_"
)

func setup() error {
	DryRun = envBool("dry_run", true)
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
