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
	"sort"
)

// OverlayEnv works on a collection of variable backing stores
type OverlayEnv struct {
	envs []Expander
}

// NewOverlayEnv creates an empty stack of environment stores
func NewOverlayEnv(envs ...Expander) *OverlayEnv {
	retval := OverlayEnv{
		envs: envs,
	}

	// all done
	return &retval
}

// GetEnvByID returns the environment you want to work with
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

// Clearenv deletes all entries
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

// Environ returns a copy of all entries in the form "key=value".
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
			key := getKeyFromPair(pair)
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

// Getenv returns the value of the variable named by the key.
//
// If the key is not found, an empty string is returned.
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

// LookupEnv returns the value of the variable named by the key.
//
// If the key is not found, an empty string is returned, and the returned
// boolean is false.
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

// LookupHomeDir retrieves the given user's home directory, or false if
// that cannot be found
func (e *OverlayEnv) LookupHomeDir(username string) (string, bool) {
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
