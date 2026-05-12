package donation_company_oracle

import (
	"database/sql"
	"time"

	"cbe-super-app-cps-action/internal/constants/dto/donation_company"
)

const donationCompanySelectCols = `RAWTOHEX(ID), COMPANY_NAME, COMPANY_CODE, COMPANY_LOGO, COMPANY_DESCRIPTION,
	ADDRESS, PHONE_NUMBER, EMAIL,
	IS_DELETED, ENABLED, CREATED_AT, LAST_MODIFIED_AT`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanDonationCompanyRow(s rowScanner) (*donation_company.DonationCompanyListResponse, error) {
	var (
		id                                                                  string
		companyName, companyCode                                            string
		companyLogo, companyDescription, address, phone, email              sql.NullString
		isDeleted, enabled                                                  int
		createdAt, lastModifiedAt                                           time.Time
	)

	if err := s.Scan(
		&id,
		&companyName,
		&companyCode,
		&companyLogo,
		&companyDescription,
		&address,
		&phone,
		&email,
		&isDeleted,
		&enabled,
		&createdAt,
		&lastModifiedAt,
	); err != nil {
		return nil, err
	}

	return &donation_company.DonationCompanyListResponse{
		ID:                 id,
		CompanyName:        companyName,
		CompanyCode:        companyCode,
		CompanyLogo:        companyLogo.String,
		CompanyDescription: companyDescription.String,
		Address:            address.String,
		PhoneNumber:        phone.String,
		Email:              email.String,
		IsDeleted:          isDeleted == 1,
		Enabled:            enabled == 1,
		CreatedAt:          createdAt.UTC().Format(time.RFC3339),
		LastModifiedAt:     lastModifiedAt.UTC().Format(time.RFC3339),
	}, nil
}
