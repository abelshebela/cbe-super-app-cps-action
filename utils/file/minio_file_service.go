// internal/service/file/minio_file_service.go
package file

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type fileService struct {
	minioClient config.MinioClientInterface
	logger      utils.Logger
}

func NewFileService(minioClient config.MinioClientInterface, logger utils.Logger) FileService {
	return &fileService{
		minioClient: minioClient,
		logger:      logger,
	}
}

func (m *fileService) UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string) (*UploadResult, error) {
	tempFile, err := os.CreateTemp("", "upload-*.tmp")
	if err != nil {
		m.logger.Errorf("failed to create temp file: %v", err)
		return nil, fmt.Errorf("UPLOAD_FAILED")
	}
	defer func() {
		if err := tempFile.Close(); err != nil {
			m.logger.Errorf("failed to close temp file: %v", err)
		}
		if err := os.Remove(tempFile.Name()); err != nil {
			m.logger.Errorf("failed to remove temp file: %v", err)
		}
	}()

	if _, err := io.Copy(tempFile, file); err != nil {
		m.logger.Errorf("failed to copy to temp file: %v", err)
		return nil, fmt.Errorf("UPLOAD_FAILED")
	}

	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return nil, fmt.Errorf("INVALID_FILE_TYPE")
	}

	bucket := folder
	objectName := fmt.Sprintf("%s/%s", folder, filepath.Base(header.Filename))

	exists, err := m.minioClient.BucketExist(ctx, bucket)
	if err != nil {
		m.logger.Errorf("bucket existence check failed: %v", err)
		return nil, fmt.Errorf("UPLOAD_FAILED")
	}
	if !exists {
		if _, err := m.minioClient.MakeBucket(ctx, bucket); err != nil {
			m.logger.Errorf("bucket creation failed: %v", err)
			return nil, fmt.Errorf("UPLOAD_FAILED")
		}
	}

	resp, err := m.minioClient.SaveObject(ctx, config.SaveObjectBody{
		BucketName:  bucket,
		ObjectName:  objectName,
		File:        tempFile.Name(),
		ContentType: config.ContentType(contentType),
	})
	if err != nil {
		m.logger.Errorf("failed to save object: %v", err)
		return nil, fmt.Errorf("UPLOAD_FAILED")
	}

	return &UploadResult{
		ObjectName: objectName,
		Bucket:     resp.Bucket,
		Key:        resp.Key,
		Name:       header.Filename,
	}, nil
}

func (s *fileService) DeleteImage(ctx context.Context, bucket, object string) error {
	_, err := s.minioClient.DeleteObject(ctx, config.DeleteObjectBody{
		BucketName: bucket,
		ObjectName: object,
	})
	if err != nil {
		s.logger.Errorf("delete object failed: %v", err)
		return fmt.Errorf("failed to delete object")
	}
	return nil
}
