package core

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
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
	exist, err := uploader.BucketExist(ctx, bucketName)
	if err != nil {
		logger.Errorf("failed to check bucket '%s': %v", bucketName, err)
		return "", fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exist {
		created, err := uploader.MakeBucket(ctx, bucketName)
		if err != nil || !created {
			logger.Errorf("failed to create bucket '%s': %v", bucketName, err)
			return "", fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	file, err := fileHeader.Open()
	if err != nil {
		logger.Errorf("failed to open file: %v", err)
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileName := fmt.Sprintf("%s-%d-%s", prefix, time.Now().UnixNano(), fileHeader.Filename)

	saveObj, err := uploader.SaveObjectN(ctx, config.SaveObjectBodyN{
		BucketName:  bucketName,
		ObjectName:  fileName,
		Reader:      file,
		Size:        fileHeader.Size,
		ContentType: config.ContentType(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		logger.Errorf("failed to upload file to MinIO: %v", err)
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	url := fmt.Sprintf("%s/%s/%s", minioEndpoint, saveObj.Bucket, saveObj.Key)
	return url, nil
}
