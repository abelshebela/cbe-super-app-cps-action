package bank_dto

import "time"

type BankResponse struct {
	ID             string    `json:"_id" bson:"_id"`
	Name           string    `json:"name" bson:"name"`
	Logo           string    `json:"logo" bson:"logo"`
	IsCBE          bool      `json:"is_cbe" bson:"is_cbe"`
	Code           string    `json:"code" bson:"code"`
	BIC            string    `json:"bic"  bson:"bic"`
	Enabled        bool      `json:"enabled" bson:"enabled"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time `json:"last_modified_at" bson:"last_modified_at"`
}
type BankOracleResponse struct {
	ID              string `sqlx:"id" json:"id"`
	BankName        string `sqlx:"bank_name" json:"bank_name"`
	Logo            string `sqlx:"logo" json:"logo"`
	BICCode         string `sqlx:"bic_code" json:"bic_code"`
	IsEnabled       bool   `sqlx:"enabled" json:"enabled"`
	IsCBE           bool   `sqlx:"is_cbe" json:"is_cbe"`
	AccountLength   int    `sqlx:"account_length" json:"account_length"`
	HasAlphaNumeric bool   `sqlx:"has_alpha_numeric" json:"has_alpha_numeric"`
	CreateAt        string `sqlx:"create_at" json:"create_at"`
	UpdateAt        string `sqlx:"update_at" json:"update_at"`
}
