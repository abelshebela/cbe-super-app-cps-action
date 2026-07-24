package model

import "time"

type AccountProduct struct {
	ID                    string     `bson:"id" json:"id"`
	CBSProductCode        string     `bson:"cbs_product_code" json:"cbs_product_code"`
	ProductName           string     `bson:"product_name" json:"product_name"`
	ProductTagLine        string     `bson:"product_tag_line" json:"product_tag_line"`
	ProductLine           string     `bson:"product_line" json:"product_line"`
	AccountCategoryID     string     `bson:"account_category_id" json:"account_category_id"`
	CategoryName          string     `bson:"category_name" json:"category_name"`
	CBSCategoryCode       string     `bson:"cbs_category_code" json:"cbs_category_code"`
	AccountType           string     `bson:"-" json:"account_type"`
	AccountCurrency       string     `bson:"account_currency" json:"account_currency"`
	MinimumOpeningBalance float64    `bson:"minimum_opening_balance" json:"minimum_opening_balance"`
	MinimumMaintenanceFee float64    `bson:"minimum_maintenance_fee" json:"minimum_maintenance_fee"`
	InterestFee           float64    `bson:"interest_fee" json:"interest_fee"`
	FaqURL                string     `bson:"faq_url" json:"faq_url"`
	ProductFeatures       []string     `bson:"product_features" json:"product_features"`
	ProductDescription     string                  `json:"product_description"`
	HasPhysicalCard       bool       `bson:"has_physical_card" json:"has_physical_card"`
	HasVirtualCard        bool       `bson:"has_virtual_card" json:"has_virtual_card"`
	ProductIcon           string     `bson:"product_icon" json:"product_icon"`
	ProductCoverImage     string     `bson:"product_cover_image" json:"product_cover_image"`
	IsAvailableForOnbording bool                 `json:"is_available_for_onbording"`
	InterestRate            float64               `json:"interest_rate"`
	IsEnabled             bool       `bson:"is_enabled" json:"is_enabled"`
	IsDeleted             bool       `bson:"is_deleted" json:"is_deleted"`
	CreatedAt             time.Time  `bson:"created_at" json:"created_at"`
	LastModifiedAt        time.Time  `bson:"last_modified_at" json:"last_modified_at"`
	DeletedAt             *time.Time `bson:"deleted_at" json:"deleted_at"`
	Eligibility           string      `json:"eligibility"`
}
