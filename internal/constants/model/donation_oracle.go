package model

import (
	"time"

	shared_types "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"
)


type DonationOracle struct {
	ID                  string
	DonationCode        string
	ServiceID           string
	CompanyID           string
	CategoryID          string
	Title               string
	IsFeatured          bool
	Target              string
	CurrentAmount       string
	DonationDescription string
	DonationImages      []shared_types.DonationImage
	CoverImage          string
	StartDate           time.Time
	EndDate             time.Time
	Enabled             bool
	IsDeleted           bool
	CreatedAt           time.Time
	LastModifiedAt      time.Time
}
