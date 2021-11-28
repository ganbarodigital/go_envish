# Welcome To Envish

## Introduction

Envish is a Golang library. It helps you emulate UNIX-like program environments in Golang packages.

It is released under the 3-clause New BSD license. See [LICENSE.md](LICENSE.md) for details.

```go
import envish "github.com/ganbarodigital/go_envish/v4"

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
	- [ShellEnv](#shellenv)
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
	- [LocalEnv.ReplacePositionalParams()](#localenvreplacepositionalparams)
	- [LocalEnv.ResetPositionalParams()](#localenvresetpositionalparams)
	- [LocalEnv.SetPositionalParams()](#localenvsetpositionalparams)
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
- [OverlayEnv](#overlayenv)
	- [NewOverlayEnv()](#newoverlayenv)
	- [OverlayEnv.GetEnvByID()](#overlayenvgetenvbyid)
	- [OverlayEnv.Clearenv()](#overlayenvclearenv)
	- [OverlayEnv.Environ()](#overlayenvenviron)
	- [OverlayEnv.Expand()](#overlayenvexpand)
	- [OverlayEnv.Getenv()](#overlayenvgetenv)
	- [OverlayEnv.IsExporter()](#overlayenvisexporter)
	- [OverlayEnv.LookupEnv()](#overlayenvlookupenv)
	- [OverlayEnv.LookupHomeDir()](#overlayenvlookuphomedir)
	- [MatchVarNames](#matchvarnames)
	- [OverlayEnv.Setenv()](#overlayenvsetenv)
	- [OverlayEnv.Unsetenv()](#overlayenvunsetenv)

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
import envish "github.com/ganbarodigital/go_envish/v4"
```

__Don't forget that `v4` on the end of the import, or you'll get an older version of this package!__

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

_Envish_ publishes five interfaces:

Interface                              | Description
---------------------------------------|--------------------------------------
[`envish.ShellEnv`](#env)              | For emulating a UNIX shell environment
[`envish.Reader`](#reader)             | For retrieving variables
[`envish.Writer`](#writer)             | For creating and updating variables
[`envish.ReaderWriter`](#readerwriter) | Full read/write support
[`envish.Expander`](#expander)         | For UNIX shell string expansion

### ShellEnv

```golang
// ShellEnv is a list of operations needed by a UNIX shell, or an emulation
// such as Scriptish.
type ShellEnv interface {
	Expander
	ReaderWriter

	// GetPositionalParamCount returns the value of the UNIX shell special
	// parameter $#.
	//
	// If $# is not set, it returns 0.
	GetPositionalParamCount() int

	// ReplacePositionalParams sets $1, $2 etc etc to the given values.
	//
	// Any existing positional parameters are deleted.
	//
	// Use SetPositionalParams instead, if you want to preserve any of
	// the existing positional params.
	//
	// It also sets the special parameter $#. The value of $# is returned.
	ReplacePositionalParams(values ...string) int

	// ResetPositionalParams deletes $1, $2 etc etc from the environment.
	//
	// It also sets the special parameter $# to 0.
	ResetPositionalParams()

	// SetPositionalParams sets $1, $2 etc etc to the given values.
	//
	// Any existing positional parameters are overwritten, up to len(values).
	// For example, the positional parameter $10 is *NOT* overwritten if
	// you only pass in nine positional parameters.
	//
	// Use ReplacePositionalParams instead, if you want `values` to be the
	// only positional parameters set.
	//
	// It also updates the special parameter $# if needed. The (possibly new)
	// value of $# is returned.
	SetPositionalParams(values ...string) int
}
```

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

### LocalEnv.ReplacePositionalParams()

```golang
func (e *LocalEnv) ReplacePositionalParams(values ...string) int
```

`ReplacePositionalParams()` sets `$1`, `$2` etc etc to the given values.

Any existing positional parameters are deleted.

Use `SetPositionalParams()` instead, if you want to preserve any of the existing positional params.

It also sets the special parameter `$#`.

The value of `$#` is returned.

```golang
// create an environment store
localEnv := envish.NewLocalEnv()

// set the positional parameters
//
// NOTE that $0 is NOT a positional parameter
localEnv.SetPositionalParameters("go", "test", "-cover")

// has the value 3
hash := localEnv.Getenv("$#")

// replace the positional parameters
//
// posCount will have the value 2
posCount := localEnv.ReplacePositionalParameters("npm", "test")

// has the value 2
hash = localEnv.Getenv("$#")
```

### LocalEnv.ResetPositionalParams()

```golang
func (e *LocalEnv) ResetPositionalParams()
```

`ResetPositionalParams()` deletes `$1`, `$2` etc etc from the environment.

It also sets the special parameter `$#` to `0`.

### LocalEnv.SetPositionalParams()

```golang
func (e *LocalEnv) SetPositionalParams(values ...string) int
```

`SetPositionalParams()` sets `$1`, `$2` etc etc to the given values.

Any existing positional parameters are overwritten, up to `len(values)`. For example, the positional parameter $10 is *NOT* overwritten if you only pass in nine positional parameters.

Use `ReplacePositionalParams()` instead, if you want `values` to be the only positional parameters set.

It also updates the special parameter `$#` if needed. The (possibly new) value of `$#` is returned.

```golang
// create an environment store
localEnv := envish.NewLocalEnv()

// set the positional parameters
//
// NOTE that $0 is NOT a positional parameter
localEnv.SetPositionalParameters("go", "fish")

// has the value "go"
param1 := localEnv.Getenv("$1")

// has the value "fish"
param2 := localEnv.Getenv("$2")

// has the value "2"
param2 := localEnv.Getenv("$#")
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

## OverlayEnv

Use an `OverlayEnv` to combine one (or more) [`LocalEnv`](#localenv) and a single [`ProgramEnv`](#programenv) into a single logical environment.

We do this in [Scriptish](https://github.com/ganbarodigital/go_scriptish) to provide local variable support.

### NewOverlayEnv()

```golang
func NewOverlayEnv(envs ...Expander) *OverlayEnv
```

Call `NewOverlayEnv()` to build a single logical environment from a set of environments.

```golang
localVars := envish.NewLocalEnv()
progVars := envish.NewLocalEnv(SetAsExporter)
progEnv := envish.NewProgramEnv()

// combine them
env := envish.NewOverlayEnv(localVars, progVars, progEnv)
```

`NewOverlayEnv()` returns a pointer to an `OverlayEnv` struct. You can use the methods of the `OverlayStruct` to read from and write to each of these environments.

The order of the arguments to `NewOverlayEnv()` matters. `OverlayEnv`'s methods will read from / write to the underlying environments in the order you've put them in, from left to right.

### OverlayEnv.GetEnvByID()

```golang
func (e *OverlayEnv) GetEnvByID(id int) (Expander, bool)
```

`GetEnvByID()` returns the requested environment from the `OverlayEnv`. ID `0` is the first environment you passed into `NewOverlayEnv()`, ID `1` is the second environment, and so on.

If you request an ID that the `OverlayEnv` does not have, it returns `nil, false`.

This is handy if you don't want to keep separate references to the environments you've combined into a single overlay:

```golang
// the IDs of each environment in the overlay
const(
    LocalVars = iota
    ProgramVars
    ProgramEnv
)

// createEnvironment() is the kind of code that you'd normally call
// during your program's bootstrap sequence
func createEnvironment() envish.Expander {
    localVars := envish.NewLocalEnv()
    progVars := envish.NewLocalEnv(SetAsExporter)
    progEnv := envish.NewProgramEnv()

    // combine them
    return envish.NewOverlayEnv(localVars, progVars, progEnv)
}

// setLocalVar() sets a local variable
func setLocalVar(env envish.Expander, key, value string) error {
    // NOTE how we use the constant defined earlier to find the right
    // environment to set this variable in
    localVars := env.GetEnvByID(LocalVars)

    // this will now be available to read and update in the `env` elsewhere
    return localVars.Setenv(key, value)
}
```

### OverlayEnv.Clearenv()

```golang
func (e *OverlayEnv) Clearenv()
```

`Clearenv()` deletes all variables __in every environment in the overlay env__. If your overlay env includes a [`ProgramEnv`](#programenv), this _will_ delete all of your program's environment variables.

Use with caution!

### OverlayEnv.Environ()

```golang
func (e *OverlayEnv) Environ() []string
```

`Environ()` returns a copy of all of the variables in your `OverlayEnv` in the form `key=value`. This format is compatible with Golang's built-in packages.

When it builds the list, it follows these rules:

* it searches the environments in the order you provided them to [`NewOverlayEnv()`](#newoverlayenv)
* it only includes variables from environments where `IsExporter()` returns `true`
* if the same variable is set in multiple environments, it uses the first value it finds

```golang
localVars := envish.NewLocalEnv()
progVars := envish.NewLocalEnv(SetAsExporter)
progEnv := envish.NewProgramEnv()

// combine them
env := envish.NewOverlayEnv(localVars, progVars, progEnv)

// export their variables
//
// it will pick up variables from:
//
// - progVars (because it was called with the SetAsExporter functional option)
// - progEnv (because ProgramEnv.IsExporter() ALWAYS returns `true`)
//
// if a variable is set in both `progVars` and `progEnv`, it will use the
// value from `progVars`
environ := env.Environ()

// pass it into run a child process
cmd := exec.Command(...)
cmd.Env = environ
cmd.Start()
```

### OverlayEnv.Expand()

```golang
func (e *OverlayEnv) Expand(fmt string) string
```

`Expand()` will replace `${key}` and `$key` entries in a format string.

* it uses the overlay env's [`LookupEnv()`](#overlayenvlookupenv) to find the values of variables
* it uses the overlay env's [`Setenv()`](#overlayenvsetenv) to set the values of variables
* it uses the overlay env's [`MatchVarNames()`](#overlayenvmatchvarnames) to expand variable name prefixes
* it uses the overlay env's [`LookupHomeDir()`](#overlayenvlookuphomedir) to expand `~` (tilde)

Hopefully, we've got the logic right, and you'll find that your expansions just work the way you'd naturally expect.

```golang
localVars := envish.NewLocalEnv()
progVars := envish.NewLocalEnv(SetAsExporter)
progEnv := envish.NewProgramEnv()

// combine them
env := envish.NewOverlayEnv(localVars, progVars, progEnv)

// show what we have
fmt.Printf(env.Expand("HOME is ${HOME}\n"))
```

`Expand()` uses the [ShellExpand package](https://github.com/ganbarodigital/go_shellexpand) to do the expansion. It supports the vast majority of UNIX shell string expansion operations.

### OverlayEnv.Getenv()

```golang
func (e *OverlayEnv) Getenv(key string) string
```

`Getenv()` returns the value of the given variable. If the variable does not exist, it returns `""` (empty string).

* it searches the environments in the order you provided them to [`NewOverlayEnv()`](#newoverlayenv)
* if the same variable is set in multiple environments, it uses the first value it finds

```golang
localVars := envish.NewLocalEnv()
progVars := envish.NewLocalEnv(SetAsExporter)
progEnv := envish.NewProgramEnv()

// combine them
env := envish.NewOverlayEnv(localVars, progVars, progEnv)

// show what we have
fmt.Printf("HOME is %s\n", env.Getenv("HOME"))
```

### OverlayEnv.IsExporter()

```golang
func (e *OverlayEnv) IsExporter() bool
```

`IsExporter()` returns `true` if (and only if) any of the environments in the overlay env are exporters too.

```golang
localVars := envish.NewLocalEnv()
progVars := envish.NewLocalEnv(SetAsExporter)
progEnv := envish.NewProgramEnv()

// combine them
env := envish.NewOverlayEnv(localVars, progVars, progEnv)

// `true`, because `progVars.IsExporter()` and `progEnv.IsExporter()` will
// have returned `true`
exporting := env.IsExporter()
```

### OverlayEnv.LookupEnv()

```golang
func (e *OverlayEnv) LookupEnv(key string) (string, bool)
```

`LookupEnv()` returns the value of the given variable. If the variable does not exist, it returns `"", false`.

* it searches the environments in the order you provided them to [`NewOverlayEnv()`](#newoverlayenv)
* if the same variable is set in multiple environments, it uses the first value it finds

```golang
localVars := envish.NewLocalEnv()
progVars := envish.NewLocalEnv(SetAsExporter)
progEnv := envish.NewProgramEnv()

// combine them
env := envish.NewOverlayEnv(localVars, progVars, progEnv)

// `home` will hold whatever value we could find
// `ok` will be `true` if the variable exists in any of the environments
home, ok := env.LookupEnv("HOME")
```

### OverlayEnv.LookupHomeDir()

```golang
func (e *OverlayEnv) LookupHomeDir(username string) (string, bool)
```

`LookupHomeDir()` returns the full path to the given user's home directory, or `false` if it cannot be retrieved for any reason.

```golang
localVars := envish.NewLocalEnv()
progVars := envish.NewLocalEnv(SetAsExporter)
progEnv := envish.NewProgramEnv()

// combine them
env := envish.NewOverlayEnv(localVars, progVars, progEnv)

// find a user's home directory
homeDir, ok := env.LookupHomeDir("root")
```

If you pass an empty string into `LookupHomeDir()`, it will look up the current user's home directory.

```golang
localVars := envish.NewLocalEnv()
progVars := envish.NewLocalEnv(SetAsExporter)
progEnv := envish.NewProgramEnv()

// combine them
env := envish.NewOverlayEnv(localVars, progVars, progEnv)

// find the current user's home directory
homeDir, ok := env.LookupHomeDir("")

// homeDir should be same as `os.UserHomeDir()`
```

This method doesn't call any methods on any of the environments in your overlay env. We might revisit that in a future release.

### MatchVarNames

```golang
func (e *OverlayEnv) MatchVarNames(prefix string) []string
```

`MatchVarNames()` returns a list of variable names that start with the given prefix.

It's a feature needed for the `${!prefix*}` string expansion syntax.

```golang
localVars := envish.NewLocalEnv()
progVars := envish.NewLocalEnv(SetAsExporter)
progEnv := envish.NewProgramEnv()

// combine them
env := envish.NewOverlayEnv(localVars, progVars, progEnv)

// find all ANSIBLE variables
vars := env.MatchVarNames("ANSIBLE_")
```

### OverlayEnv.Setenv()

```golang
func (e *OverlayEnv) Setenv(key, value string) error
```

`Setenv()` creates a variable (if it doesn't already exist) or updates its value (if it does exist).

* it searches the environments in the order you provided them to [`NewOverlayEnv()`](#newoverlayenv)
* if the same variable is set in multiple environments, it updates the first variable it finds
* if the variable does not exist, it is always created in the first environment you passed into [`NewOverlayEnv()`](#newoverlayenv)

```golang
localVars := envish.NewLocalEnv()
progVars := envish.NewLocalEnv(SetAsExporter)
progEnv := envish.NewProgramEnv()

// combine them
env := envish.NewOverlayEnv(localVars, progVars, progEnv)

// some example data to show how Setenv() works
localVars.Setenv("LOCAL", "100")
progVars.Setenv("PROG", "200")
progEnv.Setenv("ENV", "300")

// this will update the variable "PROG" in the `progVars` environment,
// because it already exists
env.Setenv("PROG", "250")

// this will create the variable "NEWVAR" in the `localVars` environment,
// because:
//
// a) NEWVAR does not yet exist, and
// b) `localVars` was the first environment passed into `NewOverlayEnv()`
env.Setenv("NEWVAR", "101")
```

### OverlayEnv.Unsetenv()

```golang
func (e *OverlayEnv) Unsetenv(key string)
```

`Unsetenv()` will delete the given variable from __all environments__ in the overlay env.

```golang
localVars := envish.NewLocalEnv()
progVars := envish.NewLocalEnv(SetAsExporter)
progEnv := envish.NewProgramEnv()

// combine them
env := envish.NewOverlayEnv(localVars, progVars, progEnv)

// some example data to show how Unsetenv() works
localVars.Setenv("VAR", "100")
progVars.Setenv("VAR", "200")
progEnv.Setenv("VAR", "300")

// `value` is "100"
value := env.Getenv("VAR")

// delete it
env.Unsetenv("VAR")

// `ok1` is false; "VAR" has been deleted from here
_, ok1 := localVars.LookupEnv("VAR")

// `ok2` is also false; "VAR" has been deleted from here as well
_, ok2 := progVars.LookupEnv("VAR")

// `ok3` is also false; "VAR" has been deleted from here as well
_, ok3 := progEnv.LookupEnv("VAR")
```

Why do we delete the variable from all environments? It would be very confusing if [`Getenv()`](#overlayenvgetenv) et al continued to return values after you've tried to delete a variable.