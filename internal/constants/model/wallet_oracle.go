package model

import (
	"time"
)

type WalletOracle struct {
	ID                  string     `bson:"id" json:"id"`
	Name                string     `bson:"name" json:"name"`
	UniqueCode          string     `bson:"unique_code" json:"unique_code"`
	Avatar              string     `bson:"avatar" json:"avatar"`
	Enabled             bool       `bson:"enabled" json:"enabled"`
	Self                bool       `bson:"self" json:"self"`
	Other               bool       `bson:"other" json:"other"`
	Agent               bool       `bson:"agent" json:"agent"`
	SelfServiceID       string     `bson:"self_service_id" json:"self_service_id"`
	OtherServiceID      string     `bson:"other_service_id" json:"other_service_id"`
	AgentServiceID      string     `bson:"agent_service_id" json:"agent_service_id"`
	SelfServiceEnabled  int        `bson:"self_service_enabled" json:"self_service_enabled"`
	OtherServiceEnabled int        `bson:"other_service_enabled" json:"other_service_enabled"`
	AgentServiceEnabled int        `bson:"agent_service_enabled" json:"agent_service_enabled"`
	SelfServiceCode     string     `bson:"self_service_code" json:"self_service_code"`
	OtherServiceCode    string     `bson:"other_service_code" json:"other_service_code"`
	AgentServiceCode    string     `bson:"agent_service_code" json:"agent_service_code"`
	IsDeleted           bool       `bson:"is_deleted" json:"is_deleted"`
	CreatedAt           time.Time  `bson:"created_at" json:"created_at"`
	LastModifiedAt      time.Time  `bson:"last_modified_at" json:"last_modified_at"`
	DeletedAt           *time.Time `bson:"deleted_at" json:"deleted_at"`
}

type WalletService struct {
	ID             string     `bson:"id" json:"id"`
	ServiceID      string     `bson:"service_id" json:"service_id"`
	WalletID       string     `bson:"wallet_id" json:"wallet_id"`
	ServiceType    string     `bson:"service_type" json:"service_type"`
	IsEnabled      int        `bson:"is_enabled" json:"is_enabled"`
	IsDeleted      int        `bson:"is_deleted" json:"is_deleted"`
	CreatedAt      time.Time  `bson:"created_at" json:"created_at"`
	LastModifiedAt time.Time  `bson:"last_modified_at" json:"last_modified_at"`
	DeletedAt      *time.Time `bson:"deleted_at" json:"deleted_at"`
}
