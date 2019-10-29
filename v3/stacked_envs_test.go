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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStackedEnvsReturnsStackOfEnvironments(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	localEnv := NewLocalEnv()
	progEnv := NewProgramEnv()

	// ----------------------------------------------------------------
	// perform the change

	stack := NewStackedEnvs(localEnv, progEnv)

	// ----------------------------------------------------------------
	// test the results

	assert.Same(t, localEnv, stack.envs[0])
	assert.Same(t, progEnv, stack.envs[1])
}

func TestStackedEnvsGetEnvByIDReturnsFromTheStack(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	localEnv := NewLocalEnv()
	progEnv := NewProgramEnv()

	stack := NewStackedEnvs(localEnv, progEnv)

	// ----------------------------------------------------------------
	// perform the change

	stack0, ok0 := stack.GetEnvByID(0)
	stack1, ok1 := stack.GetEnvByID(1)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok0)
	assert.Same(t, localEnv, stack0)
	assert.True(t, ok1)
	assert.Same(t, progEnv, stack1)
}

func TestStackedEnvsGetEnvByIDReturnsNilIfIndexTooLarge(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	localEnv := NewLocalEnv()
	progEnv := NewProgramEnv()

	stack := NewStackedEnvs(localEnv, progEnv)

	// ----------------------------------------------------------------
	// perform the change

	stack2, ok2 := stack.GetEnvByID(2)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok2)
	assert.Nil(t, stack2)
}

func TestStackedEnvsGetEnvByIDReturnsNilIfIndexTooSmall(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	localEnv := NewLocalEnv()
	progEnv := NewProgramEnv()

	stack := NewStackedEnvs(localEnv, progEnv)

	// ----------------------------------------------------------------
	// perform the change

	stack2, ok2 := stack.GetEnvByID(-2)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok2)
	assert.Nil(t, stack2)
}

func TestStackedEnvsGetEnvByIDReturnsCopesWithNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	var stack *StackedEnvs = nil

	// ----------------------------------------------------------------
	// perform the change

	stack2, ok2 := stack.GetEnvByID(0)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok2)
	assert.Nil(t, stack2)
}

func TestStackedEnvsGetEnvByIDReturnsCopesWithEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	stack := StackedEnvs{}

	// ----------------------------------------------------------------
	// perform the change

	stack2, ok2 := stack.GetEnvByID(0)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok2)
	assert.Nil(t, stack2)
}
