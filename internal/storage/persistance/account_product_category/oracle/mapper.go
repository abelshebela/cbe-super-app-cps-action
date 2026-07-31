package account_product_category_oracle

import (
	"database/sql"
	"time"

	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
)

const apcSelectCols = `
	RAWTOHEX(ID), ACCOUNT_TYPE, CATEGORY_NAME, CBS_CATEGORY_CODE,
	DESCRIPTION, IS_ENABLED, IS_DELETED, CREATED_AT, LAST_MODIFIED_AT
`

const apcFromTable = ` FROM ACCOUNT_CATEGORIES `

type rowScanner interface {
	Scan(dest ...any) error
}

func scanAPCRow(s rowScanner) (*imodel.AccountProductCategory, error) {
	var (
		id, accountType, categoryName, cbsCategoryCode string
		description                                     sql.NullString
		isEnabled, isDeleted                            int
		createdAt, lastModifiedAt                       time.Time
	)

	if err := s.Scan(
		&id, &accountType, &categoryName, &cbsCategoryCode,
		&description, &isEnabled, &isDeleted, &createdAt, &lastModifiedAt,
	); err != nil {
		return nil, err
	}

	return &imodel.AccountProductCategory{
		ID:              id,
		AccountType:     accountType,
		CategoryName:    categoryName,
		CBSCategoryCode: cbsCategoryCode,
		Description:     description.String,
		IsEnabled:       isEnabled == 1,
		IsDeleted:       isDeleted == 1,
		CreatedAt:       createdAt,
		LastModifiedAt:  lastModifiedAt,
	}, nil
}
