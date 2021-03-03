// +build tools

package main

import (
	_ "github.com/goreleaser/goreleaser"
	_ "github.com/sassoftware/gotagger/cmd/gotagger"
	_ "gotest.tools/gotestsum"
	_ "github.com/wfscheper/stentor/cmd/stentor"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)
