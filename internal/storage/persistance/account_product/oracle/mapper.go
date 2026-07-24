package account_product_oracle

import (
	"database/sql"
	"time"

	imodel "cbe-super-app-cps-action/internal/constants/model"
)

const apSelectCols = `
	RAWTOHEX(ap.ID), ap.CBS_PRODUCT_CODE, ap.PRODUCT_NAME, ap.PRODUCT_TAG_LINE,
	RAWTOHEX(ap.ACCOUNT_CATEGORY_ID), NVL(ac.CATEGORY_NAME, ''), NVL(ac.CBS_CATEGORY_CODE, ''), NVL(ac.ACCOUNT_TYPE, ''), ap.ACCOUNT_CURRENCY,
	ap.MINIMUM_OPENING_BALANCE, ap.MINIMUM_MAINTENANCE_FEE, ap.INTEREST_FEE,
	ap.FAQ_URL, ap.PRODUCT_FEATURES, ap.HAS_PHYSICAL_CARD, ap.HAS_VIRTUAL_CARD,
	ap.PRODUCT_ICON, ap.PRODUCT_COVER_IMAGE, ap.IS_ENABLED, ap.IS_DELETED,
	ap.CREATED_AT, ap.LAST_MODIFIED_AT,ap.IS_AVAILABE_FOR_ONBORDING,ap.PRODUCT_DESCRIPTION,ap.ELIGIBILITY
`

const apFromTable = ` FROM ACCOUNT_PRODUCTS ap LEFT JOIN ACCOUNT_CATEGORIES ac ON ap.ACCOUNT_CATEGORY_ID = ac.ID `

type rowScanner interface {
	Scan(dest ...any) error
}

func scanAPRow(s rowScanner) (*imodel.AccountProduct, error) {
	var (
		id, cbsCode, productName, tagLine                             string
		accountCategoryID, categoryName, cbsCategoryCode, accountType string
		accountCurrency                                               string
		minOpeningBalance, minMaintenanceFee                          float64
		interestFee                                                   float64
		faqURL, productFeatures, product_description,eligibility                 sql.NullString
		productIcon, productCoverImage                                sql.NullString
		hasPhysicalCard, hasVirtualCard                               int
		isEnabled, isDeleted, isAvailableForOnbording                 int
		createdAt, lastModifiedAt                                     time.Time
	)

	if err := s.Scan(
		&id, &cbsCode, &productName, &tagLine,
		&accountCategoryID, &categoryName, &cbsCategoryCode, &accountType, &accountCurrency,
		&minOpeningBalance, &minMaintenanceFee, &interestFee,
		&faqURL, &productFeatures,
		&hasPhysicalCard, &hasVirtualCard,
		&productIcon, &productCoverImage,
		&isEnabled, &isDeleted,
		&createdAt, &lastModifiedAt, &isAvailableForOnbording, &product_description,&eligibility,
	); err != nil {
		return nil, err
	}

	product_features_list, err := ParseOverviewHTML(productFeatures.String)
	if err != nil {
		return nil, err
	}

	// if err := json.Unmarshal([]byte(productFeatures.String), &product_features) ; err != nil{
	// 	return nil,err
	// }

	return &imodel.AccountProduct{
		ID:                      id,
		CBSProductCode:          cbsCode,
		ProductName:             productName,
		ProductTagLine:          tagLine,
		AccountCategoryID:       accountCategoryID,
		CategoryName:            categoryName,
		CBSCategoryCode:         cbsCategoryCode,
		AccountType:             accountType,
		AccountCurrency:         accountCurrency,
		MinimumOpeningBalance:   minOpeningBalance,
		MinimumMaintenanceFee:   minMaintenanceFee,
		InterestFee:             interestFee,
		FaqURL:                  faqURL.String,
		ProductFeatures:         product_features_list,
		HasPhysicalCard:         hasPhysicalCard == 1,
		HasVirtualCard:          hasVirtualCard == 1,
		ProductIcon:             productIcon.String,
		ProductCoverImage:       productCoverImage.String,
		IsEnabled:               isEnabled == 1,
		IsDeleted:               isDeleted == 1,
		CreatedAt:               createdAt,
		LastModifiedAt:          lastModifiedAt,
		IsAvailableForOnbording: isAvailableForOnbording == 1,
		ProductDescription:      product_description.String,
		Eligibility: eligibility.String,
		// InterestRate: interestRate,
	}, nil
}
