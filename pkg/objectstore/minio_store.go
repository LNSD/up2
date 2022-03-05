package objectstore

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"net/url"
	"time"
)

type MinioConfig struct {
	Endpoint          string
	AccessKeyID       string
	SecretAccessKey   string
	Bucket            string
	ObjectNamePrefix  string
	DefaultExpiration time.Duration
}

type minioObjectStore struct {
	config MinioConfig
}

func (s minioObjectStore) Connect() error {
	// Initialize minio client object
	_, err := minio.New(s.config.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(s.config.AccessKeyID, s.config.SecretAccessKey, ""),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s minioObjectStore) GetUploadURL(key string) (*PreSignedURL, error) {
	return s.GetUploadURLWithExpiration(key, s.config.DefaultExpiration)
}

func (s minioObjectStore) GetUploadURLWithExpiration(key string, expiration time.Duration) (*PreSignedURL, error) {
	// Initialize minio client object
	// TODO: Share minio client instance between calls
	minioClient, err := minio.New(s.config.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(s.config.AccessKeyID, s.config.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}

	// Generate pre-signed PUT Object URL
	preSignedURL, err := minioClient.PresignedPutObject(context.Background(), s.config.Bucket, key, expiration)
	if err != nil {
		return nil, err
	}

	return &PreSignedURL{preSignedURL.String(), expiration}, nil
}

func (s minioObjectStore) GetDownloadURL(key string) (*PreSignedURL, error) {
	return s.GetDownloadURLWithExpiration(key, s.config.DefaultExpiration)
}

func (s minioObjectStore) GetDownloadURLWithExpiration(key string, expiration time.Duration) (*PreSignedURL, error) {
	// Initialize minio client object
	// TODO: Share minio client instance between calls
	minioClient, err := minio.New(s.config.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(s.config.AccessKeyID, s.config.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}

	// Generate pre-signed GET Object URL
	preSignedURL, err := minioClient.PresignedGetObject(context.Background(), s.config.Bucket, key, expiration, url.Values{})
	if err != nil {
		return nil, err
	}

	return &PreSignedURL{preSignedURL.String(), expiration}, nil
}

func NewMinioObjectStore(config MinioConfig) ObjectStore {
	return &minioObjectStore{config}
}
