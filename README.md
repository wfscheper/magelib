# Magelib

`magelib` is a collection of [mage](https://github.com/magefile/mage) targets.
Currently, `magelib` intended for go development.
Other languages may be added when needed.

## Usage

Follow the usual method for importing mage targets:

```golang
// mage:import
_ "github.com/wfscheper/magelib"
```

## Configuration

`magelib` has several configurable options:

### Executable name

Set the `magelib.ExecName` variable.
This is best done in an `init` function:

```golang
func init() {
    magelib.ExecName = "mycommand"
}
```

### Build Tools

`magelib` itself uses
[golangci-lint](https://golangci-lint.run/),
[gotestsum](https://github.com/gotestyourself/gotestsum#documentation),
and [goreleaser](https://goreleaser.com/intro/).

As a user of `magelib`,
you control the version of the built-in tools,
and you can add your own tools via the `ProjecTools` variable.

As a convenience,
`magelib` exposes tool function factories for the built-in tools,
as well as a generic go tool factory.

An example setting the version of
`golangci-lint`,
`gotestsum`,
and a custom tool [rice](https://github.com/GeertJohan/go.rice)

```golang
const moduleRice = "github.com/GeertJohan/go.rice/rice"

var getRice = magelib.GetGoTool(moduleRice, "rice", "v1.0.0")

func init() {
    // map module import path to tool build function
    magelib.ProjectTools = magelib.ToolMap{
        magelib.ModuleGolangciLint: magelib.GetGolangciLint("v1.26.0"),
        magelib.ModuleGotestsum:    magelib.GetGotestsum("v0.4.1"),
        moduleRice:                 getRice,
    }
}
```

### Custom target deps

You can add your own mage targets
as dependencies of the built-in targets
via the "\<target>Deps" slices.

An example of building `rice`
as a `go:generate` dependency.
This example uses the `getRice` tool function
defined in the pervious section.

```golang
func init() {
    magelib.GenerateDeps = []interface{}{
        func(ctx context.Context) error { return getRice(ctx) },
    }
}
```

Additionally,
the `go:generate` target has a "rebuild" function
that provides a hook for determining
if the target needs to run or not.

```golang
func init() {
    magelib.GenerateRebuild = func(ctx context.Context) (bool, error) {
        return target.Dir(
            filepath.Join("internal", "templates", "rice-box.go"),
            filepath.Join("internal", "templates", "templates"),
        )
    }
}
```

### Complete example

```golang
// +build mage

package main

import (
    "context"

    "github.com/magefile/mage/mg"

    // mage:import
    "github.com/wfscheper/magelib"
)

var (
    // Default mage target
    Default = All

    getGolangciLint = magelib.GetGolangciLint("v1.32.2")
    getGoreleaser   = magelib.GetGoreleaser("v0.146.0")
    getGotestsum    = magelib.GetGotestsum("v0.6.0")
)

func init() {
    magelib.LintDeps = []interface{}{
        func(ctx context.Context) error { return getGolangciLint(ctx) },
    }
    magelib.ReleaseDeps = []interface{}{
        func(ctx context.Context) error { return getGoreleaser(ctx) },
    }
    magelib.TestDeps = []interface{}{
        func(ctx context.Context) error { return getGotestsum(ctx) },
    }

    magelib.ProjectTools = magelib.ToolMap{
        magelib.ModuleGolangciLint: getGolangciLint,
        magelib.ModuleGoreleaser:   getGoreleaser,
        magelib.ModuleGotestsum:    getGotestsum,
    }
}

// All runs format, lint, vet, build, and test targets
func All(ctx context.Context) {
    mg.SerialCtxDeps(ctx, magelib.Go.Lint, magelib.Go.Build, magelib.Go.Test)
}
```
