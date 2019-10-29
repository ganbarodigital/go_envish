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
	"os/user"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEnvReturnsAnEmptyEnvironmentStore(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	testData := "this is my value"

	os.Setenv(testKey, testData)

	// clean up after ourselves
	defer os.Unsetenv(testKey)

	// ----------------------------------------------------------------
	// perform the change

	env := NewEnv()
	actualResult := env.Getenv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Empty(t, actualResult)
}

func TestNewEnvRunsAnySuppliedOptionFunctions(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	expectedResult := "this is my value"

	assert.Empty(t, os.Getenv(testKey))

	op := func(e *Env) {
		e.Setenv(testKey, expectedResult)
	}

	// ----------------------------------------------------------------
	// perform the change

	env := NewEnv(op)
	actualResult := env.Getenv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestEnvGetenvReturnsFromTheEnvNotTheProgramEnv(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	testData := "this is test data"
	expectedResult := "this is my value"

	os.Setenv(testKey, testData)

	// clean up after ourselves
	defer os.Unsetenv(testKey)

	env := NewEnv()
	env.Setenv(testKey, expectedResult)

	// now remove this from the program's environment
	os.Unsetenv(testKey)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.Getenv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestEnvGetenvCopesWithEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	expectedResult := ""

	env := Env{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.Getenv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestEnvGetenvCopesWithNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	expectedResult := ""

	var env *Env = nil

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.Getenv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestEnvGetenvSetenvLookupEnvSupportDollarVars(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "$#"
	expectedResult := "5"

	env := NewEnv()

	// ----------------------------------------------------------------
	// perform the change

	env.Setenv(testKey, expectedResult)
	actualResult1 := env.Getenv(testKey)
	actualResult2, ok := env.LookupEnv(testKey)
	actualEnviron := env.Environ()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult1)
	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult2)
	assert.Equal(t, []string{"$#=5"}, actualEnviron)
}

func TestEnvSetenvDoesNotChangeTheProgramEnv(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	expectedResult := "this is my value"

	env := NewEnv()

	// make sure this key does not exist in the program environment
	os.Unsetenv(testKey)

	// ----------------------------------------------------------------
	// perform the change

	err := env.Setenv(testKey, expectedResult)
	envResult := os.Getenv(testKey)
	actualResult := env.Getenv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Empty(t, envResult)
	assert.Equal(t, expectedResult, actualResult)
}

func TestEnvSetenvReturnsErrorForZeroLengthKey(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := ""
	testData := "this is a test"
	env := NewEnv()

	// ----------------------------------------------------------------
	// perform the change

	ok := env.Setenv(testKey, testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Error(t, ok)
}

func TestEnvSetenvReturnsErrorForKeyThatOnlyHasWhitespace(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "     "
	testData := "this is a test"
	env := NewEnv()

	// ----------------------------------------------------------------
	// perform the change

	ok := env.Setenv(testKey, testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Error(t, ok)
}

func TestEnvSetenvCopesWithEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	testData := "hello world"

	env := Env{}

	// ----------------------------------------------------------------
	// perform the change

	err := env.Setenv(testKey, testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
}

func TestEnvSetenvCopesWithNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	testData := "hello world"

	var env *Env = nil

	// ----------------------------------------------------------------
	// perform the change

	err := env.Setenv(testKey, testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Error(t, err)
}

func TestEnvClearenvDeletesAllVariables(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	testData := "this is my value"

	env := NewEnv()
	env.Setenv(testKey, testData)

	// ----------------------------------------------------------------
	// perform the change

	env.Clearenv()

	// ----------------------------------------------------------------
	// test the results

	assert.Empty(t, env.Environ())
	assert.Empty(t, env.Getenv(testKey))
}

func TestEnvClearenvCopesWithEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := Env{}

	// ----------------------------------------------------------------
	// perform the change

	env.Clearenv()

	// ----------------------------------------------------------------
	// test the results

}

func TestEnvClearenvCopesWithNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	var env *Env = nil

	// ----------------------------------------------------------------
	// perform the change

	env.Clearenv()

	// ----------------------------------------------------------------
	// test the results

}

func TestEnvLookupEnvReturnsTrueIfTheVariableExists(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	expectedResult := "this is my value"

	env := NewEnv()
	env.Setenv(testKey, expectedResult)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := env.LookupEnv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.True(t, ok)
}

func TestEnvLookupEnvReturnsFalseIfTheVariableDoesNotExist(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	expectedResult := ""

	env := NewEnv()
	env.Unsetenv(testKey)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := env.LookupEnv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.False(t, ok)
}

func TestEnvLookupEnvCopesWithEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	expectedResult := ""

	env := Env{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := env.LookupEnv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.False(t, ok)
}

func TestEnvLookupEnvCopesWithNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	expectedResult := ""

	var env *Env = nil

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := env.LookupEnv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.False(t, ok)
}

func TestEnvUnsetenvDeletesAVariable(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	testData := "this is a test"
	expectedResult := ""

	env := NewEnv()
	env.Setenv(testKey, testData)

	origLen := env.Length()

	// ----------------------------------------------------------------
	// perform the change

	env.Unsetenv(testKey)

	actualResult, ok := env.LookupEnv(testKey)
	actualEnviron := env.Environ()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.False(t, ok)
	assert.Equal(t, origLen-1, len(actualEnviron))

	// make sure it isn't in the environ too
	prefix := testKey + "="
	for _, pair := range actualEnviron {
		assert.False(t, strings.HasPrefix(pair, prefix))
	}
}

func TestEnvUnsetenvDoesNotChangeProgramEnviron(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	testData := "this is a test"
	expectedResult := testData

	os.Setenv(testKey, testData)

	// clean up after ourselves
	defer os.Unsetenv(testKey)

	env := NewEnv()

	// ----------------------------------------------------------------
	// perform the change

	env.Unsetenv(testKey)

	actualResult, ok := os.LookupEnv(testKey)

	// ----------------------------------------------------------------
	// test the results

	// still in the environment
	assert.Equal(t, expectedResult, actualResult)
	assert.True(t, ok)

	// but gone from our Env
	assert.Equal(t, "", env.Getenv(testKey))
}

func TestEnvUnsetenvSupportsDollarVars(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "$#"
	expectedResult := ""

	env := NewEnv()
	env.Setenv(testKey, expectedResult)

	// ----------------------------------------------------------------
	// perform the change

	env.Unsetenv(testKey)
	actualResult, ok := env.LookupEnv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.False(t, ok)
}

func TestEnvUnsetenvCopesWithEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"

	env := Env{}

	// ----------------------------------------------------------------
	// perform the change

	env.Unsetenv(testKey)

	// ----------------------------------------------------------------
	// test the results
}

func TestEnvUnsetenvCopesWithNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"

	var env *Env = nil

	// ----------------------------------------------------------------
	// perform the change

	env.Unsetenv(testKey)

	// ----------------------------------------------------------------
	// test the results

}

func TestEnvEntriesFromProgramEnvironmentCanBeUpdated(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	testData1 := "this is a test"
	testData2 := "this is another test"
	expectedResult := testData2

	os.Setenv(testKey, testData1)

	// clean up after ourselves
	defer os.Unsetenv(testKey)

	env := NewEnv(CopyProgramEnv)

	// ----------------------------------------------------------------
	// perform the change

	env.Setenv(testKey, testData2)

	actualResult, ok := env.LookupEnv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.True(t, ok)
}

func TestEnvUpdatedEntriesCanBeUnset(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey1 := "TestNewEnv1"
	testKey2 := "TestNewEnv2"
	testData1 := "this is a test"
	testData2 := "this is another test"

	env := NewEnv(CopyProgramEnv)

	env.Setenv(testKey1, testData1)
	env.Setenv(testKey2, testData1)

	env.Setenv(testKey1, testData2)
	testValue, ok := env.LookupEnv(testKey1)
	assert.Equal(t, testData2, testValue)
	assert.True(t, ok)

	// ----------------------------------------------------------------
	// perform the change

	env.Unsetenv(testKey1)

	actualResult, ok := env.LookupEnv(testKey1)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Empty(t, actualResult)

	// prove the 2nd entry hasn't been lost
	testValue, ok = env.LookupEnv(testKey2)
	assert.True(t, ok)
	assert.Equal(t, testData1, testValue)
}

func TestEnvEnvironCopesWithEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := Env{}
	expectedResult := []string(nil)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.Environ()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestEnvEnvironCopesWithNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	var env *Env = nil
	expectedResult := []string{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.Environ()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)

}

func TestEnvExpandCopesWithNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	var env *Env = nil

	// ----------------------------------------------------------------
	// perform the change

	env.Expand("hello ${HOME}")

	// ----------------------------------------------------------------
	// test the results

	// as long as it didn't crash, we're good
}

func TestEnvExpandCopesWithEmptyStruct(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var env Env

	// ----------------------------------------------------------------
	// perform the change

	env.Expand("hello ${HOME}")

	// ----------------------------------------------------------------
	// test the results

	// as long as it didn't crash, we're good
}

func TestEnvExpandUsesEntriesInTheTemporaryEnvironment(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestSequenceKey"
	testValue1 := "this is a test"
	testValue2 := "this is another test"
	os.Setenv(testKey, testValue1)

	// clean up after ourselves
	defer os.Unsetenv(testKey)

	expectedResult := "hello this is another test"

	env := NewEnv()
	env.Setenv(testKey, testValue2)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.Expand("hello ${TestSequenceKey}")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestEnvLengthCopesWithEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := Env{}
	expectedResult := 0

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.Length()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestEnvLengthCopesWithNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	var env *Env = nil
	expectedResult := 0

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.Length()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)

}

func TestEnvMatchVarNamesReturnsAnEmptyListWhenEnvIsEmpty(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testPrefix := "TestEnvMatchVarNames"
	env := NewEnv()

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.MatchVarNames(testPrefix)

	// ----------------------------------------------------------------
	// test the results

	assert.Empty(t, actualResult)
}

func TestEnvMatchVarNamesIsCaseSensitive(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testPrefix := "TestMatchVarName"
	env := NewEnv()
	env.Setenv("testnewenv", "dummy")
	env.Setenv(testPrefix+"sOkay", "dummy")
	env.Setenv(testPrefix+"sAnotherSuffix", "dummy")

	expectedResult := []string{
		"TestMatchVarNamesOkay",
		"TestMatchVarNamesAnotherSuffix",
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.MatchVarNames(testPrefix)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestEnvMatchVarNamesOnlyMatchesPrefixes(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testPrefix := "TestMatchVarName"
	env := NewEnv()
	env.Setenv(testPrefix+"sOkay", "dummy")
	env.Setenv("Dummy"+testPrefix+"sAnotherSuffix", "dummy")

	expectedResult := []string{
		"TestMatchVarNamesOkay",
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.MatchVarNames(testPrefix)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestEnvMatchVarNamesMatchesIfKeyEqualsPrefix(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testPrefix := "TestMatchVarNames"
	env := NewEnv()
	env.Setenv("testnewenv", "dummy")
	env.Setenv(testPrefix, "dummy")
	env.Setenv(testPrefix+"Okay", "dummy")

	expectedResult := []string{
		"TestMatchVarNames",
		"TestMatchVarNamesOkay",
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.MatchVarNames(testPrefix)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestEnvMatchVarNamesCopesWithNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testPrefix := "TestMatchVarNames"
	var env *Env = nil

	expectedResult := []string{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.MatchVarNames(testPrefix)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestEnvMatchVarNamesCopesWithEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testPrefix := "TestMatchVarNames"
	env := Env{}

	expectedResult := []string{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.MatchVarNames(testPrefix)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestEnvLookupHomeDirReturnsCurrentUserHomeDir(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	expectedResult, err := os.UserHomeDir()
	assert.Nil(t, err)

	env := NewEnv()

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := env.LookupHomeDir("")

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestEnvLookupHomeDirReturnsRootUserHomeDir(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	details, err := user.Lookup("root")
	assert.Nil(t, err)
	expectedResult := details.HomeDir

	env := NewEnv()

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := env.LookupHomeDir("root")

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestEnvLookupHomeDirReturnsFalseIfUserDoesNotExist(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := NewEnv()
	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := env.LookupHomeDir("this user does not exist")

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}
