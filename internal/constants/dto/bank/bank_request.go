package bank_dto

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)


type CreateBankRequest struct {
	Name string                `form:"name" json:"name"`
	Logo *multipart.FileHeader `form:"logo" json:"logo"`
	Code string                `form:"code" json:"code"`
	BIC  string                `form:"bic" json:"bic"`
}

func (c CreateBankRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required.Error(localization.MsgBankNameRequired), validation.Length(3, 25)),
		validation.Field(&c.Code, validation.Required.Error(localization.MsgBankCodeRequired)),
		validation.Field(&c.BIC, validation.Required.Error(localization.MsgBankBICRequired)),
		validation.Field(&c.Logo, validation.By(func(value any) error {
			file, ok := value.(*multipart.FileHeader)
			if !ok {
				return fmt.Errorf(localization.MsgInvalidInput)
			}

			// Check file size (2 MB max)
			if file.Size > (2 << 20) {
				return validation.NewError("logo", localization.MsgFileTooLarge)
			}

			// Check the file type
			src, err := file.Open()
			if err != nil {
				return fmt.Errorf(localization.MsgServiceUnhandledServerError)
			}
			defer src.Close()

			buffer := make([]byte, 512)
			_, err = src.Read(buffer)
			if err != nil {
				return fmt.Errorf(localization.MsgServiceUnhandledServerError)
			}

			// Reset file pointer to the beginning
			if seeker, ok := src.(io.Seeker); ok {
				_, _ = seeker.Seek(0, io.SeekStart)
			}

			contentType := http.DetectContentType(buffer)
			switch contentType {
			case "image/jpeg", "image/png", "image/gif", "image/webp":
				return nil
			default:
				return fmt.Errorf(localization.MsgFileInvalidType)
			}
		})),
	)
}

type UpdateBankRequest struct {
	ID   string `json:"_id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Code string `json:"code" bson:"code"`
	BIC  string `json:"bic" bson:"bic"`
}

func (u UpdateBankRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name,
			validation.NilOrNotEmpty,
			validation.Length(3, 10),
			is.Alpha,
		),
	)
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

type UpdateLogo struct {
	ID   string                `json:"_id" bson:"_id"`
	Logo *multipart.FileHeader `form:"logo" json:"logo"`
}

func (u UpdateLogo) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Logo, validation.By(func(value interface{}) error {
			file, ok := value.(*multipart.FileHeader)
			if !ok {
				return validation.NewError("logo", localization.MsgInvalidInput)
			}
			if file.Size > maxFileSize {
				return validation.NewError("logo", localization.MsgFileTooLarge)
			}

			// Check MIME type is image/*
			fileObj, err := file.Open()
			if err != nil {
				return validation.NewError("logo", localization.MsgFileInvalidType)
			}
			defer fileObj.Close()
			buffer := make([]byte, 512)
			_, err = fileObj.Read(buffer)
			if err != nil {
				return validation.NewError("logo", localization.MsgFileInvalidType)
			}
			contentType := http.DetectContentType(buffer)
			if !allowedMIMETypes[contentType] {
				return validation.NewError("logo", localization.MsgFileInvalidType)
			}
			if contentType == "" || len(contentType) < 6 || contentType[:6] != "image/" {
				return validation.NewError("logo", localization.MsgFileInvalidType)
			}

			return nil
		})),
	)
}
