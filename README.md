# Welcome To Envish

__This is the documentation for Envish v2. You can find the documentation for v1 in [README.v1.md](README.v1.md).

## Introduction

Envish is a Golang library. It helps you emulate UNIX-like program environments in Golang packages.

It is released under the 3-clause New BSD license. See [LICENSE.md](LICENSE.md) for details.

```go
import envish "github.com/ganbarodigital/go_envish"

env := envish.NewEnv()

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
- [Package Docs](#package-docs)
  - [NewEnv()](#newenv)
  - [NewEnv() With Functional Options](#newenv-with-functional-options)
  - [Clearenv()](#clearenv)
  - [Environ()](#environ)
  - [Expand()](#expand)
  - [Getenv()](#getenv)
  - [Length()](#length)
  - [LookupEnv()](#lookupenv)
  - [Setenv()](#setenv)
  - [Unsetenv()](#unsetenv)

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
import envish "github.com/ganbarodigital/go_envish/v2"
```

__Don't forget that `v2` on the end of the import, or you'll get an older version of this package!__

Create a copy of your program's environment:

```golang
localEnv := envish.NewEnv(envish.CopyProgramEnv)
```

or simply start with an empty environment store:

```golang
localVars := envish.NewEnv()
```

Get and set variables in the environment store as needed:

```golang
home := localEnv.Getenv()
localEnv.Setenv("DEBIAN_FRONTEND", "noninteractive")
```

## UNIX Shell Special Variables

UNIX shells build and maintain a group of special variables, such as `$#`, `$*`, `$?` and so on. They're set automatically by the shell, and you cannot set them yourself.

In Envish, you can `Getenv()` and `Setenv()` them like any other variable:

```golang
localEnv.Setenv("$#", "5")
paramsCount := localEnv.Getenv("$#")
```

and `Expand()` will let you use them in strings:

```golang
fmt.Printf(localEnv.Expand("we have $# positional parameters"))
```

To make this work (and be readable!), we can't use the built-in `os.Expand()`. The built-in `os.Expand()` matches `$#` et al correctly, but doesn't send through the leading `$` symbol to the mapping function.

Future releases will add support for more UNIX shell substitution patterns.

## Package Docs

Envish provides an API that's compatible with Golang's standard `os` environment functions. The only difference is that they work on the key/value pairs stored in the environment store, rather than on your program's environment.

### NewEnv()

`NewEnv` creates an empty environment store:

```golang
localEnv := envish.NewEnv()
```

It returns a `envish.Env()` struct.  The struct does not export any public fields.

### NewEnv() With Functional Options

You can pass [functional options](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis) into `NewEnv()` to change the environment store before it is returned to you.

`CopyProgramEnv` is a functional option that will populate the environment store with a copy of your program's environment:

```golang
localEnv := envish.NewEnv(CopyProgramEnv)
```

### Clearenv()

`Clearenv()` deletes all entries from the given environment store. The program's environment remains unchanged.

```golang
// create a environment store
localEnv := envish.NewEnv(CopyProgramEnv)

// empty the environment store completely
localEnv.Clearenv()
```

### Environ()

`Environ()` returns a copy of all entries in the form `key=value`. This is compatible with any Golang standard library, such as `exec`.

```golang
// create an environment store
localEnv := envish.NewEnv()

// get a copy to pass into `exec`
cmd := exec.Command(...)
cmd.Env = localEnv.Environ()
```

### Expand()

`Expand()` will replace `${key}` and `$key` entries in a format string.

```golang
// create an environment store
localEnv := envish.NewEnv()

// show what we have
fmt.Printf(localEnv.Expand("HOME is ${HOME}\n"))
```

### Getenv()

`Getenv()` returns the value of the variable named by the key. If the key is not found, an empty string is returned.

If you want to find out if a key exists, use [`LookupEnv()`](#lookupenv) instead.

```golang
// create an environment store
localEnv := envish.NewEnv()

// get a variable from the environment store
home := localEnv.Getenv("HOME")
```

### Length()

`Length()` returns the number of key/value pairs stored in the environment store.

```golang
// create an environment store
localEnv := envish.NewEnv()

// find out how many variables it contains
fmt.Printf("environment has %d entries\n", localEnv.Length())
```

### LookupEnv()

`LookupEnv()` returns the value of the variable named by the key.

If the key is not found, an empty string is returned, and the returned boolean is false.

```golang
// create an environment store
localEnv := envish.NewEnv()

// find out if a key exists
value, ok := localEnv.LookupEnv("HOME")
```

### Setenv()

`Setenv()` sets the value of the variable named by the key. The program's environment remains unchanged.

```golang
// create an environment store
localEnv := envish.NewEnv()

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
localEnv := envish.NewEnv()

// delete an entry from the environment store
localEnv.Unsetenv("HOME")
```