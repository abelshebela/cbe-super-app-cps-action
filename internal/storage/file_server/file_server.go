package file_server

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type S3Persistence struct {
	s3Client *s3.Client
	bucket   string
	baseURL  string
	logger   utils.Logger
}

type FileServerServices struct {
	S3 *S3Persistence
}

func InitFileServerServices(s3Client *s3.Client, bucketName string, baseURL string, logger utils.Logger) *FileServerServices {
	return &FileServerServices{
		S3: &S3Persistence{
			s3Client: s3Client,
			bucket:   strings.Trim(bucketName, "/"),
			baseURL:  strings.TrimRight(baseURL, "/"),
			logger:   logger,
		},
	}
}

// BuildURL creates the full object URL base/bucket/key
func (p *S3Persistence) BuildURL(key string, bucketPrefix string) string {
	k := key
	if bucketPrefix != "" {
		k = path.Join(bucketPrefix, key)
	}
	return fmt.Sprintf("%s/%s/%s", p.baseURL, p.bucket, strings.TrimLeft(k, "/"))
}

// Get fetches an object metadata and body stream from S3
func (p *S3Persistence) Get(ctx context.Context, key string, bucketPrefix string) (*s3.GetObjectOutput, error) {
	k := key
	if bucketPrefix != "" {
		k = path.Join(bucketPrefix, key)
	}
	return p.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(k),
	})
}

// UploadMultipart uploads a multipart file part to S3 and returns the full URL
func (p *S3Persistence) UploadMultipart(ctx context.Context, key string, file multipart.File, header *multipart.FileHeader, bucketPrefix string) (string, error) {
	k := key
	if k == "" && header != nil {
		k = header.Filename
	}
	if bucketPrefix != "" {
		k = path.Join(bucketPrefix, k)
	}
	_, err := p.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(p.bucket),
		Key:         aws.String(k),
		Body:        file,
		ContentType: aws.String(header.Header.Get("Content-Type")),
	})
	
	if err != nil {
		p.logger.Errorf("failed to upload multipart to s3: %v", err)
		return "", err
	}
	return p.BuildURL(k, ""), nil
}

// UploadRaw uploads a raw byte slice to S3 and returns the full URL
func (p *S3Persistence) UploadRaw(ctx context.Context, key string, body []byte, contentType string, bucketPrefix string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("missing key")
	}
	k := key
	if bucketPrefix != "" {
		k = path.Join(bucketPrefix, key)
	}
	ct := contentType
	if ct == "" {
		ct = "application/octet-stream"
	}
	_, err := p.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(p.bucket),
		Key:         aws.String(k),
		Body:        bytes.NewReader(body),
		ContentType: aws.String(ct),
	})
	if err != nil {
		p.logger.Errorf("failed to upload raw to s3: %v", err)
		return "", err
	}
	return p.BuildURL(k, ""), nil
}
