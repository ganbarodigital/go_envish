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

// ShellEnv is a list of operations needed by a UNIX shell, or an emulation
// such as Scriptish.
type ShellEnv interface {
	Expander
	ReaderWriter

	// GetPositionalParamCount returns the value of the UNIX shell special
	// parameter $#.
	//
	// If $# is not set, it returns 0.
	GetPositionalParamCount() int

	// GetPositionalParams returns the (emulated) value of the UNIX
	// shell special parameter $@.
	//
	// It ignores any $@ that has been set in the environment, and builds
	// the list up using the value of $#.
	GetPositionalParams() []string

	// ReplacePositionalParams sets $1, $2 etc etc to the given values.
	//
	// Any existing positional parameters are deleted.
	//
	// Use SetPositionalParams instead, if you want to preserve any of
	// the existing positional params.
	//
	// It also sets the special parameter $#. The value of $# is returned.
	ReplacePositionalParams(values ...string) int

	// ResetPositionalParams deletes $1, $2 etc etc from the environment.
	//
	// It also sets the special parameter $# to 0.
	ResetPositionalParams()

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
	SetPositionalParams(values ...string) int

	// ShiftPositionalParams removes the first amount of positional params
	// from the environment.
	//
	// For example, if you call ShiftPositionalParams(1), then $3 becomes
	// $2, $2 becomes $1, and the original $1 is discarded.
	ShiftPositionalParams(amount int)
}
