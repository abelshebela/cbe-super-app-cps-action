package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Bank struct {
	ID             string    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name           string    `json:"name" bson:"name"`
	Logo           string    `json:"logo" bson:"logo"`
	Code           string    `json:"code" bson:"code"`
	BIC            string    `json:"bic" bson:"bic"`
	Enabled        bool      `json:"enabled" bson:"enabled"`
	IsDeleted      bool      `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time `json:"created_at,omitzero" bson:"created_at"`
	LastModifiedAt time.Time `json:"last_modified_at,omitzero" bson:"last_modified_at"`
}

type UpdateLogo struct {
	ID   string `json:"_id" bson:"_id"`
	Logo string `json:"logo" bson:"logo"`
}

type BankResponse struct {
	Page  int     `json:"page"`
	Banks []*Bank `json:"bank"`
	Limit int     `json:"limit"`
	Total int64   `json:"total"`
}

type BankDocument struct {
	ID             bson.ObjectID `bson:"_id,omitempty"`
	Name           string        `bson:"name"`
	Logo           string        `bson:"logo"`
	Code           string        `bson:"code"`
	BIC            string        `bson:"bic"`
	Enabled        bool          `bson:"enabled"`
	IsDeleted      bool          `bson:"is_deleted"`
	CreatedAt      time.Time     `bson:"created_at"`
	LastModifiedAt time.Time     `bson:"last_modified_at"`
}
