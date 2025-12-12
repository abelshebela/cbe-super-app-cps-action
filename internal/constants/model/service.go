package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Cap struct {
	KYCLevel           string `bson:"kyc_level" json:"kyc_level"`
	SingleCap          int64  `bson:"single_cap" json:"single_cap"`
	MinimumTransferCap int64  `bson:"minimum_transfer_cap" json:"minimum_transfer_cap"`
}

// ServiceCode = ProductCode
// ChargeCode(service fee account) - flat or percent, amount
// VatCode()

// 100
// 110 - 10 (Service fee) ChargeCode
// 111.5 - 1.5 (vat fee) VatCode

type FeeType string

const (
	Flat    FeeType = "FLAT"
	Percent FeeType = "PERCENT"
)

type Tier struct {
	FeeType   FeeType `bson:"fee_type" json:"fee_type"`
	FeeAmount uint8   `bson:"fee_amount" json:"fee_amount"`
	Min       int64   `bson:"min" json:"min"`
	Max       int64   `bson:"max" json:"max"`
}

type Services struct {
	ID             bson.ObjectID `bson:"_id" json:"_id"`
	ServiceCode    string        `bson:"service_code" json:"service_code"`
	ServiceName    string        `bson:"service_name" json:"service_name"`
	ServiceType    string        `bson:"service_type" json:"service_type"`
	Key            string        `bson:"key" json:"key"`
	ChargeCode     string        `bson:"charge_code" json:"charge_code"`
	CommissionCode string        `bson:"commission_code" json:"commission_code"`
	// SingleCap int64  `bson:"single_cap" json:"single_cap"`
	Cap                  Cap          `bson:"cap" json:"cap"`
	CbeProductCodes      ProductCodes `bson:"cbe_product_codes" json:"cbe_product_codes"`
	CbeIfbProductCodes   ProductCodes `bson:"cbe_ifb_product_codes" json:"cbe_ifb_product_codes"`
	AboveAmount          int64        `bson:"above_amount" json:"above_amount"`
	AboveServiceFee      int64        `bson:"above_service_fee" json:"above_service_fee"`
	PaymentType          string       `bson:"payment_type" json:"payment_type"`
	Tiers                []Tier       `bson:"tiers" json:"tiers"`
	CbeGLProductAccount  string       `bson:"cbe_gl_product_account" json:"cbe_gl_product_account"`
	CbeIFBProductAccount string       `bson:"cbe_ifb_product_account" json:"cbe_ifb_product_account"`
	Enabled              bool         `bson:"enabled" json:"enabled"`
	IsDeleted            bool         `bson:"is_deleted" json:"is_deleted"`
	CreatedAt            time.Time    `bson:"created_at" json:"created_at"`
	LastModifiedAt       time.Time    `bson:"last_modified_at" json:"last_modified_at"`
	DeletedAt            *time.Time   `bson:"deleted_at" json:"deleted_at"`
}
