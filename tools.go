package magelib

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/target"
)

const (
	BinGolangciLint = "golangci-lint"
	BinGotestsum    = "gotestsum"
	BinGoreleaser   = "goreleaser"

	ModuleGolangciLint = "github.com/golangci/golangci-lint/cmd/golangci-lint"
	ModuleGotestsum    = "gotest.tools/gotestsum"
	ModuleGoreleaser   = "github.com/goreleaser/goreleaser"

	toolsDir  = "tools"
	toolsData = `// +build tools

package main

import (%s
)
`
)

// ToolFunc is a function that installs a tool.
type ToolFunc func(context.Context) error

type ToolMap map[string]ToolFunc

var (
	// ProjectTools maps tool names to a ToolFunc that installs the tool
	ProjectTools ToolMap
)

// Tools is the naespace for all build tool targets.
type Tools mg.Namespace

// Build gets and compiles project tools
func (Tools) Build(ctx context.Context) error {
	for _, fn := range ProjectTools {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	return nil
}

// Init initializes the tools sub-module
func (Tools) Init(ctx context.Context) error {
	if len(ProjectTools) == 0 {
		// no tools, so exit
		return nil
	}

	if rebuild, _ := target.Path(filepath.Join(toolsDir, "go.mod")); rebuild {
		if err := os.MkdirAll(toolsDir, 0755); err != nil {
			return err
		}

		cmd := exec.CommandContext(ctx, "go", "mod", "init", "tools")
		cmd.Dir = toolsDir
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	toolsGo := filepath.Join(toolsDir, "tools.go")
	var imports string
	for module := range ProjectTools {
		imports += "\n\t_ \"" + module + "\""
	}

	return ioutil.WriteFile(toolsGo, []byte(fmt.Sprintf(toolsData, imports)), 0644)
}

// GetGolangciLint returns a ToolFunc that uses `go get` to install a specific version of golangci-lint .
func GetGolangciLint(version string) ToolFunc {
	return GetGoTool(ModuleGolangciLint, BinGolangciLint, version)
}

// GetGotestsum returns a ToolFunc that uses `go get` to install a specific version of gotestsum.
func GetGotestsum(version string) ToolFunc {
	return GetGoTool(ModuleGotestsum, BinGotestsum, version)
}

// GetGoTool returns a ToolFunc that uses `go get` to install a specific version of a module as name.
func GetGoTool(module, name, version string) ToolFunc {
	return func(ctx context.Context) error {
		rebuild, err := target.Glob(filepath.Join(toolsBinDir, name), filepath.Join(toolsDir, "go.*"))
		if err == nil && rebuild {
			Say("building %s@%s", module, version)
			return goGet(ctx, module+"@"+version)
		}

		Say("%s@%s up-to-date", module, version)
		return err
	}
}

// GetGoreleaser returns a ToolFunc that uses `go get` to install a specific versino of goreleaser.
func GetGoreleaser(version string) ToolFunc {
	return GetGoTool(ModuleGoreleaser, BinGoreleaser, version)
}

func goGet(ctx context.Context, s string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, "go", "get", s)
	cmd.Dir = filepath.Join(wd, toolsDir)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOBIN="+filepath.Join(wd, toolsBinDir))
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
