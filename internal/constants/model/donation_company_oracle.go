package model

import "time"

type DonationCompanyOracle struct {
	ID                 string    `json:"id,omitempty" bson:"id"`
	CompanyName        string    `json:"company_name" bson:"company_name"`
	CompanyCode        string    `json:"company_code" bson:"company_code"`
	CompanyLogo        string    `json:"company_logo" bson:"company_logo"`
	CompanyDescription string    `json:"company_description" bson:"company_description"`
	Address            string    `json:"address" bson:"address"`
	PhoneNumber        string    `json:"phone_number" bson:"phone_number"`
	Email              string    `json:"email" bson:"email"`
	Enabled            bool      `json:"enabled" bson:"enabled"`
	IsDeleted          bool      `json:"is_deleted" bson:"is_deleted"`
	CreatedAt          time.Time `json:"created_at,omitempty" bson:"created_at"`
	LastModifiedAt     time.Time `json:"last_modified_at,omitempty" bson:"last_modified_at"`
}
