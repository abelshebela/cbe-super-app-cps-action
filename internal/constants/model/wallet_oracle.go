package model

import (
	"time"
)

// WalletOracle maps to WALLETS: wallet_code → UniqueCode; services_self/other/agent → Self/Other/Agent.
type WalletOracle struct {
	ID             string     `json:"id" bson:"_id"`
	Name           string     `json:"name" bson:"name"`
	UniqueCode     string     `json:"unique_code" bson:"unique_code"`
	ServiceCode    string     `json:"service_code" bson:"service_code"`
	ServiceKey     string     `json:"service_key" bson:"service_key"`
	ServiceID      string     `json:"service_id" bson:"service_id"`
	Avatar         string     `json:"avatar" bson:"avatar"`
	Enabled        bool       `json:"enabled" bson:"enabled"`
	Self           bool       `json:"self" bson:"self"`
	Other          bool       `json:"other" bson:"other"`
	Agent          bool       `json:"agent" bson:"agent"`
	IsDeleted      bool       `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time  `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time  `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt      *time.Time `json:"deleted_at" bson:"deleted_at"`
}
