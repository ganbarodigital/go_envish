# Welcome to envish!

envish is a library to help you emulate UNIX-like program environments
in Golang packages.

It is released under the 3-clause New BSD license. See LICENSE.md for details.

```go
import envish "github.com/ganbarodigital/go_envish"

env := envish.NewLocalEnv()

// add to this temporary environment
// WITHOUT changing your program's environment
env.Setenv("EXAMPLE_KEY", "EXAMPLE VALUE")

// pass it into run a child process
cmd := exec.Command(...)
cmd.Env = env.Environ()
cmd.Start()
```

- [Why Use Envish](#why-use-envish)
- [Why A Separate Package](#why-a-separate-package)
- [Getting Started](#getting-started)
- [Functions](#functions)
	- [func CopyProgramEnv](#func-copyprogramenv)
	- [func GetKeyFromPair](#func-getkeyfrompair)
	- [func GetValueFromPair](#func-getvaluefrompair)
	- [func LookupHomeDir](#func-lookuphomedir)
	- [func SetAsExporter](#func-setasexporter)
- [Types](#types)
	- [type ErrEmptyKey](#type-erremptykey)
	- [type ErrEmptyOverlayEnv](#type-erremptyoverlayenv)
	- [type ErrNilPointer](#type-errnilpointer)
	- [type Expander](#type-expander)
	- [type LocalEnv](#type-localenv)
	- [WithFunctionalOptions1](#withfunctionaloptions1)
	- [WithFunctionalOptions2](#withfunctionaloptions2)
	- [ErrEmptyKey](#erremptykey)
	- [ErrNilPointer](#errnilpointer)
	- [type OverlayEnv](#type-overlayenv)
	- [NoExporters](#noexporters)
	- [ProgramEnvIsAlwaysAnExporter](#programenvisalwaysanexporter)
	- [type ProgramEnv](#type-programenv)
	- [type Reader](#type-reader)
	- [type ReaderWriter](#type-readerwriter)
	- [type Writer](#type-writer)

## Why Use Envish

We've built Envish for anyone who needs to emulate a UNIX-like environment in
their own Golang packages. Or anyone who just needs a simple key/value store
with a familiar API.

We're using it ourselves for our [Pipe](https://github.com/ganbarodigital/go_pipe)
and [Scriptish](https://github.com/ganbarodigital/go_scriptish) packages.

## Why A Separate Package

Golang's `os` package provides support for working with your program's
environment. But what if you want to make temporary changes to that
environment, just to pass environment variables into child processes?

This is a very common pattern used in UNIX shell script programming:

```go
DEBIAN_FRONTEND=noninteractive apt-get install -y mysql
```

In the example above, the environment variable `DEBIAN_FRONTEND` is only set
for the child process `apt-get`.

## Getting Started

Import Envish into your Golang code:

```go
go get github.com/ganbarodigital/go_envish@v4
```

Create a copy of your program's environment:

```go
localEnv := envish.NewLocalEnv(envish.CopyProgramEnv)
```

or simply start with an empty environment store:

```go
localVars := envish.NewLocalEnv()
```

Get and set variables in the environment store as needed:

```go
home := localEnv.Getenv()
localEnv.Setenv("DEBIAN_FRONTEND", "noninteractive")
```

## Functions

### func [CopyProgramEnv](/option_copy_env.go#L44)

`func CopyProgramEnv(e *LocalEnv)`

CopyProgramEnv copies your program's environment into the given
environment store.

It replaces any existing variables in the environment store.

### func [GetKeyFromPair](/utils.go#L43)

`func GetKeyFromPair(pair string) string`

GetKeyFromPair returns the `KEY` from a string `KEY=VALUE`.

### func [GetValueFromPair](/utils.go#L50)

`func GetValueFromPair(pair string, key string) string`

GetValueFromPair returns the `VALUE` from a string `KEY=VALUE`.

### func [LookupHomeDir](/util_lookuphomedir.go#L46)

`func LookupHomeDir(username string) (string, bool)`

LookupHomeDir retrieves the given user's home directory, or false if
that cannot be found.

It does not use the value of $HOME at all.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// ask the operating system where your home directory is
	homedir, ok := envish.LookupHomeDir("")

	// print the results
	fmt.Printf("ok is %v", ok)
	fmt.Printf("your homedir is: %s", homedir)
}

```

### func [SetAsExporter](/localenv_options.go#L40)

`func SetAsExporter(e *LocalEnv)`

SetAsExporter sets a flag so that the EnvStack will include its contents
when building an environ to export to Golang's exec package.

## Types

### type [ErrEmptyKey](/errors.go#L42)

`type ErrEmptyKey struct{ ... }`

ErrEmptyKey is returned whenever we're given a key that is zero-length
or only contains whitespace

#### func (ErrEmptyKey) [Error](/errors.go#L44)

`func (e ErrEmptyKey) Error() string`

### type [ErrEmptyOverlayEnv](/errors.go#L49)

`type ErrEmptyOverlayEnv struct { ... }`

ErrEmptyOverlayEnv is returned whenever you call a method on an empty EnvStack

#### func (ErrEmptyOverlayEnv) [Error](/errors.go#L53)

`func (e ErrEmptyOverlayEnv) Error() string`

### type [ErrNilPointer](/errors.go#L59)

`type ErrNilPointer struct { ... }`

ErrNilPointer is returned whenever you call a method on the Env struct
with a nil pointer

#### func (ErrNilPointer) [Error](/errors.go#L63)

`func (e ErrNilPointer) Error() string`

### type [Expander](/expander.go#L40)

`type Expander interface { ... }`

Expander is the interface that wraps a variable backing store that
also supports string expansion.

### type [LocalEnv](/localenv.go#L43)

`type LocalEnv struct { ... }`

LocalEnv holds a list key/value pairs.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// create a local environment
	localEnv := envish.NewLocalEnv()

	// it starts as an empty environment
	fmt.Print(localEnv.Getenv("$USER"))
}

```

#### func [NewLocalEnv](/localenv.go#L71)

`func NewLocalEnv(options ...func(*LocalEnv)) *LocalEnv`

NewLocalEnv creates an empty environment store.

You can pass [functional options](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)
into NewLocalEnv to change the environment store before it is returned to you.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// create a local environment
	localEnv := envish.NewLocalEnv()

	// it starts as an empty environment
	fmt.Print(localEnv.Getenv("$USER"))
}

```

### WithFunctionalOptions1

CopyProgramEnv is a functional option that will populate the environment
store with a copy of your program's environment.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// create a local environment
	//
	// it will start with a copy of your program's environment
	localEnv := envish.NewLocalEnv(envish.CopyProgramEnv)

	// on UNIX-like systems, this will print the name of the user
	// who is running the program
	fmt.Print(localEnv.Getenv("$USER"))
}

```

### WithFunctionalOptions2

SetAsExporter is a functional option that will tell the OverlayEnv to
include your LocalEnv's contents in any call to the OverlayEnv's
Environ function.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// create a local environment
	//
	// if you add this to an OverlayEnv, the OverlayEnv will include
	// its contents when you call the OverlayEnv's Environ method.
	localEnv := envish.NewLocalEnv(envish.SetAsExporter)

	// this environment is now an exporter
	fmt.Print(localEnv.IsExporter())
}

```

 Output:

```

true
```

#### func (*LocalEnv) [Clearenv](/localenv.go#L186)

`func (e *LocalEnv) Clearenv()`

Clearenv deletes all entries from the given LocalEnv. The program's
environment remains unchanged.

```golang
package main

import (
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// create a environment store
	localEnv := envish.NewLocalEnv(envish.CopyProgramEnv)

	// empty the environment store completely
	localEnv.Clearenv()
}

```

#### func (*LocalEnv) [Environ](/localenv.go#L94)

`func (e *LocalEnv) Environ() []string`

Environ returns a copy of all entries in the form "key=value".
This is compatible with any Golang standard library, such as `os/exec`.

```golang
package main

import (
	envish "github.com/ganbarodigital/go_envish"
	"os/exec"
)

func main() {
	// create an environment store
	localEnv := envish.NewLocalEnv()

	// get a copy to pass into `os/exec`
	cmd := exec.Command("go", "doc")
	cmd.Env = localEnv.Environ()
}

```

#### func (*LocalEnv) [Expand](/localenv.go#L266)

`func (e *LocalEnv) Expand(fmt string) string`

Expand replaces ${var} or $var in the input string.

Internally, it uses [https://github.com/ganbarodigital/go_shellexpand](https://github.com/ganbarodigital/go_shellexpand)
to do the expansion.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// create an environment store
	localEnv := envish.NewLocalEnv()

	// show what we have
	fmt.Print(localEnv.Expand("USER is ${USER}\n"))
}

```

#### func (*LocalEnv) [Getenv](/localenv.go#L107)

`func (e *LocalEnv) Getenv(key string) string`

Getenv returns the value of the variable named by the key.

If the key is not found, an empty string is returned.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// create an environment store
	localEnv := envish.NewLocalEnv()

	// get a variable from the environment store
	user := localEnv.Getenv("USER")
	fmt.Print(user)
}

```

#### func (*LocalEnv) [IsExporter](/localenv.go#L129)

`func (e *LocalEnv) IsExporter() bool`

IsExporter returns true if this backing store holds variables that
should be exported to external programs.

It is used by OverlayEnv.Environ to work out which keys and values
the OverlayEnv should include in its output.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// create an environment store
	localEnv := envish.NewLocalEnv()

	// by default, the environment store is NOT an exporter
	//
	// if you want to change this hint, use envish.SetAsExporter
	exporting := localEnv.IsExporter()
	fmt.Print(exporting)
}

```

 Output:

```

false
```

#### func (*LocalEnv) [Length](/localenv.go#L285)

`func (e *LocalEnv) Length() int`

Length returns the number of key/value pairs stored in the LocalEnv.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// create an environment store
	localEnv := envish.NewLocalEnv()

	// find out how many variables it contains
	//
	// a new LocalEnv starts with no entries
	fmt.Printf("environment has %d entries\n", localEnv.Length())
}

```

 Output:

```
environment has 0 entries
```

#### func (*LocalEnv) [LookupEnv](/localenv.go#L137)

`func (e *LocalEnv) LookupEnv(key string) (string, bool)`

LookupEnv returns the value of the variable named by the key.

If the key is not found, an empty string is returned, and the returned
boolean is false.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// create an environment store
	localEnv := envish.NewLocalEnv()

	// find out if a key exists
	value, ok := localEnv.LookupEnv("USER")
	fmt.Printf("key exists: %v", ok)
	fmt.Printf("value of key: %s", value)
}

```

#### func (*LocalEnv) [LookupHomeDir](/localenv.go#L274)

`func (e *LocalEnv) LookupHomeDir(username string) (string, bool)`

LookupHomeDir returns the full path to the given user's home directory,
or false if that cannot be found.

Note: it does not use the value of $HOME at all.

#### func (*LocalEnv) [MatchVarNames](/localenv.go#L158)

`func (e *LocalEnv) MatchVarNames(prefix string) []string`

MatchVarNames returns a list of variable names that start with the
given prefix.

It's a feature needed for `${!prefix*}` string expansion syntax.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// create an environment store
	localEnv := envish.NewLocalEnv()

	// find all variables that begin with 'ANSIBLE_'
	for _, key := range localEnv.MatchVarNames("ANSIBLE_") {
		fmt.Printf("%s = %s", key, localEnv.Getenv(key))
	}
}

```

#### func (*LocalEnv) [Setenv](/localenv.go#L199)

`func (e *LocalEnv) Setenv(key, value string) error`

Setenv sets the value of the variable named by the key. The program's
environment remains unchanged.

```golang
package main

import (
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// create an environment store
	localEnv := envish.NewLocalEnv()

	// set a value in the environment store
	localEnv.Setenv("DEBIAN_FRONTEND", "noninteractive")
}

```

### ErrEmptyKey

Setenv will return a `envish.ErrEmptyKey` error if you pass in a key that
is either empty, or contains only whitespace.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// create an environment store
	localEnv := envish.NewLocalEnv()

	// try to create an invalid environment variable
	err := localEnv.Setenv("", "key-is-invalid")
	fmt.Print(err)
}

```

 Output:

```
zero-length key, or key only contains whitespace
```

### ErrNilPointer

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	var localEnv *envish.LocalEnv = nil
	err := localEnv.Setenv("valid-key", "valid-value")
	fmt.Print(err)
}

```

 Output:

```
nil pointer to environment store passed to LocalEnv.Setenv
```

#### func (*LocalEnv) [Unsetenv](/localenv.go#L225)

`func (e *LocalEnv) Unsetenv(key string)`

Unsetenv deletes the variable named by the key.

```golang
package main

import (
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// create an environment store
	localEnv := envish.NewLocalEnv()

	// delete an entry from the environment store
	localEnv.Unsetenv("$#")
}

```

### type [OverlayEnv](/overlayenv.go#L49)

`type OverlayEnv struct { ... }`

OverlayEnv works on a collection of variable backing stores.

Use an OverlayEnv to combine one (or more) LocalEnv and a single
ProgramEnv into a single logical environment.

We do this in [https://github.com/ganbarodigital/go_scriptish](https://github.com/ganbarodigital/go_scriptish) to
emulate local variable support.

#### func [NewOverlayEnv](/overlayenv.go#L69)

`func NewOverlayEnv(envs ...Expander) *OverlayEnv`

NewOverlayEnv builds a single logical environments from the
given set of underlying environments.

NewOverlayEnv returns a pointer to an OverlayEnv struct. You can use the
methods of the OverlayStruct to read from and write to each of these
environments.

The order of the arguments to NewOverlayEnv matters. The returned
OverlayEnv's methods will read from / write to the underlying environments
in the order you've given.

```golang
package main

import (
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// create individual environments
	localVars := envish.NewLocalEnv()
	progVars := envish.NewLocalEnv(envish.SetAsExporter)
	progEnv := envish.NewProgramEnv()

	// combine them
	env := envish.NewOverlayEnv(localVars, progVars, progEnv)

	// you can now treat them as a single environment
	env.Setenv("$1", "go")
}

```

#### func (*OverlayEnv) [Clearenv](/overlayenv.go#L256)

`func (e *OverlayEnv) Clearenv()`

Clearenv deletes all variables in every environment in the OverlayEnv.
If your overlay env includes a ProgramEnv, this *WILL* delete all of
your program's environment variables.

Use with extreme caution!

#### func (*OverlayEnv) [Environ](/overlayenv.go#L98)

`func (e *OverlayEnv) Environ() []string`

Environ() returns a copy of all of the variables in your `OverlayEnv`
in the form `key=value`. This format is compatible with Golang's
built-in packages.

When it builds the list, it follows these rules:

* it searches the environments in the order you provided them to
NewOverlayEnv

* it only includes variables from environments where the IsExporter
method returns `true`

* if the same variable is set in multiple environments, it uses the first
value it finds

```golang
package main

import (
	envish "github.com/ganbarodigital/go_envish"
	"os/exec"
)

func main() {
	// create our independent environments
	localVars := envish.NewLocalEnv()
	progVars := envish.NewLocalEnv(envish.SetAsExporter)
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
	cmd := exec.Command("godoc")
	cmd.Env = environ

	// you can now call cmd.Start()
}

```

#### func (*OverlayEnv) [Expand](/overlayenv.go#L340)

`func (e *OverlayEnv) Expand(fmt string) string`

Expand will replace `${key}` and `$key` entries in a format string,
by looking up values from the environments contained within the OverlayEnv.

* it uses the given OverlayEnv's LookupEnv to find the values of variables

* it uses the given OverlayEnv's Setenv to set the values of variables

* it uses the given OverlayEnv's MatchVarNames to expand variable name prefixes

* it uses the given OverlayEnv's LookupHomeDir to expand `~` (tilde)

Hopefully, we've got the logic right, and you'll find that your expansions
just work the way you'd naturally expect.

Internally, it uses [https://github.com/ganbarodigital/go_shellexpand](https://github.com/ganbarodigital/go_shellexpand) to
do the shell expansion. It supports the vast majority of UNIX shell
string expansion operations.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// create individual environments
	localVars := envish.NewLocalEnv()
	progVars := envish.NewLocalEnv(envish.SetAsExporter)
	progEnv := envish.NewProgramEnv()

	// combine them
	env := envish.NewOverlayEnv(localVars, progVars, progEnv)

	// use UNIX shell expansion to see what we have
	fmt.Print(env.Expand("USER is ${USER}\n"))
}

```

#### func (*OverlayEnv) [GetEnvByID](/overlayenv.go#L365)

`func (e *OverlayEnv) GetEnvByID(id int) (Expander, bool)`

GetEnvByID returns the requested environment from the given OverlayEnv.
ID `0` is the first environment you passed into NewOverlayEnv, ID `1`
is the second environment, and so on.

If you request an ID that the given OverlayEnv does not have, it returns
`nil, false`.

GetEnvByID is handy if you don't want to keep separate references to the
environments after you've combined them into the OverlayEnv.

```golang
// Envish is a library to help you emulate UNIX-like program environments
// in Golang packages
//
// Copyright 2019-present Ganbaro Digital Ltd
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
//   * Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer.
//
//   * Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer in
//     the documentation and/or other materials provided with the
//     distribution.
//
//   * Neither the names of the copyright holders nor the names of his
//     contributors may be used to endorse or promote products derived
//     from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
// FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE
// COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN
// ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package main

import (
	envish "github.com/ganbarodigital/go_envish"
)

const (
	LocalVars = iota
	ProgramVars
	ProgramEnv
)

func main() {
	// build an environment stack without keeping a reference
	// to any of the individual environments
	env := envish.NewOverlayEnv(
		envish.NewLocalEnv(),
		envish.NewLocalEnv(envish.SetAsExporter),
		envish.NewProgramEnv(),
	)

	// NOTE how we use the constant defined earlier to find the right
	// environment to set this variable in
	localVars, _ := env.GetEnvByID(LocalVars)

	// this will now be available to read and update in the `env` elsewhere
	localVars.Setenv("$#", "2")
}

```

#### func (*OverlayEnv) [GetTopMostEnv](/overlayenv.go#L384)

`func (e *OverlayEnv) GetTopMostEnv() (Expander, error)`

GetTopMostEnv returns the envish environment that's at the top of the
overlay stack.

If we don't have that environment, we return a suitable error.

```golang
package main

import (
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// build an environment stack without keeping a reference
	// to any of the individual environments
	env := envish.NewOverlayEnv(
		envish.NewLocalEnv(),
		envish.NewLocalEnv(envish.SetAsExporter),
		envish.NewProgramEnv(),
	)

	// now, imagine we want to set some local variables
	localVars, _ := env.GetTopMostEnv()
	localVars.Setenv("$#", "2")
}

```

#### func (*OverlayEnv) [Getenv](/overlayenv.go#L145)

`func (e *OverlayEnv) Getenv(key string) string`

Getenv returns the value of the given variable. If the variable does not
exist, it returns `""` (empty string).

* it searches the environments in the order you provided them to NewOverlayEnv

* if the same variable is set in multiple environments, it uses the first
value it finds

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// build an environment stack without keeping a reference
	// to any of the individual environments
	env := envish.NewOverlayEnv(
		envish.NewLocalEnv(),
		envish.NewLocalEnv(envish.SetAsExporter),
		envish.NewProgramEnv(),
	)

	// show what we have
	fmt.Printf("USER is %s`n", env.Getenv("USER"))
}

```

#### func (*OverlayEnv) [IsExporter](/overlayenv.go#L166)

`func (e *OverlayEnv) IsExporter() bool`

IsExporter returns `true` if (and only if) any of the environments in
the overlay env hold variables that should be exported to external
programs.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// build an environment stack without keeping a reference
	// to any of the individual environments
	env := envish.NewOverlayEnv(
		envish.NewLocalEnv(),
		envish.NewLocalEnv(envish.SetAsExporter),
		envish.NewProgramEnv(),
	)

	fmt.Print(env.IsExporter())
}

```

 Output:

```
true
```

### NoExporters

If none of the environments passed to NewOverlayEnv are exporters,
IsExporter will return `false`.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// build an environment stack without keeping a reference
	// to any of the individual environments
	env := envish.NewOverlayEnv(
		envish.NewLocalEnv(),
		envish.NewLocalEnv(),
	)

	fmt.Print(env.IsExporter())
}

```

 Output:

```
false
```

### ProgramEnvIsAlwaysAnExporter

If you have a ProgramEnv anywhere in your OverlayEnv, IsExporter
will always return `true`.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// build an environment stack without keeping a reference
	// to any of the individual environments
	env := envish.NewOverlayEnv(
		envish.NewLocalEnv(),
		envish.NewProgramEnv(),
	)

	fmt.Print(env.IsExporter())
}

```

 Output:

```
true
```

#### func (*OverlayEnv) [LookupEnv](/overlayenv.go#L190)

`func (e *OverlayEnv) LookupEnv(key string) (string, bool)`

LookupEnv returns the value of the given variable. If the variable does
not exist, it returns `"", false`.

* it searches the environments in the order you provided them to NewOverlayEnv

* if the same variable is set in multiple environments, it uses the first
value it finds

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// build an environment stack without keeping a reference
	// to any of the individual environments
	env := envish.NewOverlayEnv(
		envish.NewLocalEnv(),
		envish.NewLocalEnv(envish.SetAsExporter),
		envish.NewProgramEnv(),
	)

	home, ok := env.LookupEnv("HOME")
	fmt.Printf("did we find $HOME?: %v", ok)
	fmt.Printf("what value did we find?: %v", home)
}

```

#### func (*OverlayEnv) [LookupHomeDir](/overlayenv.go#L346)

`func (e *OverlayEnv) LookupHomeDir(username string) (string, bool)`

LookupHomeDir retrieves the given user's home directory, or false if
that cannot be found.

#### func (*OverlayEnv) [MatchVarNames](/overlayenv.go#L216)

`func (e *OverlayEnv) MatchVarNames(prefix string) []string`

MatchVarNames returns a list of variable names that start with the
given prefix.

It's a feature needed for `${!prefix*}` string expansion syntax.

* it searches the environments in the order you provided them to NewOverlayEnv

* if the same key is found in multiple environments, it only returns
the key once (ie, results are deduped before they are returned)

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// build an environment stack without keeping a reference
	// to any of the individual environments
	env := envish.NewOverlayEnv(
		envish.NewLocalEnv(),
		envish.NewLocalEnv(envish.SetAsExporter),
		envish.NewProgramEnv(),
	)

	// print out all the variables with the prefix ANSIBLE_
	for _, key := range env.MatchVarNames("ANSIBLE_") {
		fmt.Printf("%s = %s", key, env.Getenv(key))
	}
}

```

#### func (*OverlayEnv) [Setenv](/overlayenv.go#L280)

`func (e *OverlayEnv) Setenv(key, value string) error`

Setenv creates a variable (if it doesn't already exist) or updates its
value (if it does exist).

* it searches the environments in the order you provided to NewOverlayEnv

* if the same variable is set in multiple environments, it updates the
first variable it finds

* if the variable does not exist, it is always created in the first
environment you provided to NewOverlayEnv

```golang
package main

import (
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// create individual environments
	localVars := envish.NewLocalEnv()
	progVars := envish.NewLocalEnv(envish.SetAsExporter)
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
}

```

#### func (*OverlayEnv) [Unsetenv](/overlayenv.go#L306)

`func (e *OverlayEnv) Unsetenv(key string)`

Unsetenv deletes the variable named by the key.

It will be deleted from all the environments in the stack.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// create individual environments
	localVars := envish.NewLocalEnv()
	progVars := envish.NewLocalEnv(envish.SetAsExporter)
	progEnv := envish.NewProgramEnv()

	// combine them
	env := envish.NewOverlayEnv(localVars, progVars, progEnv)

	// some example data to show how Unsetenv() works
	localVars.Setenv("VAR", "100")
	progVars.Setenv("VAR", "200")
	progEnv.Setenv("VAR", "300")

	// `value` is "100"
	value := env.Getenv("VAR")
	fmt.Printf("value is: %s\n", value)

	// delete it
	env.Unsetenv("VAR")

	// `ok1` is false; "VAR" has been deleted from here
	_, ok1 := localVars.LookupEnv("VAR")
	fmt.Printf("ok1 is: %v\n", ok1)

	// `ok2` is also false; "VAR" has been deleted from here as well
	_, ok2 := progVars.LookupEnv("VAR")
	fmt.Printf("ok2 is: %v\n", ok2)

	// `ok3` is also false; "VAR" has been deleted from here as well
	_, ok3 := progEnv.LookupEnv("VAR")
	fmt.Printf("ok3 is: %v\n", ok3)

}

```

 Output:

```
value is: 100
ok1 is: false
ok2 is: false
ok3 is: false
```

### type [ProgramEnv](/programenv.go#L45)

`type ProgramEnv struct { ... }`

ProgramEnv puts helper wrapper functions around your program's
environment.

#### func [NewProgramEnv](/programenv.go#L56)

`func NewProgramEnv() *ProgramEnv`

NewProgramEnv returns an envish environment that works directly with
your program's environment.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	env := envish.NewProgramEnv()

	fmt.Printf("USER is %s", env.Getenv("USER"))
}

```

#### func (*ProgramEnv) [Clearenv](/programenv.go#L126)

`func (e *ProgramEnv) Clearenv()`

Clearenv deletes all entries from your program's environment.
Use with extreme caution!

#### func (*ProgramEnv) [Environ](/programenv.go#L71)

`func (e *ProgramEnv) Environ() []string`

Environ returns a copy of all entries in the form "key=value".
This format is compatible with Golang's built-in packages.

```golang
package main

import (
	envish "github.com/ganbarodigital/go_envish"
	"os/exec"
)

func main() {
	// get access to our program environment
	env := envish.NewProgramEnv()

	// get a list of all entries in the environment
	environ := env.Environ()

	// pass it into run a child process
	cmd := exec.Command("godoc")
	cmd.Env = environ

	// you can now call cmd.Start()
}

```

#### func (*ProgramEnv) [Expand](/programenv.go#L155)

`func (e *ProgramEnv) Expand(fmt string) string`

Expand replaces ${var} or $var in the input string, by looking up
values from your program's environment.

Internally, it uses [https://github.com/ganbarodigital/go_shellexpand](https://github.com/ganbarodigital/go_shellexpand)
to do the shell expansion. It supports the vast majority of UNIX shell
string expansion operations.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// get access to our program environment
	env := envish.NewProgramEnv()

	// show what we have
	fmt.Print(env.Expand("USER is ${USER}"))
}

```

#### func (*ProgramEnv) [Getenv](/programenv.go#L78)

`func (e *ProgramEnv) Getenv(key string) string`

Getenv returns the value of the variable named by the key.

If the key is not found, an empty string is returned.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// get access to our program environment
	env := envish.NewProgramEnv()

	// show what we have
	fmt.Printf("USER is %s", env.Getenv("USER"))
}

```

#### func (*ProgramEnv) [IsExporter](/programenv.go#L86)

`func (e *ProgramEnv) IsExporter() bool`

IsExporter always returns `true`.

It is used by OverlayEnv.Environ to work out which keys and values
the OverlayEnv should include in its output.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// get access to our program environment
	env := envish.NewProgramEnv()

	// a ProgramEnv is always an environment exporter
	fmt.Printf("env.IsExporter() is %v", env.IsExporter())

}

```

 Output:

```
env.IsExporter() is true
```

#### func (*ProgramEnv) [LookupEnv](/programenv.go#L94)

`func (e *ProgramEnv) LookupEnv(key string) (string, bool)`

LookupEnv returns the value of the variable named by the key.

If the key is not found, an empty string is returned, and the returned
boolean is false.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// get access to our program environment
	env := envish.NewProgramEnv()

	value, ok := env.LookupEnv("USER")

	fmt.Printf("ok is: %v\n", ok)
	fmt.Printf("value is: %v\n", value)
}

```

#### func (*ProgramEnv) [LookupHomeDir](/programenv.go#L161)

`func (e *ProgramEnv) LookupHomeDir(username string) (string, bool)`

LookupHomeDir retrieves the given user's home directory, or false if
that cannot be found.

#### func (*ProgramEnv) [MatchVarNames](/programenv.go#L102)

`func (e *ProgramEnv) MatchVarNames(prefix string) []string`

MatchVarNames returns a list of variable names that start with the
given prefix.

It's a feature needed for `${!prefix*}` string expansion syntax.

```golang
package main

import (
	"fmt"
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// get access to your program environment
	env := envish.NewProgramEnv()

	// print out all the variables with the prefix ANSIBLE_
	for _, key := range env.MatchVarNames("ANSIBLE_") {
		fmt.Printf("%s = %s", key, env.Getenv(key))
	}
}

```

#### func (*ProgramEnv) [RestoreEnvironment](/programenv.go#L178)

`func (e *ProgramEnv) RestoreEnvironment(pairs []string)`

RestoreEnvironment writes the given "key=value" pairs into your
program's environment.

It does *not* empty your program's environment first!

It was originally added so that our unit tests could put the 'go test'
program environment back in place after each test had run.

```golang
package main

import (
	envish "github.com/ganbarodigital/go_envish"
)

func main() {
	// get access to your program environment
	env := envish.NewProgramEnv()

	// take a backup of the whole environment
	backup := env.Environ()

	// nuke the environment
	env.Clearenv()

	// restore from backup
	env.RestoreEnvironment(backup)
}

```

#### func (*ProgramEnv) [Setenv](/programenv.go#L131)

`func (e *ProgramEnv) Setenv(key, value string) error`

Setenv sets the value of the variable named by the key.

```golang
package main

import (
	envish "github.com/ganbarodigital/go_envish"
	"os/exec"
)

func main() {
	// get access to your program environment
	env := envish.NewProgramEnv()

	// set a new environment variable
	//
	// this will be picked up the next time you call a child process
	env.Setenv("DEBIAN_FRONTEND", "noninteractive")

	// pass it into run a child process
	cmd := exec.Command("apt-get", "install", "mysql-server")
	cmd.Env = env.Environ()

	// you can now call cmd.Start()
}

```

#### func (*ProgramEnv) [Unsetenv](/programenv.go#L139)

`func (e *ProgramEnv) Unsetenv(key string)`

Unsetenv deletes the variable named by the key.

This will remove the given variable from your program's environment.
Use with caution!

```golang
package main

import (
	envish "github.com/ganbarodigital/go_envish"
	"os/exec"
)

func main() {
	// get access to your program environment
	env := envish.NewProgramEnv()

	// set a new environment variable
	//
	// this will be picked up the next time you call a child process
	env.Setenv("DEBIAN_FRONTEND", "noninteractive")

	// pass it into run a child process
	cmd := exec.Command("apt-get", "install", "mysql-server")
	cmd.Env = env.Environ()

	// you can now call cmd.Start()

	// afterwards, undo what you originally set
	env.Unsetenv("DEBIAN_FRONTEND")
}

```

### type [Reader](/reader.go#L40)

`type Reader interface { ... }`

Reader is the interface that wraps a basic, read-only
variable backing store

### type [ReaderWriter](/readerwriter.go#L40)

`type ReaderWriter interface { ... }`

ReaderWriter is the interface that groups the basic Read and Write methods
for a variable backing store

### type [Writer](/writer.go#L40)

`type Writer interface { ... }`

Writer is the interface that wraps a basic, write-only
variable backing store

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
