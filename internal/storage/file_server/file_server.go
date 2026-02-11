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

func sanitizeKey(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("invalid key")
	}
	s = strings.ReplaceAll(s, "\\", "/")
	s = strings.TrimLeft(s, "/")
	parts := strings.Split(s, "/")
	cleaned := make([]string, 0, len(parts))
	for _, p := range parts {
		if p == "" || p == "." || p == ".." {
			return "", fmt.Errorf("invalid key segment")
		}
		cleaned = append(cleaned, p)
	}
	joined := strings.Join(cleaned, "/")
	if joined == "" {
		return "", fmt.Errorf("invalid key")
	}
	return joined, nil
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
	if cleaned, err := sanitizeKey(k); err == nil {
		k = cleaned
	}
	if bucketPrefix != "" {
		if pref, err := sanitizeKey(bucketPrefix); err == nil {
			k = path.Join(pref, k)
		} else {
			k = path.Join(bucketPrefix, k)
		}
	}
	return fmt.Sprintf("%s/%s/%s", p.baseURL, p.bucket, strings.TrimLeft(k, "/"))
}

// Get fetches an object metadata and body stream from S3
func (p *S3Persistence) Get(ctx context.Context, key string, bucketPrefix string) (*s3.GetObjectOutput, error) {
	cleanedKey, err := sanitizeKey(key)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %w", err)
	}
	k := cleanedKey
	if bucketPrefix != "" {
		pref, err := sanitizeKey(bucketPrefix)
		if err != nil {
			return nil, fmt.Errorf("invalid prefix: %w", err)
		}
		k = path.Join(pref, k)
	}
	return p.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(k),
	})
}

// UploadMultipart uploads a multipart file part to S3 and returns the full URL
func (p *S3Persistence) UploadMultipart(ctx context.Context, key string, file multipart.File, header *multipart.FileHeader, bucketPrefix string) (string, error) {
	k := strings.TrimSpace(key)
	if k == "" && header != nil {
		name := header.Filename
		name = strings.ReplaceAll(name, "\\", "/")
		name = path.Base(name)
		k = name
	}
	cleanedKey, err := sanitizeKey(k)
	if err != nil {
		return "", fmt.Errorf("invalid key: %w", err)
	}
	finalKey := cleanedKey
	if bucketPrefix != "" {
		pref, err := sanitizeKey(bucketPrefix)
		if err != nil {
			return "", fmt.Errorf("invalid prefix: %w", err)
		}
		finalKey = path.Join(pref, finalKey)
	}
	_, err = p.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(p.bucket),
		Key:         aws.String(finalKey),
		Body:        file,
		ContentType: aws.String(header.Header.Get("Content-Type")),
	})

	if err != nil {
		p.logger.Errorf("failed to upload multipart to s3: %v", err)
		return "", err
	}
	return p.BuildURL(finalKey, ""), nil
}

// UploadRaw uploads a raw byte slice to S3 and returns the full URL
func (p *S3Persistence) UploadRaw(ctx context.Context, key string, body []byte, contentType string, bucketPrefix string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("missing key")
	}
	cleanedKey, err := sanitizeKey(key)
	if err != nil {
		return "", fmt.Errorf("invalid key: %w", err)
	}
	finalKey := cleanedKey
	if bucketPrefix != "" {
		pref, err := sanitizeKey(bucketPrefix)
		if err != nil {
			return "", fmt.Errorf("invalid prefix: %w", err)
		}
		finalKey = path.Join(pref, finalKey)
	}
	ct := contentType
	if ct == "" {
		ct = "application/octet-stream"
	}
	_, err = p.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(p.bucket),
		Key:         aws.String(finalKey),
		Body:        bytes.NewReader(body),
		ContentType: aws.String(ct),
	})
	if err != nil {
		p.logger.Errorf("failed to upload raw to s3: %v", err)
		return "", err
	}
	return p.BuildURL(finalKey, ""), nil
}
