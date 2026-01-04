package service

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"
	"time"

	cg "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"cbe-super-app-mini-app/internal/constants"
	"cbe-super-app-mini-app/internal/domain/entity"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type MediaServiceIntf interface {
	DeleteFile(ctx context.Context, fileName string) *constants.AppError
	GetFile(ctx context.Context, fileName string) (io.ReadCloser, *constants.AppError)
	GetFileInfo(ctx context.Context, fileName string) (*entity.FileInfo, *constants.AppError)
	GetFileRange(ctx context.Context, fileName string, start int64, end int64) (io.ReadCloser, *constants.AppError)
	GetPresignedURL(ctx context.Context, fileName string, expiry time.Duration) (string, *constants.AppError)
	UploadFile(ctx context.Context, fileName string, contentType string, size int64, reader io.Reader) (string, *constants.AppError)
	UploadMultipleFiles(ctx context.Context, files []entity.FileUpload) ([]string, *constants.AppError)
}

type MediaService struct {
	client     *s3.Client
	uploader   *manager.Uploader
	bucketName string
	publicURL  string
	logger     shared_utils.Logger
	endpoint   string
}

// NewMediaService creates a new MediaService instance
func NewMediaService(cfg *cg.VaultConfig, logger shared_utils.Logger) MediaServiceIntf {
	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion("us-east-1"),
		awsConfig.WithBaseEndpoint(cfg.S3BucketURL),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.S3AccessKeyID,
				cfg.S3SecretAccessKey,
				"",
			),
		),
	)

	if err != nil {
		logger.Errorf("failed to load AWS config: %v", err)
		return nil
	}

	// Create S3 client with path-style addressing for MinIO
	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.DisableLogOutputChecksumValidationSkipped = true
	})

	return &MediaService{
		client:     s3Client,
		uploader:   manager.NewUploader(s3Client),
		bucketName: cfg.S3BucketName,
		publicURL:  cfg.MinioPublicEndPoint,
		endpoint:   cfg.S3BucketURL,
		logger:     logger,
	}
}

// UploadFile uploads a file to MinIO (via S3 SDK) and returns the URL
func (m *MediaService) UploadFile(ctx context.Context, fileName string, contentType string, size int64, reader io.Reader) (string, *constants.AppError) {
	ext := filepath.Ext(fileName)
	genName := uuid.New().String() + ext

	// Determine the folder based on content type
	folder := m.getFolderByContentType(contentType)
	objectName := fmt.Sprintf("%s/%s", folder, genName)

	putInput := &s3.PutObjectInput{
		Bucket:      aws.String(m.bucketName),
		Key:         aws.String(objectName),
		Body:        reader,
		ContentType: aws.String(contentType),
	}

	_, err := m.uploader.Upload(ctx, putInput)
	if err != nil {
		m.logger.Errorf("failed to upload file %s: %v", objectName, err)
		return "", constants.ErrRequestFailed
	}

	url := fmt.Sprintf("%s/%s/%s", m.publicURL, folder, genName)
	return url, nil
}

// DeleteFile deletes a file from MinIO (via S3 SDK)
func (m *MediaService) DeleteFile(ctx context.Context, fileName string) *constants.AppError {
	// Extract object name from filename (remove public URL prefix if present)
	objectName := fileName
	if after, ok := strings.CutPrefix(fileName, m.publicURL+"/"); ok {
		objectName = after
	}

	_, err := m.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(m.bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		m.logger.Errorf("failed to delete file %s: %v", objectName, err)
		return constants.ErrRequestFailed
	}

	return nil
}

// GetPresignedURL generates a presigned URL for temporary access to a file
func (m *MediaService) GetPresignedURL(ctx context.Context, fileName string, expiry time.Duration) (string, *constants.AppError) {
	objectName := fileName
	if after, ok := strings.CutPrefix(fileName, m.publicURL+"/"); ok {
		objectName = after
	}

	// Create presign client
	presigner := s3.NewPresignClient(m.client)
	req, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(m.bucketName),
		Key:    aws.String(objectName),
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		m.logger.Errorf("failed to generate presigned URL for %s: %v", objectName, err)
		return "", constants.ErrRequestFailed
	}

	return req.URL, nil
}

// getFolderByContentType determines the folder structure based on content type
func (m *MediaService) getFolderByContentType(contentType string) string {
	switch {
	case strings.HasPrefix(contentType, "mini-app/image/"):
		return "images"
	case strings.HasPrefix(contentType, "mini-app/video/"):
		return "videos"
	case strings.HasPrefix(contentType, "mini-app/audio/"):
		return "audio"
	case contentType == "text/markdown" || contentType == "application/markdown":
		return "mini-app/markdown"
	default:
		return "mini-app/misc"
	}
}

func (m *MediaService) cleanUp(ctx context.Context, uploadedURLs []string) {
	var delWg sync.WaitGroup
	for _, url := range uploadedURLs {
		delWg.Add(1)
		go func(u string) {
			defer delWg.Done()
			_ = m.DeleteFile(ctx, u)
		}(url)
	}
	delWg.Wait()
}

// UploadMultipleFiles uploads multiple files concurrently and returns their URLs
func (m *MediaService) UploadMultipleFiles(ctx context.Context, files []entity.FileUpload) ([]string, *constants.AppError) {
	if len(files) == 0 {
		return nil, nil
	}

	const maxConcurrent = 4
	sem := make(chan struct{}, maxConcurrent)
	results := make(chan struct {
		url   string
		err   *constants.AppError
		index int
	}, len(files))
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr *constants.AppError
	var uploadedURLs []string

	for i, file := range files {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, f entity.FileUpload) {
			defer wg.Done()
			defer func() { <-sem }()

			subCtx, cancel := context.WithCancel(ctx)
			defer cancel()

			url, err := m.UploadFile(subCtx, f.FileName, f.ContentType, f.Size, f.Reader)
			mu.Lock()
			if firstErr == nil && err != nil {
				firstErr = constants.ErrRequestFailed
			} else if err == nil {
				uploadedURLs = append(uploadedURLs, url)
			}
			mu.Unlock()

			results <- struct {
				url   string
				err   *constants.AppError
				index int
			}{url: url, err: err, index: idx}
		}(i, file)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	urls := make([]string, len(files))
	for res := range results {
		if res.err != nil && firstErr == nil {
			firstErr = res.err
		}
		if firstErr == nil {
			urls[res.index] = res.url
		}
	}

	if firstErr != nil {
		go m.cleanUp(ctx, uploadedURLs)
		return nil, firstErr
	}

	return urls, nil
}

// GetFileInfo retrieves file information from MinIO (via S3 SDK)
func (m *MediaService) GetFileInfo(ctx context.Context, fileName string) (*entity.FileInfo, *constants.AppError) {
	// Extract object name from filename
	objectName := fileName
	if after, ok := strings.CutPrefix(fileName, m.publicURL+"/"); ok {
		objectName = after
	}

	headOut, err := m.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(m.bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		m.logger.Errorf("failed to get file info for %s: %v", objectName, err)
		return nil, constants.ErrRequestFailed
	}

	return &entity.FileInfo{
		Size:         *headOut.ContentLength,
		ContentType:  *headOut.ContentType,
		LastModified: *headOut.LastModified,
		ETag:         *headOut.ETag,
	}, nil
}

// GetFileRange retrieves a specific range of bytes from a file for streaming
func (m *MediaService) GetFileRange(ctx context.Context, fileName string, start, end int64) (io.ReadCloser, *constants.AppError) {
	// Extract object name from filename
	objectName := fileName
	if after, ok := strings.CutPrefix(fileName, m.publicURL+"/"); ok {
		objectName = after
	}

	// Get object with range
	getInput := &s3.GetObjectInput{
		Bucket: aws.String(m.bucketName),
		Key:    aws.String(objectName),
		Range:  aws.String(fmt.Sprintf("bytes=%d-%d", start, end)),
	}
	out, err := m.client.GetObject(ctx, getInput)
	if err != nil {
		m.logger.Errorf("failed to get file range %d-%d for %s: %v", start, end, objectName, err)
		return nil, constants.ErrRequestFailed
	}

	return out.Body, nil
}

// GetFile retrieves the entire file
func (m *MediaService) GetFile(ctx context.Context, fileName string) (io.ReadCloser, *constants.AppError) {
	// Extract object name from filename
	objectName := fileName
	if after, ok := strings.CutPrefix(fileName, m.publicURL+"/"); ok {
		objectName = after
	}

	// Get object
	getInput := &s3.GetObjectInput{
		Bucket: aws.String(m.bucketName),
		Key:    aws.String(objectName),
	}
	out, err := m.client.GetObject(ctx, getInput)
	if err != nil {
		m.logger.Errorf("failed to get file %s: %v", objectName, err)
		return nil, constants.ErrRequestFailed
	}

	return out.Body, nil
}
