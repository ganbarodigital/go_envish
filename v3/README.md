# Welcome To Envish

## Introduction

Envish is a Golang library. It helps you emulate UNIX-like program environments in Golang packages.

It is released under the 3-clause New BSD license. See [LICENSE.md](LICENSE.md) for details.

```go
import envish "github.com/ganbarodigital/go_envish/v3"

env := envish.NewLocalEnv()

// add to this temporary environment
// WITHOUT changing your program's environment
env.Setenv("EXAMPLE_KEY", "EXAMPLE VALUE")

// pass it into run a child process
cmd := exec.Command(...)
cmd.Env = env.Environ()
cmd.Start()
```

## Table of Contents <!-- omit in toc -->

- [Introduction](#introduction)
- [Why Use Envish?](#why-use-envish)
  - [Who Is Envish For?](#who-is-envish-for)
  - [Why A Separate Package?](#why-a-separate-package)
- [How Does It Work?](#how-does-it-work)
  - [Getting Started](#getting-started)
- [Interfaces](#interfaces)
  - [Reader](#reader)
  - [Writer](#writer)
  - [ReaderWriter](#readerwriter)
  - [Expander](#expander)
- [LocalEnv](#localenv)
  - [NewLocalEnv()](#newlocalenv)
  - [NewLocalEnv() With Functional Options](#newlocalenv-with-functional-options)
  - [LocalEnv.Clearenv()](#localenvclearenv)
  - [LocalEnv.Environ()](#localenvenviron)
  - [LocalEnv.Expand()](#localenvexpand)
  - [LocalEnv.Getenv()](#localenvgetenv)
  - [LocalEnv.IsExporter()](#localenvisexporter)
  - [LocalEnv.Length()](#localenvlength)
  - [LocalEnv.LookupEnv()](#localenvlookupenv)
  - [LocalEnv.LookupHomeDir()](#localenvlookuphomedir)
  - [LocalEnv.MatchVarNames()](#localenvmatchvarnames)
  - [LocalEnv.Setenv()](#localenvsetenv)
  - [LocalEnv.Unsetenv()](#localenvunsetenv)
- [ProgramEnv](#programenv)
  - [NewProgramEnv()](#newprogramenv)
  - [ProgramEnv.Clearenv()](#programenvclearenv)
  - [ProgramEnv.Environ()](#programenvenviron)
  - [ProgramEnv.Expand()](#programenvexpand)
  - [ProgramEnv.Getenv()](#programenvgetenv)
  - [ProgramEnv.IsExporter()](#programenvisexporter)
  - [ProgramEnv.LookupEnv()](#programenvlookupenv)
  - [ProgramEnv.LookupHomeDir()](#programenvlookuphomedir)
  - [ProgramEnv.MatchVarNames()](#programenvmatchvarnames)
  - [ProgramEnv.Setenv()](#programenvsetenv)
  - [ProgramEnv.Unsetenv()](#programenvunsetenv)

## Why Use Envish?

### Who Is Envish For?

We've built Envish for anyone who needs to emulate a UNIX-like environment in their own Golang packages. Or anyone who just needs a simple key/value store with a familiar API.

We're using it ourselves for our [Pipe](https://github.com/ganbarodigital/go_pipe) and [Scriptish](https://github.com/ganbarodigital/go_scriptish) packages.

### Why A Separate Package?

Golang's `os` package provides support for working with your program's environment. But what if you want to make temporary changes to that environment, just to pass environment variables into child processes?

This is a very common pattern used in UNIX shell script programming:

```bash
DEBIAN_FRONTEND=noninteractive apt-get install -y mysql
```

In the example above, the environment variable `DEBIAN_FRONTEND` is only set for the child process `apt-get`.

## How Does It Work?

### Getting Started

Import Envish into your Golang code:

```golang
import envish "github.com/ganbarodigital/go_envish/v3"
```

__Don't forget that `v3` on the end of the import, or you'll get an older version of this package!__

Create a copy of your program's environment:

```golang
localEnv := envish.NewLocalEnv(envish.CopyProgramEnv)
```

or simply start with an empty environment store:

```golang
localVars := envish.NewLocalEnv()
```

Get and set variables in the environment store as needed:

```golang
home := localEnv.Getenv()
localEnv.Setenv("DEBIAN_FRONTEND", "noninteractive")
```

## Interfaces

_Envish_ publishes three interfaces:

Interface                              | Description
---------------------------------------|--------------------------------------
[`envish.Reader`](#reader)             | For retrieving variables
[`envish.Writer`](#writer)             | For creating and updating variables
[`envish.ReaderWriter`](#readerwriter) | Full read/write support
[`envish.Expander`](#expander)         | For UNIX shell string expansion

### Reader

```golang
// Reader is the interface that wraps a basic, read-only
// variable backing store
type Reader interface {
	// Environ returns a copy of all entries in the form "key=value".
	Environ() []string

	// Getenv returns the value of the variable named by the key.
	//
	// If the key is not found, an empty string is returned.
	Getenv(key string) string

	// IsExporter returns true if this backing store holds variables that
	// should be exported to external programs
	IsExporter() bool

	// LookupEnv returns the value of the variable named by the key.
	//
	// If the key is not found, an empty string is returned, and the returned
	// boolean is false.
	LookupEnv(key string) (string, bool)

	// MatchVarNames returns a list of variable names that start with the
	// given prefix.
	//
	// This is very useful if you want to support `${PARAM:=word}` shell
	// expansion in your own code.
	MatchVarNames(prefix string) []string
}
```

### Writer

```golang
// Writer is the interface that wraps a basic, write-only
// variable backing store
type Writer interface {
	// Clearenv deletes all entries
	Clearenv()

	// Setenv sets the value of the variable named by the key.
	Setenv(key, value string) error

	// Unsetenv deletes the variable named by the key.
	Unsetenv(key string)
}
```

### ReaderWriter

```golang
// ReaderWriter is the interface that groups the basic Read and Write methods
// for a variable backing store
type ReaderWriter interface {
	Reader
	Writer
}
```

### Expander

```golang
// Expander is the interface that wraps a variable backing store that
// also supports string expansion
type Expander interface {
	Reader
	Writer

	// Expand replaces ${var} or $var in the input string.
	Expand(fmt string) string

	// LookupHomeDir retrieves the given user's home directory, or false if
	// that cannot be found
	LookupHomeDir(username string) (string, bool)
}
```

## LocalEnv

`LocalEnv` is what we call _a variable backing store_. It emulates an environment and environment variables. It never reads from, or writes to, your program's environment.

### NewLocalEnv()

```golang
func NewLocalEnv(options ...func(*LocalEnv)) *LocalEnv
```

`NewLocalEnv` creates an empty environment store:

```golang
localEnv := envish.NewLocalEnv()
```

It returns a `envish.Env()` struct.  The struct does not export any public fields.

### NewLocalEnv() With Functional Options

You can pass [functional options](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis) into `NewLocalEnv()` to change the environment store before it is returned to you.

`CopyProgramEnv` is a functional option that will populate the environment store with a copy of your program's environment:

```golang
localEnv := envish.NewLocalEnv(CopyProgramEnv)
```

`SetAsExporter` is a functional option that will tell the `OverlayEnv` to include your `LocalEnv`'s contents in any call to its [`Environ()`](#overlayenvenviron) function.

```golang
localEnv := envish.NewLocalEnv(SetAsExporter)
```

### LocalEnv.Clearenv()

```golang
func (e *LocalEnv) Clearenv()
```

`Clearenv()` deletes all entries from the given environment store. The program's environment remains unchanged.

```golang
// create a environment store
localEnv := envish.NewLocalEnv(CopyProgramEnv)

// empty the environment store completely
localEnv.Clearenv()
```

### LocalEnv.Environ()

```golang
func (e *LocalEnv) Environ() []string
```

`Environ()` returns a copy of all entries in the form `key=value`. This is compatible with any Golang standard library, such as `exec`.

```golang
// create an environment store
localEnv := envish.NewLocalEnv()

// get a copy to pass into `exec`
cmd := exec.Command(...)
cmd.Env = localEnv.Environ()
```

### LocalEnv.Expand()

```golang
func (e *LocalEnv) Expand(fmt string) string
```

`Expand()` will replace `${key}` and `$key` entries in a format string.

```golang
// create an environment store
localEnv := envish.NewLocalEnv()

// show what we have
fmt.Printf(localEnv.Expand("HOME is ${HOME}\n"))
```

`Expand()` uses the [ShellExpand package](https://github.com/ganbarodigital/go_shellexpand) to do the expansion. It supports the vast majority of UNIX shell string expansion operations.

### LocalEnv.Getenv()

```golang
func (e *LocalEnv) Getenv(key string) string
```

`Getenv()` returns the value of the variable named by the key. If the key is not found, an empty string is returned.

If you want to find out if a key exists, use [`LocalEnv.LookupEnv()`](#localenvlookupenv) instead.

```golang
// create an environment store
localEnv := envish.NewLocalEnv()

// get a variable from the environment store
home := localEnv.Getenv("HOME")
```

### LocalEnv.IsExporter()

```golang
func (e *LocalEnv) IsExporter() bool
```

`IsExporter()` returns true if the environment store's contents should be exported to external programs.

```golang
// create an environment store
localEnv := envish.NewLocalEnv()

// by default, the environment store is NOT an exporter
exporting := localEnv.IsExporter()
```

If you want to change this hint, use the option function `SetAsExporter`:

```golang
// create an environment store
localEnv := envish.NewLocalEnv(envish.SetAsExporter)

// this will now return TRUE
exporting := localEnv.IsExporter()
```

### LocalEnv.Length()

```golang
func (e *LocalEnv) Length() int
```

`Length()` returns the number of key/value pairs stored in the environment store.

```golang
// create an environment store
localEnv := envish.NewLocalEnv()

// find out how many variables it contains
fmt.Printf("environment has %d entries\n", localEnv.Length())
```

### LocalEnv.LookupEnv()

```golang
func (e *LocalEnv) LookupEnv(key string) (string, bool)
```

`LookupEnv()` returns the value of the variable named by the key.

If the key is not found, an empty string is returned, and the returned boolean is false.

```golang
// create an environment store
localEnv := envish.NewLocalEnv()

// find out if a key exists
value, ok := localEnv.LookupEnv("HOME")
```

### LocalEnv.LookupHomeDir()

```golang
func (e *LocalEnv) LookupHomeDir(username string) (string, bool)
```

`LookupHomeDir()` returns the full path to the given user's home directory, or `false` if it cannot be retrieved for any reason.

```golang
localEnv := envish.NewLocalEnv()
homeDir, ok := localEnv.LookupHomeDir("root")
```

If you pass an empty string into `LookupHomeDir()`, it will look up the current user's home directory.

```golang
localEnv := envish.NewLocalEnv()
homeDir, ok := localEnv.LookupHomeDir("")

// homeDir should be same as `os.UserHomeDir()`
```

### LocalEnv.MatchVarNames()

```golang
func (e *LocalEnv) MatchVarNames(prefix string) []string
```

`MatchVarNames()` returns a list of keys that begin with the given prefix.

```golang
// create an environment store
localEnv := envish.NewLocalEnv()

// find all variables that begin with 'ANSIBLE_'
keys := localEnv.MatchVarNames("ANSIBLE_")
```

### LocalEnv.Setenv()

```golang
func (e *LocalEnv) Setenv(key, value string) error
```

`Setenv()` sets the value of the variable named by the key. The program's environment remains unchanged.

```golang
// create an environment store
localEnv := envish.NewLocalEnv()

// set a value in the environment store
localEnv.Setenv("DEBIAN_FRONTEND", "noninteractive")
```

`Setenv()` will return a `envish.ErrEmptyKey` error if you pass in a key that is either empty, or contains only whitespace.

```golang
err := localEnv.Setenv("", "key-is-invalid")
```

`Setenv()` will return an `envish.ErrNilPointer` error if you call `Setenv()` with a nil pointer to the environment store:

```golang
var localEnv *Env = nil
err := localEnv.Setenv("valid-key", "valid-value")
```

Other errors may be added in future releases.

### LocalEnv.Unsetenv()

```golang
func (e *LocalEnv) Unsetenv(key string)
```

`Unsetenv()` deletes the variable named by the key. The program's environment remains completely unchanged.

`Unsetenv()` does not return an error if the key is not found.

```golang
// create an environment store
localEnv := envish.NewLocalEnv()

// delete an entry from the environment store
localEnv.Unsetenv("HOME")
```

## ProgramEnv

ProgramEnv gives you the same API as [LocalEnv](#localenv), only it works directly on your program's environment instead.

### NewProgramEnv()

```golang
func NewProgramEnv() *ProgramEnv
```

`NewProgramEnv` creates an empty environment store:

```golang
progEnv := envish.NewProgramEnv()
```

It returns a `envish.ProgramEnv()` struct.  The struct does not export any public fields.

### ProgramEnv.Clearenv()

```golang
func (e *ProgramEnv) Clearenv()
```

`Clearenv()` deletes all entries from your program's environment. Use with caution.

```golang
progEnv := envish.NewProgramEnv()
progEnv.Clearenv()
```

### ProgramEnv.Environ()

```golang
func (e *ProgramEnv) Environ() []string
```

`Environ()` returns a copy of all entries in the form `key=value`. This is compatible with any Golang standard library, such as `exec`.

```golang
progEnv := envish.NewProgramEnv()

// get a copy to pass into `exec`
cmd := exec.Command(...)
cmd.Env = progEnv.Environ()
```

### ProgramEnv.Expand()

```golang
func (e *ProgramEnv) Expand(fmt string) string
```

`Expand()` will replace `${key}` and `$key` entries in a format string.

```golang
progEnv := envish.NewProgramEnv()

// show what we have
fmt.Printf(progEnv.Expand("HOME is ${HOME}\n"))
```

`Expand()` uses the [ShellExpand package](https://github.com/ganbarodigital/go_shellexpand) to do the expansion. It supports the vast majority of UNIX shell string expansion operations.

### ProgramEnv.Getenv()

```golang
func (e *ProgramEnv) Getenv(key string) string
```

`Getenv()` returns the value of the variable named by the key. If the key is not found, an empty string is returned.

If you want to find out if a key exists, use [`LookupEnv()`](#lookupenv) instead.

```golang
progEnv := envish.NewProgramEnv()

// get a variable from your program's environment
home := progEnv.Getenv("HOME")
```

### ProgramEnv.IsExporter()

```golang
func (e *ProgramEnv) IsExporter() bool
```

`IsExporter()` always returns `true`.

```golang
progEnv := envish.NewProgramEnv()
// always TRUE
exporting := progEnv.IsExporter()
```

### ProgramEnv.LookupEnv()

```golang
func (e *ProgramEnv) LookupEnv(key string) (string, bool)
```

`LookupEnv()` returns the value of the variable named by the key.

If the key is not found, an empty string is returned, and the returned boolean is false.

```golang
progEnv := envish.NewProgramEnv()

// find out if a key exists
value, ok := progEnv.LookupEnv("HOME")
```

### ProgramEnv.LookupHomeDir()

```golang
func (e *ProgramEnv) LookupHomeDir(username string) (string, bool) {
```

`LookupHomeDir()` returns the full path to the given user's home directory, or `false` if it cannot be retrieved for any reason.

```golang
progEnv := envish.NewProgramEnv()
homeDir, ok := progEnv.LookupHomeDir("root")
```

If you pass an empty string into `LookupHomeDir()`, it will look up the current user's home directory.

```golang
progEnv := envish.NewProgramEnv()
homeDir, ok := progEnv.LookupHomeDir("")

// homeDir should be same as `os.UserHomeDir()`
```

### ProgramEnv.MatchVarNames()

```golang
func (e *ProgramEnv) MatchVarNames(prefix string) []string
```

`MatchVarNames()` returns a list of keys that begin with the given prefix.

```golang
progEnv := envish.NewProgramEnv()

// find all variables that begin with 'ANSIBLE_'
keys := progEnv.MatchVarNames("ANSIBLE_")
```

### ProgramEnv.Setenv()

```golang
func (e *ProgramEnv) Setenv(key, value string) error
```

`Setenv()` sets the value of the variable named by the key. This is published into your program's environment immediately.

```golang
env := envish.NewProgramEnv()

// create/update an environment variable in your program's environment
progEnv.Setenv("DEBIAN_FRONTEND", "noninteractive")
```

`Setenv` will return an error if something went wrong.

### ProgramEnv.Unsetenv()

```golang
func (e *ProgramEnv) Unsetenv(key string)
```

`Unsetenv()` deletes the variable from your program's environment, if it exists.

```golang
progEnv := envish.NewProgramEnv()

// delete an entry from your program's environment
progEnv.Unsetenv("HOME")
```