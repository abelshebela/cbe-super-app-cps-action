package model

import "time"

type AccountProduct struct {
	ID                       string
	CBSProductCode           string
	ProductName              string
	ProductTagLine           string
	ProductLine              string
	AccountCategoryID        string
	AccountCurrency          string
	MinimumOpeningBalance    float64
	MinimumMaintenanceFee    float64
	InterestFee              float64
	FaqURL                   string
	ProductFeatures          string
	HasPhysicalCard          bool
	HasVirtualCard           bool
	ProductIcon              string
	ProductCoverImage        string
	IsEnabled                bool
	IsDeleted                bool
	CreatedAt                time.Time
	LastModifiedAt           time.Time
	DeletedAt                *time.Time
}
