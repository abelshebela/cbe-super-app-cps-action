package avatar

import (
	"fmt"
	"mime/multipart"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const maxFileSize = 2 * 1024 * 1024 // 2MB

var allowedMIMETypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

var IsValidImage = func(fileHeader *multipart.FileHeader) bool {
	file, err := fileHeader.Open()
	if err != nil {
		return false
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return false
	}
	contentType := http.DetectContentType(buffer)
	return allowedMIMETypes[contentType]
}

type CreateAvatar struct {
	Label  string                `form:"label"`
	Avatar *multipart.FileHeader `form:"avatar"`
}

func (c CreateAvatar) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Label, validation.Required),
		validation.Field(&c.Avatar, validation.Required, validation.By(func(value interface{}) error {
			file, ok := value.(*multipart.FileHeader)
			if !ok {
				return fmt.Errorf("invalid file")
			}
			if file.Size > maxFileSize {
				return fmt.Errorf("file size should be less than 2MB")
			}

			if !IsValidImage(file) {
				return fmt.Errorf("invalid file content")
			}

			return nil
		})),
	)
}

type UpdateAvatar struct {
	Avatar *multipart.FileHeader `form:"avatar"`
}

func (u UpdateAvatar) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Avatar, validation.Required, validation.By(func(value interface{}) error {
			file, ok := value.(*multipart.FileHeader)
			if !ok {
				return fmt.Errorf("invalid file")
			}
			if file.Size > maxFileSize {
				return fmt.Errorf("file size should be less than 2MB")
			}

			if !IsValidImage(file) {
				return fmt.Errorf("invalid file content")
			}

			return nil
		})),
	)
}

type AvatarResponse struct {
	Page    int       `json:"page"`
	Avatars []*Avatar `json:"avatars"`
	Limit   int       `json:"limit"`
	Total   int64     `json:"total"`
}
