package bank_dto

import (
	"mime/multipart"
)

type CreateBankRequest struct {
	Name            string                `form:"name" json:"name" binding:"required"`
	Logo            *multipart.FileHeader `form:"logo" json:"logo" binding:"required" swaggertype:"string" format:"binary"`
	BICCode         string                `form:"bic_code" json:"bic_code" binding:"required"`
	HasAlphaNumeric *bool                 `form:"has_alpha_numeric" json:"has_alpha_numeric"`
	AccountLength   int                   `form:"account_length" json:"account_length"`
}

type UpdateBankRequest struct {
	ID              string                `json:"_id" bson:"_id"`
	Logo            *multipart.FileHeader `form:"logo" json:"logo" swaggertype:"string" format:"binary"`
	Name            string                `json:"name" bson:"name"`
	BICCode         string                `json:"bic_code" bson:"bic_code"`
	HasAlphaNumeric *bool                 `form:"has_alpha_numeric" json:"has_alpha_numeric"`
	AccountLength   int                   `form:"account_length" json:"account_length"`
}

type UpdateLogo struct {
	ID   string                `json:"_id" bson:"_id"`
	Logo *multipart.FileHeader `form:"logo" json:"logo" swaggertype:"string" format:"binary"`
}
