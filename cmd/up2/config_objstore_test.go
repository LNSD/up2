package main

import (
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
	"testing"
	"upload-presigned-url-provider/internal/testutil"
)

func TestLoadMinioObjectStoreConfig(t *testing.T) {
	// Given
	config := new(ObjectStoreConfig)
	secretAccessKey := faker.New().Hash().MD5()
	testutil.Setenv(t, "UP2_MINIO_SECRET_ACCESS_KEY", secretAccessKey)

	// When
	err := config.LoadConfig()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 300, config.DefaultExpiration)
	assert.NotNil(t, config.Minio)
	assert.Equal(t, secretAccessKey, config.Minio.SecretAccessKey)
}

func TestValidateIncompleteObjectStoreConfig(t *testing.T) {
	// Given
	config := new(ObjectStoreConfig)

	// When
	err := config.LoadAndValidateConfig()

	// Then
	assert.Error(t, err)
}

func TestValidateCorrectObjectStoreConfig(t *testing.T) {
	// Given
	config := new(ObjectStoreConfig)
	secretAccessKey := faker.New().Hash().MD5()
	testutil.Setenv(t, "UP2_AWS_SECRET_ACCESS_KEY", secretAccessKey)

	// When
	err := config.LoadAndValidateConfig()

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, config.Aws)
	assert.Equal(t, secretAccessKey, config.Aws.SecretAccessKey)
}
