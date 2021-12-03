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

package envish_test

import (
	"fmt"
	"os/exec"

	envish "github.com/ganbarodigital/go_envish/v3"
)

func ExampleLocalEnv() {
	// create a local environment
	localEnv := envish.NewLocalEnv()

	// it starts as an empty environment
	fmt.Print(localEnv.Getenv("$USER"))
}

func ExampleNewLocalEnv() {
	// create a local environment
	localEnv := envish.NewLocalEnv()

	// it starts as an empty environment
	fmt.Print(localEnv.Getenv("$USER"))
}

// CopyProgramEnv is a functional option that will populate the environment
// store with a copy of your program's environment.
func ExampleNewLocalEnv_withFunctionalOptions1() {
	// create a local environment
	//
	// it will start with a copy of your program's environment
	localEnv := envish.NewLocalEnv(envish.CopyProgramEnv)

	// on UNIX-like systems, this will print the name of the user
	// who is running the program
	fmt.Print(localEnv.Getenv("$USER"))
}

// SetAsExporter is a functional option that will tell the OverlayEnv to
// include your LocalEnv's contents in any call to the OverlayEnv's
// Environ function.
func ExampleNewLocalEnv_withFunctionalOptions2() {
	// create a local environment
	//
	// if you add this to an OverlayEnv, the OverlayEnv will include
	// its contents when you call the OverlayEnv's Environ method.
	localEnv := envish.NewLocalEnv(envish.SetAsExporter)

	// this environment is now an exporter
	fmt.Print(localEnv.IsExporter())
	// Output:
	//
	// true
}

// ================================================================
//
// Reader interface examples
//
// ----------------------------------------------------------------

func ExampleLocalEnv_Environ() {
	// create an environment store
	localEnv := envish.NewLocalEnv()

	// get a copy to pass into `os/exec`
	cmd := exec.Command("go", "doc")
	cmd.Env = localEnv.Environ()
}

func ExampleLocalEnv_Getenv() {
	// create an environment store
	localEnv := envish.NewLocalEnv()

	// get a variable from the environment store
	user := localEnv.Getenv("USER")
	fmt.Print(user)
}

func ExampleLocalEnv_IsExporter() {
	// create an environment store
	localEnv := envish.NewLocalEnv()

	// by default, the environment store is NOT an exporter
	//
	// if you want to change this hint, use envish.SetAsExporter
	exporting := localEnv.IsExporter()
	fmt.Print(exporting)
	// Output:
	//
	// false
}

func ExampleLocalEnv_LookupEnv() {
	// create an environment store
	localEnv := envish.NewLocalEnv()

	// find out if a key exists
	value, ok := localEnv.LookupEnv("USER")
	fmt.Printf("key exists: %v", ok)
	fmt.Printf("value of key: %s", value)
}

func ExampleLocalEnv_MatchVarNames() {
	// create an environment store
	localEnv := envish.NewLocalEnv()

	// find all variables that begin with 'ANSIBLE_'
	for _, key := range localEnv.MatchVarNames("ANSIBLE_") {
		fmt.Printf("%s = %s", key, localEnv.Getenv(key))
	}
}

// ================================================================
//
// Writer interface examples
//
// ----------------------------------------------------------------

func ExampleLocalEnv_Clearenv() {
	// create a environment store
	localEnv := envish.NewLocalEnv(envish.CopyProgramEnv)

	// empty the environment store completely
	localEnv.Clearenv()
}

func ExampleLocalEnv_Setenv() {
	// create an environment store
	localEnv := envish.NewLocalEnv()

	// set a value in the environment store
	localEnv.Setenv("DEBIAN_FRONTEND", "noninteractive")
}

// Setenv will return a `envish.ErrEmptyKey` error if you pass in a key that
// is either empty, or contains only whitespace.
func ExampleLocalEnv_Setenv_errEmptyKey() {
	// create an environment store
	localEnv := envish.NewLocalEnv()

	// try to create an invalid environment variable
	err := localEnv.Setenv("", "key-is-invalid")
	fmt.Print(err)
	// Output:
	// zero-length key, or key only contains whitespace
}

// Setenv will return an `envish.ErrNilPointer` error if you call Setenv with
// a nil pointer to the environment store.

func ExampleLocalEnv_Setenv_errNilPointer() {
	var localEnv *envish.LocalEnv = nil
	err := localEnv.Setenv("valid-key", "valid-value")
	fmt.Print(err)
	// Output:
	// nil pointer to environment store passed to LocalEnv.Setenv
}

func ExampleLocalEnv_Unsetenv() {
	// create an environment store
	localEnv := envish.NewLocalEnv()

	// delete an entry from the environment store
	localEnv.Unsetenv("$#")
}

// ================================================================
//
// Expander interface examples
//
// ----------------------------------------------------------------

func ExampleLocalEnv_Expand() {
	// create an environment store
	localEnv := envish.NewLocalEnv()

	// show what we have
	fmt.Print(localEnv.Expand("USER is ${USER}\n"))
}

// ================================================================
//
// Unique methods examples
//
// ----------------------------------------------------------------

func ExampleLocalEnv_Length() {
	// create an environment store
	localEnv := envish.NewLocalEnv()

	// find out how many variables it contains
	//
	// a new LocalEnv starts with no entries
	fmt.Printf("environment has %d entries\n", localEnv.Length())
	// Output:
	// environment has 0 entries
}
