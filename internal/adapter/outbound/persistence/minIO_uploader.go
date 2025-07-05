package persistence

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url" // <-- add this
	"time"    // <-- add this

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOUploaderImpl struct {
	Client     *minio.Client
	BucketName string
	Endpoint   string
	UseSSL     bool
}

func NewMinIOUploader(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinIOUploaderImpl, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	return &MinIOUploaderImpl{
		Client:     client,
		BucketName: bucket,
		Endpoint:   endpoint,
		UseSSL:     useSSL,
	}, nil
}

func (u *MinIOUploaderImpl) UploadProfileImage(ctx context.Context, objectName string, file multipart.File, contentType string, userID string) (string, error) {
	defer file.Close()
	_, err := u.Client.PutObject(
		ctx,
		u.BucketName,
		objectName,
		file,
		-1,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", err
	}
	protocol := "http"
	if u.UseSSL {
		protocol = "https"
	}
	url := fmt.Sprintf("%s://%s/%s/%s", protocol, u.Endpoint, u.BucketName, objectName)
	return url, nil
}

func GeneratePresignedURL(client *minio.Client, bucketName, objectName string) (string, error) {
	expiry := time.Minute * 15 // Valid for 15 minutes

	reqParams := make(url.Values)
	presignedURL, err := client.PresignedGetObject(context.Background(), bucketName, objectName, expiry, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
