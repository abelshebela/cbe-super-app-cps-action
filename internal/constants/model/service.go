package model

import (
	"cbe-super-app-cps-action/internal/constants"
	"time"
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

type Cap struct {
	Source             constants.SourceApp `json:"source" bson:"source"`
	Currency           string              `json:"currency" bson:"currency"`
	SingleCap          string              `json:"single_cap" bson:"single_cap"`
	MinimumTransferCap string              `json:"minimum_transfer_cap" bson:"minimum_transfer_cap"`
	CreatedAt          time.Time           `json:"created_at" bson:"created_at"`
	LastModifiedAt     time.Time           `json:"last_modified_at" bson:"last_modified_at"`
}

type ServiceKey struct {
	ID             string     `json:"id,omitempty" bson:"_id,omitempty"`
	ServiceName    string     `json:"service_name" bson:"service_name"`
	ServiceKey     string     `json:"service_key" bson:"service_key"`
	IsUSSDEnabled  bool       `json:"is_ussd_enabled" bson:"is_ussd_enabled"`
	IsEnabled      bool       `json:"is_enabled" bson:"is_enabled"`
	IsDeleted      bool       `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time  `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time  `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type Service struct {
	ID                       string    `json:"id" bson:"_id"`
	ServiceKeyId             string    `json:"service_key_id" bson:"service_key_id"`
	ServiceCode              string    `json:"service_code" bson:"service_code"`
	Cap                      []Cap     `json:"cap" bson:"cap"`
	MinimumFraudAmount       string    `json:"minimum_fraud_amount" bson:"minimum_fraud_amount"`
	ProductGlAccount         string    `json:"product_gl_account" bson:"product_gl_account"`
	ProductGlAccountCurrency string    `json:"product_gl_account_currency" bson:"product_gl_account_currency"`
	CreatedAt                time.Time `json:"created_at" bson:"created_at"`
	LastModifiedAt           time.Time `json:"last_modified_at" bson:"last_modified_at"`
}
