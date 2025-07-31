package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service struct {
	ID                 bson.ObjectID `json:"_id" bson:"_id"`
	ServiceCode        string        `json:"service_code" bson:"service_code"`
	ServiceName        string        `json:"service_name" bson:"service_name"`
	ServiceType        string        `json:"service_type" bson:"service_type"`
	Key                string        `json:"key" bson:"key"`
	Cap                Cap           `json:"cap" bson:"cap"`
	CBEProductCodes    ProductCodes  `json:"cbe_product_codes" bson:"cbe_product_codes"`
	CBEIFBProductCodes ProductCodes  `json:"cbe_ifb_product_codes" bson:"cbe_ifb_product_codes"`
	AboveAmount        uint64        `json:"above_amount" bson:"above_amount"`
	AboveServiceFee    uint64        `json:"above_service_fee" bson:"above_service_fee"`
	PaymentType        string        `json:"payment_type" bson:"payment_type"`
	Tiers              []Tier        `json:"tiers" bson:"tiers"`
	CBEGLEntry         GLEntry       `json:"cbe_gl_entry" bson:"cbe_gl_entry"`
	CBEIFBGLEntry      GLEntry       `json:"cbe_ifb_gl_entry" bson:"cbe_ifb_gl_entry"`
	Enabled            bool          `json:"enabled" bson:"enabled"`
	IsDeleted          bool          `json:"is_deleted" bson:"is_deleted"`
	CreatedAt          time.Time     `json:"created_at" bson:"created_at"`
	LastModifiedAt     time.Time     `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt          time.Time     `json:"deleted_at" bson:"deleted_at"`
}

type ProductCodes struct {
	PRD    string `bson:"prd"`
	VATPRD string `bson:"vatprd"`
	SFPRD  string `bson:"sfprd"`
	TRXN   string `bson:"trxn"`
}

type GLEntry struct {
	ProductAccount    string `bson:"product_account"`
	ProductBranchCode string `bson:"product_branch_code"`
	ServiceAccount    string `bson:"service_account"`
	ServiceBranchCode string `bson:"service_branch_code"`
	VatAccount        string `bson:"vat_account"`
	VatBranchCode     string `bson:"vat_branch_code"`
}

type KYCLevel string

const (
	KYCLevelZero KYCLevel = "ZERO"
	KYCLevelOne  KYCLevel = "ONE"
	KYCLevelTwo  KYCLevel = "TWO"
)
