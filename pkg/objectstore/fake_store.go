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

func (s fakeObjectStore) getFakeUrlWithExpiration(key string, expiration time.Duration) (*PreSignedUrl, error) {
	fake := faker.New()

	fakeBucket := s.config.Bucket
	fakeObjectPrefix := s.config.ObjectNamePrefix
	fakeDomain := s.config.Domain
	fakeToken := fake.Hash().SHA512()
	fakeUrl := fmt.Sprintf("https://%s/%s/%s%s?token=%s", fakeDomain, fakeBucket, fakeObjectPrefix, key, fakeToken)

	return &PreSignedUrl{fakeUrl, expiration}, nil
}

func (s fakeObjectStore) GetUploadUrl(key string) (*PreSignedUrl, error) {
	return s.GetUploadUrlWithExpiration(key, s.config.DefaultExpiration)
}

func (s fakeObjectStore) GetUploadUrlWithExpiration(key string, expiration time.Duration) (*PreSignedUrl, error) {
	return s.getFakeUrlWithExpiration(key, expiration)
}

func (s fakeObjectStore) GetDownloadUrl(key string) (*PreSignedUrl, error) {
	return s.getFakeUrlWithExpiration(key, s.config.DefaultExpiration)
}

func (s fakeObjectStore) GetDownloadUrlWithExpiration(key string, expiration time.Duration) (*PreSignedUrl, error) {
	return s.getFakeUrlWithExpiration(key, expiration)
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
