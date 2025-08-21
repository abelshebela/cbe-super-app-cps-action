package wallet

import (
	"mime/multipart"
	"net/http"

	utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	// error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type WalletRequest struct {
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

func (c WalletRequest) Validate(isCreate bool) error {
	var rules = []*validation.FieldRules{}

	if isCreate {
		rules = append(rules, validation.Field(&c.Name,
			validation.Required.Error("name is required"),
			validation.Length(3, 10),
			validation.By(utils.NoSpecialChars),
		))
	} else if c.Name != "" {
		rules = append(rules, validation.Field(&c.Name,
			validation.Length(3, 10),
			validation.By(utils.NoSpecialChars),
			validation.By(func(value interface{}) error {
				if str, ok := value.(string); ok {
					return common_util.ValidateInputNoSpecialChars(str)
				}
				return nil
			}),
		))
		rules = append(rules, validation.Field(&c.Code,
			validation.By(utils.NoSpecialChars)))
	}

	if isCreate {
		// Avatar is required on create
		rules = append(rules, validation.Field(&c.Avatar, validation.Required, validation.By(func(value any) error {
			file, ok := value.(*multipart.FileHeader)
			if !ok {
				return validation.NewError("avatar", utils.InvalidInput)
			}
			if file.Size > maxFileSize {
				return validation.NewError("avatar", utils.FileTooLarge)
			}

			if !IsValidImage(file) {
				return validation.NewError("avatar", utils.InvalidFileType)
			}
			return nil
		})))
	} else if c.Avatar != nil {
		// On update, if Avatar provided, validate file
		rules = append(rules, validation.Field(&c.Avatar, validation.By(func(value any) error {
			file, ok := value.(*multipart.FileHeader)
			if !ok {
				return validation.NewError("avatar", utils.InvalidInput)
			}
			if file.Size > maxFileSize {
				return validation.NewError("avatar", utils.FileTooLarge)
			}

			if !IsValidImage(file) {
				return validation.NewError("avatar", utils.InvalidFileType)
			}
			return nil
		})))
	}

	return validation.ValidateStruct(&c, rules...)
}

func IsEmpty(value WalletRequest) bool {
	return value.Name == "" && value.Code == "" && value.Avatar == nil
}
