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
	"sort"
)

// OverlayEnv works on a collection of variable backing stores.
//
// Use an OverlayEnv to combine one (or more) LocalEnv and a single
// ProgramEnv into a single logical environment.
//
// We do this in https://github.com/ganbarodigital/go_scriptish to
// emulate local variable support.
type OverlayEnv struct {
	envs []Expander
}

// ================================================================
//
// Constructors
//
// ----------------------------------------------------------------

// NewOverlayEnv builds a single logical environments from the
// given set of underlying environments.
//
// NewOverlayEnv returns a pointer to an OverlayEnv struct. You can use the
// methods of the OverlayStruct to read from and write to each of these
// environments.
//
// The order of the arguments to NewOverlayEnv matters. The returned
// OverlayEnv's methods will read from / write to the underlying environments
// in the order you've given.
func NewOverlayEnv(envs []Expander) *OverlayEnv {
	retval := OverlayEnv{
		envs: envs,
	}

	// all done
	return &retval
}

// ================================================================
//
// Reader interface
//
// ----------------------------------------------------------------

// Environ() returns a copy of all of the variables in your `OverlayEnv`
// in the form `key=value`. This format is compatible with Golang's
// built-in packages.
//
// When it builds the list, it follows these rules:
//
// * it searches the environments in the order you provided them to
// NewOverlayEnv
//
// * it only includes variables from environments where the IsExporter
// method returns `true`
//
// * if the same variable is set in multiple environments, it uses the first
// value it finds
func (e *OverlayEnv) Environ() []string {
	// our return value
	retval := []string{}

	// do we have a stack to work with?
	if e == nil {
		return retval
	}

	// we need somewhere to keep track of the variables we are exporting
	foundPairs := make(map[string]string)

	for _, env := range e.envs {
		if !env.IsExporter() {
			continue
		}

		pairs := env.Environ()
		for _, pair := range pairs {
			key := GetKeyFromPair(pair)
			_, ok := foundPairs[key]
			if !ok {
				foundPairs[key] = pair
			}
		}
	}

	// at this point, foundPairs needs to be flattened
	for _, pair := range foundPairs {
		retval = append(retval, pair)
	}

	// sort the results, otherwise they're non-deterministic
	// which makes tests unreliable
	sort.Strings(retval)

	// all done
	return retval
}

// Getenv returns the value of the given variable. If the variable does not
// exist, it returns `""` (empty string).
//
// * it searches the environments in the order you provided them to NewOverlayEnv
//
// * if the same variable is set in multiple environments, it uses the first
// value it finds
func (e *OverlayEnv) Getenv(key string) string {
	// do we have a variable backing store to work with?
	if e == nil {
		return ""
	}

	// search for this variable
	for _, env := range e.envs {
		value, ok := env.LookupEnv(key)
		if ok {
			return value
		}
	}

	// if we get here, then it doesn't exist
	return ""
}

// IsExporter returns `true` if (and only if) any of the environments in
// the overlay env hold variables that should be exported to external
// programs.
func (e *OverlayEnv) IsExporter() bool {
	// do we have an overlay to work with?
	if e == nil {
		return false
	}

	// do any of our overlays export anything?
	for _, env := range e.envs {
		if env.IsExporter() {
			return true
		}
	}

	// sadly, they do not
	return false
}

// LookupEnv returns the value of the given variable. If the variable does
// not exist, it returns `"", false`.
//
// * it searches the environments in the order you provided them to NewOverlayEnv
//
// * if the same variable is set in multiple environments, it uses the first
// value it finds
func (e *OverlayEnv) LookupEnv(key string) (string, bool) {
	// do we have a stack?
	if e == nil {
		return "", false
	}

	for _, env := range e.envs {
		value, ok := env.LookupEnv(key)
		if ok {
			return value, true
		}
	}

	// no joy
	return "", false
}

// MatchVarNames returns a list of variable names that start with the
// given prefix.
//
// It's a feature needed for `${!prefix*}` string expansion syntax.
//
// * it searches the environments in the order you provided them to NewOverlayEnv
//
// * if the same key is found in multiple environments, it only returns
// the key once (ie, results are deduped before they are returned)
func (e *OverlayEnv) MatchVarNames(prefix string) []string {
	// our return value
	retval := []string{}

	// do we have a stack to work with?
	if e == nil {
		return retval
	}

	// let's go and find things
	foundKeys := make(map[string]bool)
	for _, env := range e.envs {
		keys := env.MatchVarNames(prefix)
		for _, key := range keys {
			foundKeys[key] = true
		}
	}

	for key := range foundKeys {
		retval = append(retval, key)
	}

	// sort it
	sort.Strings(retval)

	// all done
	return retval
}

// ================================================================
//
// Writer interface
//
// ----------------------------------------------------------------

// Clearenv deletes all variables in every environment in the OverlayEnv.
// If your overlay env includes a ProgramEnv, this *WILL* delete all of
// your program's environment variables.
//
// Use with extreme caution!
func (e *OverlayEnv) Clearenv() {
	// do we have a stack to work with?
	if e == nil {
		return
	}

	// wipe them out ... all of them ...
	for i := range e.envs {
		e.envs[i].Clearenv()
	}

	// all done
}

// Setenv creates a variable (if it doesn't already exist) or updates its
// value (if it does exist).
//
// * it searches the environments in the order you provided to NewOverlayEnv
//
// * if the same variable is set in multiple environments, it updates the
// first variable it finds
//
// * if the variable does not exist, it is always created in the first
// environment you provided to NewOverlayEnv
func (e *OverlayEnv) Setenv(key, value string) error {
	// do we have a stack?
	if e == nil {
		return ErrNilPointer{"OverlayEnv.Setenv"}
	}

	// do we have any environments in the stack?
	if len(e.envs) == 0 {
		return ErrEmptyOverlayEnv{"OverlayEnv.Setenv"}
	}

	// are we updating an existing variable?
	for _, env := range e.envs {
		_, ok := env.LookupEnv(key)
		if ok {
			return env.Setenv(key, value)
		}
	}

	// nope, it's a brand new variable
	return e.envs[0].Setenv(key, value)
}

// Unsetenv deletes the variable named by the key.
//
// It will be deleted from all the environments in the stack.
func (e *OverlayEnv) Unsetenv(key string) {
	// do we have a stack?
	if e == nil {
		return
	}

	for _, env := range e.envs {
		env.Unsetenv(key)
	}
}

// ================================================================
//
// Expander interface
//
// ----------------------------------------------------------------

// Expand will replace `${key}` and `$key` entries in a format string,
// by looking up values from the environments contained within the OverlayEnv.
//
// * it uses the given OverlayEnv's LookupEnv to find the values of variables
//
// * it uses the given OverlayEnv's Setenv to set the values of variables
//
// * it uses the given OverlayEnv's MatchVarNames to expand variable name prefixes
//
// * it uses the given OverlayEnv's LookupHomeDir to expand `~` (tilde)
//
// Hopefully, we've got the logic right, and you'll find that your expansions
// just work the way you'd naturally expect.
//
// Internally, it uses https://github.com/ganbarodigital/go_shellexpand to
// do the shell expansion. It supports the vast majority of UNIX shell
// string expansion operations.
func (e *OverlayEnv) Expand(fmt string) string {
	return expand(e, fmt)
}

// ================================================================
//
// Struct-unique functions
//
// ----------------------------------------------------------------

// GetEnvByID returns the requested environment from the given OverlayEnv.
// ID `0` is the first environment you passed into NewOverlayEnv, ID `1`
// is the second environment, and so on.
//
// If you request an ID that the given OverlayEnv does not have, it returns
// `nil, false`.
//
// GetEnvByID is handy if you don't want to keep separate references to the
// environments after you've combined them into the OverlayEnv.
func (e *OverlayEnv) GetEnvByID(id int) (Expander, bool) {
	// do we have a stack to work with?
	if e == nil {
		return nil, false
	}

	// do we have the environment that has been requested?
	if id >= len(e.envs) || id < 0 {
		return nil, false
	}

	// yes, we do
	return e.envs[id], true
}

// GetTopMostEnv returns the envish environment that's at the top of the
// overlay stack.
//
// If we don't have that environment, we return a suitable error.
func (e *OverlayEnv) GetTopMostEnv() (Expander, error) {
	// do we have an OverlayEnv to work with?
	if e == nil {
		return nil, ErrNilPointer{"OverlayEnv.GetTopMostEnv"}
	}

	// do we have a stack to work with?
	if len(e.envs) == 0 {
		return nil, ErrEmptyOverlayEnv{}
	}

	// yes we do
	return e.envs[0], nil
}

// Export emulates UNIX shell `export XXX=YYY` behaviour. It returns an
// error if the environment variable cannot be set.
//
// * It makes no changes at all if the given OverlayEnv does not have any
// environments that are exporters.
//
// * It searches each environment inside the given OverlayEnv in the order
// that you passed them into NewOverlayEnv.
//
// * It overwrites any existing value for the given key, in every environment
// that it searches (including environments that are not exporters). This
// should ensure consistent results whenever you call Getenv on the given
// OverlayEnv.
//
// * It stops once it has set the environment variable inside an environment
// that is an exporter.
func (e *OverlayEnv) Export(key, value string) error {
	// do we have an OverlayEnv to work with?
	if e == nil {
		return ErrNilPointer{"OverlayEnv.Export"}
	}

	// do we have a stack to work with?
	if len(e.envs) == 0 {
		return ErrEmptyOverlayEnv{}
	}

	// do we have any exporters in the stack?
	hasExporter := false
	for _, env := range e.envs {
		if env.IsExporter() {
			hasExporter = true
		}
	}
	if !hasExporter {
		return ErrNoExporterEnv{"OverlayEnv.Export"}
	}

	// work through the stack
	for _, env := range e.envs {
		// shorthand
		isExporter := env.IsExporter()
		_, hasKey := env.LookupEnv(key)

		if isExporter || hasKey {
			err := env.Setenv(key, value)
			if err != nil {
				// we have to bail
				return err
			}
		}

		// are we done?
		if isExporter {
			break
		}
	}

	// all done
	return nil
}
