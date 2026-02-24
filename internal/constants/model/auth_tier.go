package model

import (
	"cbe-super-app-cps-action/internal/constants"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuthTier struct {
	ID           bson.ObjectID          `json:"id" bson:"_id"`
	Currency     constants.CurrencyType `json:"currency" bson:"currency"`
	MinAmount    uint64                 `json:"min_amount" bson:"min_amount"`
	MaxAmount    uint64                 `json:"max_amount" bson:"max_amount"`
	Method       constants.Method       `json:"method" bson:"method"`
	Enabled      bool                   `json:"enabled" bson:"enabled"`
	IsDeleted    bool                   `json:"is_deleted" bson:"is_deleted"`
	CreatedAt    time.Time              `json:"created_at" bson:"created_at"`
	LastModified time.Time              `json:"last_modified" bson:"last_modified"`
}
