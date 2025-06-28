package config

import (
	"context"
	"regexp"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ContentType string

const (
	ContentTypeJSON              ContentType = "application/json"
	ContentTypeXML               ContentType = "application/xml"
	ContentTypePDF               ContentType = "application/pdf"
	ContentTypeZIP               ContentType = "application/zip"
	ContentTypeOctetStream       ContentType = "application/octet-stream"
	ContentTypeMSWord            ContentType = "application/msword"
	ContentTypeExcel             ContentType = "application/vnd.ms-excel"
	ContentTypeExcelOpenXML      ContentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	ContentTypePlainText         ContentType = "text/plain"
	ContentTypeHTML              ContentType = "text/html"
	ContentTypeCSV               ContentType = "text/csv"
	ContentTypeJPEG              ContentType = "image/jpeg"
	ContentTypePNG               ContentType = "image/png"
	ContentTypeGIF               ContentType = "image/gif"
	ContentTypeWEBP              ContentType = "image/webp"
	ContentTypeMP4               ContentType = "video/mp4"
	ContentTypeMP3               ContentType = "audio/mpeg"
	ContentTypeMultipartFormData ContentType = "multipart/form-data"
)

type MinioClientInterface interface {
	ListBuckets(ctx context.Context) ([]minio.BucketInfo, error)
	MakeBucket(ctx context.Context, bucket_name string) (bool, error)
	BucketExist(ctx context.Context, bucket_name string) (bool, error)
	SaveObject(ctx context.Context, obj SaveObjectBody) (*SaveObjectResponse, error)
	DeleteObject(ctx context.Context, obj DeleteObjectBody) (bool, error)
}

type DeleteObjectBody struct {
	BucketName string
	ObjectName string
}

type SaveObjectBody struct {
	BucketName  string
	ObjectName  string
	File        string
	ContentType ContentType
}

type SaveObjectResponse struct {
	Key       string
	Bucket    string
	ETag      string
	VersionID string
	Size      int64
}

type MinioClient struct {
	client *minio.Client
}

func NewMinioClient(env *VaultConfig) (MinioClientInterface, error) {
	client, err := minio.New(env.MinioEndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(env.MinioAccessKey, env.MinioSecretKey, ""),
		Secure: false,
	})

	if err != nil {
		return nil, err
	}

	return &MinioClient{
		client: client,
	}, nil
}

func (m *MinioClient) ListBuckets(ctx context.Context) ([]minio.BucketInfo, error) {
	bucket, err := m.client.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}

func (m *MinioClient) MakeBucket(ctx context.Context, bucket_name string) (bool, error) {
	regex := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	bucket := regex.ReplaceAllString(bucket_name, "-")
	err := m.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})

	if err != nil {
		return false, err
	}

	found, err := m.BucketExist(ctx, bucket)
	if err != nil {
		return false, nil
	}
	return found, nil
}

func (m *MinioClient) BucketExist(ctx context.Context, bucket_name string) (bool, error) {
	exists, err := m.client.BucketExists(ctx, bucket_name)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (m *MinioClient) SaveObject(ctx context.Context, obj SaveObjectBody) (*SaveObjectResponse, error) {
	ui, err := m.client.FPutObject(
		ctx,
		obj.BucketName,
		obj.ObjectName,
		obj.File,
		minio.PutObjectOptions{
			ContentType: string(obj.ContentType),
		},
	)

	if err != nil {
		return nil, err
	}

	return &SaveObjectResponse{
		Bucket:    ui.Bucket,
		Key:       ui.Key,
		ETag:      ui.ETag,
		VersionID: ui.VersionID,
		Size:      ui.Size,
	}, nil
}

func (m *MinioClient) DeleteObject(ctx context.Context, obj DeleteObjectBody) (bool, error) {
	if err := m.client.RemoveObject(
		ctx,
		obj.BucketName,
		obj.ObjectName,
		minio.RemoveObjectOptions{},
	); err != nil {
		return false, err
	}

	return true, nil
}
