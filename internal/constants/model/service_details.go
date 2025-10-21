package model

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ServiceDetails struct {
	ID                 bson.ObjectID      `json:"id,omitempty" bson:"_id,omitempty"`
	ServiceCode        string             `json:"service_code" bson:"service_code"`
	ServiceName        string             `json:"service_name" bson:"service_name"`
	ServiceType        string             `json:"service_type" bson:"service_type"`
	Key                string             `json:"key" bson:"key"`
	Cap                types.Cap          `json:"cap" bson:"cap"`
	CBEProductCodes    types.ProductCodes `json:"cbe_product_codes" bson:"cbe_product_codes"`
	CBEIFBProductCodes types.ProductCodes `json:"cbe_ifb_product_codes" bson:"cbe_ifb_product_codes"`
	AboveAmount        uint64             `json:"above_amount" bson:"above_amount"`
	AboveServiceFee    uint64             `json:"above_service_fee" bson:"above_service_fee"`
	PaymentType        string             `json:"payment_type" bson:"payment_type"`
	Tiers              []types.Tier       `json:"tiers" bson:"tiers"`
	CBEGLEntry         types.GLEntry      `json:"cbe_gl_entry" bson:"cbe_gl_entry"`
	CBEIFBGLEntry      types.GLEntry      `json:"cbe_ifb_gl_entry" bson:"cbe_ifb_gl_entry"`
	Enabled            bool               `json:"enabled" bson:"enabled"`
	IsDeleted          bool               `json:"is_deleted" bson:"is_deleted"`
	CreatedAt          time.Time          `json:"created_at" bson:"created_at"`
	LastModifiedAt     time.Time          `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt          time.Time          `json:"deleted_at" bson:"deleted_at"`
}

type Cap struct {
	ISingleCap         uint64 `json:"individual_single_cap" bson:"individual_single_cap"`
	IDailyCap          uint64 `json:"individual_daily_cap" bson:"individual_daily_cap"`
	CorporateSingleCap uint64 `json:"corporate_single_cap" bson:"corporate_single_cap"`
	CorporateDailyCap  uint64 `json:"corporate_daily_cap" bson:"corporate_daily_cap"`
	MinAmount          uint64 `json:"min_amount" bson:"min_amount"`
}
