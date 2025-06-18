package entity

import "time"

type Bank struct {
	ID             string    `json:"id,omitempty" bson:"id"`
	Name           string    `json:"name" bson:"name"`
	Logo           string    `json:"logo" bson:"logo"`
	Code           string    `json:"code" bson:"code"`
	BIC            string    `json:"bic" bson:"bic"`
	Enabled        bool      `json:"enabled" bson:"enabled"`
	IsDeleted      bool      `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time `json:"created_at,omitzero" bson:"created_at"`
	LastModifiedAt time.Time `json:"last_modified_at,omitzero" bson:"last_modified_at"`
}

type BankResponse struct {
	Page  int     `json:"page"`
	Banks []*Bank `json:"bank"`
	Limit int     `json:"limit"`
	Total int64   `json:"total"`
}
