package avatar

import (
	"fmt"
	"mime/multipart"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var allowedMIMETypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

type CreateAvatar struct {
	Label  string                `form:"label"`
	Avatar *multipart.FileHeader `form:"avatar"`
}

func (c CreateAvatar) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Label, validation.Required),
		validation.Field(&c.Avatar, validation.By(func(value interface{}) error {
			file, ok := value.(*multipart.FileHeader)
			if !ok {
				return fmt.Errorf("invalid file")
			}
			if file.Size > (2 << 20) {
				return fmt.Errorf("file size should be less than 2MB")
			}

			if !isValidImage(file) {
				return fmt.Errorf("invalid file content")
			}

			return nil
		})),
	)
}

func isValidImage(fileHeader *multipart.FileHeader) bool {
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
