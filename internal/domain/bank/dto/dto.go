package dto

import (
	"fmt"
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
		validation.Field(&c.Name, validation.Required.Error("name is required"), validation.Length(3, 10), is.Alpha),
		validation.Field(&c.Code, validation.Required.Error("code is required")),
		validation.Field(&c.BIC, validation.Required.Error("bank identifier code is required")),
		validation.Field(&c.Logo, validation.By(func(value interface{}) error {
			file, ok := value.(*multipart.FileHeader)

			if !ok {
				return fmt.Errorf(error_codes.InvalidInput)
			}
			// Check file size (2 MB max)
			if file.Size > (2 << 20) {
				return fmt.Errorf(error_codes.FileTooLarge)
			}

			// check the file type
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
	ID   string `json:"id"`
	Name string `json:"name"`
	// Logo *string
	Code string `json:"code"`
	BIC  string `json:"bic"`
}
