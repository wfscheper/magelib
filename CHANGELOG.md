# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Changes for the next release can be found in the [".stentor.d" directory](./.stentor.d).

<!-- stentor output starts -->
## [v0.7.0] - 2023-08-19

### Added

- `magelib` now supports versioning tools via a module.

  This means consumers can take advantage of
  Go module checksums for security
  and use the standard Go tools for updating versions.
  [#21](https://github.com/wfscheper/magelib/issues/21)


[v0.7.0]: https://github.com/wfscheper/magelib/compare/v0.6.0...v0.7.0


----


## [v0.6.0] - 2021-09-26

### Deprecated

- go 1.15 is no longer supported,
  so `magelib` is dropping support for it.
  [#8](https://github.com/wfscheper/magelib/issues/8)
- Now that `go install` supports specifying a version,
  `magelib` no longer needs a tools module to track build to versions.
  [#8](https://github.com/wfscheper/magelib/issues/8)


[v0.6.0]: https://github.com/wfscheper/magelib/compare/v0.5.0...v0.6.0


----


## [v0.5.0] - 2021-03-02

### Added

- Adds a `CleanPaths` slice
  for adding additional files to clean.
  [#6](https://github.com/wfscheper/magelib/issues/6)
- Adds [gotagger](https://github.com/sassoftware/gotagger)
  as the built-in versioning tool.
  [#6](https://github.com/wfscheper/magelib/issues/6)
- Adds a `go:format` target
  to the golang package.

  This target runs `golangci-lint` in "fix" mode.
  [#7](https://github.com/wfscheper/magelib/issues/7)


[v0.5.0]: https://github.com/wfscheper/magelib/compare/v0.4.1...v0.5.0


----


## [v0.4.1] - 2021-01-30

### Fixed

- Fixed `changelog` target to print output.
  [#5](https://github.com/wfscheper/magelib/issues/5)


[v0.4.1]: https://github.com/wfscheper/magelib/compare/v0.4.0...v0.4.1


----


## [v0.4.0] - 2021-01-27

### Added

- Adds a `changelog` target
  that uses [stentor](https://github.com/wfscheper/stentor)
  to manage the changelog.
  [#4](https://github.com/wfscheper/magelib/issues/4)


[v0.4.0]: https://github.com/wfscheper/magelib/compare/v0.3.0...v0.4.0


----


## [v0.3.0] - 2020-11-27

### Changed

- The type of `ProjectTools` is now `ToolMap`,
  which is a type for a `map[string]ToolFunc`.
  [#1](https://github.com/wfscheper/magelib/issues/1)


### Added

- `magelib` now defines a type, `ToolMap` for the `ProjectTools` map.
  [#1](https://github.com/wfscheper/magelib/issues/1)
- Add `TestTimeout` to control how long tests can run.
  [#2](https://github.com/wfscheper/magelib/issues/2)


[v0.3.0]: https://github.com/wfscheper/magelib/compare/v0.2.0...v0.3.0


----


## [v0.2.0] - 2020-09-20

### Added

- Added `go:release` target to run [goreleaser](https://goreleaser.com/)

[v0.2.0]: https://github.com/wfscheper/magelib/compare/v0.1.0...v0.2.0


----

## v0.1.0 - 2020-09-19

Initial release
