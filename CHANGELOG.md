# CHANGELOG

## develop

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