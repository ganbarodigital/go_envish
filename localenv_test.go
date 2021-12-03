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
	"os"
	"strings"
	"testing"

	envish "github.com/ganbarodigital/go_envish/v4"
	"github.com/stretchr/testify/assert"
)

// ================================================================
//
// Constructors
//
// ----------------------------------------------------------------

func TestNewLocalEnvReturnsAnEmptyEnvironmentStore(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	testData := "this is my value"

	os.Setenv(testKey, testData)

	// clean up after ourselves
	defer os.Unsetenv(testKey)

	// ----------------------------------------------------------------
	// perform the change

	env := envish.NewLocalEnv()
	actualResult := env.Getenv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Empty(t, actualResult)
}

func TestNewLocalEnvRunsAnySuppliedOptionFunctions(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewLocalEnv"
	expectedResult := "this is my value"

	assert.Empty(t, os.Getenv(testKey))

	op := func(e *envish.LocalEnv) {
		e.Setenv(testKey, expectedResult)
	}

	// ----------------------------------------------------------------
	// perform the change

	env := envish.NewLocalEnv(op)
	actualResult := env.Getenv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

// ================================================================
//
// Interface compatibility
//
// ----------------------------------------------------------------

func TestLocalEnvImplementsReader(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	unit := envish.NewLocalEnv()
	var i interface{} = unit

	// ----------------------------------------------------------------
	// perform the change

	_, ok := i.(envish.Reader)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
}

func TestLocalEnvImplementsWriter(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	unit := envish.NewLocalEnv()
	var i interface{} = unit

	// ----------------------------------------------------------------
	// perform the change

	_, ok := i.(envish.Writer)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
}

func TestLocalEnvImplementsReaderWriter(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	unit := envish.NewLocalEnv()
	var i interface{} = unit

	// ----------------------------------------------------------------
	// perform the change

	_, ok := i.(envish.ReaderWriter)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
}

func TestLocalEnvImplementsExpander(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	unit := envish.NewLocalEnv()
	var i interface{} = unit

	// ----------------------------------------------------------------
	// perform the change

	_, ok := i.(envish.Expander)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
}

// ================================================================
//
// Reader Compatibility
//
// ----------------------------------------------------------------

func TestLocalEnvEnvironCopesWithEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := envish.LocalEnv{}
	expectedResult := []string(nil)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.Environ()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestLocalEnvEnvironCopesWithNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	var env *envish.LocalEnv = nil
	expectedResult := []string{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.Environ()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)

}

func TestLocalEnvGetenvReturnsFromTheEnvNotTheProgramEnv(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewLocalEnv"
	testData := "this is test data"
	expectedResult := "this is my value"

	os.Setenv(testKey, testData)

	// clean up after ourselves
	defer os.Unsetenv(testKey)

	env := envish.NewLocalEnv()
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

func TestLocalEnvGetenvCopesWithEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewLocalEnv"
	expectedResult := ""

	env := envish.LocalEnv{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.Getenv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestLocalEnvGetenvCopesWithNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewLocalEnv"
	expectedResult := ""

	var env *envish.LocalEnv = nil

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.Getenv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestLocalEnvIsNotAnExporterByDefault(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := envish.NewLocalEnv()
	expectedResult := false

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.IsExporter()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestLocalEnvCanBeCreatedAsExporter(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := envish.NewLocalEnv(envish.SetAsExporter)
	expectedResult := true

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.IsExporter()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestLocalEnvLookupEnvReturnsTrueIfTheVariableExists(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewLocalEnv"
	expectedResult := "this is my value"

	env := envish.NewLocalEnv()
	env.Setenv(testKey, expectedResult)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := env.LookupEnv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.True(t, ok)
}

func TestLocalEnvLookupEnvReturnsFalseIfTheVariableDoesNotExist(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewLocalEnv"
	expectedResult := ""

	env := envish.NewLocalEnv()
	env.Unsetenv(testKey)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := env.LookupEnv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.False(t, ok)
}

func TestLocalEnvLookupEnvCopesWithEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewLocalEnv"
	expectedResult := ""

	env := envish.LocalEnv{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := env.LookupEnv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.False(t, ok)
}

func TestLocalEnvLookupEnvCopesWithNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewLocalEnv"
	expectedResult := ""

	var env *envish.LocalEnv = nil

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := env.LookupEnv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.False(t, ok)
}

func TestLocalEnvMatchVarNamesReturnsAnEmptyListWhenEnvIsEmpty(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testPrefix := "TestLocalEnvMatchVarNames"
	env := envish.NewLocalEnv()

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.MatchVarNames(testPrefix)

	// ----------------------------------------------------------------
	// test the results

	assert.Empty(t, actualResult)
}

func TestLocalEnvMatchVarNamesIsCaseSensitive(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testPrefix := "TestMatchVarName"
	env := envish.NewLocalEnv()
	env.Setenv("testNewLocalEnv", "dummy")
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

func TestLocalEnvMatchVarNamesOnlyMatchesPrefixes(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testPrefix := "TestMatchVarName"
	env := envish.NewLocalEnv()
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

func TestLocalEnvMatchVarNamesMatchesIfKeyEqualsPrefix(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testPrefix := "TestMatchVarNames"
	env := envish.NewLocalEnv()
	env.Setenv("testNewLocalEnv", "dummy")
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

func TestLocalEnvMatchVarNamesCopesWithNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testPrefix := "TestMatchVarNames"
	var env *envish.LocalEnv = nil

	expectedResult := []string{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.MatchVarNames(testPrefix)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestLocalEnvMatchVarNamesCopesWithEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testPrefix := "TestMatchVarNames"
	env := envish.LocalEnv{}

	expectedResult := []string{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.MatchVarNames(testPrefix)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

// ================================================================
//
// Writer interface
//
// ----------------------------------------------------------------

func TestLocalEnvClearenvDeletesAllVariables(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewLocalEnv"
	testData := "this is my value"

	env := envish.NewLocalEnv()
	env.Setenv(testKey, testData)

	// ----------------------------------------------------------------
	// perform the change

	env.Clearenv()

	// ----------------------------------------------------------------
	// test the results

	assert.Empty(t, env.Environ())
	assert.Empty(t, env.Getenv(testKey))
}

func TestLocalEnvClearenvCopesWithEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := envish.LocalEnv{}

	// ----------------------------------------------------------------
	// perform the change

	env.Clearenv()

	// ----------------------------------------------------------------
	// test the results

}

func TestLocalEnvClearenvCopesWithNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	var env *envish.LocalEnv = nil

	// ----------------------------------------------------------------
	// perform the change

	env.Clearenv()

	// ----------------------------------------------------------------
	// test the results

}

func TestLocalEnvGetenvSetenvLookupEnvSupportDollarVars(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "$#"
	expectedResult := "5"

	env := envish.NewLocalEnv()

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

func TestLocalEnvSetenvDoesNotChangeTheProgramEnv(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewLocalEnv"
	expectedResult := "this is my value"

	env := envish.NewLocalEnv()

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

func TestLocalEnvSetenvReturnsErrorForZeroLengthKey(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := ""
	testData := "this is a test"
	env := envish.NewLocalEnv()

	// ----------------------------------------------------------------
	// perform the change

	ok := env.Setenv(testKey, testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Error(t, ok)
}

func TestLocalEnvSetenvReturnsErrorForKeyThatOnlyHasWhitespace(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "     "
	testData := "this is a test"
	env := envish.NewLocalEnv()

	// ----------------------------------------------------------------
	// perform the change

	ok := env.Setenv(testKey, testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Error(t, ok)
}

func TestLocalEnvSetenvCopesWithEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewLocalEnv"
	testData := "hello world"

	env := envish.LocalEnv{}

	// ----------------------------------------------------------------
	// perform the change

	err := env.Setenv(testKey, testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
}

func TestLocalEnvSetenvCopesWithNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewLocalEnv"
	testData := "hello world"

	var env *envish.LocalEnv = nil

	// ----------------------------------------------------------------
	// perform the change

	err := env.Setenv(testKey, testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Error(t, err)
}

func TestLocalEnvUnsetenvDeletesAVariable(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewLocalEnv"
	testData := "this is a test"
	expectedResult := ""

	env := envish.NewLocalEnv()
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

func TestLocalEnvUnsetenvDoesNotChangeProgramEnviron(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewLocalEnv"
	testData := "this is a test"
	expectedResult := testData

	os.Setenv(testKey, testData)

	// clean up after ourselves
	defer os.Unsetenv(testKey)

	env := envish.NewLocalEnv()

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

func TestLocalEnvUnsetenvSupportsDollarVars(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "$#"
	expectedResult := ""

	env := envish.NewLocalEnv()
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

func TestLocalEnvUnsetenvCopesWithEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewLocalEnv"

	env := envish.LocalEnv{}

	// ----------------------------------------------------------------
	// perform the change

	env.Unsetenv(testKey)

	// ----------------------------------------------------------------
	// test the results
}

func TestLocalEnvUnsetenvCopesWithNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewLocalEnv"

	var env *envish.LocalEnv = nil

	// ----------------------------------------------------------------
	// perform the change

	env.Unsetenv(testKey)

	// ----------------------------------------------------------------
	// test the results

}

func TestLocalEnvEntriesFromProgramEnvironmentCanBeUpdated(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewLocalEnv"
	testData1 := "this is a test"
	testData2 := "this is another test"
	expectedResult := testData2

	os.Setenv(testKey, testData1)

	// clean up after ourselves
	defer os.Unsetenv(testKey)

	env := envish.NewLocalEnv(envish.CopyProgramEnv)

	// ----------------------------------------------------------------
	// perform the change

	env.Setenv(testKey, testData2)

	actualResult, ok := env.LookupEnv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.True(t, ok)
}

func TestLocalEnvUpdatedEntriesCanBeUnset(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey1 := "TestNewLocalEnv1"
	testKey2 := "TestNewLocalEnv2"
	testData1 := "this is a test"
	testData2 := "this is another test"

	env := envish.NewLocalEnv(envish.CopyProgramEnv)

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

// ================================================================
//
// Expander interface
//
// ----------------------------------------------------------------

func TestLocalEnvExpandCopesWithNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	var env *envish.LocalEnv = nil

	// ----------------------------------------------------------------
	// perform the change

	env.Expand("hello ${HOME}")

	// ----------------------------------------------------------------
	// test the results

	// as long as it didn't crash, we're good
}

func TestLocalEnvExpandCopesWithEmptyStruct(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var env envish.LocalEnv

	// ----------------------------------------------------------------
	// perform the change

	env.Expand("hello ${HOME}")

	// ----------------------------------------------------------------
	// test the results

	// as long as it didn't crash, we're good
}

func TestLocalEnvExpandUsesEntriesInTheTemporaryEnvironment(t *testing.T) {
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

	env := envish.NewLocalEnv()
	env.Setenv(testKey, testValue2)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.Expand("hello ${TestSequenceKey}")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestLocalEnvExpandReturnsOriginalStringIfExpansionFails(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	env := envish.NewLocalEnv()

	// the search pattern is invalid, and this will trigger an error
	expectedResult := "hello ${TestSequenceKey#abc[}"

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.Expand(expectedResult)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestLocalEnvLengthCopesWithEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := envish.LocalEnv{}
	expectedResult := 0

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.Length()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestLocalEnvLengthCopesWithNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	var env *envish.LocalEnv = nil
	expectedResult := 0

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.Length()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)

}
