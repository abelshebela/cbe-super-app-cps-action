package bank_dto

import "time"

type BankResponse struct {
	ID             string    `json:"_id" bson:"_id"`
	Name           string    `json:"name" bson:"name"`
	Logo           string    `json:"logo" bson:"logo"`
	Code           string    `json:"code" bson:"code"`
	BIC            string    `json:"bic"  bson:"bic"`
	Enabled        bool      `json:"enabled" bson:"enabled"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time `json:"last_modified_at" bson:"last_modified_at"`
}
