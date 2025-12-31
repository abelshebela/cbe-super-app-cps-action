package model

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Wallet struct {
	ID             bson.ObjectID  `json:"id" bson:"_id"`
	Name           string         `json:"name" bson:"name"`
	UniqueCode     string         `json:"unique_code" bson:"unique_code"`
	ServiceID      string         `json:"service_id" bson:"service_id"`
	Avatar         string         `json:"avatar" bson:"avatar"`
	Enabled        bool           `json:"enabled" bson:"enabled"`
	Services       types.Services `json:"services" bson:"services"`
	IsDeleted      bool           `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time      `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time      `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt      time.Time      `json:"deleted_at" bson:"deleted_at"`
}
type GRPCWallet struct {
	ID             bson.ObjectID  `json:"id" bson:"_id"`
	Name           string         `json:"name" bson:"name"`
	UniqueCode     string         `json:"unique_code" bson:"unique_code"`
	ServiceID      string         `json:"service_id" bson:"service_id"`
	ServiceCode    string         `json:"service_code" bson:"service_code"`
	ServiceKey     string         `json:"service_key" bson:"service_key"`
	Avatar         string         `json:"avatar" bson:"avatar"`
	Enabled        bool           `json:"enabled" bson:"enabled"`
	Services       types.Services `json:"services" bson:"services"`
	IsDeleted      bool           `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time      `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time      `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt      time.Time      `json:"deleted_at" bson:"deleted_at"`
}
