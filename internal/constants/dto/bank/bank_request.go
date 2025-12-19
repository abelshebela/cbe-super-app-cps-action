package bank_dto

import (
	"cbe-super-app-cps-action/internal/constants"
	"mime/multipart"
)

type CreateBankRequest struct {
	Name string                             `form:"name" json:"name" binding:"required"`
	Logo *multipart.FileHeader              `form:"logo" json:"logo" binding:"required"`
	Code string                             `form:"code" json:"code" binding:"required"`
	BIC  string                             `form:"bic" json:"bic" binding:"required"`
	Type constants.FinancialInstitutionType `form:"type" json:"type" binding:"required"`
}

type UpdateBankRequest struct {
	ID   string                             `json:"_id" bson:"_id"`
	Logo *multipart.FileHeader              `form:"logo" json:"logo"`
	Name string                             `json:"name" bson:"name"`
	Code string                             `json:"code" bson:"code"`
	BIC  string                             `json:"bic" bson:"bic"`
	Type constants.FinancialInstitutionType `form:"type" json:"type"`
}

type UpdateLogo struct {
	ID   string                `json:"_id" bson:"_id"`
	Logo *multipart.FileHeader `form:"logo" json:"logo"`
}
