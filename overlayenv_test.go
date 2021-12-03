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
	"os/user"
	"testing"

	envish "github.com/ganbarodigital/go_envish/v3"
	"github.com/stretchr/testify/assert"
)

// ================================================================
//
// Constructors
//
// ----------------------------------------------------------------

func TestNewOverlayEnvReturnsStackOfEnvironments(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	localEnv := envish.NewLocalEnv()
	progEnv := envish.NewProgramEnv()

	// ----------------------------------------------------------------
	// perform the change

	stack := envish.NewOverlayEnv(localEnv, progEnv)

	// ----------------------------------------------------------------
	// test the results

	stack0, _ := stack.GetEnvByID(0)
	stack1, _ := stack.GetEnvByID(1)
	assert.Same(t, localEnv, stack0)
	assert.Same(t, progEnv, stack1)
}

// ================================================================
//
// Interface compatibility
//
// ----------------------------------------------------------------

func TestOverlayEnvImplementsReader(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	localEnv := envish.NewLocalEnv()
	progEnv := envish.NewProgramEnv()
	unit := envish.NewOverlayEnv(localEnv, progEnv)

	var i interface{} = unit

	// ----------------------------------------------------------------
	// perform the change

	_, ok := i.(envish.Reader)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
}

func TestOverlayEnvImplementsWriter(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	localEnv := envish.NewLocalEnv()
	progEnv := envish.NewProgramEnv()
	unit := envish.NewOverlayEnv(localEnv, progEnv)

	var i interface{} = unit

	// ----------------------------------------------------------------
	// perform the change

	_, ok := i.(envish.Writer)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
}

func TestOverlayEnvImplementsReaderWriter(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	localEnv := envish.NewLocalEnv()
	progEnv := envish.NewProgramEnv()
	unit := envish.NewOverlayEnv(localEnv, progEnv)

	var i interface{} = unit

	// ----------------------------------------------------------------
	// perform the change

	_, ok := i.(envish.ReaderWriter)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
}

func TestOverlayEnvImplementsExpander(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	localEnv := envish.NewLocalEnv()
	progEnv := envish.NewProgramEnv()
	unit := envish.NewOverlayEnv(localEnv, progEnv)

	var i interface{} = unit

	// ----------------------------------------------------------------
	// perform the change

	_, ok := i.(envish.Expander)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
}

func TestOverlayEnvImplementsShellEnv(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	localEnv := envish.NewLocalEnv()
	progEnv := envish.NewProgramEnv()
	unit := envish.NewOverlayEnv(localEnv, progEnv)

	var i interface{} = unit

	// ----------------------------------------------------------------
	// perform the change

	_, ok := i.(envish.Writer)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
}

// ================================================================
//
// Unique Methods
//
// ----------------------------------------------------------------

func TestOverlayEnvGetEnvByIDReturnsFromTheStack(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	localEnv := envish.NewLocalEnv()
	progEnv := envish.NewProgramEnv()

	stack := envish.NewOverlayEnv(localEnv, progEnv)

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

func TestOverlayEnvGetEnvByIDReturnsNilIfIndexTooLarge(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	localEnv := envish.NewLocalEnv()
	progEnv := envish.NewProgramEnv()

	stack := envish.NewOverlayEnv(localEnv, progEnv)

	// ----------------------------------------------------------------
	// perform the change

	stack2, ok2 := stack.GetEnvByID(2)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok2)
	assert.Nil(t, stack2)
}

func TestOverlayEnvGetEnvByIDReturnsNilIfIndexTooSmall(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	localEnv := envish.NewLocalEnv()
	progEnv := envish.NewProgramEnv()

	stack := envish.NewOverlayEnv(localEnv, progEnv)

	// ----------------------------------------------------------------
	// perform the change

	stack2, ok2 := stack.GetEnvByID(-2)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok2)
	assert.Nil(t, stack2)
}

func TestOverlayEnvGetEnvByIDReturnsCopesWithNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	var stack *envish.OverlayEnv = nil

	// ----------------------------------------------------------------
	// perform the change

	stack2, ok2 := stack.GetEnvByID(0)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok2)
	assert.Nil(t, stack2)
}

func TestOverlayEnvGetEnvByIDReturnsCopesWithEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	stack := envish.OverlayEnv{}

	// ----------------------------------------------------------------
	// perform the change

	stack2, ok2 := stack.GetEnvByID(0)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok2)
	assert.Nil(t, stack2)
}

// ================================================================
//
// Reader Compatibility
//
// ----------------------------------------------------------------

func TestOverlayEnvironReturnsAllVariablesFromEnvsThatAreExporting(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we don't use a program environment here because its contents are
	// unpredictable
	env1 := envish.NewLocalEnv(envish.SetAsExporter)
	env1.Setenv("PARAM1.1", "hello")
	env1.Setenv("PARAM1.2", "world")
	env2 := envish.NewLocalEnv(envish.SetAsExporter)
	env2.Setenv("PARAM2.1", "trout")
	env2.Setenv("PARAM2.2", "haddock")

	expectedResult := []string{
		"PARAM1.1=hello",
		"PARAM1.2=world",
		"PARAM2.1=trout",
		"PARAM2.2=haddock",
	}
	stack := envish.NewOverlayEnv(env1, env2)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := stack.Environ()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvironSkipsOverEnvsThatAreNotExporting(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we don't use a program environment here because its contents are
	// unpredictable
	env1 := envish.NewLocalEnv()
	env1.Setenv("PARAM1.1", "hello")
	env1.Setenv("PARAM1.2", "world")
	env2 := envish.NewLocalEnv(envish.SetAsExporter)
	env2.Setenv("PARAM2.1", "trout")
	env2.Setenv("PARAM2.2", "haddock")

	expectedResult := []string{
		"PARAM2.1=trout",
		"PARAM2.2=haddock",
	}
	stack := envish.NewOverlayEnv(env1, env2)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := stack.Environ()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvironGivesPrecidenceToEarlierStackEntries(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we don't use a program environment here because its contents are
	// unpredictable
	env1 := envish.NewLocalEnv(envish.SetAsExporter)
	env1.Setenv("PARAM1.1", "hello")
	env1.Setenv("PARAM1.2", "world")
	env2 := envish.NewLocalEnv(envish.SetAsExporter)
	env2.Setenv("PARAM1.1", "trout")
	env2.Setenv("PARAM2.1", "haddock")

	expectedResult := []string{
		"PARAM1.1=hello",
		"PARAM1.2=world",
		"PARAM2.1=haddock",
	}
	stack := envish.NewOverlayEnv(env1, env2)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := stack.Environ()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvironReturnsEmptyListIfNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	expectedResult := []string{}
	var stack *envish.OverlayEnv = nil

	// ----------------------------------------------------------------
	// perform the change

	actualResult := stack.Environ()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvironReturnsEmptyListIfEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	expectedResult := []string{}
	stack := envish.OverlayEnv{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult := stack.Environ()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvGetenvSearchesTheStack(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we don't use a program environment here because its contents are
	// unpredictable
	env1 := envish.NewLocalEnv(envish.SetAsExporter)
	env1.Setenv("PARAM1.1", "hello")
	env1.Setenv("PARAM1.2", "world")
	env2 := envish.NewLocalEnv(envish.SetAsExporter)
	env2.Setenv("PARAM2.1", "trout")
	env2.Setenv("PARAM2.2", "haddock")

	expectedResult := "haddock"
	stack := envish.NewOverlayEnv(env1, env2)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := stack.Getenv("PARAM2.2")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvGetenvGivesPrecidenceToEarlierStackEntries(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we don't use a program environment here because its contents are
	// unpredictable
	env1 := envish.NewLocalEnv(envish.SetAsExporter)
	env1.Setenv("PARAM1.1", "hello")
	env1.Setenv("PARAM1.2", "world")
	env2 := envish.NewLocalEnv(envish.SetAsExporter)
	env2.Setenv("PARAM1.1", "trout")
	env2.Setenv("PARAM1.2", "haddock")

	expectedResult := "world"
	stack := envish.NewOverlayEnv(env1, env2)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := stack.Getenv("PARAM1.2")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvGetenvReturnsEmptyStringIfVariableNotFound(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we don't use a program environment here because its contents are
	// unpredictable
	env1 := envish.NewLocalEnv(envish.SetAsExporter)
	env1.Setenv("PARAM1.1", "hello")
	env1.Setenv("PARAM1.2", "world")
	env2 := envish.NewLocalEnv(envish.SetAsExporter)
	env2.Setenv("PARAM1.1", "trout")
	env2.Setenv("PARAM1.2", "haddock")

	expectedResult := ""
	stack := envish.NewOverlayEnv(env1, env2)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := stack.Getenv("PARAM")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvGetenvReturnsEmptyStringIfNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""
	var stack *envish.OverlayEnv = nil

	// ----------------------------------------------------------------
	// perform the change

	actualResult := stack.Getenv("PARAM")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvGetenvReturnsEmptyStringIfEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""
	stack := envish.OverlayEnv{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult := stack.Getenv("PARAM")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvIsExporterReturnsTrueIfAnyEnvIsAnExporter(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	localEnv := envish.NewLocalEnv()
	progEnv := envish.NewProgramEnv()

	stack := envish.NewOverlayEnv(localEnv, progEnv)

	// ----------------------------------------------------------------
	// perform the change

	ok := stack.IsExporter()

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
}

func TestOverlayEnvIsExporterReturnsFalseIfAllEnvsAreNotExporters(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	localEnv := envish.NewLocalEnv()
	progEnv := envish.NewLocalEnv()

	stack := envish.NewOverlayEnv(localEnv, progEnv)

	// ----------------------------------------------------------------
	// perform the change

	ok := stack.IsExporter()

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
}

func TestOverlayEnvIsExporterReturnsFalseIfNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	var stack *envish.OverlayEnv = nil

	// ----------------------------------------------------------------
	// perform the change

	ok := stack.IsExporter()

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
}

func TestOverlayEnvIsExporterReturnsFalseIfEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	stack := envish.OverlayEnv{}

	// ----------------------------------------------------------------
	// perform the change

	ok := stack.IsExporter()

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
}

func TestOverlayEnvLookupEnvSearchesTheStack(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we don't use a program environment here because its contents are
	// unpredictable
	env1 := envish.NewLocalEnv(envish.SetAsExporter)
	env1.Setenv("PARAM1.1", "hello")
	env1.Setenv("PARAM1.2", "world")
	env2 := envish.NewLocalEnv(envish.SetAsExporter)
	env2.Setenv("PARAM2.1", "trout")
	env2.Setenv("PARAM2.2", "haddock")

	expectedResult := "haddock"
	stack := envish.NewOverlayEnv(env1, env2)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := stack.LookupEnv("PARAM2.2")

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvLookupEnvGivesPrecidenceToEarlierStackEntries(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we don't use a program environment here because its contents are
	// unpredictable
	env1 := envish.NewLocalEnv(envish.SetAsExporter)
	env1.Setenv("PARAM1.1", "hello")
	env1.Setenv("PARAM1.2", "world")
	env2 := envish.NewLocalEnv(envish.SetAsExporter)
	env2.Setenv("PARAM1.1", "trout")
	env2.Setenv("PARAM1.2", "haddock")

	expectedResult := "world"
	stack := envish.NewOverlayEnv(env1, env2)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := stack.LookupEnv("PARAM1.2")

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvLookupEnvReturnsEmptyStringIfVariableNotFound(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we don't use a program environment here because its contents are
	// unpredictable
	env1 := envish.NewLocalEnv(envish.SetAsExporter)
	env1.Setenv("PARAM1.1", "hello")
	env1.Setenv("PARAM1.2", "world")
	env2 := envish.NewLocalEnv(envish.SetAsExporter)
	env2.Setenv("PARAM1.1", "trout")
	env2.Setenv("PARAM1.2", "haddock")

	expectedResult := ""
	stack := envish.NewOverlayEnv(env1, env2)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := stack.LookupEnv("PARAM")

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvLookupEnvReturnsEmptyStringIfNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""
	var stack *envish.OverlayEnv = nil

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := stack.LookupEnv("PARAM")

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvLookupEnvReturnsEmptyStringIfEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""
	stack := envish.OverlayEnv{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := stack.LookupEnv("PARAM")

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvMatchVarNamesSearchesTheStack(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we don't use a program environment here because its contents are
	// unpredictable
	env1 := envish.NewLocalEnv(envish.SetAsExporter)
	env1.Setenv("PARAM1.1", "hello")
	env1.Setenv("PARAM1.2", "world")
	env2 := envish.NewLocalEnv(envish.SetAsExporter)
	env2.Setenv("PARAM2.1", "trout")
	env2.Setenv("PARAM2.2", "haddock")

	expectedResult := []string{
		"PARAM1.1",
		"PARAM1.2",
		"PARAM2.1",
		"PARAM2.2",
	}
	stack := envish.NewOverlayEnv(env1, env2)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := stack.MatchVarNames("PARAM")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvMatchVarNamesGivesPrecidenceToEarlierStackEntries(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we don't use a program environment here because its contents are
	// unpredictable
	env1 := envish.NewLocalEnv(envish.SetAsExporter)
	env1.Setenv("PARAM1.1", "hello")
	env1.Setenv("PARAM1.2", "world")
	env2 := envish.NewLocalEnv(envish.SetAsExporter)
	env2.Setenv("PARAM1.1", "trout")
	env2.Setenv("PARAM1.2", "haddock")

	expectedResult := []string{
		"PARAM1.1",
		"PARAM1.2",
	}
	stack := envish.NewOverlayEnv(env1, env2)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := stack.MatchVarNames("PARAM")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvMatchVarNamesReturnsEmptyListIfNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	expectedResult := []string{}
	var stack *envish.OverlayEnv = nil

	// ----------------------------------------------------------------
	// perform the change

	actualResult := stack.MatchVarNames("PARAM")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvMatchVarNamesReturnsEmptyListIfEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	expectedResult := []string{}
	stack := envish.OverlayEnv{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult := stack.MatchVarNames("PARAM")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

// ================================================================
//
// Writer Compatibility
//
// ----------------------------------------------------------------

func TestOverlayEnvClearenvEmptiesAllEnvironments(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	localEnv := envish.NewLocalEnv()
	progEnv := envish.NewProgramEnv()
	stack := envish.NewOverlayEnv(localEnv, progEnv)

	// we'll need to put the program's environment back afterwards!
	origEnviron := os.Environ()
	defer progEnv.RestoreEnvironment(origEnviron)

	// ----------------------------------------------------------------
	// perform the change

	stack.Clearenv()

	// ----------------------------------------------------------------
	// test the results

	localEnviron := localEnv.Environ()
	progEnviron := progEnv.Environ()

	assert.Empty(t, localEnviron)
	assert.Empty(t, progEnviron)
}

func TestOverlayEnvClearenvDoesNothingIfNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	var stack *envish.OverlayEnv = nil

	// ----------------------------------------------------------------
	// perform the change

	stack.Clearenv()

	// ----------------------------------------------------------------
	// test the results

	osEnviron := os.Environ()

	assert.NotEmpty(t, osEnviron)
}

func TestOverlayEnvClearenvDoesNothingIfEmptyStruct(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	stack := envish.OverlayEnv{}

	// ----------------------------------------------------------------
	// perform the change

	stack.Clearenv()

	// ----------------------------------------------------------------
	// test the results

	osEnviron := os.Environ()

	assert.NotEmpty(t, osEnviron)
}

func TestOverlayEnvSetenvUpdatesExistingEntriesInStackOrder(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we don't use a program environment here because its contents are
	// unpredictable
	env1 := envish.NewLocalEnv(envish.SetAsExporter)
	env1.Setenv("PARAM1.1", "hello")
	env1.Setenv("PARAM1.2", "world")
	env2 := envish.NewLocalEnv(envish.SetAsExporter)
	env2.Setenv("PARAM1.1", "trout")
	env2.Setenv("PARAM1.2", "haddock")

	expectedResult1 := "there"
	expectedResult2 := "haddock"
	stack := envish.NewOverlayEnv(env1, env2)

	// ----------------------------------------------------------------
	// perform the change

	stack.Setenv("PARAM1.2", expectedResult1)
	actualResult1 := env1.Getenv("PARAM1.2")
	actualResult2 := env2.Getenv("PARAM1.2")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult1, actualResult1)
	assert.Equal(t, expectedResult2, actualResult2)
}

func TestOverlayEnvSetenvCreatesNewEntriesInFirstStackMember(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we don't use a program environment here because its contents are
	// unpredictable
	env1 := envish.NewLocalEnv(envish.SetAsExporter)
	env1.Setenv("PARAM1.1", "hello")
	env1.Setenv("PARAM1.2", "world")
	env2 := envish.NewLocalEnv(envish.SetAsExporter)
	env2.Setenv("PARAM1.1", "trout")
	env2.Setenv("PARAM1.2", "haddock")

	expectedResult := "there"
	stack := envish.NewOverlayEnv(env1, env2)

	// ----------------------------------------------------------------
	// perform the change

	stack.Setenv("PARAM1.3", expectedResult)
	actualResult := env1.Getenv("PARAM1.3")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvSetenvReturnsErrorIfNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	var stack *envish.OverlayEnv = nil

	expectedResult := "nil pointer to environment store passed to OverlayEnv.Setenv"

	// ----------------------------------------------------------------
	// perform the change

	actualResult := stack.Setenv("PARAM1.3", "DOES NOT MATTER")

	// ----------------------------------------------------------------
	// test the results

	assert.Error(t, actualResult)
	assert.Equal(t, expectedResult, actualResult.Error())
}

func TestOverlayEnvSetenvReturnsErrorIfEmptyStack(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	stack := envish.OverlayEnv{}

	expectedResult := "overlay env is empty; OverlayEnv.Setenv"

	// ----------------------------------------------------------------
	// perform the change

	actualResult := stack.Setenv("PARAM1.3", "DOES NOT MATTER")

	// ----------------------------------------------------------------
	// test the results

	assert.Error(t, actualResult)
	assert.Equal(t, expectedResult, actualResult.Error())
}

func TestOverlayEnvUnsetenvUpdatesEveryEntryInTheStack(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we don't use a program environment here because its contents are
	// unpredictable
	env1 := envish.NewLocalEnv(envish.SetAsExporter)
	env1.Setenv("PARAM1.1", "hello")
	env1.Setenv("PARAM1.2", "world")
	env2 := envish.NewLocalEnv(envish.SetAsExporter)
	env2.Setenv("PARAM1.1", "trout")
	env2.Setenv("PARAM1.2", "haddock")

	stack := envish.NewOverlayEnv(env1, env2)

	// ----------------------------------------------------------------
	// perform the change

	stack.Unsetenv("PARAM1.2")
	actualResult1 := env1.Getenv("PARAM1.2")
	actualResult2 := env2.Getenv("PARAM1.2")

	// ----------------------------------------------------------------
	// test the results

	assert.Empty(t, actualResult1)
	assert.Empty(t, actualResult2)
}

func TestOverlayEnvUnsetenvDoesNothingIfNilPointer(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	var stack *envish.OverlayEnv = nil

	// ----------------------------------------------------------------
	// perform the change

	assert.NotPanics(t, func() { stack.Unsetenv("PARAM") })

	// ----------------------------------------------------------------
	// test the results

	// as long as it didn't crash, all is well
}

func TestOverlayEnvUnsetenvDoesNothingIfEmptyStack(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	stack := envish.OverlayEnv{}

	// ----------------------------------------------------------------
	// perform the change

	assert.NotPanics(t, func() { stack.Unsetenv("PARAM") })

	// ----------------------------------------------------------------
	// test the results
}

// ================================================================
//
// Expander Compatibility
//
// ----------------------------------------------------------------

func TestOverlayExpandSearchesTheStack(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we don't use a program environment here because its contents are
	// unpredictable
	env1 := envish.NewLocalEnv(envish.SetAsExporter)
	env1.Setenv("PARAM1_1", "hello")
	env1.Setenv("PARAM1_2", "world")
	env2 := envish.NewLocalEnv(envish.SetAsExporter)
	env2.Setenv("PARAM2_1", "trout")
	env2.Setenv("PARAM1_2", "haddock")

	expectedResult := "hello, world trout"
	stack := envish.NewOverlayEnv(env1, env2)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := stack.Expand("${PARAM1_1}, ${PARAM1_2} ${PARAM2_1}")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayExpandReturnsOriginalStringIfExpansionFails(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we don't use a program environment here because its contents are
	// unpredictable
	env1 := envish.NewLocalEnv(envish.SetAsExporter)
	env1.Setenv("PARAM1_1", "hello")
	env1.Setenv("PARAM1_2", "world")
	env2 := envish.NewLocalEnv(envish.SetAsExporter)
	env2.Setenv("PARAM1_1", "trout")
	env2.Setenv("PARAM1_2", "haddock")

	expectedResult := "${PARAM1_1#abc[}, ${PARAM1_2}"
	stack := envish.NewOverlayEnv(env1, env2)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := stack.Expand(expectedResult)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvLookupHomeDirReturnsCurrentUserHomeDir(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	expectedResult, err := os.UserHomeDir()
	assert.Nil(t, err)

	env := envish.NewOverlayEnv()

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := env.LookupHomeDir("")

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvLookupHomeDirReturnsRootUserHomeDir(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	details, err := user.Lookup("root")
	assert.Nil(t, err)
	expectedResult := details.HomeDir

	env := envish.NewOverlayEnv()

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := env.LookupHomeDir("root")

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestOverlayEnvLookupHomeDirReturnsFalseIfUserDoesNotExist(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	env := envish.NewOverlayEnv()
	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := env.LookupHomeDir("this user does not exist")

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}
