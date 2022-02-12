package objectstore

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"time"
)

type MinioConfig struct {
	Endpoint        string
	AccessKeyId     string
	SecretAccessKey string
	BucketName      string
	BaseObjectName  string
}

type MinioObjectStore struct {
	config MinioConfig
}

func (s MinioObjectStore) GetUrl(key string, exp time.Duration) (*PreSignedUrl, error) {
	// Initialize minio client object
	minioClient, err := minio.New(s.config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s.config.AccessKeyId, s.config.SecretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	// Generate pre-signed PUT Object URL
	objectName := fmt.Sprintf("%s/%s", s.config.BucketName, key)
	presignedURL, err := minioClient.PresignedPutObject(context.Background(), s.config.BucketName, objectName, exp*time.Second)
	if err != nil {
		return nil, err
	}

	return &PreSignedUrl{Url: presignedURL.String(), Expiration: exp}, nil
}

func NewMinioObjectStore(config MinioConfig) ObjectStore {
	return &MinioObjectStore{config}
}
