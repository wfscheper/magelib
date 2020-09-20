// +build tools

package main

import (
	_ "github.com/GeertJohan/go.rice/rice"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "gotest.tools/gotestsum"
)
