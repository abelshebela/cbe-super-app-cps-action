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
	Key                string             `bson:"key"`
	Cap                types.Cap          `bson:"cap"`
	CBEProductCodes    types.ProductCodes `bson:"cbe_product_codes"`
	CBEIFBProductCodes types.ProductCodes `bson:"cbe_ifb_product_codes"`
	AboveAmount        uint64             `bson:"above_amount"`
	AboveServiceFee    uint64             `bson:"above_service_fee"`
	PaymentType        string             `bson:"payment_type"`
	Tiers              []types.Tier       `bson:"tiers"`
	CBEGLEntry         types.GLEntry      `bson:"cbe_gl_entry"`
	CBEIFBGLEntry      types.GLEntry      `bson:"cbe_ifb_gl_entry"`
	Enabled            bool               `bson:"enabled"`
	IsDeleted          bool               `bson:"is_deleted"`
	CreatedAt          time.Time          `bson:"created_at"`
	LastModifiedAt     time.Time          `bson:"last_modified_at"`
	DeletedAt          time.Time          `bson:"deleted_at"`
}
