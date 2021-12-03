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

	envish "github.com/ganbarodigital/go_envish/v4"
)

func ExampleNewProgramEnv() {
	env := envish.NewProgramEnv()

	fmt.Printf("USER is %s", env.Getenv("USER"))
}

func ExampleProgramEnv_Environ() {
	// get access to our program environment
	env := envish.NewProgramEnv()

	// get a list of all entries in the environment
	environ := env.Environ()

	// pass it into run a child process
	cmd := exec.Command("godoc")
	cmd.Env = environ

	// you can now call cmd.Start()
}

func ExampleProgramEnv_Expand() {
	// get access to our program environment
	env := envish.NewProgramEnv()

	// show what we have
	fmt.Print(env.Expand("USER is ${USER}"))
}

func ExampleProgramEnv_Getenv() {
	// get access to our program environment
	env := envish.NewProgramEnv()

	// show what we have
	fmt.Printf("USER is %s", env.Getenv("USER"))
}

func ExampleProgramEnv_IsExporter() {
	// get access to our program environment
	env := envish.NewProgramEnv()

	// a ProgramEnv is always an environment exporter
	fmt.Printf("env.IsExporter() is %v", env.IsExporter())

	// Output:
	// env.IsExporter() is true
}

func ExampleProgramEnv_LookupEnv() {
	// get access to our program environment
	env := envish.NewProgramEnv()

	value, ok := env.LookupEnv("USER")

	fmt.Printf("ok is: %v\n", ok)
	fmt.Printf("value is: %v\n", value)
}

func ExampleProgramEnv_MatchVarNames() {
	// get access to your program environment
	env := envish.NewProgramEnv()

	// print out all the variables with the prefix ANSIBLE_
	for _, key := range env.MatchVarNames("ANSIBLE_") {
		fmt.Printf("%s = %s", key, env.Getenv(key))
	}
}

func ExampleProgramEnv_RestoreEnvironment() {
	// get access to your program environment
	env := envish.NewProgramEnv()

	// take a backup of the whole environment
	backup := env.Environ()

	// nuke the environment
	env.Clearenv()

	// restore from backup
	env.RestoreEnvironment(backup)
}

func ExampleProgramEnv_Setenv() {
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

func ExampleProgramEnv_Unsetenv() {
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
