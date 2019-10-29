# Welcome To Envish

__This is the documentation for Envish v2. You can find the documentation for v1 in [README.v1.md](README.v1.md).

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
- [LocalEnv](#localenv)
  - [NewLocalEnv()](#newlocalenv)
  - [NewLocalEnv() With Functional Options](#newlocalenv-with-functional-options)
  - [Clearenv()](#clearenv)
  - [Environ()](#environ)
  - [Expand()](#expand)
  - [Getenv()](#getenv)
  - [IsExporter()](#isexporter)
  - [Length()](#length)
  - [LookupEnv()](#lookupenv)
  - [LookupHomeDir()](#lookuphomedir)
  - [MatchVarNames()](#matchvarnames)
  - [Setenv()](#setenv)
  - [Unsetenv()](#unsetenv)
- [ProgramEnv](#programenv)
  - [NewProgramEnv()](#newprogramenv)
  - [Clearenv()](#clearenv-1)
  - [Environ()](#environ-1)
  - [Expand()](#expand-1)
  - [Getenv()](#getenv-1)
  - [IsExporter()](#isexporter-1)
  - [Length()](#length-1)
  - [LookupEnv()](#lookupenv-1)
  - [LookupHomeDir()](#lookuphomedir-1)
  - [MatchVarNames()](#matchvarnames-1)
  - [Setenv()](#setenv-1)
  - [Unsetenv()](#unsetenv-1)

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

## LocalEnv

Envish provides an API that's compatible with Golang's standard `os` environment functions. The only difference is that they work on the key/value pairs stored in the environment store, rather than on your program's environment.

### NewLocalEnv()

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

### Clearenv()

`Clearenv()` deletes all entries from the given environment store. The program's environment remains unchanged.

```golang
// create a environment store
localEnv := envish.NewLocalEnv(CopyProgramEnv)

// empty the environment store completely
localEnv.Clearenv()
```

### Environ()

`Environ()` returns a copy of all entries in the form `key=value`. This is compatible with any Golang standard library, such as `exec`.

```golang
// create an environment store
localEnv := envish.NewLocalEnv()

// get a copy to pass into `exec`
cmd := exec.Command(...)
cmd.Env = localEnv.Environ()
```

### Expand()

`Expand()` will replace `${key}` and `$key` entries in a format string.

```golang
// create an environment store
localEnv := envish.NewLocalEnv()

// show what we have
fmt.Printf(localEnv.Expand("HOME is ${HOME}\n"))
```

`Expand()` uses the [ShellExpand package](https://github.com/ganbarodigital/go_shellexpand) to do the expansion. It supports the vast majority of UNIX shell string expansion operations.

### Getenv()

`Getenv()` returns the value of the variable named by the key. If the key is not found, an empty string is returned.

If you want to find out if a key exists, use [`LookupEnv()`](#lookupenv) instead.

```golang
// create an environment store
localEnv := envish.NewLocalEnv()

// get a variable from the environment store
home := localEnv.Getenv("HOME")
```

### IsExporter()

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

### Length()

`Length()` returns the number of key/value pairs stored in the environment store.

```golang
// create an environment store
localEnv := envish.NewLocalEnv()

// find out how many variables it contains
fmt.Printf("environment has %d entries\n", localEnv.Length())
```

### LookupEnv()

`LookupEnv()` returns the value of the variable named by the key.

If the key is not found, an empty string is returned, and the returned boolean is false.

```golang
// create an environment store
localEnv := envish.NewLocalEnv()

// find out if a key exists
value, ok := localEnv.LookupEnv("HOME")
```

### LookupHomeDir()

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

### MatchVarNames()

`MatchVarNames()` returns a list of keys that begin with the given prefix.

```golang
// create an environment store
localEnv := envish.NewLocalEnv()

// find all variables that begin with 'ANSIBLE_'
keys := localEnv.MatchVarNames("ANSIBLE_")
```

### Setenv()

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

### Unsetenv()

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

`NewProgramEnv` creates an empty environment store:

```golang
progEnv := envish.NewProgramEnv()
```

It returns a `envish.ProgramEnv()` struct.  The struct does not export any public fields.

### Clearenv()

`Clearenv()` deletes all entries from your program's environment. Use with caution.

```golang
progEnv := envish.NewProgramEnv()
progEnv.Clearenv()
```

### Environ()

`Environ()` returns a copy of all entries in the form `key=value`. This is compatible with any Golang standard library, such as `exec`.

```golang
progEnv := envish.NewProgramEnv()

// get a copy to pass into `exec`
cmd := exec.Command(...)
cmd.Env = progEnv.Environ()
```

### Expand()

`Expand()` will replace `${key}` and `$key` entries in a format string.

```golang
progEnv := envish.NewProgramEnv()

// show what we have
fmt.Printf(progEnv.Expand("HOME is ${HOME}\n"))
```

`Expand()` uses the [ShellExpand package](https://github.com/ganbarodigital/go_shellexpand) to do the expansion. It supports the vast majority of UNIX shell string expansion operations.

### Getenv()

`Getenv()` returns the value of the variable named by the key. If the key is not found, an empty string is returned.

If you want to find out if a key exists, use [`LookupEnv()`](#lookupenv) instead.

```golang
progEnv := envish.NewProgramEnv()

// get a variable from your program's environment
home := progEnv.Getenv("HOME")
```

### IsExporter()

`IsExporter()` always returns `true`.

```golang
progEnv := envish.NewProgramEnv()
// always TRUE
exporting := progEnv.IsExporter()
```

### Length()

`Length()` returns the number of key/value pairs stored in your program's environment.

```golang
env := envish.NewProgramEnv()

// find out how many variables it contains
fmt.Printf("environment has %d entries\n", env.Length())
```

### LookupEnv()

`LookupEnv()` returns the value of the variable named by the key.

If the key is not found, an empty string is returned, and the returned boolean is false.

```golang
progEnv := envish.NewProgramEnv()

// find out if a key exists
value, ok := progEnv.LookupEnv("HOME")
```

### LookupHomeDir()

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

### MatchVarNames()

`MatchVarNames()` returns a list of keys that begin with the given prefix.

```golang
progEnv := envish.NewProgramEnv()

// find all variables that begin with 'ANSIBLE_'
keys := progEnv.MatchVarNames("ANSIBLE_")
```

### Setenv()

`Setenv()` sets the value of the variable named by the key. This is published into your program's environment immediately.

```golang
env := envish.NewProgramEnv()

// create/update an environment variable in your program's environment
progEnv.Setenv("DEBIAN_FRONTEND", "noninteractive")
```

`Setenv` will return an error if something went wrong.

### Unsetenv()

`Unsetenv()` deletes the variable from your program's environment.

`Unsetenv()` does not return an error if the key is not found.

```golang
progEnv := envish.NewProgramEnv()

// delete an entry from your program's environment
progEnv.Unsetenv("HOME")
```