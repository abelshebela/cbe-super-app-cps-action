package bank_dto

import (
	"mime/multipart"
)

type CreateBankRequest struct {
	Name string                `form:"name" json:"name"`
	Logo *multipart.FileHeader `form:"logo" json:"logo"`
	Code string                `form:"code" json:"code"`
	BIC  string                `form:"bic" json:"bic"`
}

type UpdateBankRequest struct {
	ID   string                `json:"_id" bson:"_id"`
	Logo *multipart.FileHeader `form:"logo" json:"logo"`
	Name string                `json:"name" bson:"name"`
	Code string                `json:"code" bson:"code"`
	BIC  string                `json:"bic" bson:"bic"`
}

type UpdateLogo struct {
	ID   string                `json:"_id" bson:"_id"`
	Logo *multipart.FileHeader `form:"logo" json:"logo"`
}
