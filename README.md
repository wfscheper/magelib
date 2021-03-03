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

If your project produces a binary,
then set the `magelib.ExeName` variable
to the name of your executable.
Setting this is best done in an `init` function:

```golang
func init() {
    magelib.ExecName = "mycommand"
}
```


### Main package

If your main package is not at the root of your project,
then set the `magelib.MainPackage` variable
to the full path of your main package.
Setting this is best done in an `init` function:

```golang
func init() {
    magelib.MainPackage = "example.com/foo/cmd/foo"
}
```


### Dry run mode

Some `magelib` targets support a "dry run" mode,
where they only report what would be done.
This is best set via an environment variable:

```bash
MAGELIB_DRY_RUN=true mage
```


## Project Versioning

`magelib` uses [gotagger] to version projects
based on their commit history and tags.
This requires using [conventional commits].
If your project does not wish to adopt conventional commits,
or if you want direct control over your project's version,
then set the `magelib.ProjectVersion` variable.
Setting this is best done in an `init` function:

```golang
func init() {
    magelib.ProjectVersion = "1.2.3"
}
```

For situations where the full git history is not available,
you can also set the project version
via the `MAGELIB_VERSION` environment variable.

```bash
MAGELIB_VERSION=1.2.3
```


### Ignore go modules

For `magelib` projects that are not themselves written in golang,
you will want to disable strict module versioning rules
by setting the `magelib.IgnoreModules` variable.
Setting this is best done in an `init` function:

```golang
func init() {
    magelib.IgnoreModules = true
}
```


## Changelog Management

`magelib` uses [stentor] for changelog management.


## Tools

`magelib` itself uses several external tools:

- [golangci-lint] for linting
- [goreleaser] for creating github releases
- [gotagger] for project versioning
- [gotestsum] as a test driver
- [stentor] for changelog management

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
    magelib.ExeName = "example"
    magelib.MainPackage = "example.com/example/cmd/example"

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

// All runs format, lint, build, and test targets
func All(ctx context.Context) {
    mg.SerialCtxDeps(ctx, magelib.Go.Format, magelib.Go.Lint, magelib.Go.Build, magelib.Go.Test)
}
```

[golangci-lint]: https://golangci-lint.run/
[goreleaser]: https://goreleaser.com/intro/
[gotagger]: https://github.com/sassoftware/gotagger
[gotestsum]: https://github.com/gotestyourself/gotestsum#documentation
[stentor]: https://github.com/wfscheper/stentor
