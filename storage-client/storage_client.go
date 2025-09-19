package storageclient

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
)

const credsEnvVar = "GOOGLE_APPLICATION_CREDENTIALS"

type StorageClient struct {
	client         *storage.Client
	projectID      string
	bucketName     string
	timeoutSeconds int
}

//go:generate mockery --name IStorageClient
type IStorageClient interface {
	ReplaceFile(path string, destObject string) error
	UploadNewFile(path string, destObject string) error
}

type BackupFile struct {
	Name    string
	Updated time.Time
}

func New() *StorageClient {
	projectID := os.Getenv("GCS_PROJECT_ID")
	if projectID == "" {
		panic("GCS_PROJECT_ID must be set")
	}
	bucketName := os.Getenv("GCS_BUCKET_NAME")
	if bucketName == "" {
		panic("GCS_BUCKET_NAME must be set")
	}
	gcsTimeoutSeconds := os.Getenv("GCS_TIMEOUT_SECONDS")
	if gcsTimeoutSeconds == "" {
		panic("GCS_TIMEOUT_SECONDS must be set")
	}
	gcsTimeoutSecondsInt, err := strconv.Atoi(gcsTimeoutSeconds)
	if err != nil {
		panic(fmt.Sprintf("could not convert %v to int", gcsTimeoutSeconds))
	}
	gcsCredsFile := os.Getenv("GCS_CREDS_FILE")
	if gcsCredsFile == "" {
		panic("GCS_CREDS_FILE must be set")
	}

	os.Setenv(credsEnvVar, gcsCredsFile)

	client, err := storage.NewClient(context.Background())
	if err != nil {
		panic(err)
	}

	return &StorageClient{
		client:         client,
		bucketName:     bucketName,
		projectID:      projectID,
		timeoutSeconds: gcsTimeoutSecondsInt,
	}
}

func (sc *StorageClient) ReplaceFile(path string, destObject string) error {
	blobFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer blobFile.Close()

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(sc.timeoutSeconds))
	defer cancel()

	obj := sc.client.Bucket(sc.bucketName).Object(destObject)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to replace the file is aborted
	// if the object's generation number does not match your precondition.
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("error getting object attributes: %v", err)
	}
	obj = obj.If(storage.Conditions{GenerationMatch: attrs.Generation})

	wc := obj.NewWriter(ctx)

	if _, err := io.Copy(wc, blobFile); err != nil {
		return fmt.Errorf("error on Copy to bucket %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("error on Close during bucket upload: %v", err)
	}

	return nil
}

func (sc *StorageClient) UploadNewFile(path string, destObject string) error {
	blobFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer blobFile.Close()

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(sc.timeoutSeconds))
	defer cancel()

	obj := sc.client.Bucket(sc.bucketName).Object(destObject)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	// For an object that does not yet exist, set the DoesNotExist precondition.
	obj = obj.If(storage.Conditions{DoesNotExist: true})

	wc := obj.NewWriter(ctx)

	if _, err := io.Copy(wc, blobFile); err != nil {
		return fmt.Errorf("error on Copy to bucket %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("error on Close during bucket upload: %v", err)
	}

	return nil
}
