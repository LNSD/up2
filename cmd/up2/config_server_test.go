package main

import (
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"upload-presigned-url-provider/internal/testutil"
)

func TestLoadServerConfig(t *testing.T) {
	// Given
	config := new(ServerConfig)
	serverPort := int(faker.New().UInt16())
	testutil.Setenv(t, "UP2_PORT", strconv.Itoa(serverPort))

	// When
	err := config.LoadConfig()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "", config.Host)
	assert.Equal(t, serverPort, config.Port)
}

func TestValidateInvalidServerConfig(t *testing.T) {
	// Given
	config := new(ServerConfig)
	testutil.Setenv(t, "UP2_PORT", strconv.Itoa(1_000_000))

	// When
	err := config.LoadAndValidateConfig()

	// Then
	assert.Error(t, err)
}

func TestValidateCorrectServerConfig(t *testing.T) {
	// Given
	config := new(ServerConfig)
	serverPort := int(faker.New().UInt16())
	testutil.Setenv(t, "UP2_PORT", strconv.Itoa(serverPort))

	// When
	err := config.LoadAndValidateConfig()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, serverPort, config.Port)
}
