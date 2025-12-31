package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type FeeType string

const (
	Flat    FeeType = "FLAT"
	Percent FeeType = "PERCENT"
)

type Tier struct {
	FeeType   FeeType `bson:"fee_type" json:"fee_type"`
	FeeAmount string  `bson:"fee_amount" json:"fee_amount"`
	Min       string  `bson:"min" json:"min"`
	Max       string  `bson:"max" json:"max"`
}

type ServiceLists struct {
	ServiceName             string `bson:"service_name" json:"service_name"`
	ServiceKey              string `bson:"service_key" json:"service_key"`
	OverideCap              Cap    `bson:"overide_cap" json:"overide_cap"`
	OverideProductGlAccount string `bson:"overide_product_gl_account" json:"overide_product_gl_account"`
	OverideTiers            []Tier `bson:"overide_tiers" json:"overide_tiers"`
	IsEnabled               bool   `bson:"is_enabled" json:"is_enabled"`
}
type Service struct {
	ID               bson.ObjectID  `bson:"_id,omitempty" json:"id"`
	ServiceCode      string         `bson:"service_code" json:"service_code"`
	ServiceKey       string         `bson:"service_key" json:"service_key"`
	ServiceName      string         `bson:"service_name" json:"service_name"`
	ServiceList      []ServiceLists `bson:"service_list" json:"service_list"`
	Cap              Cap            `bson:"cap" json:"cap"`
	Tiers            []Tier         `bson:"tiers" json:"tiers"`
	ProductGlAccount string         `bson:"product_gl_account" json:"product_gl_account"`
	Enabled          bool           `bson:"enabled" json:"enabled"`
	IsDeleted        bool           `bson:"is_deleted" json:"is_deleted"`
	CreatedAt        time.Time      `bson:"created_at" json:"created_at"`
	LastModifiedAt   time.Time      `bson:"last_modified_at" json:"last_modified_at"`
	DeletedAt        *time.Time     `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type Cap struct {
	SingleCap          string `bson:"single_cap" json:"single_cap"`
	MinimumTransferCap string `bson:"minimum_transfer_cap" json:"minimum_transfer_cap"`
}
