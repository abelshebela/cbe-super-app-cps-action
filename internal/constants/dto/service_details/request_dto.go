package dto

type UpdateServiceFeeRequest struct {
	ID        string `json:"tire_id" bson:"tire_id"`
	MinAmount uint64 `json:"min_amount" bson:"min_amount"`
	MaxAmount uint64 `json:"max_amount" bson:"max_amount"`
}

type PaymentType string

const (
	PaymentTypeFlatFee    PaymentType = "flat_fee"
	PaymentTypePercentage PaymentType = "percentage"
)

// TierDTO represents a fee tier structure
type TierDTO struct {
	Min       uint64 `json:"min" bson:"min"`
	Max       uint64 `json:"max" bson:"max"`
	FeeAmount uint64 `json:"fee_amount" bson:"fee_amount"`
}

// CBglEntryDTO represents CBE GL entry details
type CBglEntryDTO struct {
	ProductAccount    string `json:"cbgl_product_account" bson:"cbgl_product_account"`
	ProductBranchCode string `json:"cbgl_product_branchcode" bson:"cbgl_product_branch_code"`
	ServiceAccount    string `json:"cbgl_service_account" bson:"cbgl_service_account"`
	ServiceBranchCode string `json:"cbgl_service_branchcode" bson:"cbgl_service_branch_code"`
	VatAccount        string `json:"cbgl_vat_account" bson:"cbgl_vat_account"`
	VatBranchCode     string `json:"cbgl_vat_branch_code" bson:"cbgl_vat_branch_code"`
}

// IFBglEntryDTO represents IFB GL entry details
type IFBglEntryDTO struct {
	ProductAccount    string `json:"ifbgl_product_account" bson:"ifbgl_product_account"`
	ProductBranchCode string `json:"ifbgl_product_branchcode" bson:"ifbgl_product_branch_code"`
	ServiceAccount    string `json:"ifbgl_service_account" bson:"ifbgl_service_account"`
	ServiceBranchCode string `json:"ifbgl_service_branchcode" bson:"ifbgl_service_branch_code"`
	VatAccount        string `json:"ifbgl_vat_account" bson:"ifbgl_vat_account"`
	VatBranchCode     string `json:"ifbgl_vat_branch_code" bson:"ifbgl_vat_branch_code"`
}

// ServiceFeeDetailDTO represents the complete service fee structure
type ServiceFeeDetailDTO struct {
	ServiceType       string        `json:"service_type" bson:"service_type"`
	PaymentType       PaymentType   `json:"payment_type" bson:"payment_type"`
	SingleCapLevelOne uint64        `json:"single_cap_level_one" bson:"single_cap_level_one"`
	DailyCapLevelOne  uint64        `json:"daily_cap_level_one" bson:"daily_cap_level_one"`
	MinAmountVIRTUAL  uint64        `json:"min_amount_virtual" bson:"min_amount_virtual"`
	AboveAmount       uint64        `json:"above_amount" bson:"above_amount"`
	AboveServiceFee   uint64        `json:"above_service_fee" bson:"above_service_fee"`
	Tiers             []TierDTO     `json:"tiers" bson:"tiers"`
	CBglEntry         CBglEntryDTO  `json:"cbgl_entry" bson:"cbgl_entry"`
	IFBglEntry        IFBglEntryDTO `json:"ifbgl_entry" bson:"ifbgl_entry"`
}

type SingleMaxTransferRequest struct {
	ISingleCap uint64 `json:"individual_single_cap" bson:"individual_single_cap"`
	IDailyCap  uint64 `json:"individual_daily_cap" bson:"individual_daily_cap"`
	CSingleCap uint64 `json:"corporate_single_cap" bson:"corporate_single_cap"`
	CDailyCap  uint64 `json:"corporate_daily_cap" bson:"corporate_daily_cap"`
}

// TotalMaxTransferUpdateRequest represents a request to update the total transfer cap
type TotalMaxTransferUpdateRequest struct {
	TotalTransferLimit uint64 `json:"total_cap" bson:"total_cap"`
}

// MinimumTransferUpdateRequest represents a request to update the minimum transfer amount
type MinimumTransferUpdateRequest struct {
	Minimum uint64 `json:"min_amount" bson:"min_amount"`
}

// DeleteServiceFeeTireRequest represents a request to delete a service fee tier
type DeleteServiceFeeTireRequest struct {
	ID string `json:"tire_id" bson:"tire_id"`
}
