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

package envish_test

import (
	"testing"

	envish "github.com/ganbarodigital/go_envish/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetKeyFromPairReturnsKey(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testData := "PARAM1=foo"
	expectedResult := "PARAM1"

	// ----------------------------------------------------------------
	// perform the change

	actualResult := envish.GetKeyFromPair(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestGetKeyFromPairSupportsVarsThatStartWithEqualsSign(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testData := "=PARAM1=foo"
	expectedResult := "=PARAM1"

	// ----------------------------------------------------------------
	// perform the change

	actualResult := envish.GetKeyFromPair(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestGetValueFromPairReturnsValue(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testData := "PARAM1=foo"
	expectedResult := "foo"

	// ----------------------------------------------------------------
	// perform the change

	actualResult := envish.GetValueFromPair(testData, "PARAM1")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestGetValueFromPairSupportsEmptyValue(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testData := "PARAM1="
	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	actualResult := envish.GetValueFromPair(testData, "PARAM1")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}
