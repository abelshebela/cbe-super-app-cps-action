package entities

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/pkgs/entities/type_definition"
)

type Service struct {
	ID                 *string                      `json:"id" bson:"id"`
	ServiceCode        string                       `json:"service_code" bson:"service_code"`
	ServiceName        string                       `json:"service_name" bson:"service_name"`
	ServiceType        string                       `json:"service_type" bson:"service_type"`
	Key                string                       `json:"key" bson:"key"`
	Cap                type_definition.Cap          `json:"cap" bson:"cap"`
	CBEProductCodes    type_definition.ProductCodes `json:"cbe_product_codes" bson:"cbe_product_codes"`
	CBEIFBProductCodes type_definition.ProductCodes `json:"cbe_ifb_product_codes" bson:"cbe_ifb_product_codes"`
	AboveAmount        uint64                       `json:"above_amount" bson:"above_amount"`
	AboveServiceFee    uint64                       `json:"above_service_fee" bson:"above_service_fee"`
	PaymentType        string                       `json:"payment_type" bson:"payment_type"`
	Tiers              []type_definition.Tier       `json:"tiers" bson:"tiers"`
	CBEGLEntry         type_definition.GLEntry      `json:"cbe_gl_entry" bson:"cbe_gl_entry"`
	CBEIFBGLEntry      type_definition.GLEntry      `json:"cbe_ifb_gl_entry" bson:"cbe_ifb_gl_entry"`
	Enabled            bool                         `json:"enabled" bson:"enabled"`
	IsDeleted          bool                         `json:"is_deleted" bson:"is_deleted"`
	CreatedAt          time.Time                    `json:"created_at" bson:"created_at"`
	LastModifiedAt     time.Time                    `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt          time.Time                    `json:"deleted_at" bson:"deleted_at"`
}
