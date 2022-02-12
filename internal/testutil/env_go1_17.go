//go:build go1.17
// +build go1.17

package testutil

import (
	"testing"
)

// Based on https://github.com/spf13/viper/blob/master/internal/testutil/env_go1_16.go
// Licensed under the MIT license
// Copyright (c) 2014 Steve Francia

// Setenv sets an environment variable to a temporary value for the
// duration of the test.
//
// This shim can be removed once support for Go <1.17 is dropped.
func Setenv(t *testing.T, name, val string) {
	t.Helper()

	t.Setenv(name, val)
}
