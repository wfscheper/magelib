package magelib

import (
	"context"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

const (
	BinGolangciLint = "golangci-lint"
	BinGoreleaser   = "goreleaser"
	BinGotagger     = "gotagger"
	BinGotestsum    = "gotestsum"
	BinStentor      = "stentor"

	ModuleGolangciLint = "github.com/golangci/golangci-lint/cmd/golangci-lint"
	ModuleGoreleaser   = "github.com/goreleaser/goreleaser"
	ModuleGotagger     = "github.com/sassoftware/gotagger/cmd/gotagger"
	ModuleGotestsum    = "gotest.tools/gotestsum"
	ModuleStentor      = "github.com/wfscheper/stentor/cmd/stentor"
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

// GetGoTool returns a ToolFunc that uses `go install` to install a specific version of a module as name.
//
// If no `version` is specified,
// then just the module name is used.
// This is useful if you just want the latest version,
// or you want to manage the version via a go module.
func GetGoTool(module, name string, versions ...string) ToolFunc {
	if len(versions) > 0 {
		module += "@" + versions[0]
	}

	return func(ctx context.Context) error {
		rebuild, err := target.Glob(filepath.Join(toolsBinDir, name))
		if err == nil && rebuild {
			Say("building %s", module)
			return goInstall(ctx, module)
		}

		return err
	}
}

// GetGolangciLint returns a ToolFunc that uses `go get` to install a specific version of golangci-lint .
func GetGolangciLint(version ...string) ToolFunc {
	return GetGoTool(ModuleGolangciLint, BinGolangciLint, version...)
}

// GetGoreleaser returns a ToolFunc that uses `go get` to install a specific versino of goreleaser.
func GetGoreleaser(version ...string) ToolFunc {
	return GetGoTool(ModuleGoreleaser, BinGoreleaser, version...)
}

// GetGotagger returns a ToolFunc that uses `go get` to install a specific versin of gotagger.
func GetGotagger(version ...string) ToolFunc {
	return GetGoTool(ModuleGotagger, BinGotagger, version...)
}

// GetGotestsum returns a ToolFunc that uses `go get` to install a specific version of gotestsum.
func GetGotestsum(version ...string) ToolFunc {
	return GetGoTool(ModuleGotestsum, BinGotestsum, version...)
}

// GetStentor returns a ToolFunc that uses `go get` to install a specific version of stentor.
func GetStentor(version ...string) ToolFunc {
	return GetGoTool(ModuleStentor, BinStentor, version...)
}

func goInstall(_ context.Context, module string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := os.Chdir(filepath.Join(wd, "mage")); err != nil {
		return err
	}
	defer func() {
		_ = os.Chdir(wd)
	}()

	return sh.RunWith(map[string]string{"GOBIN": filepath.Join(wd, toolsBinDir)}, "go", "install", module)
}
