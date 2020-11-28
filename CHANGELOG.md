# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Changes for the next release can be found in the [".stentor.d" directory](./.stentor.d).

<!-- stentor output starts -->

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
