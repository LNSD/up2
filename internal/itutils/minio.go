package itutils

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"testing"
)

const testContainerRegion = "eu-central-1"
const testContainerBucket = "test"
const testContainerAccessKeyId = "minio"
const testContainerSecretAccessKey = "miniostorage"

type MinioContainer struct {
	testcontainers.Container
	Endpoint        string
	ConsoleEndpoint string
	Region          string
	Bucket          string
	AccessKeyId     string
	SecretAccessKey string
}

func NewMinioContainer(t *testing.T, ctx context.Context) *MinioContainer {
	req := testcontainers.ContainerRequest{
		Image:        "quay.io/minio/minio:latest",
		Name:         "minio",
		ExposedPorts: []string{"9000/tcp", "9001/tcp"},
		Env: map[string]string{
			"MINIO_ROOT_USER":     testContainerAccessKeyId,
			"MINIO_ROOT_PASSWORD": testContainerSecretAccessKey,
		},
		Cmd:        []string{"server", "--console-address", ":9001", "/data"},
		WaitingFor: wait.ForHTTP("/minio/health/live").WithPort("9000"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	apiPort, err := container.MappedPort(ctx, "9000")
	if err != nil {
		t.Fatal(err)
	}
	endpoint := fmt.Sprintf("%s:%s", host, apiPort.Port())

	consolePort, err := container.MappedPort(ctx, "9001")
	if err != nil {
		t.Fatal(err)
	}
	consoleEndpoint := fmt.Sprintf("%s:%s", host, consolePort.Port())

	minioContainer := &MinioContainer{
		Container:       container,
		Region:          testContainerRegion,
		Bucket:          testContainerBucket,
		Endpoint:        endpoint,
		ConsoleEndpoint: consoleEndpoint,
		AccessKeyId:     testContainerAccessKeyId,
		SecretAccessKey: testContainerSecretAccessKey,
	}
	return minioContainer

}

func InitMinio(t *testing.T, container *MinioContainer) {
	minioClient, err := minio.New(container.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(container.AccessKeyId, container.SecretAccessKey, ""),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create bucket
	err = minioClient.MakeBucket(context.Background(), container.Bucket, minio.MakeBucketOptions{
		Region:        container.Region,
		ObjectLocking: false,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TerminateMinio(t *testing.T, ctx context.Context, container *MinioContainer) {
	if err := container.Terminate(ctx); err != nil {
		t.Fatal(err)
	}
}

func ListBucketObjects(t *testing.T, container *MinioContainer) []minio.ObjectInfo {
	minioClient, err := minio.New(container.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(container.AccessKeyId, container.SecretAccessKey, ""),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create bucket
	objects := make([]minio.ObjectInfo, 0)
	for object := range minioClient.ListObjects(context.Background(), container.Bucket, minio.ListObjectsOptions{MaxKeys: 20}) {
		objects = append(objects, object)
	}
	if err != nil {
		t.Fatal(err)
	}

	return objects
}

func PutObject(t *testing.T, container *MinioContainer, key string, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	minioClient, err := minio.New(container.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(container.AccessKeyId, container.SecretAccessKey, ""),
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = minioClient.PutObject(context.Background(), container.Bucket, key, file, fileStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		t.Fatal(err)
	}
}
