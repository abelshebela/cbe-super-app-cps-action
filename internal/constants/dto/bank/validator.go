package bank_dto

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (c CreateBankRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required.Error(localization.MsgBankNameRequired), validation.Length(3, 25), validation.Match(regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)).Error("Name must not be contain special char")),
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

func (u UpdateBankRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name,
			validation.NilOrNotEmpty,
			validation.Length(3, 10),
			is.Alpha,
		),
		validation.Field(&u.Code, validation.Required.Error(localization.MsgBankCodeRequired)),
		validation.Field(&u.BIC, validation.Required.Error(localization.MsgBankBICRequired)),
	)
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
