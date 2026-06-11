package model

import "time"

type AccountProductCategory struct {
	ID              string
	AccountType     string
	CategoryName    string
	CBSCategoryCode string
	Description     string
	IsEnabled       bool
	IsDeleted       bool
	CreatedAt       time.Time
	LastModifiedAt  time.Time
	DeletedAt       *time.Time
}
