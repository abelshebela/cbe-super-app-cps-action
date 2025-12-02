package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Bank struct {
	ID             bson.ObjectID `json:"id,omitempty" bson:"_id"`
	Name           string        `json:"name" bson:"name"`
	Logo           string        `json:"logo" bson:"logo"`
	Code           string        `json:"code" bson:"code"`
	BIC            string        `json:"bic" bson:"bic"`
	Enabled        bool          `json:"enabled" bson:"enabled"`
	IsDeleted      bool          `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time     `json:"created_at,omitzero" bson:"created_at"`
	LastModifiedAt time.Time     `json:"last_modified_at,omitzero" bson:"last_modified_at"`
}
