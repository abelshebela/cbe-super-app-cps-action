package dto

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type BankResponse struct {
	ID             string    `json:"id" bson:"id"`
	Name           string    `json:"name" bson:"name"`
	Logo           string    `json:"logo" bson:"logo"`
	Code           string    `json:"code" bson:"code"`
	BIC            string    `json:"bic"  bson:"bic"`
	Enabled        bool      `json:"enabled" bson:"enabled"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time `json:"last_modified_at" bson:"last_modified_at"`
}

type CreateBankRequest struct {
	Name string                `form:"name" json:"name"`
	Logo *multipart.FileHeader `form:"logo" json:"logo"`
	Code string                `form:"code" json:"code"`
	BIC  string                `form:"bic" json:"bic"`
}

func (c CreateBankRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required.Error(error_codes.MissingBankName), validation.Length(3, 25)),
		validation.Field(&c.Code, validation.Required.Error(error_codes.MissingBankCode)),
		validation.Field(&c.BIC, validation.Required.Error(error_codes.MissingBankBIC)),
		validation.Field(&c.Logo, validation.By(func(value any) error {
			file, ok := value.(*multipart.FileHeader)
			if !ok {
				return fmt.Errorf(error_codes.InvalidInput)
			}

			// Check file size (2 MB max)
			if file.Size > (2 << 20) {
				return validation.NewError("logo", error_codes.FileTooLarge)
			}

			// Check the file type
			src, err := file.Open()
			if err != nil {
				return fmt.Errorf(error_codes.UnhandledServerError)
			}
			defer src.Close()

			buffer := make([]byte, 512)
			_, err = src.Read(buffer)
			if err != nil {
				return fmt.Errorf(error_codes.UnhandledServerError)
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
				return fmt.Errorf(error_codes.InvalidFileType)
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
				return validation.NewError("logo", error_codes.InvalidInput)
			}
			if file.Size > maxFileSize {
				return validation.NewError("logo", error_codes.FileTooLarge)
			}

			// Check MIME type is image/*
			fileObj, err := file.Open()
			if err != nil {
				return validation.NewError("logo", error_codes.InvalidFileType)
			}
			defer fileObj.Close()
			buffer := make([]byte, 512)
			_, err = fileObj.Read(buffer)
			if err != nil {
				return validation.NewError("logo", error_codes.InvalidFileType)
			}
			contentType := http.DetectContentType(buffer)
			if !allowedMIMETypes[contentType] {
				return validation.NewError("logo", error_codes.InvalidFileType)
			}
			if contentType == "" || len(contentType) < 6 || contentType[:6] != "image/" {
				return validation.NewError("logo", error_codes.InvalidFileType)
			}

			return nil
		})),
	)
}
