package model

import (
	"time"
)

// WalletOracle maps to WALLETS: UNIQUE_CODE → UniqueCode; SERVICES_SELF/OTHER/AGENT → Self/Other/Agent.
type WalletOracle struct {
	ID             string     `json:"id" bson:"_id"`
	Name           string     `json:"name" bson:"name"`
	UniqueCode     string     `json:"unique_code" bson:"unique_code"`
	Avatar         string     `json:"avatar" bson:"avatar"`
	Enabled        bool       `json:"enabled" bson:"enabled"`
	Self           bool       `json:"self" bson:"self"`
	Other          bool       `json:"other" bson:"other"`
	Agent          bool       `json:"agent" bson:"agent"`
	SelfServiceID  string     `json:"self_service_id" bson:"self_service_id"`
	OtherServiceID string     `json:"other_service_id" bson:"other_service_id"`
	AgentServiceID string     `json:"agent_service_id" bson:"agent_service_id"`
	IsDeleted      bool       `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time  `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time  `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt      *time.Time `json:"deleted_at" bson:"deleted_at"`
}
