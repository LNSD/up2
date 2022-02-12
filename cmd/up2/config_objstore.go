package main

import (
	"github.com/spf13/viper"
)

type ObjectStoreConfig struct {
	// Default expiration time for the upload URL in seconds. Default: 300
	DefaultExpiration int `validate:"min=1" mapstructure:"default_expiration"`

	// Object store bucket name. Default: "up2"
	Bucket string

	// Object name prefix. Default: ""
	ObjectNamePrefix string `mapstructure:"object_name_prefix"`

	Aws   *AwsConfig   `validate:"required_without=Minio"`
	Minio *MinioConfig `validate:"required_without=Aws"`
}

type AwsConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	Region          string `mapstructure:"region"`
	AccessKeyId     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
}

type MinioConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyId     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
}

func (osc ObjectStoreConfig) newConfigParser() *viper.Viper {
	parser := newConfigParser()

	// Common configuration
	parser.BindEnv("default_expiration")
	parser.SetDefault("default_expiration", 300)

	parser.BindEnv("bucket")
	parser.SetDefault("bucket", "up2")

	parser.BindEnv("object_name_prefix")
	parser.SetDefault("object_name_prefix", "")

	// AWS S3 object store
	parser.BindEnv("aws.endpoint")
	parser.BindEnv("aws.region")
	parser.BindEnv("aws.access_key_id")
	parser.BindEnv("aws.secret_access_key")

	// MinIO object store
	parser.BindEnv("minio.endpoint")
	parser.BindEnv("minio.access_key_id")
	parser.BindEnv("minio.secret_access_key")

	return parser
}

func (osc *ObjectStoreConfig) LoadConfig() error {
	parser := osc.newConfigParser()
	if err := parser.Unmarshal(osc); err != nil {
		return err
	}
	return nil
}

func (osc ObjectStoreConfig) ValidateConfig() error {
	if err := validateConfig(osc); err != nil {
		return err
	}
	return nil
}

func (osc *ObjectStoreConfig) LoadAndValidateConfig() error {
	if err := osc.LoadConfig(); err != nil {
		return err
	}

	if err := osc.ValidateConfig(); err != nil {
		return err
	}

	return nil
}
