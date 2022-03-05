package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"
	"upload-presigned-url-provider/pkg/objectstore"
	"upload-presigned-url-provider/pkg/server"
)

const serverShutdownTimeout = 5 * time.Second

func main() {
	// Load configuration
	serverConfig := new(ServerConfig)
	if err := serverConfig.LoadAndValidateConfig(); err != nil {
		log.Fatal("Failed to load config: " + err.Error())
	}

	storeConfig := new(ObjectStoreConfig)
	if err := storeConfig.LoadAndValidateConfig(); err != nil {
		log.Fatal("Failed to load config: " + err.Error())
	}

	// Get ObjectStore client instance
	store := getObjetStoreFromConfig(storeConfig)
	// TODO: Check object store connection

	// Instantiate and create server
	srv := server.New(&store, &server.Config{
		Host: serverConfig.Host,
		Port: serverConfig.Port,
	})

	go srv.Start()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 5 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()
	srv.Shutdown(ctx)
}

func getObjetStoreFromConfig(storeConfig *ObjectStoreConfig) objectstore.ObjectStore {
	var store objectstore.ObjectStore
	if storeConfig.Aws != nil {
		store = objectstore.NewS3ObjectStore(objectstore.AwsConfig{
			Endpoint:          (*storeConfig.Aws).Endpoint,
			Region:            (*storeConfig.Aws).Region,
			AccessKeyID:       (*storeConfig.Aws).AccessKeyID,
			SecretAccessKey:   (*storeConfig.Aws).SecretAccessKey,
			Bucket:            storeConfig.Bucket,
			ObjectNamePrefix:  storeConfig.ObjectNamePrefix,
			DefaultExpiration: time.Duration(storeConfig.DefaultExpiration) * time.Second,
		})
	} else if storeConfig.Minio != nil {
		store = objectstore.NewMinioObjectStore(objectstore.MinioConfig{
			Endpoint:          (*storeConfig.Minio).Endpoint,
			AccessKeyID:       (*storeConfig.Minio).AccessKeyID,
			SecretAccessKey:   (*storeConfig.Minio).SecretAccessKey,
			Bucket:            storeConfig.Bucket,
			ObjectNamePrefix:  storeConfig.ObjectNamePrefix,
			DefaultExpiration: time.Duration(storeConfig.DefaultExpiration) * time.Second,
		})
	} else {
		log.Fatal("missing object store credentials")
	}
	return store
}
