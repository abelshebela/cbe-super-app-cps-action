// Package dto provides data transfer objects for wallet operations.
package dto

import (
	"mime/multipart"
	"net/http"

	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateWalletRequest struct {
	Name   string                `form:"name" json:"name"`
	Avatar *multipart.FileHeader `form:"avatar" json:"avatar"`
	Code   string                `form:"code" json:"code"`
}

const maxFileSize = 2 * 1024 * 1024 // 2MB

var allowedMIMETypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
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

func (c CreateWalletRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required.Error("name is required"), validation.Length(3, 10)),
		validation.Field(&c.Code, validation.Required.Error("code is required")),
		validation.Field(&c.Avatar, validation.By(func(value any) error {
			file, ok := value.(*multipart.FileHeader)
			if !ok {
				return validation.NewError("avatar", error_codes.InvalidInput)
			}
			if file.Size > maxFileSize {
				return validation.NewError("avatar", error_codes.FileTooLarge)
			}

			if !IsValidImage(file) {
				return validation.NewError("avatar", error_codes.InvalidFileType)
			}
			return nil
		})),
	)
}

type UpdateWalletRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}
