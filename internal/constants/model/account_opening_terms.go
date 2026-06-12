package model

import "time"

type AccountOpeningTerms struct {
	ID                     string     `bson:"id" json:"id"`
	AccountProductID       string     `bson:"account_product_id" json:"account_product_id"`
	ProductName            string     `bson:"product_name" json:"product_name"`
	ActivationTime         string     `bson:"activation_time" json:"activation_time"`
	VersionLabel           string     `bson:"version_label" json:"version_label"`
	TermsAndConditionsPath string     `bson:"terms_and_conditions_path" json:"terms_and_conditions_path"`
	IsEnabled              bool       `bson:"is_enabled" json:"is_enabled"`
	IsDeleted              bool       `bson:"is_deleted" json:"is_deleted"`
	CreatedAt              time.Time  `bson:"created_at" json:"created_at"`
	LastModifiedAt         time.Time  `bson:"last_modified_at" json:"last_modified_at"`
	DeletedAt              *time.Time `bson:"deleted_at" json:"deleted_at"`
}
