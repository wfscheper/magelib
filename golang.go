package magelib

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

const (
	// tests
	testDir = "tests"
)

var (
	ExeName     string
	MainPackage string

	// Configurable deps

	BuildDeps    []interface{}
	GenerateDeps []interface{}
	LintDeps     []interface{}
	ReleaseDeps  []interface{}
	TestDeps     []interface{}

	goexe = "go"

	// tests
	coverageDir     = filepath.Join(testDir, "coverage")
	coverageProfile = filepath.Join(coverageDir, "coverage.out")

	// tools
	toolsBinDir      = filepath.Join(toolsDir, "bin")
	golangcilintPath = filepath.Join(toolsBinDir, "golangci-lint")
	goreleaserPath   = filepath.Join(toolsBinDir, "goreleaser")
	gotestsumPath    = filepath.Join(toolsBinDir, "gotestsum")

	// commands
	gobuild = sh.RunCmd(goexe, "build")
	govet   = sh.RunCmd(goexe, "vet")

	// args
	gotestArgs = []string{"--", "-timeout=15s"}
)

func init() {
	// Force use of go modules
	os.Setenv("GO111MODULES", "on")
	if runtime.GOOS == "windows" {
		golangcilintPath += ".exe"
		goreleaserPath += ".exe"
		gotestsumPath += ".exe"
	}
}

type Go mg.Namespace

// Benchmark runs the benchmark suite
func (Go) Benchmark(ctx context.Context) error {
	return runTests("-run=__absolutelynothing__", "-bench")
}

// Bulid runs go build
func (Go) Build(ctx context.Context) error {
	Say("building")

	return gobuild("-v", "./...")
}

// Exec builds the main binary
func (g Go) Exec(ctx context.Context) error {
	mg.CtxDeps(ctx, append(BuildDeps, g.Generate)...)

	Say("building " + ExeName)
	version, err := sh.Output("git", "describe", "--tags", "--always", "--dirty", "--match=v*")
	if err != nil {
		return err
	}

	commit, err := sh.Output("git", "rev-parse", "HEAD")
	if err != nil {
		return err
	}

	buildDate := time.Now().UTC()
	ldflags := "-X main.version=" + version +
		" -X main.commit=" + commit +
		" -X main.buildDate=" + buildDate.Format(time.RFC3339)

	return gobuild("-v", "-o", filepath.Join("bin", ExeName), "-ldflags", ldflags, MainPackage)
}

// Coverage generates coverage reports
func (Go) Coverage(ctx context.Context) error {
	mg.CtxDeps(ctx, append(TestDeps, mkCoverageDir)...)

	mode := os.Getenv("COVERAGE_MODE")
	if mode == "" {
		mode = "atomic"
	}
	if err := runTests(
		"-cover",
		"-covermode",
		mode,
		"-coverprofile="+coverageProfile,
	); err != nil {
		return err
	}
	if err := sh.Run(
		goexe,
		"tool",
		"cover",
		"-html="+coverageProfile,
		"-o",
		filepath.Join(coverageDir, "index.html"),
	); err != nil {
		return err
	}
	return nil
}

// Clean removes generated files
func (Go) Clean(ctx context.Context) error {
	Say("cleaning files")

	var err error
	for _, path := range []string{"bin", "dist", testDir, toolsBinDir} {
		err = sh.Rm(path)
	}

	return err
}

// Generate runs go generate
func (Go) Generate(ctx context.Context) error {
	mg.CtxDeps(ctx, GenerateDeps...)

	rebuild, err := target.Dir(
		filepath.Join("internal", "templates", "rice-box.go"),
		filepath.Join("internal", "templates", "templates"),
	)
	if err == nil && rebuild {
		Say("running go generate")
		return sh.RunV(goexe, "generate", "-x", "./...")
	}
	return err
}

// Lint runs golangci-lint
func (Go) Lint(ctx context.Context) error {
	mg.CtxDeps(ctx, LintDeps...)
	Say("running pre-commit hooks")
	return sh.RunV("pre-commit", "run", "--all-files")
}

// Release runs goreleaser to create a release. Must set MAGELIB_DRY_RUN=false.
func (Go) Release(ctx context.Context) error {
	mg.CtxDeps(ctx, ReleaseDeps...)
	if dryRun, err := strconv.ParseBool(os.Getenv("MAGELIB_DRY_RUN")); err != nil || dryRun {
		// run goreleaser in snapshot mode
		Say("running goreleaser test")
		return sh.Run(goreleaserPath, "--snapshot", "--skip-publish", "--rm-dist")
	}

	// run for reals
	Say("running golreleaser")
	return sh.Run(goreleaserPath)
}

// Test runs the test suite
func (Go) Test(ctx context.Context) error {
	mg.CtxDeps(ctx, TestDeps...)
	Say("running tests")
	return runTests()
}

// TestRace runs the test suite with race detection
func (Go) TestRace(ctx context.Context) error {
	mg.CtxDeps(ctx, TestDeps...)
	Say("running race condition tests")
	return runTests("-race")
}

// TestShort runs only tests marked as short
func (Go) TestShort(ctx context.Context) error {
	mg.CtxDeps(ctx, TestDeps...)
	Say("running short tests")
	return runTests("-short")
}

// Vet runs go vet
func (Go) Vet(ctx context.Context) error {
	Say("running go vet")
	return govet("./...")
}

func mkCoverageDir(ctx context.Context) error {
	_, err := os.Stat(coverageDir)
	if os.IsNotExist(err) {
		return os.MkdirAll(coverageDir, 0755)
	}
	return err
}

func runTests(testType ...string) error {
	if update, err := strconv.ParseBool(os.Getenv("UPDATE_GOLDEN")); err == nil && update {
		testType = append(testType, "./cmd/stentor", "-update")
	} else {
		testType = append(testType, "./...")
	}
	testType = append(gotestArgs, testType...)
	return sh.RunV(gotestsumPath, testType...)
}
