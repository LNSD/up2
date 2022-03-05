package objectstore

import "time"

type ObjectStore interface {
	Connect() error
	GetUploadURL(key string) (*PreSignedURL, error)
	GetUploadURLWithExpiration(key string, expiration time.Duration) (*PreSignedURL, error)
	GetDownloadURL(key string) (*PreSignedURL, error)
	GetDownloadURLWithExpiration(key string, expiration time.Duration) (*PreSignedURL, error)
}

type PreSignedURL struct {
	URL        string
	Expiration time.Duration
}
