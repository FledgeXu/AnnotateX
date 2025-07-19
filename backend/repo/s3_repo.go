package repo

import (
	"annotate-x/config"
	"annotate-x/models"
	"annotate-x/utils"
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type IS3Repo interface {
	UploadObject(ctx context.Context, objectName string, data []byte, contentType string) error
	UploadFile(ctx context.Context, objectName, filePath string) error
	GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error)
}

type S3Repo struct {
	Client     *minio.Client
	BucketName string
}

func NewS3Repo(s3Config config.S3Config, bucketName models.BucketName) *S3Repo {
	client, err := minio.New(s3Config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3Config.AccessKey, s3Config.SecretKey, ""),
		Secure: s3Config.UseSSL,
	})
	if err != nil {
		panic(fmt.Errorf("failed to init s3 client: %w", err))
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, string(bucketName))
	if err != nil {
		panic(fmt.Errorf("failed to check bucket: %w", err))
	}
	if !exists {
		panic(fmt.Errorf("bucket: %s is not exists", bucketName))
	}

	return &S3Repo{client, string(bucketName)}
}

func (r *S3Repo) UploadObject(ctx context.Context, objectName string, data []byte, contentType string) error {
	_, err := r.Client.PutObject(ctx, r.BucketName, objectName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (r *S3Repo) UploadFile(ctx context.Context, objectName, filePath string) error {
	contentType, err := utils.DetectContentType(filePath)
	if err != nil {
		return err
	}
	_, err = r.Client.FPutObject(ctx, r.BucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	return err
}

func (r *S3Repo) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := r.Client.PresignedGetObject(ctx, r.BucketName, objectName, expiry, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
