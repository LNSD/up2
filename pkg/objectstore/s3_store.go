package objectstore

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"time"
)

type AwsConfig struct {
	Endpoint          string
	Region            string
	AccessKeyID       string
	SecretAccessKey   string
	Bucket            string
	ObjectNamePrefix  string
	DefaultExpiration time.Duration
}

type s3ObjectStore struct {
	config AwsConfig
}

func (s s3ObjectStore) Connect() error {
	_, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(s.config.Endpoint),
		Region:      aws.String(s.config.Region),
		Credentials: credentials.NewStaticCredentials(s.config.AccessKeyID, s.config.SecretAccessKey, ""),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s s3ObjectStore) GetUploadURL(key string) (*PreSignedURL, error) {
	return s.GetUploadURLWithExpiration(key, s.config.DefaultExpiration)
}

func (s s3ObjectStore) GetUploadURLWithExpiration(key string, expiration time.Duration) (*PreSignedURL, error) {
	// Initialize aws client session instance
	// TODO: Share client instance between calls
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(s.config.Endpoint),
		Region:      aws.String(s.config.Region),
		Credentials: credentials.NewStaticCredentials(s.config.AccessKeyID, s.config.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}

	// Create S3 service client
	s3Client := s3.New(sess)

	// Request a Put Object URL and pre-sign it
	reqBucket := s.config.Bucket
	reqKey := fmt.Sprintf("%s/%s", s.config.ObjectNamePrefix, key)
	req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(reqBucket),
		Key:    aws.String(reqKey),
	})

	url, err := req.Presign(expiration)
	if err != nil {
		return nil, fmt.Errorf("error presigning the request: %s", err)
	}

	return &PreSignedURL{url, expiration}, nil
}

func (s s3ObjectStore) GetDownloadURL(key string) (*PreSignedURL, error) {
	return s.GetUploadURLWithExpiration(key, s.config.DefaultExpiration)
}

func (s s3ObjectStore) GetDownloadURLWithExpiration(key string, expiration time.Duration) (*PreSignedURL, error) {
	// Initialize aws client session instance
	// TODO: Share client instance between calls
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(s.config.Endpoint),
		Region:      aws.String(s.config.Region),
		Credentials: credentials.NewStaticCredentials(s.config.AccessKeyID, s.config.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}

	// Create S3 service client
	s3Client := s3.New(sess)

	// Request a Get Object URL and pre-sign it
	reqBucket := s.config.Bucket
	reqKey := fmt.Sprintf("%s/%s", s.config.ObjectNamePrefix, key)
	req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(reqBucket),
		Key:    aws.String(reqKey),
	})

	url, err := req.Presign(expiration)
	if err != nil {
		return nil, fmt.Errorf("error presigning the request: %s", err)
	}

	return &PreSignedURL{url, expiration}, nil
}

func NewS3ObjectStore(config AwsConfig) ObjectStore {
	return s3ObjectStore{config}
}
