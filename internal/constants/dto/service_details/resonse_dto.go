package dto

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// MinimumTransferCapResponse represents the response for minimum transfer cap projection
type MinimumTransferCapResponse struct {
	ID            bson.ObjectID `json:"id" bson:"_id"`
	ServiceName   string        `json:"service_name" bson:"service_name"`
	ServiceCode   string        `json:"service_code" bson:"service_code"`
	ServiceType   string        `json:"service_type" bson:"service_type"`
	MinimumAmount uint64        `json:"min_amount" bson:"min_amount"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
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
	ID                bson.ObjectID      `json:"id" bson:"_id"`
	ServiceName       string             `json:"service_name" bson:"service_name"`
	ServiceKey        string             `json:"service_key" bson:"service_key"`
	ServiceType       string             `json:"service_type" bson:"service_type"`
	PaymentType       string             `json:"payment_type" bson:"payment_type"`
	DailyCapLevelOne  uint64             `json:"daily_cap_level_one" bson:"daily_cap_level_one"`
	MinAmountVirtual  uint64             `json:"min_amount_virtual" bson:"min_amount_virtual"`
	SingleCapLevelOne uint64             `json:"single_cap_level_one" bson:"single_cap_level_one"`
	Tiers             []types.Tier       `json:"tiers" bson:"tiers"`
	MinAmount         uint64             `json:"min_amount" bson:"min_amount"`
	ProductCodes      types.ProductCodes `json:"product_codes" bson:"product_codes"`
	IFBProductCodes   types.ProductCodes `json:"ifb_product_codes" bson:"ifb_product_codes"`
	GLEntry           types.GLEntry      `json:"gl_entry" bson:"gl_entry"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
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
	CBEIFBGLEntry      types.IFBglEntry   `json:"cbe_ifb_gl_entry" bson:"cbe_ifb_gl_entry"`
	Enabled            bool               `json:"enabled" bson:"enabled"`
	DailyCapLevelOne   uint64             `json:"daily_cap_level_one" bson:"daily_cap_level_one"`
	MinAmountVirtual   uint64             `json:"min_amount_virtual" bson:"min_amount_virtual"`
	SingleCapLevelOne  uint64             `json:"single_cap_level_one" bson:"single_cap_level_one"`
	IsDeleted          bool               `json:"is_deleted" bson:"is_deleted"`
	CreatedAt          time.Time          `json:"created_at" bson:"created_at"`
	LastModifiedAt     time.Time          `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt          time.Time          `json:"deleted_at" bson:"deleted_at"`
}
