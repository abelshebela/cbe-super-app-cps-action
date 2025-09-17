package dto

import "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"

// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/service/dto"

type UpdateServiceFeeRequest struct {
	ServiceType       string        `json:"service_type" bson:"service_type"`
	PaymentType       string        `json:"payment_type" bson:"payment_type"`
	AboveAmount       float64       `json:"above_amount" bson:"above_amount"`
	AboveServiceFee   float64       `json:"above_service_fee" bson:"above_service_fee"`
	MinAmountVIRTUAL  float64       `json:"min_amount_virtual" bson:"min_amount_virtual"`
	DailyCapLevelOne  float64       `json:"daily_cap_level_one" bson:"daily_cap_level_one"`
	SingleCapLevelOne float64       `json:"single_cap_level_one" bson:"single_cap_level_one"`
	CBglEntry         *GLEntry      `json:"cbgl_entry" bson:"cbgl_entry"`
	IFBglEntry        *GLEntry      `json:"ifbgl_entry" bson:"ifbgl_entry"`
	Tiers             []action.Tier `json:"tiers" bson:"tiers" validate:"required,min=1,dive"`
}

type GLEntry struct {
	ProductAccount    string `json:"product_account" bson:"product_account"`
	ProductBranchcode string `json:"product_branchcode" bson:"product_branchcode"`
	ServiceAccount    string `json:"service_account" bson:"service_account"`
	ServiceBranchcode string `json:"service_branchcode" bson:"service_branchcode"`
	VatAccount        string `json:"vat_account" bson:"vat_account"`
	VatBranchcode     string `json:"vat_branchcode" bson:"vat_branchcode"`
}
