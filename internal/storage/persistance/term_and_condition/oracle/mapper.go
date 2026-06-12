package term_and_condition_oracle

import (
	"time"

	imodel "cbe-super-app-cps-action/internal/constants/model"
)

const tacSelectCols = `
	RAWTOHEX(t.ID), RAWTOHEX(t.ACCOUNT_PRODUCT_ID), TO_CHAR(t.ACTIVATION_TIME, 'YYYY-MM-DD'), t.VERSION_LABEL,
	t.TERMS_AND_CONDITIONS_PATH, t.IS_ENABLED, t.IS_DELETED, t.CREATED_AT, t.LAST_MODIFIED_AT
`

const tacFromTable = ` FROM ACCOUNT_OPENING_TERMS t `

type rowScanner interface {
	Scan(dest ...any) error
}

func scanTACRow(s rowScanner) (*imodel.AccountOpeningTerms, error) {
	var (
		id, accountProductID              string
		activationTime, versionLabel      string
		termsPath                         string
		isEnabled, isDeleted              int
		createdAt, lastModifiedAt         time.Time
	)

	if err := s.Scan(
		&id, &accountProductID, &activationTime, &versionLabel,
		&termsPath, &isEnabled, &isDeleted, &createdAt, &lastModifiedAt,
	); err != nil {
		return nil, err
	}

	return &imodel.AccountOpeningTerms{
		ID:                     id,
		AccountProductID:       accountProductID,
		ActivationTime:         activationTime,
		VersionLabel:           versionLabel,
		TermsAndConditionsPath: termsPath,
		IsEnabled:              isEnabled == 1,
		IsDeleted:              isDeleted == 1,
		CreatedAt:              createdAt,
		LastModifiedAt:         lastModifiedAt,
	}, nil
}
