package entity

import "time"

type Wallet struct {
	ID             string    `json:"id" bson:"id"`
	Name           string    `json:"name" bson:"name"`
	Code           string    `json:"code" bson:"code"`
	Avatar         string    `json:"avatar" bson:"avatar"`
	Enabled        bool      `json:"enabled" bson:"enabled"`
	IsDeleted      bool      `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt      time.Time `json:"deletd_at" bson:"deleted_at"`
}

type WalletResponse struct {
	Page    int       `json:"page"`
	Wallets []*Wallet `json:"wallets"`
	Limit   int       `json:"limit"`
	Total   int64     `json:"total"`
}
