package updatedbulkservice

import "time"

type Tier struct {
	Min       float64 `json:"min" bson:"min"`
	Max       float64 `json:"max" bson:"max"`
	FeeAmount float64 `json:"fee_amount" bson:"fee_amount"`
}

type ProductCodes struct {
	PRD    string `json:"prd" bson:"prd"`
	VATPRD string `json:"vatprd" bson:"vatprd"`
	SFPRD  string `json:"sfprd" bson:"sfprd"`
	TRXN   string `json:"trxn" bson:"trxn"`
}

type CBglEntry struct {
	CBglProductAccount    string `json:"cbgl_product_account" bson:"cbgl_product_account"`
	CBglProductBranchcode string `json:"cbgl_product_branchcode" bson:"cbgl_product_branchcode"`
	CBglServiceAccount    string `json:"cbgl_service_account" bson:"cbgl_service_account"`
	CBglServiceBranchcode string `json:"cbgl_service_branchcode" bson:"cbgl_service_branchcode"`
	CBglVatAccount        string `json:"cbgl_vat_account" bson:"cbgl_vat_account"`
	CBglVatBranchcode     string `json:"cbgl_vat_branchcode" bson:"cbgl_vat_branchcode"`
}

type IFBglEntry struct {
	IFBglProductAccount    string `json:"ifbgl_product_account" bson:"ifbgl_product_account"`
	IFBglProductBranchcode string `json:"ifbgl_product_branchcode" bson:"ifbgl_product_branchcode"`
	IFBglServiceAccount    string `json:"ifbgl_service_account" bson:"ifbgl_service_account"`
	IFBglServiceBranchcode string `json:"ifbgl_service_branchcode" bson:"ifbgl_service_branchcode"`
	IFBglVatAccount        string `json:"ifbgl_vat_account" bson:"ifbgl_vat_account"`
	IFBglVatBranchcode     string `json:"ifbgl_vat_branchcode" bson:"ifbgl_vat_branchcode"`
}

type ServiceDetails struct {
	Key               string       `json:"key" bson:"key"`
	ServiceCode       string       `json:"service_code" bson:"service_code"`
	ServiceName       string       `json:"service_name" bson:"service_name"`
	SingleCap         float64      `json:"single_cap" bson:"single_cap"`
	DailyCap          float64      `json:"daily_cap" bson:"daily_cap"`
	MinAmount         float64      `json:"min_amount" bson:"min_amount"`
	ServiceType       string       `json:"service_type" bson:"service_type"`
	PaymentType       string       `json:"payment_type" bson:"payment_type"`
	VAT               float64      `json:"vat" bson:"vat"`
	Tiers             []Tier       `json:"tiers" bson:"tiers"`
	MinAmountVIRTUAL  float64      `json:"min_amount_virtual" bson:"min_amount_virtual"`
	SingleCapLevelOne float64      `json:"single_cap_level_one" bson:"single_cap_level_one"`
	DailyCapLevelOne  float64      `json:"daily_cap_level_one" bson:"daily_cap_level_one"`
	AboveAmount       float64      `json:"above_amount" bson:"above_amount"`
	AboveServiceFee   float64      `json:"above_service_fee" bson:"above_service_fee"`
	Prefix            string       `json:"prefix,omitempty" bson:"prefix,omitempty"`
	USSDEnabled       bool         `json:"ussd_enabled" bson:"ussd_enabled"`
	Enabled           bool         `json:"enabled" bson:"enabled"`
	IsDeleted         bool         `json:"is_deleted" bson:"is_deleted"`
	DateCreated       time.Time    `json:"date_created" bson:"date_created"`
	LastModified      time.Time    `json:"last_modified" bson:"last_modified"`
	CreatedDate       time.Time    `json:"created_date" bson:"created_date"`
	ProductCodes      ProductCodes `json:"product_codes" bson:"product_codes"`
	IFBProductCodes   ProductCodes `json:"ifb_product_codes" bson:"ifb_product_codes"`
	CBglEntry         CBglEntry    `json:"cbgl_entry" bson:"cbgl_entry"`
	IFBglEntry        IFBglEntry   `json:"ifbgl_entry" bson:"ifbgl_entry"`
}
