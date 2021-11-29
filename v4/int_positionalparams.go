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
	"fmt"
	"strconv"
)

func buildPositionalParamName(position int) string {
	return fmt.Sprintf("$%d", position)
}

func getPositionalParamCount(e Reader) int {
	retval, err := strconv.Atoi(e.Getenv("$#"))
	if err != nil {
		return 0
	}

	return retval
}

func getPositionalParams(e Reader) []string {
	positionalCount := getPositionalParamCount(e)

	// special case - no params to return
	if positionalCount == 0 {
		return []string{}
	}

	retval := make([]string, 0, positionalCount)
	for i := 1; i <= positionalCount; i++ {
		name := buildPositionalParamName(i)
		retval = append(retval, e.Getenv(name))
	}

	// all done
	return retval
}

// replacePositionalParams sets $1, $2 etc etc to the given values.
//
// Any existing positional parameters are deleted.
//
// Use SetPositionalParams instead, if you want to preserve any of
// the existing positional params.
//
// It also sets the special parameter $#. The value of $# is returned.
func replacePositionalParams(e ReaderWriter, values ...string) int {
	// get rid of any existing positional params
	resetPositionalParams(e)

	// set the new ones
	return setPositionalParams(e, values...)
}

func resetPositionalParams(e ReaderWriter) {
	// how many positional params do we have?
	positionalCount := getPositionalParamCount(e)

	// let's get rid of them
	for i := 1; i <= positionalCount; i++ {
		name := buildPositionalParamName(i)
		e.Unsetenv(name)
	}

	// now we need to update $# too
	e.Setenv("$#", "0")
}

// setPositionalParams sets $1, $2 etc etc to the given values.
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
func setPositionalParams(e ReaderWriter, values ...string) int {
	// set the positional values
	for i, value := range values {
		name := buildPositionalParamName(i + 1)
		e.Setenv(name, value)
	}

	return updatePositionalCount(e, len(values))
}

// ShiftPositionalParams removes the first amount of positional params
// from the environment.
//
// For example, if you call ShiftPositionalParams(1), then $3 becomes
// $2, $2 becomes $1, and the original $1 is discarded.
func shiftPositionalParams(e ReaderWriter, amount int) {
	// how many positional params do we have?
	paramCount := getPositionalParamCount(e)

	// special case - are we removing all of the remaining positional
	// parameters?
	if amount >= paramCount {
		resetPositionalParams(e)
		return
	}

	// general case
	params := getPositionalParams(e)
	replacePositionalParams(e, params[amount:]...)
}

func updatePositionalCount(e ReaderWriter, newLen int) int {
	// what is the current value of $#?
	positionalCount := getPositionalParamCount(e)

	// do we need to update $#?
	if positionalCount < newLen {
		e.Setenv("$#", strconv.Itoa(newLen))
		return newLen
	}

	// no, we do not
	return positionalCount
}
