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

// ================================================================
//
// buildPositionalParamName
//
// ----------------------------------------------------------------

func TestBuildPositionalParamNameWorks(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := []string{
		"$1",
		"$2",
		"$3",
		"$4",
		"$5",
		"$6",
		"$7",
		"$8",
		"$9",
		"$10",
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult := make([]string, 10)

	for i := 1; i <= 10; i++ {
		actualResult[i-1] = buildPositionalParamName(i)
	}

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

// ================================================================
//
// getPositionalParamCount()
//
// ----------------------------------------------------------------

func TestGetPositionalParamCountReturnsDollarHash(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	testData := []string{
		"one",
		"two",
		"three",
		"four",
		"five",
		"six",
		"seven",
		"eight",
		"nine",
		"ten",
	}

	env := NewLocalEnv()
	env.SetPositionalParams(testData...)

	expectedResult := 10

	// ----------------------------------------------------------------
	// perform the change

	actualResult := getPositionalParamCount(env)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

// ================================================================
//
// getPositionalParams()
//
// ----------------------------------------------------------------

func TestGetPositionalParamsEmulatesDollarStar(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	testData := []string{
		"one",
		"two",
		"three",
		"four",
		"five",
		"six",
		"seven",
		"eight",
		"nine",
		"ten",
	}

	env := NewLocalEnv()
	env.SetPositionalParams(testData...)

	expectedResult := testData

	// ----------------------------------------------------------------
	// perform the change

	actualResult := getPositionalParams(env)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

// ================================================================
//
// replacePositionalParams()
//
// ----------------------------------------------------------------

func TestReplacePositionalParamsSetsThePositionalParams(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	testData := []string{
		"one",
		"two",
		"three",
		"four",
		"five",
		"six",
		"seven",
		"eight",
		"nine",
		"ten",
	}

	env := NewLocalEnv()

	expectedLen := 10
	expectedResult := []string{
		"$#=10",
		"$1=one",
		"$2=two",
		"$3=three",
		"$4=four",
		"$5=five",
		"$6=six",
		"$7=seven",
		"$8=eight",
		"$9=nine",
		"$10=ten",
	}

	// ----------------------------------------------------------------
	// perform the change

	actualLen := replacePositionalParams(env, testData...)
	actualResult := env.Environ()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedLen, actualLen)
	assert.Equal(t, expectedResult, actualResult)
}

func TestReplacePositionalParamsReplacesAllExistingPositionalParams(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	seedData := []string{
		"one",
		"two",
		"three",
		"four",
		"five",
		"six",
		"seven",
		"eight",
		"nine",
		"ten",
	}

	testData := []string{
		"new one",
		"new two",
		"new three",
	}

	env := NewLocalEnv()
	env.SetPositionalParams(seedData...)

	expectedLen := 3
	expectedResult := []string{
		"$#=3",
		"$1=new one",
		"$2=new two",
		"$3=new three",
	}

	// ----------------------------------------------------------------
	// perform the change

	actualLen := replacePositionalParams(env, testData...)
	actualResult := env.Environ()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedLen, actualLen)
	assert.Equal(t, expectedResult, actualResult)
}

func TestReplacePositionalParamsUpdatesDollarHash(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	seedData := []string{
		"one",
		"two",
		"three",
	}

	testData := []string{
		"new one",
		"new two",
		"new three",
		"new four",
		"new five",
	}

	env := NewLocalEnv()
	env.SetPositionalParams(seedData...)

	origHash := env.Getenv("$#")
	assert.Equal(t, "3", origHash)

	expectedResult := "5"

	// ----------------------------------------------------------------
	// perform the change

	replacePositionalParams(env, testData...)
	actualResult := env.Getenv("$#")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

// ================================================================
//
// resetPositionalParams()
//
// ----------------------------------------------------------------

func TestResetPositionalParamsDeletesThePositionalParams(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	seedData := []string{
		"one",
		"two",
		"three",
		"four",
		"five",
		"six",
		"seven",
		"eight",
		"nine",
		"ten",
	}

	env := NewLocalEnv()
	replacePositionalParams(env, seedData...)

	expectedResult := "0"

	// ----------------------------------------------------------------
	// perform the change

	env.ResetPositionalParams()
	actualResult := env.Getenv("$#")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

// ================================================================
//
// setPositionalParams()
//
// ----------------------------------------------------------------

func TestSetPositionalParamsSetsThePositionalParams(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	testData := []string{
		"one",
		"two",
		"three",
		"four",
		"five",
		"six",
		"seven",
		"eight",
		"nine",
		"ten",
	}

	env := NewLocalEnv()

	expectedLen := 10
	expectedResult := []string{
		"$1=one",
		"$2=two",
		"$3=three",
		"$4=four",
		"$5=five",
		"$6=six",
		"$7=seven",
		"$8=eight",
		"$9=nine",
		"$10=ten",
		"$#=10",
	}

	// ----------------------------------------------------------------
	// perform the change

	actualLen := setPositionalParams(env, testData...)
	actualResult := env.Environ()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedLen, actualLen)
	assert.Equal(t, expectedResult, actualResult)
}

func TestSetPositionalParamsPreservesExistingPositionalParams(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	seedData := []string{
		"one",
		"two",
		"three",
		"four",
		"five",
		"six",
		"seven",
		"eight",
		"nine",
		"ten",
	}

	testData := []string{
		"new one",
		"new two",
		"new three",
	}

	env := NewLocalEnv()
	env.SetPositionalParams(seedData...)

	expectedLen := 10
	expectedResult := []string{
		"$1=new one",
		"$2=new two",
		"$3=new three",
		"$4=four",
		"$5=five",
		"$6=six",
		"$7=seven",
		"$8=eight",
		"$9=nine",
		"$10=ten",
		"$#=10",
	}

	// ----------------------------------------------------------------
	// perform the change

	actualLen := setPositionalParams(env, testData...)
	actualResult := env.Environ()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedLen, actualLen)
	assert.Equal(t, expectedResult, actualResult)
}

func TestSetPositionalParamsUpdatesDollarHash(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	seedData := []string{
		"one",
		"two",
		"three",
	}

	testData := []string{
		"new one",
		"new two",
		"new three",
		"new four",
		"new five",
	}

	env := NewLocalEnv()
	env.SetPositionalParams(seedData...)

	origHash := env.Getenv("$#")
	assert.Equal(t, "3", origHash)

	expectedResult := "5"

	// ----------------------------------------------------------------
	// perform the change

	setPositionalParams(env, testData...)
	actualResult := env.Getenv("$#")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

// ================================================================
//
// shiftPositionalParams()
//
// ----------------------------------------------------------------

func TestShiftPositionalParamsLeftPopsThePositionalParams(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	testData := []string{
		"one",
		"two",
		"three",
		"four",
		"five",
	}

	unit := NewLocalEnv()
	unit.SetPositionalParams(testData...)

	expectedResult := []string{
		"two",
		"three",
		"four",
		"five",
	}

	// ----------------------------------------------------------------
	// perform the change

	shiftPositionalParams(unit, 1)

	// ----------------------------------------------------------------
	// test the results

	actualResult := getPositionalParams(unit)
	assert.Equal(t, expectedResult, actualResult)
}

func TestShiftPositionalParamsSupportsAmountGreaterThanOne(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	testData := []string{
		"one",
		"two",
		"three",
		"four",
		"five",
	}

	unit := NewLocalEnv()
	unit.SetPositionalParams(testData...)

	expectedResult := []string{
		"four",
		"five",
	}

	// ----------------------------------------------------------------
	// perform the change

	shiftPositionalParams(unit, 3)

	// ----------------------------------------------------------------
	// test the results

	actualResult := getPositionalParams(unit)
	assert.Equal(t, expectedResult, actualResult)
}

func TestShiftPositionalParamsRemovesAllPositionalParamsWhenAmountGreaterThanDollarHash(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	testData := []string{
		"one",
		"two",
		"three",
		"four",
		"five",
	}

	unit := NewLocalEnv()
	unit.SetPositionalParams(testData...)

	expectedResult := []string{}

	// ----------------------------------------------------------------
	// perform the change

	shiftPositionalParams(unit, 6)

	// ----------------------------------------------------------------
	// test the results

	actualResult := getPositionalParams(unit)
	assert.Equal(t, expectedResult, actualResult)
}

func TestShiftPositionalParamsDoesNotCrashWhenEnvironmentHasNoPositionalParams(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	unit := NewLocalEnv()

	// ----------------------------------------------------------------
	// perform the change

	shiftPositionalParams(unit, 1)

	// ----------------------------------------------------------------
	// test the results

	// if we get here without crashing, we're good :)
	actualResult := getPositionalParams(unit)
	assert.Empty(t, actualResult)
}

// ================================================================
//
// updatePositionalSets()
//
// ----------------------------------------------------------------
