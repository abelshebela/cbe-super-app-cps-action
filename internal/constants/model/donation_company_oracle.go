package model

import "time"

type DonationCompanyOracle struct {
	ID                 string    `json:"id,omitempty"`
	CompanyName        string    `json:"company_name"`
	CompanyCode        string    `json:"company_code"`
	CompanyLogo        string    `json:"company_logo"`
	CompanyDescription string    `json:"company_description"`
	Address            string    `json:"address"`
	PhoneNumber        string    `json:"phone_number"`
	Email              string    `json:"email"`
	Enabled            bool      `json:"enabled"`
	IsDeleted          bool      `json:"is_deleted"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
	LastModifiedAt     time.Time `json:"last_modified_at,omitempty"`
}
