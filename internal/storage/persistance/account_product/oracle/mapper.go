package account_product_oracle

import (
	"database/sql"
	"time"

	imodel "cbe-super-app-cps-action/internal/constants/model"
)

const apSelectCols = `
	RAWTOHEX(ap.ID), ap.CBS_PRODUCT_CODE, ap.PRODUCT_NAME, ap.PRODUCT_TAG_LINE,
	RAWTOHEX(ap.ACCOUNT_CATEGORY_ID), ap.ACCOUNT_CURRENCY,
	ap.MINIMUM_OPENING_BALANCE, ap.MINIMUM_MAINTENANCE_FEE, ap.INTEREST_FEE,
	ap.FAQ_URL, ap.PRODUCT_FEATURES, ap.HAS_PHYSICAL_CARD, ap.HAS_VIRTUAL_CARD,
	ap.PRODUCT_ICON, ap.PRODUCT_COVER_IMAGE, ap.IS_ENABLED, ap.IS_DELETED,
	ap.CREATED_AT, ap.LAST_MODIFIED_AT
`

const apFromTable = ` FROM ACCOUNT_PRODUCTS ap `

type rowScanner interface {
	Scan(dest ...any) error
}

func scanAPRow(s rowScanner) (*imodel.AccountProduct, error) {
	var (
		id, cbsCode, productName, tagLine string
		accountCategoryID, accountCurrency string
		minOpeningBalance, minMaintenanceFee float64
		interestFee                          float64
		faqURL, productFeatures             sql.NullString
		productIcon, productCoverImage      sql.NullString
		hasPhysicalCard, hasVirtualCard     int
		isEnabled, isDeleted                int
		createdAt, lastModifiedAt           time.Time
	)

	if err := s.Scan(
		&id, &cbsCode, &productName, &tagLine,
		&accountCategoryID, &accountCurrency,
		&minOpeningBalance, &minMaintenanceFee, &interestFee,
		&faqURL, &productFeatures,
		&hasPhysicalCard, &hasVirtualCard,
		&productIcon, &productCoverImage,
		&isEnabled, &isDeleted,
		&createdAt, &lastModifiedAt,
	); err != nil {
		return nil, err
	}

	return &imodel.AccountProduct{
		ID:                    id,
		CBSProductCode:        cbsCode,
		ProductName:           productName,
		ProductTagLine:        tagLine,
		AccountCategoryID:     accountCategoryID,
		AccountCurrency:       accountCurrency,
		MinimumOpeningBalance: minOpeningBalance,
		MinimumMaintenanceFee: minMaintenanceFee,
		InterestFee:           interestFee,
		FaqURL:                faqURL.String,
		ProductFeatures:       productFeatures.String,
		HasPhysicalCard:       hasPhysicalCard == 1,
		HasVirtualCard:        hasVirtualCard == 1,
		ProductIcon:           productIcon.String,
		ProductCoverImage:     productCoverImage.String,
		IsEnabled:             isEnabled == 1,
		IsDeleted:             isDeleted == 1,
		CreatedAt:             createdAt,
		LastModifiedAt:        lastModifiedAt,
	}, nil
}
