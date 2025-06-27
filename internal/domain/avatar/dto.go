package avatar

import (
	"fmt"
	"mime/multipart"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

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
			return nil
		})),
	)
}
