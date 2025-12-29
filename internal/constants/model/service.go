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
	FeeAmount uint8   `bson:"fee_amount" json:"fee_amount"`
	Min       float64 `bson:"min" json:"min"`
	Max       float64 `bson:"max" json:"max"`
}

type Services struct {
	ID                  bson.ObjectID `bson:"_id" json:"_id"`
	ServiceCode         string        `bson:"service_code" json:"service_code"`
	ServiceName         string        `bson:"service_name" json:"service_name"`
	ChargeCode          string        `bson:"charge_code" json:"charge_code"`
	CommissionCode      string        `bson:"commission_code" json:"commission_code"`
	Cap                 Cap           `bson:"cap" json:"cap"`
	AboveAmount         float64       `bson:"above_amount" json:"above_amount"`
	AboveServiceFee     float64       `bson:"above_service_fee" json:"above_service_fee"`
	PaymentType         string        `bson:"payment_type" json:"payment_type"`
	Tiers               []Tier        `bson:"tiers" json:"tiers"`
	CbeGLProductAccount string        `bson:"cbe_gl_product_account" json:"cbe_gl_product_account"`
	Enabled             bool          `bson:"enabled" json:"enabled"`
	IsDeleted           bool          `bson:"is_deleted" json:"is_deleted"`
	CreatedAt           time.Time     `bson:"created_at" json:"created_at"`
	LastModifiedAt      time.Time     `bson:"last_modified_at" json:"last_modified_at"`
	DeletedAt           *time.Time    `bson:"deleted_at" json:"deleted_at"`
}

type Cap struct {
	SingleCap          float64 `bson:"single_cap" json:"single_cap"`
	MinimumTransferCap float64 `bson:"minimum_transfer_cap" json:"minimum_transfer_cap"`
}
