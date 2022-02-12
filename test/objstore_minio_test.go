//go:build integration
// +build integration

package test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
	"upload-presigned-url-provider/internal/itutils"
	"upload-presigned-url-provider/pkg/objectstore"
)

func uploadFile(t *testing.T, filename string, url string) {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", url, file)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.ContentLength = stat.Size()

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatal(res.Status)
	}
}

func downloadFile(t *testing.T, url string, filename string) {
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		t.Fatal(err)
	}
}

func calculateFileSha256(t *testing.T, filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		t.Fatal(err)
	}

	sha := hash.Sum(nil)
	return hex.EncodeToString(sha)
}

func TestIntegrationMinioObjectStoreConnect(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	/// Given
	ctx := context.Background()

	// Start the container
	minioContainer := itutils.NewMinioContainer(t, ctx)
	defer itutils.TerminateMinio(t, ctx, minioContainer)

	// Initialize MinIO instance
	itutils.InitMinio(t, minioContainer)

	// Create object store
	store := objectstore.NewMinioObjectStore(objectstore.MinioConfig{
		Endpoint:          minioContainer.Endpoint,
		Bucket:            minioContainer.Bucket,
		AccessKeyId:       minioContainer.AccessKeyId,
		SecretAccessKey:   minioContainer.SecretAccessKey,
		ObjectNamePrefix:  "",
		DefaultExpiration: 60 * time.Second,
	})

	/// When
	err := store.Connect()

	/// Then
	assert.NoError(t, err)
}

func TestIntegrationMinioObjectStoreGeUploadUrl(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	/// Given
	ctx := context.Background()

	// Start the container
	minioContainer := itutils.NewMinioContainer(t, ctx)
	defer itutils.TerminateMinio(t, ctx, minioContainer)

	// Initialize MinIO instance
	itutils.InitMinio(t, minioContainer)

	// Create object store
	store := objectstore.NewMinioObjectStore(objectstore.MinioConfig{
		Endpoint:          minioContainer.Endpoint,
		Bucket:            minioContainer.Bucket,
		AccessKeyId:       minioContainer.AccessKeyId,
		SecretAccessKey:   minioContainer.SecretAccessKey,
		ObjectNamePrefix:  "",
		DefaultExpiration: 60 * time.Second,
	})

	/// When
	uuid := faker.New().UUID().V4()
	preSignedURL, err := store.GetUploadUrl(uuid)

	uploadFile(t, "./testdata/up.gif", preSignedURL.Url)

	/// Then
	assert.NoError(t, err)
	assert.NotNil(t, preSignedURL)
	assert.NotEmpty(t, preSignedURL.Url)
	assert.Equal(t, 60*time.Second, preSignedURL.Expiration)

	objects := itutils.ListBucketObjects(t, minioContainer)
	assert.Len(t, objects, 1)
	assert.Equal(t, objects[0].Key, uuid)
}

func TestIntegrationMinioObjectStoreGeDownloadUrl(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	/// Given
	ctx := context.Background()

	// Start the container
	minioContainer := itutils.NewMinioContainer(t, ctx)
	defer itutils.TerminateMinio(t, ctx, minioContainer)

	// Initialize MinIO instance
	itutils.InitMinio(t, minioContainer)

	// Upload test file
	uuid := faker.New().UUID().V4()
	itutils.PutObject(t, minioContainer, uuid, "./testdata/up.gif")

	// Create object store
	store := objectstore.NewMinioObjectStore(objectstore.MinioConfig{
		Endpoint:          minioContainer.Endpoint,
		Bucket:            minioContainer.Bucket,
		AccessKeyId:       minioContainer.AccessKeyId,
		SecretAccessKey:   minioContainer.SecretAccessKey,
		ObjectNamePrefix:  "",
		DefaultExpiration: 60 * time.Second,
	})

	// Create temporary download directory
	dir, err := ioutil.TempDir("/tmp", "prefix")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	/// When
	preSignedURL, err := store.GetDownloadUrl(uuid)

	filename := filepath.Join(dir, "up.gif")
	downloadFile(t, preSignedURL.Url, filename)

	/// Then
	assert.NoError(t, err)
	assert.NotNil(t, preSignedURL)
	assert.NotEmpty(t, preSignedURL.Url)
	assert.Equal(t, 60*time.Second, preSignedURL.Expiration)

	checksum := calculateFileSha256(t, filename)
	assert.Equal(t, "3bb5224dbb80594443d83a460c851d350f805866289991ad4a9ca664c25c0078", checksum)
}
