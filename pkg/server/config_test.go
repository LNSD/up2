package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigAddresses(t *testing.T) {
	/// Given
	config := Config{Host: "localhost", Port: 42}

	/// When
	address := config.Address()

	/// Then
	assert.Equal(t, "localhost:42", address)
}
