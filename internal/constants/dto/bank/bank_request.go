package bank_dto

import (
	"mime/multipart"
)

type CreateBankRequest struct {
	Name    string                `form:"name" json:"name" binding:"required"`
	Logo    *multipart.FileHeader `form:"logo" json:"logo" binding:"required"`
	BICCode string                `form:"bic_code" json:"bic_code" binding:"required"`
	Type    string                `form:"type" json:"type" binding:"required"`
}

type UpdateBankRequest struct {
	ID      string                `json:"_id" bson:"_id"`
	Logo    *multipart.FileHeader `form:"logo" json:"logo"`
	Name    string                `json:"name" bson:"name"`
	BICCode string                `json:"bic_code" bson:"bic_code"`
	Type    string                `form:"type" json:"type"`
}

type UpdateLogo struct {
	ID   string                `json:"_id" bson:"_id"`
	Logo *multipart.FileHeader `form:"logo" json:"logo"`
}
