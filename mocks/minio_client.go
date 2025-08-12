package mocks

import (
	"context"
	"io"

	"github.com/stretchr/testify/mock"
)

type MinioClientInterface struct {
	mock.Mock
}

func (m *MinioClientInterface) UploadFile(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts map[string]string) (string, error) {
	args := m.Called(ctx, bucketName, objectName, reader, objectSize, opts)
	return args.String(0), args.Error(1)
}

func (m *MinioClientInterface) DeleteFile(ctx context.Context, bucketName, objectName string) error {
	args := m.Called(ctx, bucketName, objectName)
	return args.Error(0)
}

func (m *MinioClientInterface) GetFileURL(bucketName, objectName string) string {
	args := m.Called(bucketName, objectName)
	return args.String(0)
}

func (m *MinioClientInterface) FileExists(ctx context.Context, bucketName, objectName string) (bool, error) {
	args := m.Called(ctx, bucketName, objectName)
	return args.Bool(0), args.Error(1)
}
