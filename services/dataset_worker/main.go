package main

import (
	"context"
	"dataset_worker/config"
	"dataset_worker/models"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func GetS3Object(ctx context.Context, client *minio.Client, bucket, objectName string) (*minio.Object, error) {
	object, err := client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("get object %s from bucket %s: %w", objectName, bucket, err)
	}
	return object, nil
}

func PrepareLocalFilePath(rootDir, objectName string) (string, error) {
	localPath := path.Join(rootDir, objectName)
	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("create directory %s: %w", dir, err)
	}
	return localPath, nil
}

func SaveReaderToFile(r io.Reader, localPath string) error {
	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("create file %s: %w", localPath, err)
	}
	defer file.Close()

	if _, err := io.Copy(file, r); err != nil {
		return fmt.Errorf("write to file %s: %w", localPath, err)
	}
	return nil
}

func DownloadS3ObjectToLocal(ctx context.Context, client *minio.Client, bucket, objectName, rootDir string) (string, error) {
	object, err := GetS3Object(ctx, client, bucket, objectName)
	if err != nil {
		return "", err
	}
	defer object.Close()

	localPath, err := PrepareLocalFilePath(rootDir, objectName)
	if err != nil {
		return "", err
	}

	if err := SaveReaderToFile(object, localPath); err != nil {
		return "", err
	}

	return localPath, nil
}

func main() {
	minioClient, err := minio.New(config.GetConfig().S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.GetConfig().S3AccessKey, config.GetConfig().S3SecretKey, ""),
		Secure: config.GetConfig().S3UseSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	message := models.TransformDatasetMessage{
		Dataset: models.Dataset{
			ID:            1,
			ProjectID:     1,
			Name:          "HH",
			Description:   nil,
			FormatVersion: "HHH",
			Status:        "HHH",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		Keys: []string{"/1/test2/finial.zip"},
	}
	rootDir, err := os.MkdirTemp(".", "*")
	defer os.RemoveAll(rootDir)
	if err != nil {
		fmt.Println(err)
	}

	ctx := context.Background()
	for _, objectName := range message.Keys {
		DownloadS3ObjectToLocal(ctx, minioClient, config.GetConfig().S3BucketName, objectName, rootDir)
	}
}
