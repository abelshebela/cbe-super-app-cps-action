package dto

import (
	"fmt"
	"mime/multipart"
	"time"

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
				return fmt.Errorf("invalid file")
			}
			if file.Size > (2 << 20) {
				return fmt.Errorf("file size should be less than 2MB")
			}
			return nil
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
