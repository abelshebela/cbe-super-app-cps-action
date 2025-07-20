// Package entity provides domain entities for wallet management.
package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Wallet struct {
	ID             string    `json:"id,omitempty" bson:"id"`
	Name           string    `json:"name" bson:"name"`
	Code           string    `json:"code" bson:"code"`
	Avatar         string    `json:"avatar" bson:"avatar"`
	Enabled        bool      `json:"enabled" bson:"enabled"`
	IsDeleted      bool      `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt      *time.Time `json:"deletd_at,omitempty" bson:"deleted_at"`
}

type WalletDocument struct {
	ID             bson.ObjectID `json:"id" bson:"_id"`
	Name           string        `json:"name" bson:"name"`
	Code           string        `json:"code" bson:"code"`
	Avatar         string        `json:"avatar" bson:"avatar"`
	Enabled        bool          `json:"enabled" bson:"enabled"`
	IsDeleted      bool          `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time     `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt      time.Time     `json:"deletd_at" bson:"deleted_at"`
}

type WalletResponse struct {
	Page    int       `json:"page"`
	Wallets []*Wallet `json:"wallets"`
	Limit   int       `json:"limit"`
	Total   int64     `json:"total"`
}
