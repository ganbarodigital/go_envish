# Welcome To Envish

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
  - [Clearenv()](#clearenv)
  - [Environ()](#environ)
  - [Getenv()](#getenv)
  - [Length()](#length)
  - [LookupEnv()](#lookupenv)
  - [Setenv()](#setenv)
  - [Unsetenv()](#unsetenv)

## Why Use Envish?

### Who Is Envish For?

We've built Envish for anyone who needs to emulate a UNIX-like environment in their own Golang packages.

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
import envish "github.com/ganbarodigital/go_envish"
```

Create a copy of your program's environment:

```golang
localEnv := NewEnv()
```

Get and set temporary environment variables as needed:

```golang
home := localEnv.Getenv()
localEnv.Setenv("DEBIAN_FRONTEND", "noninteractive")
```

## Package Docs

Envish provides an API that's compatible with Golang's standard `os` environment functions. The only difference is that they work on a copy of the program's environment.

### NewEnv()

`NewEnv` creates a copy of your process's current environment.

```golang
localEnv := envish.NewEnv()
```

It returns a `envish.Env()` struct.  The struct does not export any public fields.

### Clearenv()

`Clearenv()` deletes all entries from the given environment. The program's environment remains unchanged.

```golang
// create a temporary environment
localEnv := envish.NewEnv()

// empty the temporary environment completely
localEnv.Clearenv()
```

### Environ()

`Environ()` returns a copy of all entries in the form `key=value`. This is compatible with any Golang standard library, such as `exec`.

```golang
// create a temporary environment
localEnv := envish.NewEnv()

// get a copy to pass into `exec`
cmd := exec.Command(...)
cmd.Env = localEnv.Environ()
```

### Getenv()

`Getenv()` returns the value of the variable named by the key. If the key is not found, an empty string is returned.

If you want to find out if a key exists, use [`LookupEnv()`](#lookupenv) instead.

```golang
// create a temporary environment
localEnv := envish.NewEnv()

// get a variable from that temporary environment
home := localEnv.Getenv("HOME")
```

### Length()

`Length()` returns the number of key/value pairs stored in the Env.

```golang
// create a temporary environment
localEnv := envish.NewEnv()

// find out how many variables it contains
fmt.Printf("environment has %d entries\n", localEnv.Length())
```

### LookupEnv()

`LookupEnv()` returns the value of the variable named by the key.

If the key is not found, an empty string is returned, and the returned boolean is false.

```golang
// create a temporary environment
localEnv := envish.NewEnv()

// find out if a key exists
value, ok := localEnv.LookupEnv("HOME")
```

### Setenv()

`Setenv()` sets the value of the variable named by the key. The program's environment remains unchanged.

```golang
// create a temporary environment
localEnv := envish.NewEnv()

// set a value in the temporary environment
localEnv.Setenv("DEBIAN_FRONTEND", "noninteractive")
```

`Setenv()` will return a `envish.ErrEmptyKey` error if you pass in a key that is either empty, or contains only whitespace.

```golang
err := localEnv.Setenv("", "key-is-invalid")
```

Other errors may be added in future releases.

### Unsetenv()

`Unsetenv()` deletes the variable named by the key. The program's environment remains completely unchanged.

`Unsetenv()` does not return an error if the key is not found.

```golang
// create a temporary environment
localEnv := envish.NewEnv()

// delete an entry from the temporary environment
localEnv.Unsetenv("HOME")
```