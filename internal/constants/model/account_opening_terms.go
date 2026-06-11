package model

import "time"

type AccountOpeningTerms struct {
	ID                     string
	AccountProductID       string
	ActivationTime         string
	VersionLabel           string
	TermsAndConditionsPath string
	IsEnabled              bool
	IsDeleted              bool
	CreatedAt              time.Time
	LastModifiedAt         time.Time
	DeletedAt              *time.Time
}
