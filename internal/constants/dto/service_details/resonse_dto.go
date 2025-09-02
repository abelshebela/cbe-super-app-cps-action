package dto

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// MinimumTransferCapResponse represents the response for minimum transfer cap projection
type MinimumTransferCapResponse struct {
	ID          bson.ObjectID `json:"id" bson:"_id"`
	ServiceName string        `json:"service_name" bson:"service_name"`
	ServiceCode string        `json:"service_code" bson:"service_code"`
	ServiceType string        `json:"service_type" bson:"service_type"`
	Cap         types.Cap     `json:"cap" bson:"cap"`
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
}

// MaximumTransferCapResponse represents the response for maximum transfer cap projection
type MaximumTransferCapResponse struct {
	ID              bson.ObjectID      `json:"id" bson:"_id"`
	ServiceName     string             `json:"service_name" bson:"service_name"`
	ServiceKey      string             `json:"service_key" bson:"service_key"`
	ServiceCode     string             `json:"service_code" bson:"service_code"`
	ServiceType     string             `json:"service_type" bson:"service_type"`
	Cap             types.Cap          `json:"cap" bson:"cap"`
	ProductCodes    types.ProductCodes `json:"product_codes" bson:"product_codes"`
	IFBProductCodes types.ProductCodes `json:"ifb_product_codes" bson:"ifb_product_codes"`
	GLEntry         types.GLEntry      `json:"gl_entry" bson:"gl_entry"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
}

// ServiceFeeResponse represents the response for service fee projection
type ServiceFeeResponse struct {
	ID              bson.ObjectID      `json:"id" bson:"_id"`
	ServiceName     string             `json:"service_name" bson:"service_name"`
	ServiceKey      string             `json:"service_key" bson:"service_key"`
	ServiceType     string             `json:"service_type" bson:"service_type"`
	PaymentType     string             `json:"payment_type" bson:"payment_type"`
	Tiers           []types.Tier       `json:"tiers" bson:"tiers"`
	MinAmount       uint64             `json:"min_amount" bson:"min_amount"`
	ProductCodes    types.ProductCodes `json:"product_codes" bson:"product_codes"`
	IFBProductCodes types.ProductCodes `json:"ifb_product_codes" bson:"ifb_product_codes"`
	GLEntry         types.GLEntry      `json:"gl_entry" bson:"gl_entry"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
}

// TotalTransferCapResponse represents the response for total transfer cap projection
type TotalTransferCapResponse struct {
	ID                bson.ObjectID `json:"id" bson:"_id"`
	TotalCap          uint64        `json:"total_cap" bson:"total_cap"`
	CreatedAt         time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAtTotalCap time.Time     `json:"updated_at_total_cap" bson:"updated_at_total_cap"`
}

// ServiceFeeDetailResponse represents the response for service fee detail projection
type ServiceFeeDetailResponse struct {
	ID                 bson.ObjectID      `json:"id" bson:"_id"`
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
