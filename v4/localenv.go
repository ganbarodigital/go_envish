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
	"os/user"
	"strings"

	shellexpand "github.com/ganbarodigital/go_shellexpand"
)

// LocalEnv holds a list key/value pairs.
type LocalEnv struct {
	// pairs is the list we'll need to pass to Golang's standard library
	// for things like running external software
	pairs []string

	// pairKeys is a lookup table into pairs
	//
	// we populate this whenever anyone does a lookup, to speed up
	// any subsequent lookups of the same variable
	pairKeys map[string]int

	// should the variables in here be made available to external programs?
	//
	// this helps our EnvStack work out which stacked environments to
	// export out
	isExporter bool
}

// ================================================================
//
// Constructors
//
// ----------------------------------------------------------------

// NewLocalEnv creates an empty environment store
func NewLocalEnv(options ...func(*LocalEnv)) *LocalEnv {
	retval := LocalEnv{}

	// set aside some space for our fast key lookup
	retval.makePairIndex()

	// apply any options that we've been given
	for _, option := range options {
		option(&retval)
	}

	// all done
	return &retval
}

// ================================================================
//
// Reader interface
//
// ----------------------------------------------------------------

// Environ returns a copy of all entries in the form "key=value".
func (e *LocalEnv) Environ() []string {
	// do we have an environment store to work with?
	if e == nil {
		return []string{}
	}

	// yes we do
	return e.pairs
}

// Getenv returns the value of the variable named by the key.
//
// If the key is not found, an empty string is returned.
func (e *LocalEnv) Getenv(key string) string {
	// do we have an environment store to work with?
	if e == nil {
		return ""
	}

	// yes we do
	i := e.findPairIndex(key)
	if i >= 0 {
		key := getKeyFromPair(e.pairs[i])
		return getValueFromPair(e.pairs[i], key)
	}

	// not found
	return ""
}

// IsExporter returns true if this backing store holds variables that
// should be exported to external programs
func (e *LocalEnv) IsExporter() bool {
	return e.isExporter
}

// LookupEnv returns the value of the variable named by the key.
//
// If the key is not found, an empty string is returned, and the returned
// boolean is false.
func (e *LocalEnv) LookupEnv(key string) (string, bool) {
	// do we have an environment store to work with?
	if e == nil {
		return "", false
	}

	// yes we do
	i := e.findPairIndex(key)
	if i >= 0 {
		key := getKeyFromPair(e.pairs[i])
		return getValueFromPair(e.pairs[i], key), true
	}

	// not found
	return "", false
}

// MatchVarNames returns a list of variable names that start with the
// given prefix.
//
// This is very useful if you want to support `${PARAM:=word}` shell
// expansion in your own code.
func (e *LocalEnv) MatchVarNames(prefix string) []string {
	// our return value
	retval := []string{}

	// do we have an environment store to work with?
	if e == nil {
		return retval
	}

	// yes we do
	for i := range e.pairs {
		if strings.HasPrefix(e.pairs[i], prefix) {
			retval = append(retval, getKeyFromPair(e.pairs[i]))
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
func (e *LocalEnv) Clearenv() {
	// do we have an environment store to work with?
	if e == nil {
		return
	}

	// yes, we do
	e.pairs = []string{}
	e.makePairIndex()
}

// Setenv sets the value of the variable named by the key.
func (e *LocalEnv) Setenv(key, value string) error {
	// do we have an environment store to work with
	if e == nil {
		return ErrNilPointer{"LocalEnv.Setenv"}
	}

	// make sure we have a key that we can work with
	if len(key) == 0 || len(strings.TrimSpace(key)) == 0 {
		return ErrEmptyKey{}
	}

	// we need to update the Golang-compatible list too
	i := e.findPairIndex(key)
	if i >= 0 {
		// we're updating an existing entry
		e.pairs[i] = key + "=" + value
	} else {
		// we have a new entry!
		e.appendPairIndex(key, value)
	}

	// all done
	return nil
}

// Unsetenv deletes the variable named by the key.
func (e *LocalEnv) Unsetenv(key string) {
	// do we have an environment store to work with?
	if e == nil {
		return
	}

	// yes we do
	//
	// but do we have this variable?
	i := e.findPairIndex(key)
	if i < 0 {
		return
	}

	// we need to shuffle up
	e.pairs = append(e.pairs[:i], e.pairs[i+1:]...)

	// and we need to rewrite our fast lookup map too
	newPairKeys := make(map[string]int, len(e.pairKeys))
	for cachedKey, cachedIndex := range e.pairKeys {
		if cachedKey == key {
			continue
		}

		if cachedIndex >= i {
			newPairKeys[cachedKey] = cachedIndex - 1
		}
	}
	e.pairKeys = newPairKeys
}

// ================================================================
//
// Expander interface
//
// ----------------------------------------------------------------

// Expand replaces ${var} or $var in the input string.
func (e *LocalEnv) Expand(fmt string) string {
	// do we have an environment to work with?
	if e == nil {
		return fmt
	}

	// yes we do
	cb := shellexpand.ExpansionCallbacks{
		AssignToVar:   e.Setenv,
		LookupHomeDir: e.LookupHomeDir,
		LookupVar:     e.LookupEnv,
		MatchVarNames: e.MatchVarNames,
	}

	// attempt full-on shell expansion
	retval, err := shellexpand.Expand(fmt, cb)

	// did it work?
	if err != nil {
		return fmt
	}

	// yes it did :)
	return retval
}

// LookupHomeDir retrieves the given user's home directory, or false if
// that cannot be found
func (e *LocalEnv) LookupHomeDir(username string) (string, bool) {
	var details *user.User
	var err error

	if username == "" {
		details, err = user.Current()
	} else {
		details, err = user.Lookup(username)
	}

	if err != nil {
		return "", false
	}

	return details.HomeDir, true
}

// ================================================================
//
// ShellEnv interface
//
// ----------------------------------------------------------------

// GetPositionalParamCount returns the value of the UNIX shell special
// parameter $#.
//
// If $# is not set, it returns 0.
func (e *LocalEnv) GetPositionalParamCount() int {
	return getPositionalParamCount(e)
}

// GetPositionalParams returns the (emulated) value of the UNIX
// shell special parameter $@.
//
// It ignores any $@ that has been set in the environment, and builds
// the list up using the value of $#.
func (e *LocalEnv) GetPositionalParams() []string {
	return getPositionalParams(e)
}

// ReplacePositionalParams sets $1, $2 etc etc to the given values.
//
// Any existing positional parameters are deleted.
//
// Use SetPositionalParams instead, if you want to preserve any of
// the existing positional params.
//
// It also sets the special parameter $#. The value of $# is returned.
func (e *LocalEnv) ReplacePositionalParams(values ...string) int {
	return replacePositionalParams(e, values...)
}

// ResetPositionalParams deletes $1, $2 etc etc from the environment.
//
// It also sets the special parameter $# to 0.
func (e *LocalEnv) ResetPositionalParams() {
	resetPositionalParams(e)
}

// SetPositionalParams sets $1, $2 etc etc to the given values.
//
// Any existing positional parameters are overwritten, up to len(values).
// For example, the positional parameter $10 is *NOT* overwritten if
// you only pass in nine positional parameters.
//
// Use ReplacePositionalParams instead, if you want `values` to be the
// only positional parameters set.
//
// It also updates the special parameter $# if needed. The (possibly new)
// value of $# is returned.
func (e *LocalEnv) SetPositionalParams(values ...string) int {
	return setPositionalParams(e, values...)
}

// ================================================================
//
// Internal helpers
//
// ----------------------------------------------------------------

// Length returns the number of key/value pairs stored in the Env
func (e *LocalEnv) Length() int {
	// do we have an environment store to work with?
	if e == nil {
		return 0
	}

	// yes we do
	return len(e.pairs)
}

func (e *LocalEnv) findPairIndex(key string) int {
	// special case - we've already got this cached
	i, ok := e.pairKeys[key]
	if ok {
		return i
	}

	// general case - we have to search the full list of pairs
	//
	// this is what we are looking for
	prefix := key + "="

	// yes, this is horrible
	for i := range e.pairs {
		if strings.HasPrefix(e.pairs[i], prefix) {
			// cache it
			e.pairKeys[key] = i

			// all done
			return i
		}
	}

	// if we get here, the key doesn't exist in the pairs
	return -1
}

func (e *LocalEnv) appendPairIndex(key, value string) {
	// do we have a map to write to?
	if e.pairKeys == nil {
		e.makePairIndex()
	}

	// add the new keys to the end of the map
	e.pairs = append(e.pairs, key+"="+value)
	e.pairKeys[key] = len(e.pairs) - 1
}

func (e *LocalEnv) makePairIndex() {
	// set aside some space to store our faster lookups
	e.pairKeys = make(map[string]int, 10)
}
