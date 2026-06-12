package model

import "time"

type AccountProductCategory struct {
	ID              string     `bson:"id" json:"id"`
	AccountType     string     `bson:"account_type" json:"account_type"`
	CategoryName    string     `bson:"category_name" json:"category_name"`
	CBSCategoryCode string     `bson:"cbs_category_code" json:"cbs_category_code"`
	Description     string     `bson:"description" json:"description"`
	IsEnabled       bool       `bson:"is_enabled" json:"is_enabled"`
	IsDeleted       bool       `bson:"is_deleted" json:"is_deleted"`
	CreatedAt       time.Time  `bson:"created_at" json:"created_at"`
	LastModifiedAt  time.Time  `bson:"last_modified_at" json:"last_modified_at"`
	DeletedAt       *time.Time `bson:"deleted_at" json:"deleted_at"`
}
