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

	envish "github.com/ganbarodigital/go_envish"
)

// ================================================================
//
// Constructor examples
//
// ----------------------------------------------------------------

func ExampleNewOverlayEnv() {
	// create individual environments
	localVars := envish.NewLocalEnv()
	progVars := envish.NewLocalEnv(envish.SetAsExporter)
	progEnv := envish.NewProgramEnv()

	// combine them
	env := envish.NewOverlayEnv(
		[]envish.Expander{
			localVars,
			progVars,
			progEnv,
		},
	)

	// you can now treat them as a single environment
	env.Setenv("$1", "go")
}

// ================================================================
//
// Reader interface examples
//
// ----------------------------------------------------------------

func ExampleOverlayEnv_Environ() {
	// create our independent environments
	localVars := envish.NewLocalEnv()
	progVars := envish.NewLocalEnv(envish.SetAsExporter)
	progEnv := envish.NewProgramEnv()

	// combine them
	env := envish.NewOverlayEnv(
		[]envish.Expander{
			localVars,
			progVars,
			progEnv,
		},
	)

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

func ExampleOverlayEnv_Getenv() {
	// build an environment stack without keeping a reference
	// to any of the individual environments
	env := envish.NewOverlayEnv(
		[]envish.Expander{
			envish.NewLocalEnv(),
			envish.NewLocalEnv(envish.SetAsExporter),
			envish.NewProgramEnv(),
		},
	)

	// show what we have
	fmt.Printf("USER is %s`n", env.Getenv("USER"))
}

func ExampleOverlayEnv_IsExporter() {
	// build an environment stack without keeping a reference
	// to any of the individual environments
	env := envish.NewOverlayEnv(
		[]envish.Expander{
			envish.NewLocalEnv(),
			envish.NewLocalEnv(envish.SetAsExporter),
			envish.NewProgramEnv(),
		},
	)

	fmt.Print(env.IsExporter())
	// Output:
	// true
}

// If none of the environments passed to NewOverlayEnv are exporters,
// IsExporter will return `false`.
func ExampleOverlayEnv_IsExporter_noExporters() {
	// build an environment stack without keeping a reference
	// to any of the individual environments
	env := envish.NewOverlayEnv(
		[]envish.Expander{
			envish.NewLocalEnv(),
			envish.NewLocalEnv(),
		},
	)

	fmt.Print(env.IsExporter())
	// Output:
	// false
}

// If you have a ProgramEnv anywhere in your OverlayEnv, IsExporter
// will always return `true`.
func ExampleOverlayEnv_IsExporter_programEnvIsAlwaysAnExporter() {
	// build an environment stack without keeping a reference
	// to any of the individual environments
	env := envish.NewOverlayEnv(
		[]envish.Expander{
			envish.NewLocalEnv(),
			envish.NewProgramEnv(),
		},
	)

	fmt.Print(env.IsExporter())
	// Output:
	// true
}

func ExampleOverlayEnv_LookupEnv() {
	// build an environment stack without keeping a reference
	// to any of the individual environments
	env := envish.NewOverlayEnv(
		[]envish.Expander{
			envish.NewLocalEnv(),
			envish.NewLocalEnv(envish.SetAsExporter),
			envish.NewProgramEnv(),
		},
	)

	home, ok := env.LookupEnv("HOME")
	fmt.Printf("did we find $HOME?: %v", ok)
	fmt.Printf("what value did we find?: %v", home)
}

func ExampleOverlayEnv_MatchVarNames() {
	// build an environment stack without keeping a reference
	// to any of the individual environments
	env := envish.NewOverlayEnv(
		[]envish.Expander{
			envish.NewLocalEnv(),
			envish.NewLocalEnv(envish.SetAsExporter),
			envish.NewProgramEnv(),
		},
	)

	// print out all the variables with the prefix ANSIBLE_
	for _, key := range env.MatchVarNames("ANSIBLE_") {
		fmt.Printf("%s = %s", key, env.Getenv(key))
	}
}

// ================================================================
//
// Writer interface examples
//
// ----------------------------------------------------------------

// there is no Clearenv example because it's too dangerous to
// include here!

func ExampleOverlayEnv_Setenv() {
	// create individual environments
	localVars := envish.NewLocalEnv()
	progVars := envish.NewLocalEnv(envish.SetAsExporter)
	progEnv := envish.NewProgramEnv()

	// combine them
	env := envish.NewOverlayEnv(
		[]envish.Expander{
			localVars,
			progVars,
			progEnv,
		},
	)

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

func ExampleOverlayEnv_Unsetenv() {
	// create individual environments
	localVars := envish.NewLocalEnv()
	progVars := envish.NewLocalEnv(envish.SetAsExporter)
	progEnv := envish.NewProgramEnv()

	// combine them
	env := envish.NewOverlayEnv(
		[]envish.Expander{
			localVars,
			progVars,
			progEnv,
		},
	)

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

	// Output:
	// value is: 100
	// ok1 is: false
	// ok2 is: false
	// ok3 is: false
}

// ================================================================
//
// Expander interface examples
//
// ----------------------------------------------------------------

func ExampleOverlayEnv_Expand() {
	// create individual environments
	localVars := envish.NewLocalEnv()
	progVars := envish.NewLocalEnv(envish.SetAsExporter)
	progEnv := envish.NewProgramEnv()

	// combine them
	env := envish.NewOverlayEnv(
		[]envish.Expander{
			localVars,
			progVars,
			progEnv,
		},
	)

	// use UNIX shell expansion to see what we have
	fmt.Print(env.Expand("USER is ${USER}\n"))
}

// ================================================================
//
// Unique method examples
//
// ----------------------------------------------------------------

func ExampleOverlayEnv_GetTopMostEnv() {
	// build an environment stack without keeping a reference
	// to any of the individual environments
	env := envish.NewOverlayEnv(
		[]envish.Expander{
			envish.NewLocalEnv(),
			envish.NewLocalEnv(envish.SetAsExporter),
			envish.NewProgramEnv(),
		},
	)

	// now, imagine we want to set some local variables
	localVars, _ := env.GetTopMostEnv()
	localVars.Setenv("$#", "2")
}
