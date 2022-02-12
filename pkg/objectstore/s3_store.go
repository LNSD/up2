package objectstore

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"time"
)

type AwsConfig struct {
	Region         string
	Bucket         string
	BaseObjectName string
}

type S3ObjectStore struct {
	config AwsConfig
}

func (s S3ObjectStore) GetUrl(key string, exp time.Duration) (*PreSignedUrl, error) {
	// Initialize a session in <Region> that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s.config.Region),
	})
	if err != nil {
		return nil, err
	}

	// Create S3 service client
	s3Client := s3.New(sess)

	// Request a Put Object URL and pre-sign it
	reqBucket := s.config.Bucket
	reqKey := fmt.Sprintf("%s/%s", s.config.BaseObjectName, key)
	req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(reqBucket),
		Key:    aws.String(reqKey),
	})

	url, err := req.Presign(exp)
	if err != nil {
		return nil, fmt.Errorf("error presigning the request: %s", err)
	}

	return &PreSignedUrl{Url: url, Expiration: exp}, nil
}

func NewS3ObjectStore(config AwsConfig) ObjectStore {
	return S3ObjectStore{config}
}
