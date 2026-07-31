package donation_category_oracle

import (
	"database/sql"
	"time"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/donation_category"
)

const donationCategorySelectCols = `RAWTOHEX(ID), CATEGORY_NAME, ICON, IS_DELETED, ENABLED, CREATED_AT, LAST_MODIFIED_AT`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanDonationCategoryRow(s rowScanner) (*donation_category.DonationCategoryListResponse, error) {
	var (
		id                       string
		categoryName             string
		icon                     sql.NullString
		isDeleted, enabled       int
		createdAt, lastModifiedAt time.Time
	)

	if err := s.Scan(&id, &categoryName, &icon, &isDeleted, &enabled, &createdAt, &lastModifiedAt); err != nil {
		return nil, err
	}

	return &donation_category.DonationCategoryListResponse{
		ID:             id,
		CategoryName:   categoryName,
		Icon:           icon.String,
		IsDeleted:      isDeleted == 1,
		Enabled:        enabled == 1,
		CreatedAt:      createdAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		LastModifiedAt: lastModifiedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
