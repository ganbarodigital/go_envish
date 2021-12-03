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

import "fmt"

// ErrEmptyKey is returned whenever we're given a key that is zero-length
// or only contains whitespace
type ErrEmptyKey struct{}

func (e ErrEmptyKey) Error() string {
	return "zero-length key, or key only contains whitespace"
}

// ErrEmptyOverlayEnv is returned whenever you call a method on an empty EnvStack
type ErrEmptyOverlayEnv struct {
	Method string
}

func (e ErrEmptyOverlayEnv) Error() string {
	return fmt.Sprintf("overlay env is empty; %s", e.Method)
}

// ErrNilPointer is returned whenever you call a method on the Env struct
// with a nil pointer
type ErrNilPointer struct {
	Method string
}

func (e ErrNilPointer) Error() string {
	return fmt.Sprintf("nil pointer to environment store passed to %s", e.Method)
}

type ErrNoExporterEnv struct {
	Method string
}

func (e ErrNoExporterEnv) Error() string {
	return fmt.Sprintf("no exporting environment in OverlayEnv passed to %s", e.Method)
}
