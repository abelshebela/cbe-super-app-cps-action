package ussd_merchant_dto

import "mime/multipart"

type CreateUssdMerchantRequest struct {
	SettlementMethod string                `json:"settlement_method"`
	Name             string                `json:"name"`
	PhoneNumber      string                `json:"phone_number"`
	Email            string                `json:"email"`
	Service          string                `json:"service"`
	AccountNumber    string                `json:"account_number"`
	Logo             *multipart.FileHeader `json:"logo"`
}

type UpdateUssdMerchantRequest struct {
	SettlementMethod string                `json:"settlement_method"`
	Name             string                `json:"name"`
	PhoneNumber      string                `json:"phone_number"`
	Email            string                `json:"email"`
	Service          string                `json:"service"`
	AccountNumber    string                `json:"account_number"`
	Logo             *multipart.FileHeader `json:"logo"`
}
