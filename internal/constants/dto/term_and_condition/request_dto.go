package term_and_condition_dto

import "mime/multipart"

type CreateTACRequest struct {
	AccountProductID    string                `json:"product"`
	TermAndCondition    *multipart.FileHeader `json:"term_and_condition"`
	ActivationTime      string                `json:"activation_time"`
	VersionLabel        string                `json:"version_label"`
}

type CPSTACRequest struct {
	ID                     string `json:"id"`
	AccountProductID       string `json:"account_product_id"`
	ActivationTime         string `json:"activation_time"`
	VersionLabel           string `json:"version_label"`
	TermsAndConditionsPath string `json:"terms_and_conditions_path"`
	IsEnabled              bool   `json:"is_enabled"`
	IsDeleted              bool   `json:"is_deleted"`
	CreatedAt              string `json:"created_at"`
	LastModifiedAt         string `json:"last_modified_at"`
}

type TACResponse struct {
	ID                     string `json:"id"`
	AccountProductID       string `json:"account_product_id"`
	ProductName            string `json:"product_name"`
	ActivationTime         string `json:"activation_time"`
	VersionLabel           string `json:"version_label"`
	TermsAndConditionsPath string `json:"terms_and_conditions_path"`
	IsEnabled              bool   `json:"is_enabled"`
	IsDeleted              bool   `json:"is_deleted"`
	CreatedAt              string `json:"created_at"`
	LastModifiedAt         string `json:"last_modified_at"`
}
