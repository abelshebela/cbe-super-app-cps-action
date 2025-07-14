package file

import (
	"context"
	"mime/multipart"
)

type UploadResult struct {
	ObjectName string
	Bucket     string
	Key        string
	Name       string
}

type FileService interface {
	UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string) (*UploadResult, error)
	DeleteImage(ctx context.Context, bucketName string, objectName string) error
}
