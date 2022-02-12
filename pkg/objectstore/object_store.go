package objectstore

import "time"

type ObjectStore interface {
	Connect() error
	GetUploadUrl(key string) (*PreSignedUrl, error)
	GetUploadUrlWithExpiration(key string, expiration time.Duration) (*PreSignedUrl, error)
	GetDownloadUrl(key string) (*PreSignedUrl, error)
	GetDownloadUrlWithExpiration(key string, expiration time.Duration) (*PreSignedUrl, error)
}

type PreSignedUrl struct {
	Url        string
	Expiration time.Duration
}
