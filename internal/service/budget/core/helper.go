package core

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// SaveIconToMinio saves the icon to MinIO and returns the saved object response and error.
func SaveIconToMinio(ctx context.Context, fileHeader *multipart.FileHeader, bucketName string, minio config.MinioClientInterface, logger utils.Logger) (*config.SaveObjectResponse, error) {
	imageFileHeader := fileHeader

	file, err := imageFileHeader.Open()
	if err != nil {
		logger.Errorf("failed to open uploaded image: %v", err)
		return nil, fmt.Errorf("failed to open uploaded image")
	}
	defer file.Close()

	exist, err := minio.BucketExist(ctx, bucketName)
	if err != nil {
		logger.Errorf("failed to check icon bucket: %v", err)
		return nil, fmt.Errorf("failed to check icon bucket")
	}

	if !exist {
		created, err := minio.MakeBucket(ctx, bucketName)
		if !created || err != nil {
			logger.Errorf("failed to create icon bucket: %v", err)
			return nil, fmt.Errorf("failed to create icon bucket")
		}
	}

	fileName := fmt.Sprintf("budget-icon-%d-%s", time.Now().UnixNano(), imageFileHeader.Filename)
	dir, err := os.Getwd()
	if err != nil {
		logger.Errorf("failed to get current working directory", err)
		return nil, fmt.Errorf("failed to get current working directory")
	}

	filePath := filepath.Join(dir, fileName)

	tempFile, err := os.Create(filePath)
	if err != nil {
		logger.Errorf("failed to create temp file: %v", err)
		return nil, fmt.Errorf("failed to create temp file")
	}
	defer func() {
		tempFile.Close()
		os.Remove(filePath)
	}()

	// Copy the uploaded file to the temp file
	if _, err := file.Seek(0, 0); err != nil {
		logger.Errorf("failed to seek file: %v", err)
		return nil, fmt.Errorf("failed to seek file")
	}
	if _, err := tempFile.ReadFrom(file); err != nil {
		logger.Errorf("failed to copy file to temp file: %v", err)
		return nil, fmt.Errorf("failed to copy file to temp file")
	}

	saveObj, err := minio.SaveObject(ctx, config.SaveObjectBody{
		BucketName: bucketName,
		ObjectName: fileName,
		File:       filePath,
	})
	if err != nil {
		logger.Errorf("failed to save object to MinIO: %v", err)
		return nil, fmt.Errorf("failed to save object to MinIO")
	}

	return saveObj, nil
}
