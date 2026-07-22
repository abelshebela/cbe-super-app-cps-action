package account_product_dto

import "mime/multipart"

type CreateAPRequest struct {
	CBSProductCode        string                `json:"cps_product_code"`
	ProductName           string                `json:"product_name"`
	ProductTagLine        string                `json:"product_tagline"`
	AccountCategoryID     string                `json:"account_category"`
	AccountCurrency       string                `json:"account_currency"`
	MinimumOpeningBalance float64               `json:"minimum_opening_balance"`
	MinimumMaintenanceFee float64               `json:"minimum_balance_to_maintain_account"`
	InterestFee           float64               `json:"interest_fee"`
	FaqURL                string                `json:"faq_url"`
	ProductFeatures       string                `json:"product_features"`
	HasPhysicalCard       bool                  `json:"has_atm_and_debit_card"`
	HasVirtualCard        bool                  `json:"has_virtual_debit_card"`
	Icon                  *multipart.FileHeader `json:"icon"`
	CoverImage            *multipart.FileHeader `json:"cover_image"`
	IsAvailableForOnbording bool                 `json:"is_available_for_onbording"`
	InterestRate            float64               `json:"interest_rate"`
}

type UpdateAPRequest struct {
	CBSProductCode        string                `json:"cps_product_code"`
	ProductName           string                `json:"product_name"`
	ProductTagLine        string                `json:"product_tagline"`
	AccountCategoryID     string                `json:"account_category"`
	AccountCurrency       string                `json:"account_currency"`
	MinimumOpeningBalance float64               `json:"minimum_opening_balance"`
	MinimumMaintenanceFee float64               `json:"minimum_balance_to_maintain_account"`
	InterestFee           float64               `json:"interest_fee"`
	FaqURL                string                `json:"faq_url"`
	ProductFeatures       string                `json:"product_features"`
	HasPhysicalCard       *bool                 `json:"has_atm_and_debit_card"`
	HasVirtualCard        *bool                 `json:"has_virtual_debit_card"`
	Icon                  *multipart.FileHeader `json:"icon"`
	CoverImage            *multipart.FileHeader `json:"cover_image"`
	IsAvailableForOnbording *bool                 `json:"is_available_for_onbording"`
	InterestRate            float64               `json:"interest_rate"`
}

type CPSAPRequest struct {
	ID                    string  `json:"id"`
	CBSProductCode        string  `json:"cbs_product_code"`
	ProductName           string  `json:"product_name"`
	ProductTagLine        string  `json:"product_tag_line"`
	ProductLine           string  `json:"product_line"`
	AccountCategoryID     string  `json:"account_category_id"`
	AccountCurrency       string  `json:"account_currency"`
	MinimumOpeningBalance float64 `json:"minimum_opening_balance"`
	MinimumMaintenanceFee float64 `json:"minimum_maintenance_fee"`
	InterestFee           float64 `json:"interest_fee"`
	FaqURL                string  `json:"faq_url"`
	ProductFeatures       string  `json:"product_features"`
	HasPhysicalCard       bool    `json:"has_physical_card"`
	HasVirtualCard        bool    `json:"has_virtual_card"`
	ProductIcon           string  `json:"product_icon"`
	ProductCoverImage     string  `json:"product_cover_image"`
	IsAvailableForOnbording bool                 `json:"is_available_for_onbording"`
	InterestRate            float64               `json:"interest_rate"`
	IsEnabled             bool    `json:"is_enabled"`
	IsDeleted             bool    `json:"is_deleted"`
	CreatedAt             string  `json:"created_at"`
	LastModifiedAt        string  `json:"last_modified_at"`
}

type APAccountCategory struct {
	ID              string `json:"id"`
	CategoryName    string `json:"category_name"`
	CBSCategoryCode string `json:"cbs_category_code"`
	AccountType     string `json:"account_type"`
}

type APResponse struct {
	ID                    string            `json:"id"`
	CBSProductCode        string            `json:"cbs_product_code"`
	ProductName           string            `json:"product_name"`
	ProductTagLine        string            `json:"product_tag_line"`
	AccountCategory       APAccountCategory `json:"account_category"`
	AccountCurrency       string            `json:"account_currency"`
	MinimumOpeningBalance float64           `json:"minimum_opening_balance"`
	MinimumMaintenanceFee float64           `json:"minimum_maintenance_fee"`
	InterestFee           float64           `json:"interest_fee"`
	FaqURL                string            `json:"faq_url"`
	ProductFeatures       string            `json:"product_features"`
	HasPhysicalCard       bool              `json:"has_physical_card"`
	HasVirtualCard        bool              `json:"has_virtual_card"`
	ProductIcon           string            `json:"product_icon"`
	ProductCoverImage     string            `json:"product_cover_image"`
	IsAvailableForOnbording bool                 `json:"is_available_for_onbording"`
	InterestRate            float64               `json:"interest_rate"`
	IsEnabled             bool              `json:"is_enabled"`
	IsDeleted             bool              `json:"is_deleted"`
	CreatedAt             string            `json:"created_at"`
	LastModifiedAt        string            `json:"last_modified_at"`
}
