package outbound

import (
	"context"
	"mime/multipart"
)

type MinIOUploader interface {
	UploadProfileImage(ctx context.Context, objectName string, file multipart.File, contentType string, userID string) (string, error)
}
