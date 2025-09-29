package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Foo string `env:"FOO" envDefault:"foo-default"`
	Bar string `env:"BAR" envDefault:"bar-default"`
}

type nestedStruct struct {
	Outer string `env:"OUTER" envDefault:"outer-default"`
	Inner testStruct
}

const (
	fooEnv       = "foo-env"
	barEnv       = "bar-env"
	fooDefault   = "foo-default"
	barDefault   = "bar-default"
	outerEnv     = "outer-env"
	outerDefault = "outer-default"
)

func TestReflectStructEnvValue(t *testing.T) {
	t.Setenv("FOO", fooEnv)
	t.Setenv("BAR", barEnv)
	s := &testStruct{}

	reflectStruct(s)

	assert.Equal(t, fooEnv, s.Foo)
	assert.Equal(t, barEnv, s.Bar)
}

func TestReflectStructDefaultValue(t *testing.T) {
	os.Unsetenv("FOO")
	os.Unsetenv("BAR")
	s := &testStruct{}

	reflectStruct(s)

	assert.Equal(t, fooDefault, s.Foo)
	assert.Equal(t, barDefault, s.Bar)
}

func TestReflectStructNestedStruct(t *testing.T) {
	t.Setenv("FOO", fooEnv)
	t.Setenv("BAR", barEnv)
	t.Setenv("OUTER", outerEnv)
	ns := &nestedStruct{}

	reflectStruct(ns)

	assert.Equal(t, outerEnv, ns.Outer)
	assert.Equal(t, fooEnv, ns.Inner.Foo)
	assert.Equal(t, barEnv, ns.Inner.Bar)
}

func TestReflectStructNoEnvTags(t *testing.T) {
	type noEnv struct {
		A string
		B int
	}
	s := &noEnv{A: "a", B: 1}

	assert.NotPanics(t, func() { reflectStruct(s) })

	assert.Equal(t, "a", s.A)
	assert.Equal(t, 1, s.B)
}
