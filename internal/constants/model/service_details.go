package model

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ServiceDetails struct {
	ID                 bson.ObjectID      `json:"id,omitempty" bson:"_id,omitempty"`
	ServiceCode        string             `bson:"service_code"`
	ServiceName        string             `bson:"service_name"`
	ServiceType        string             `bson:"service_type"`
	Key                string             `bson:"key omitempty"`
	Cap                types.Cap          `bson:"cap omitempty"`
	CBEProductCodes    types.ProductCodes `bson:"cbe_product_codes omitempty"`
	CBEIFBProductCodes types.ProductCodes `bson:"cbe_ifb_product_codes omitempty"`
	AboveAmount        uint64             `bson:"above_amount omitempty"`
	AboveServiceFee    uint64             `bson:"above_service_fee omitempty"`
	PaymentType        string             `bson:"payment_type omitempty"`
	Tiers              []types.Tier       `bson:"tiers omitempty"`
	CBEGLEntry         types.GLEntry      `bson:"cbe_gl_entry omitempty"`
	CBEIFBGLEntry      types.GLEntry      `bson:"cbe_ifb_gl_entry omitempty"`
	Enabled            bool               `bson:"enabled"`
	IsDeleted          bool               `bson:"is_deleted"`
	CreatedAt          time.Time          `bson:"created_at omitempty"`
	LastModifiedAt     time.Time          `bson:"last_modified_at omitempty"`
	DeletedAt          time.Time          `bson:"deleted_at omitempty"`
}

type Cap struct {
	ISingleCap         uint64 `json:"individual_single_cap" bson:"individual_single_cap"`
	IDailyCap          uint64 `json:"individual_daily_cap" bson:"individual_daily_cap"`
	CorporateSingleCap uint64 `json:"corporate_single_cap" bson:"corporate_single_cap"`
	CorporateDailyCap  uint64 `json:"corporate_daily_cap" bson:"corporate_daily_cap"`
	MinAmount          uint64 `json:"min_amount" bson:"min_amount"`
}