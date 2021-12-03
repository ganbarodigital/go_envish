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

package envish

import (
	"os"
	"strings"
)

// ProgramEnv works directly on the program's environment
type ProgramEnv struct {
}

// ================================================================
//
// Constructors
//
// ----------------------------------------------------------------

// NewProgramEnv creates an empty environment store, and populates
// it with the contents of your program's environment.
func NewProgramEnv() *ProgramEnv {
	retval := ProgramEnv{}

	// all done
	return &retval
}

// ================================================================
//
// Reader interface
//
// ----------------------------------------------------------------

// Environ returns a copy of all entries in the form "key=value".
func (e *ProgramEnv) Environ() []string {
	return os.Environ()
}

// Getenv returns the value of the variable named by the key.
//
// If the key is not found, an empty string is returned.
func (e *ProgramEnv) Getenv(key string) string {
	return os.Getenv(key)
}

// IsExporter returns true if this backing store holds variables that
// should be exported to external programs
func (e *ProgramEnv) IsExporter() bool {
	return true
}

// LookupEnv returns the value of the variable named by the key.
//
// If the key is not found, an empty string is returned, and the returned
// boolean is false.
func (e *ProgramEnv) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

// MatchVarNames returns a list of variable names that start with the
// given prefix.
//
// This is very useful if you want to support `${PARAM:=word}` shell
// expansion in your own code.
func (e *ProgramEnv) MatchVarNames(prefix string) []string {
	// our return value
	retval := []string{}

	// the current, full environment
	pairs := os.Environ()
	for i := range pairs {
		if strings.HasPrefix(pairs[i], prefix) {
			retval = append(retval, getKeyFromPair(pairs[i]))
		}
	}

	// all done
	return retval
}

// ================================================================
//
// Writer interface
//
// ----------------------------------------------------------------

// Clearenv deletes all entries
func (e *ProgramEnv) Clearenv() {
	os.Clearenv()
}

// Setenv sets the value of the variable named by the key.
func (e *ProgramEnv) Setenv(key, value string) error {
	return os.Setenv(key, value)
}

// Unsetenv deletes the variable named by the key.
func (e *ProgramEnv) Unsetenv(key string) {
	_ = os.Unsetenv(key)
}

// ================================================================
//
// Expand interface
//
// ----------------------------------------------------------------

// Expand replaces ${var} or $var in the input string.
func (e *ProgramEnv) Expand(fmt string) string {
	return expand(e, fmt)
}

// LookupHomeDir retrieves the given user's home directory, or false if
// that cannot be found
func (e *ProgramEnv) LookupHomeDir(username string) (string, bool) {
	return lookupHomeDir(username)
}

// RestoreEnvironment writes the given "key=value" pairs into your
// program's environment
//
// It does *not* empty the environment.
func (e *ProgramEnv) RestoreEnvironment(pairs []string) {
	for _, pair := range pairs {
		key := getKeyFromPair(pair)
		value := getValueFromPair(pair, key)

		e.Setenv(key, value)
	}
}
