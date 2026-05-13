package donation_oracle

import (
	"database/sql"
	"time"

	donation_dto "cbe-super-app-cps-action/internal/constants/dto/donation"
	"cbe-super-app-cps-action/internal/constants/types"
)

const donationListSelectCols = `
	RAWTOHEX(d.ID), d.DONATION_CODE, d.TITLE, d.IS_FEATURED,
	d.TARGET, d.CURRENT_AMOUNT, d.DONATION_DESCRIPTION, d.COVER_IMAGE,
	d.START_DATE, d.END_DATE, d.IS_DELETED, d.ENABLED,
	d.CREATED_AT, d.LAST_MODIFIED_AT,
	RAWTOHEX(d.SERVICE_ID), al.NAME as SERVICE_NAME, al.SERVICE_KEY as SERVICE_KEY,
	RAWTOHEX(d.COMPANY_ID), c.COMPANY_NAME, c.COMPANY_LOGO, c.ENABLED,
	s.PRODUCT_GL_ACCOUNT_NUMBER as ACCOUNT_NUMBER,
	RAWTOHEX(d.CATEGORY_ID), cat.CATEGORY_NAME, cat.ICON
`

// INNER JOIN: donations with a dangling service/company/category FK are excluded.
const donationListFromJoin = `
	FROM DONATIONS d
	INNER JOIN SERVICES s            ON s.ID   = d.SERVICE_ID
	INNER JOIN ACCESS_LISTS al            ON   al.ID = s.ACCESS_LIST_ID 
	INNER JOIN DONATION_COMPANIES c  ON c.ID   = d.COMPANY_ID
	INNER JOIN DONATION_CATEGORIES cat ON cat.ID = d.CATEGORY_ID
`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanDonationListRow(s rowScanner) (*donation_dto.DonationListResponse, error) {
	var (
		id, donationCode, title        string
		isFeatured, isDeleted, enabled int
		target, currentAmount          sql.NullString
		description, coverImage        sql.NullString
		startDate                      time.Time
		endDate                        sql.NullTime
		createdAt, lastModifiedAt      time.Time

		serviceID, serviceName, serviceKey sql.NullString

		companyID, companyName sql.NullString
		companyLogo            sql.NullString
		companyEnabled         sql.NullInt64

		categoryID, categoryName sql.NullString
		categoryIcon             sql.NullString
	)

	if err := s.Scan(
		&id, &donationCode, &title, &isFeatured,
		&target, &currentAmount, &description, &coverImage,
		&startDate, &endDate, &isDeleted, &enabled,
		&createdAt, &lastModifiedAt,
		&serviceID, &serviceName, &serviceKey,
		&companyID, &companyName, &companyLogo, &companyEnabled,
		&categoryID, &categoryName, &categoryIcon,
	); err != nil {
		return nil, err
	}

	endDateStr := ""
	if endDate.Valid {
		endDateStr = endDate.Time.UTC().Format(time.RFC3339)
	}

	return &donation_dto.DonationListResponse{
		ID:           id,
		DonationCode: donationCode,
		Service: donation_dto.Service{
			ID:          serviceID.String,
			ServiceName: serviceName.String,
			ServiceKey:  serviceKey.String,
		},
		Company: donation_dto.Company{
			ID:          companyID.String,
			CompanyName: companyName.String,
			CompanyLogo: companyLogo.String,
			Enabled:     companyEnabled.Valid && companyEnabled.Int64 == 1,
		},
		Category: donation_dto.Category{
			ID:           categoryID.String,
			CategoryName: categoryName.String,
			Icon:         categoryIcon.String,
		},
		Title:               title,
		IsFeatured:          isFeatured == 1,
		Target:              target.String,
		CurrentAmount:       currentAmount.String,
		DonationDescription: description.String,
		DonationImages:      []types.DonationImage{},
		CoverImage:          coverImage.String,
		EndDate:             endDateStr,
		StartDate:           startDate.UTC().Format(time.RFC3339),
		IsDeleted:           isDeleted == 1,
		CreatedAt:           createdAt.UTC().Format(time.RFC3339),
		LastModifiedAt:      lastModifiedAt.UTC().Format(time.RFC3339),
		Enabled:             enabled == 1,
	}, nil
}
