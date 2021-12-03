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

/*
envish is a library to help you emulate UNIX-like program environments
in Golang packages.

It is released under the 3-clause New BSD license. See LICENSE.md for details.

  import envish "github.com/ganbarodigital/go_envish/v3"

  env := envish.NewLocalEnv()

  // add to this temporary environment
  // WITHOUT changing your program's environment
  env.Setenv("EXAMPLE_KEY", "EXAMPLE VALUE")

  // pass it into run a child process
  cmd := exec.Command(...)
  cmd.Env = env.Environ()
  cmd.Start()


Why Use Envish

We've built Envish for anyone who needs to emulate a UNIX-like environment in
their own Golang packages. Or anyone who just needs a simple key/value store
with a familiar API.

We're using it ourselves for our (Pipe) https://github.com/ganbarodigital/go_pipe
and (Scriptish) https://github.com/ganbarodigital/go_scriptish packages.


Why A Separate Package

Golang's `os` package provides support for working with your program's
environment. But what if you want to make temporary changes to that
environment, just to pass environment variables into child processes?

This is a very common pattern used in UNIX shell script programming:

  DEBIAN_FRONTEND=noninteractive apt-get install -y mysql

In the example above, the environment variable `DEBIAN_FRONTEND` is only set
for the child process `apt-get`.


Getting Started

Import Envish into your Golang code:

  import envish "github.com/ganbarodigital/go_envish/v3"


Don't forget that `v3` on the end of the import, or you'll get an older
version of this package!

Create a copy of your program's environment:

  localEnv := envish.NewLocalEnv(envish.CopyProgramEnv)

or simply start with an empty environment store:

  localVars := envish.NewLocalEnv()

Get and set variables in the environment store as needed:

  home := localEnv.Getenv()
  localEnv.Setenv("DEBIAN_FRONTEND", "noninteractive")

*/
package envish
