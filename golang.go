package magelib

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	// CleanPaths are additional paths to delete when running the clean target.
	CleanPaths []string
	// ExeName sets the filename of a CLI executable.
	ExeName string
	// MainPackage sets the full package name of the main package
	// This is needed if main.go isn't at the root.
	MainPackage string
	// BuildDeps is a slice of targets that are dependencies for go:build
	BuildDeps []interface{}
	// GenerateDeps is a slice of targets that are dependencies for go:generate
	GenerateDeps []interface{}
	// LintDeps is a slice of targets that are dependencies for go:lint
	LintDeps []interface{}
	// ReleaseDeps is a slice of targets that are dependencies for go:release
	ReleaseDeps []interface{}
	// TestDeps is a slice of targets that are dependencies for go:test
	TestDeps []interface{}
	// TestTimeout sets the duration tests are allowed to run.
	TestTimeout = 15 * time.Second
	// GenerateRebuild is a function that returns whether go:generate needs to run or not.
	GenerateRebuild = func(_ context.Context) (bool, error) { return true, nil }

	goexe = "go"

	// tests
	coverageProfile = "coverage.out"

	// tools
	golangcilintPath = filepath.Join(toolsBinDir, "golangci-lint")
	goreleaserPath   = filepath.Join(toolsBinDir, "goreleaser")
	gotestsumPath    = filepath.Join(toolsBinDir, "gotestsum")

	// commands
	gobuild = sh.RunCmd(goexe, "build")
	govet   = sh.RunCmd(goexe, "vet")

	// default paths to clean
	cleanPaths = []string{"bin", "dist", toolsBinDir, coverageProfile}
)

type Go mg.Namespace

// Benchmark runs the benchmark suite
func (Go) Benchmark(ctx context.Context) error {
	return runTests(ctx, "-run=__absolutelynothing__", "-bench")
}

// Bulid runs go build
func (Go) Build(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, setup)

	Say("building")
	return gobuild("-v", "./...")
}

// Clean removes generated files
func (Go) Clean(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, setup)

	cleanPaths = append(cleanPaths, CleanPaths...)

	var err error
	Say("cleaning files")
	for _, path := range cleanPaths {
		err = sh.Rm(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	if err != nil {
		return errors.New("failed to clean some files")
	}

	return nil
}

// Coverage generates coverage reports
func (Go) Coverage(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, setup)
	mg.CtxDeps(ctx, TestDeps...)

	mode := envString("coverage_mode", "atomic")
	return runTests(ctx, "-cover", "-covermode", mode, "-coverprofile="+coverageProfile)
}

// Exec builds the main binary
func (g Go) Exec(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, setup)
	mg.CtxDeps(ctx, append(BuildDeps, g.Generate)...)

	if ExeName == "" {
		return errors.New("no executable name set")
	}

	commit, err := sh.Output("git", "rev-parse", "HEAD")
	if err != nil {
		return err
	}

	buildDate := time.Now().UTC()
	ldflags := fmt.Sprintf(
		"-X main.version=%s -X main.commit=%s -X main.buildDate=%s",
		ProjectVersion,
		commit,
		buildDate.Format(time.RFC3339),
	)

	exe := filepath.Join("bin", execName(ExeName))
	Say("building " + exe)
	return gobuild("-v", "-o", exe, "-ldflags", ldflags, MainPackage)
}

// Format runs golangci-lint in "fix" mode
func (Go) Format(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, setup)
	mg.CtxDeps(ctx, LintDeps...)

	Say("running golangci-lint fix")
	return sh.RunV(execName(golangcilintPath), "run", "--fix")
}

// Generate runs go generate
func (Go) Generate(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, setup)
	mg.CtxDeps(ctx, GenerateDeps...)

	rebuild, err := GenerateRebuild(ctx)
	if err == nil && rebuild {
		Say("running go generate")
		return sh.RunV(goexe, "generate", "-x", "./...")
	}
	return err
}

// Lint runs golangci-lint
func (Go) Lint(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, setup)
	mg.CtxDeps(ctx, LintDeps...)

	Say("running golangci-lint")
	return sh.RunV(execName(golangcilintPath), "run")
}

// Release runs goreleaser to create a release. Must set MAGELIB_DRY_RUN=false.
func (Go) Release(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, setup)
	mg.CtxDeps(ctx, ReleaseDeps...)

	exe := execName(goreleaserPath)
	if DryRun {
		// run goreleaser in snapshot mode
		Say("running goreleaser dry run")
		return sh.RunV(exe, "--snapshot", "--skip-publish", "--rm-dist")
	}

	// run for reals
	Say("running golreleaser")
	return sh.RunV(exe)
}

// Test runs the test suite
func (Go) Test(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, setup)
	mg.CtxDeps(ctx, TestDeps...)

	Say("running tests")
	return runTests(ctx)
}

// TestRace runs the test suite with race detection
func (Go) TestRace(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, setup)
	mg.CtxDeps(ctx, TestDeps...)

	Say("running race condition tests")
	return runTests(ctx, "-race")
}

// TestShort runs only tests marked as short
func (Go) TestShort(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, setup)
	mg.CtxDeps(ctx, TestDeps...)

	Say("running short tests")
	return runTests(ctx, "-short")
}

// Vet runs go vet
func (Go) Vet(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, setup)
	Say("running go vet")
	return govet("./...")
}

func runTests(_ context.Context, args ...string) error {
	args = append(args, "./...")
	args = append([]string{"--", fmt.Sprintf("-timeout=%s", TestTimeout)}, args...)

	return sh.RunV(execName(gotestsumPath), args...)
}
