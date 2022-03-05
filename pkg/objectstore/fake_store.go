package objectstore

import (
	"fmt"
	"github.com/jaswdr/faker"
	"time"
)

const DefaultFakeDomain = "obj.example.com"
const DefaultFakeBucket = "up2"
const DefaultFakeObjectNamePrefix = ""
const DefaultFakeExpiration = 5 * time.Minute

type FakeObjectStoreConfig struct {
	Domain            string
	Bucket            string
	ObjectNamePrefix  string
	DefaultExpiration time.Duration
}

type fakeObjectStore struct {
	config FakeObjectStoreConfig
}

func (s fakeObjectStore) Connect() error {
	return nil
}

func (s fakeObjectStore) getFakeURLWithExpiration(key string, expiration time.Duration) (*PreSignedURL, error) {
	fake := faker.New()

	fakeBucket := s.config.Bucket
	fakeObjectPrefix := s.config.ObjectNamePrefix
	fakeDomain := s.config.Domain
	fakeToken := fake.Hash().SHA512()
	fakeURL := fmt.Sprintf("https://%s/%s/%s%s?token=%s", fakeDomain, fakeBucket, fakeObjectPrefix, key, fakeToken)

	return &PreSignedURL{fakeURL, expiration}, nil
}

func (s fakeObjectStore) GetUploadURL(key string) (*PreSignedURL, error) {
	return s.GetUploadURLWithExpiration(key, s.config.DefaultExpiration)
}

func (s fakeObjectStore) GetUploadURLWithExpiration(key string, expiration time.Duration) (*PreSignedURL, error) {
	return s.getFakeURLWithExpiration(key, expiration)
}

func (s fakeObjectStore) GetDownloadURL(key string) (*PreSignedURL, error) {
	return s.getFakeURLWithExpiration(key, s.config.DefaultExpiration)
}

func (s fakeObjectStore) GetDownloadURLWithExpiration(key string, expiration time.Duration) (*PreSignedURL, error) {
	return s.getFakeURLWithExpiration(key, expiration)
}

func NewFakeObjectStore(config *FakeObjectStoreConfig) ObjectStore {
	if config == nil {
		config = &FakeObjectStoreConfig{
			Domain:            DefaultFakeDomain,
			Bucket:            DefaultFakeBucket,
			ObjectNamePrefix:  DefaultFakeObjectNamePrefix,
			DefaultExpiration: DefaultFakeExpiration,
		}
	}
	return &fakeObjectStore{*config}
}
