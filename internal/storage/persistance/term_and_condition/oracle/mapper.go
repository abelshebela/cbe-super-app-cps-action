package term_and_condition_oracle

import (
	"time"

	imodel "cbe-super-app-cps-action/internal/constants/model"
)

const tacSelectCols = `
	RAWTOHEX(t.ID), RAWTOHEX(t.ACCOUNT_PRODUCT_ID), NVL(ap.PRODUCT_NAME, ''),
	TO_CHAR(t.ACTIVATION_TIME, 'YYYY-MM-DD'), t.VERSION_LABEL,
	t.TERMS_AND_CONDITIONS_PATH, t.IS_ENABLED, t.IS_DELETED, t.CREATED_AT, t.LAST_MODIFIED_AT
`

const tacFromTable = ` FROM ACCOUNT_OPENING_TERMS t LEFT JOIN ACCOUNT_PRODUCTS ap ON t.ACCOUNT_PRODUCT_ID = ap.ID `

type rowScanner interface {
	Scan(dest ...any) error
}

func scanTACRow(s rowScanner) (*imodel.AccountOpeningTerms, error) {
	var (
		id, accountProductID      string
		productName               string
		activationTime, versionLabel string
		termsPath                 string
		isEnabled, isDeleted      int
		createdAt, lastModifiedAt time.Time
	)

	if err := s.Scan(
		&id, &accountProductID, &productName,
		&activationTime, &versionLabel,
		&termsPath, &isEnabled, &isDeleted, &createdAt, &lastModifiedAt,
	); err != nil {
		return nil, err
	}

	return &imodel.AccountOpeningTerms{
		ID:                     id,
		AccountProductID:       accountProductID,
		ProductName:            productName,
		ActivationTime:         activationTime,
		VersionLabel:           versionLabel,
		TermsAndConditionsPath: termsPath,
		IsEnabled:              isEnabled == 1,
		IsDeleted:              isDeleted == 1,
		CreatedAt:              createdAt,
		LastModifiedAt:         lastModifiedAt,
	}, nil
}
