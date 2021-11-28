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
	"testing"

	"github.com/stretchr/testify/assert"
)

// ================================================================
//
// Constructors
//
// ----------------------------------------------------------------

func TestNewProgramEnvReturnsAnEmptyEnvironmentStore(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// ----------------------------------------------------------------
	// perform the change

	env := NewProgramEnv()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, env)
}

func TestProgramEnvClearenvEmptiesYourProgramsEnvironment(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := NewProgramEnv()

	origEnviron := env.Environ()
	defer env.RestoreEnvironment(origEnviron)

	// ----------------------------------------------------------------
	// perform the change

	env.Clearenv()

	// ----------------------------------------------------------------
	// test the results

	assert.Empty(t, os.Environ())
}

// ================================================================
//
// Interface compatibility
//
// ----------------------------------------------------------------

func TestProgramEnvImplementsReader(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	unit := NewProgramEnv()
	var i interface{} = unit

	// ----------------------------------------------------------------
	// perform the change

	_, ok := i.(Reader)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
}

func TestProgramEnvImplementsWriter(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	unit := NewProgramEnv()
	var i interface{} = unit

	// ----------------------------------------------------------------
	// perform the change

	_, ok := i.(Writer)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
}

func TestProgramEnvImplementsReaderWriter(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	unit := NewProgramEnv()
	var i interface{} = unit

	// ----------------------------------------------------------------
	// perform the change

	_, ok := i.(ReaderWriter)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
}

func TestProgramEnvImplementsExpander(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	unit := NewProgramEnv()
	var i interface{} = unit

	// ----------------------------------------------------------------
	// perform the change

	_, ok := i.(Expander)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
}

func TestProgramEnvImplementsShellEnv(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	unit := NewProgramEnv()
	var i interface{} = unit

	// ----------------------------------------------------------------
	// perform the change

	_, ok := i.(Writer)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
}

// ================================================================
//
// Reader Compatibility
//
// ----------------------------------------------------------------

func TestProgramEnvExpandPerformsStringExpansion(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := NewProgramEnv()
	env.Setenv("PARAM1", "foo")
	expectedResult := "FOO"

	// clean up after ourselves
	defer os.Unsetenv("PARAM1")

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.Expand("${PARAM1^^}")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestProgramEnvExpandReturnsOriginalStringIfExpansionGeneratesAnError(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := NewProgramEnv()

	// we need to set this, to make sure an attempt is made to compile
	// the invalid pattern
	env.Setenv("PARAM1", "foo")

	// clean up after ourselves
	defer os.Unsetenv("PARAM1")

	expectedResult := "${PARAM1##abc[}"

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.Expand(expectedResult)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestProgramEnvGetenvRetrievesAValue(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := NewProgramEnv()
	env.Setenv("PARAM1", "foo")
	expectedResult := "foo"

	// clean up after ourselves
	defer os.Unsetenv("PARAM1")

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.Getenv("PARAM1")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestProgramEnvGetenvReturnsEmptyStringWhenVariableNotSet(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := NewProgramEnv()
	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.Getenv("PARAM1_DOES_NOT_EXIST")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestProgramEnvIsAnExporter(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := NewProgramEnv()

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.IsExporter()

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, actualResult)
}

func TestProgramEnvLookupEnvReturnsTrueWhenVarExists(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := NewProgramEnv()
	env.Setenv("PARAM1", "foo")
	expectedResult := "foo"

	// clean up after ourselves
	defer os.Unsetenv("PARAM1")

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := env.LookupEnv("PARAM1")

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestProgramEnvLookupEnvReturnsFalseWhenVarDoesNotExist(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := NewProgramEnv()
	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := env.LookupEnv("PARAM1_DOES_NOT_EXIST")

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestProgramEnvLookupHomeDirReturnsCurrentUserHomeDir(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	expectedResult, err := os.UserHomeDir()
	assert.Nil(t, err)

	env := NewProgramEnv()

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := env.LookupHomeDir("")

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestProgramEnvLookupHomeDirReturnsRootUserHomeDir(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	details, err := user.Lookup("root")
	assert.Nil(t, err)
	expectedResult := details.HomeDir

	env := NewProgramEnv()

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := env.LookupHomeDir("root")

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestProgramEnvLookupHomeDirReturnsFalseIfUserDoesNotExist(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := NewProgramEnv()
	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := env.LookupHomeDir("this user does not exist")

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestProgramEnvMatchVarNamesReturnsListOfKeys(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := NewProgramEnv()
	env.Setenv("TestProgramEnv_PARAM1", "foo")
	env.Setenv("TestProgramEnv_PARAM2", "bar")
	expectedResult := []string{
		"TestProgramEnv_PARAM1",
		"TestProgramEnv_PARAM2",
	}

	// clean up after ourselves
	defer func() {
		for _, key := range expectedResult {
			os.Unsetenv(key)
		}
	}()

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.MatchVarNames("TestProgramEnv_")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestProgramEnvSetenvUpdatesTheProgramEnvironment(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := NewProgramEnv()
	expectedResult := "foo"

	// clean up after ourselves
	defer os.Unsetenv("PARAM1")

	// ----------------------------------------------------------------
	// perform the change

	env.Setenv("PARAM1", "foo")
	actualResult, ok := os.LookupEnv("PARAM1")

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestProgramEnvUnsetenvUpdatesTheProgramEnvironment(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := NewProgramEnv()
	expectedResult := ""

	env.Setenv("PARAM1", "foo")

	// clean up after ourselves
	defer os.Unsetenv("PARAM1")

	// ----------------------------------------------------------------
	// perform the change

	env.Unsetenv("PARAM1")
	actualResult, ok := os.LookupEnv("PARAM1")

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}
