package objectstore

import "time"

type ObjectStore interface {
	GetUrl(key string, exp time.Duration) (*PreSignedUrl, error)
}

type PreSignedUrl struct {
	Url        string
	Expiration time.Duration
}
