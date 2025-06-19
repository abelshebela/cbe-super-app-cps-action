package dto

import (
	"fmt"
	"mime/multipart"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateWalletRequest struct {
	Name   string                `form:"name" json:"name"`
	Avatar *multipart.FileHeader `form:"avatar" json:"avatar"`
	Code   string                `form:"code" json:"code"`
}

func (c CreateWalletRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required.Error("name is required"), validation.Length(3, 10), is.Alpha),
		validation.Field(&c.Code, validation.Required.Error("code is required")),
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

type UpdateWalletRequest struct {
	ID    string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}
