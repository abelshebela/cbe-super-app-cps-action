package dto

// import (
// 	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/service/dto"
// )

type UpdateServiceFeeRequest struct {
	ServiceType       string    `json:"serviceType"`
	PaymentType       string    `json:"paymentType"`
	AboveAmount       float64   `json:"aboveAmount"`
	AboveServiceFee   float64   `json:"aboveServiceFee"`
	MinAmountVIRTUAL  float64   `json:"minAmountVIRTUAL"`
	DailyCapLevelOne  float64   `json:"dailyCapLevelOne"`
	SingleCapLevelOne float64   `json:"singleCapLevelOne"`
	CBglEntry         *GLEntry  `json:"CBglEntry"`
	IFBglEntry        *GLEntry  `json:"IFBglEntry"`
	Tiers             []TierDTO `json:"tiers" validate:"required,min=1,dive"`
}

type TierDTO struct {
	MinAmount  float64 `json:"minAmount"`
	MaxAmount  float64 `json:"maxAmount"`
	ServiceFee float64 `json:"serviceFee"`
	DailyCap   float64 `json:"dailyCap"`
	SingleCap  float64 `json:"singleCap"`
}

type GLEntry struct {
	ProductAccount    string `json:"productAccount"`
	ProductBranchcode string `json:"productBranchcode"`
	ServiceAccount    string `json:"serviceAccount"`
	ServiceBranchcode string `json:"serviceBranchcode"`
	VatAccount        string `json:"vatAccount"`
	VatBranchcode     string `json:"vatBranchcode"`
}
