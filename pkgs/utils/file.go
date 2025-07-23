package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	config "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
)

func UploadFileToMinio(
	ctx context.Context,
	uploader config.MinioClientInterface,
	bucketName string,
	fileHeader *multipart.FileHeader,
	prefix string,
	minioEndpoint string,
	logger interface {
		Errorf(format string, args ...any)
	},
) (string, error) {
	// Ensure bucket exists
	exist, err := uploader.BucketExist(ctx, bucketName)
	if err != nil {
		logger.Errorf("failed to check bucket '%s': %v", bucketName, err)
		return "", fmt.Errorf(UnhandledServerError)
	}

	if !exist {
		created, err := uploader.MakeBucket(ctx, bucketName)
		if err != nil || !created {
			logger.Errorf("failed to create bucket '%s': %v", bucketName, err)
			return "", fmt.Errorf(UnhandledServerError)
		}
	}

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		logger.Errorf("failed to open file: %v", err)
		return "", fmt.Errorf(UnhandledServerError)
	}
	defer file.Close()

	// Generate file name
	fileName := fmt.Sprintf("%s-%d-%s", prefix, time.Now().UnixNano(), fileHeader.Filename)

	// Upload file
	saveObj, err := uploader.SaveObjectN(ctx, config.SaveObjectBodyN{
		BucketName:  bucketName,
		ObjectName:  fileName,
		Reader:      file,
		Size:        fileHeader.Size,
		ContentType: config.ContentType(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		logger.Errorf("failed to upload file to MinIO: %v", err)
		return "", fmt.Errorf(UnhandledServerError)
	}

	// Return full URL
	url := fmt.Sprintf("%s/%s/%s", minioEndpoint, saveObj.Bucket, saveObj.Key)
	return url, nil
}
