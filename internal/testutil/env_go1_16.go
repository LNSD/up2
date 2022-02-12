//go:build !go1.17
// +build !go1.17

package testutil

import (
	"os"
	"testing"
)

// Based on https://github.com/spf13/viper/blob/master/internal/testutil/env_go1_16.go
// Licensed under the MIT license
// Copyright (c) 2014 Steve Francia

// Setenv sets an environment variable to a temporary value for the
// duration of the test.
//
// At the end of the test (see "Deferred execution" in the package docs), the
// environment variable is returned to its original value.
func Setenv(t *testing.T, name, val string) {
	setenv(t, name, val, true)
}

// setenv sets or unsets an environment variable to a temporary value for the
// duration of the test
func setenv(t *testing.T, name, val string, valOK bool) {
	oldVal, oldOK := os.LookupEnv(name)
	if valOK {
		os.Setenv(name, val)
	} else {
		os.Unsetenv(name)
	}
	t.Cleanup(func() {
		if oldOK {
			os.Setenv(name, oldVal)
		} else {
			os.Unsetenv(name)
		}
	})
}
