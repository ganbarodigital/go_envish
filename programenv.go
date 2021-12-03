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

// ProgramEnv puts helper wrapper functions around your program's
// environment.
type ProgramEnv struct {
}

// ================================================================
//
// Constructors
//
// ----------------------------------------------------------------

// NewProgramEnv returns an envish environment that works directly with
// your program's environment.
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
// This format is compatible with Golang's built-in packages.
func (e *ProgramEnv) Environ() []string {
	return os.Environ()
}

// Getenv returns the value of the variable named by the key.
//
// If the key is not found, an empty string is returned.
func (e *ProgramEnv) Getenv(key string) string {
	return os.Getenv(key)
}

// IsExporter always returns `true`.
//
// It is used by OverlayEnv.Environ to work out which keys and values
// the OverlayEnv should include in its output.
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
// It's a feature needed for `${!prefix*}` string expansion syntax.
func (e *ProgramEnv) MatchVarNames(prefix string) []string {
	// our return value
	retval := []string{}

	// the current, full environment
	pairs := os.Environ()
	for i := range pairs {
		if strings.HasPrefix(pairs[i], prefix) {
			retval = append(retval, GetKeyFromPair(pairs[i]))
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

// Clearenv deletes all entries from your program's environment.
// Use with extreme caution!
func (e *ProgramEnv) Clearenv() {
	os.Clearenv()
}

// Setenv sets the value of the variable named by the key.
func (e *ProgramEnv) Setenv(key, value string) error {
	return os.Setenv(key, value)
}

// Unsetenv deletes the variable named by the key.
//
// This will remove the given variable from your program's environment.
// Use with caution!
func (e *ProgramEnv) Unsetenv(key string) {
	_ = os.Unsetenv(key)
}

// ================================================================
//
// Expand interface
//
// ----------------------------------------------------------------

// Expand replaces ${var} or $var in the input string, by looking up
// values from your program's environment.
//
// Internally, it uses https://github.com/ganbarodigital/go_shellexpand
// to do the shell expansion. It supports the vast majority of UNIX shell
// string expansion operations.
func (e *ProgramEnv) Expand(fmt string) string {
	return expand(e, fmt)
}

// ================================================================
//
// Unique methods
//
// ----------------------------------------------------------------

// RestoreEnvironment writes the given "key=value" pairs into your
// program's environment.
//
// It does *not* empty your program's environment first!
//
// It was originally added so that our unit tests could put the 'go test'
// program environment back in place after each test had run.
func (e *ProgramEnv) RestoreEnvironment(pairs []string) {
	for _, pair := range pairs {
		key := GetKeyFromPair(pair)
		value := GetValueFromPair(pair, key)

		e.Setenv(key, value)
	}
}
