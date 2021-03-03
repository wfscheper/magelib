package magelib

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_envBool(t *testing.T) {
	teardown := setEnv("MAGELIB_FOO", "true")
	defer teardown()

	assert.Equal(t, true, envBool("foo", false))
}

func Test_envBool_invalid(t *testing.T) {
	teardown := setEnv("MAGELIB_FOO", "foo")
	defer teardown()

	assert.Equal(t, true, envBool("foo", true))
	assert.Equal(t, false, envBool("foo", false))
}

func Test_envBool_default(t *testing.T) {
	teardown := unsetEnv("MAGELIB_FOO")
	defer teardown()

	assert.Equal(t, false, envBool("foo", false))
}

func Test_envString(t *testing.T) {
	teardown := setEnv("MAGELIB_FOO", "foo")
	defer teardown()

	assert.Equal(t, "foo", envString("foo", ""))
}

func Test_envString_default(t *testing.T) {
	teardown := unsetEnv("MAGELIB_FOO")
	defer teardown()

	assert.Equal(t, "missing", envString("foo", "missing"))
}

func setEnv(key, value string) func() {
	old, ok := os.LookupEnv(key)
	os.Setenv(key, value)
	if ok {
		return func() { os.Setenv(key, old) }
	}

	return func() { os.Unsetenv(key) }
}

func unsetEnv(key string) func() {
	if old, ok := os.LookupEnv(key); ok {
		os.Unsetenv(key)
		return func() { os.Setenv(key, old) }
	}

	return func() {}
}
