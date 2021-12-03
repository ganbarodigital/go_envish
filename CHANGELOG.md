# CHANGELOG

## develop

## v4.0.1

Released Friday, 3rd December 2021.

### Fixes

* Resolve compatibility with 'go get'

## v4.0.0

Released Friday, 3rd December 2021.

### Backwards-Compatibility Breaks

* Removed `LookupHomeDir()` from the various interfaces and structs
  - removed `Expander.LookupHomeDir()`
  - removed `LocalEnv.LookupHomeDir()`
  - removed `OverlayEnv.LookupHomeDir()`
  - removed `ProgramEnv.LookupHomeDir()`
* NewOverlayEnv() now requires a `[]envish.Expander`
  - no longer supports rest parameters for adding environments

### New

* Added LookupHomeDir() as a standalone util
  - we may also publish this as a separate package one day
* Added OverlayEnv.Export()

### Docs

* Added `godoc` compatibility
* README is now generated from the `godoc` output

## v3.0.1

Released Wednesday, 30th October 2019.

### Fix

* `OverlayEnv.IsExporter()` added, for interface compatibility

## v3.0.0

Released Wednesday, 30th October 2019.

### Backwards-Compatibility Breaks

* `envish.Env` is now `envish.LocalEnv`
* `envish.NewEnv` is now `envish.NewLocalEnv`

### New

* Added `ProgramEnv` and its methods
* Added `NewProgramEnv()`
* Added `ErrEmptyOverlayEnv` error
* Added `OverlayEnv` and its methods
* Added `NewOverlayEnv()`

## v2.3.0

Released Tuesday, 29th October 2019.

### New

* `Env.Expand()` now uses ShellExpand, supports a lot more string expansion features than before
* `Env.IsExporter()` added
* `Env.LookupHomeDir()` added
* `Env.SetAsExporter()` option function added
* `Reader` interface added
* `Writer` interface added
* `ReaderWriter` interface added
* `Expander` interface added

### Deps

* Added `go_shellexpand` v0.1.0

## v2.2.0

Released Tuesday, 29th October 2019.

### New

* `Env.MatchVarNames()` added

## v2.1.0

Released Tuesday, 8th October 2019.

### New

* `Env.Expand()` now supports UNIX shell special variable names

## v2.0.1

Released Monday, 7th October 2019.

### New

* Added `ErrNilPointer` error

### Fixes

* Added nil pointer checks on all Env methods

## v2.0.0

Released Monday, 7th October 2019.

### B/C Breaks

We are making some changes to `go_envish` to make it easier to reuse.

* `NewEnv()` now returns an empty environment store by default

Other B/C breaks:

* Package name is now `envish` instead of `pipe`

### New

* Added NewEnv option: `envish.CopyProgramEnv`

### Fixes

* `Env.Unsetenv()` now works when trying to unset the first variable in the environment store.

## v1.1.0

Released Sunday, 29th September 2019.

### New

* Added `Expand()`

## v1.0.0

Released Sunday, 29th September 2019.

### New

* New: imported code from `go_pipe`
