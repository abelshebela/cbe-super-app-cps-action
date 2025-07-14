package dto

import (
	"fmt"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/service/dto"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UpdateServiceFeeRequest struct {
	ServiceType       string     `json:"serviceType"`
	PaymentType       string     `json:"paymentType"`
	AboveAmount       float64    `json:"aboveAmount"`
	AboveServiceFee   float64    `json:"aboveServiceFee"`
	MinAmountVIRTUAL  float64    `json:"minAmountVIRTUAL"`
	DailyCapLevelOne  float64    `json:"dailyCapLevelOne"`
	SingleCapLevelOne float64    `json:"singleCapLevelOne"`
	CBglEntry         *GLEntry   `json:"CBglEntry"`
	IFBglEntry        *GLEntry   `json:"IFBglEntry"`
	Tiers             []dto.Tier `json:"tiers" validate:"required,min=1,dive"`
}

type GLEntry struct {
	ProductAccount    string `json:"productAccount"`
	ProductBranchcode string `json:"productBranchcode"`
	ServiceAccount    string `json:"serviceAccount"`
	ServiceBranchcode string `json:"serviceBranchcode"`
	VatAccount        string `json:"vatAccount"`
	VatBranchcode     string `json:"vatBranchcode"`
}

// ... existing code ...
type CreateServiceRequest struct {
	ServiceCode        string       `json:"service_code"`
	ServiceName        string       `json:"service_name"`
	ServiceType        string       `json:"service_type"`
	Key                string       `json:"key"`
	Cap                Cap          `json:"cap"`
	CBEProductCodes    ProductCodes `json:"cbe_product_codes"`
	CBEIFBProductCodes ProductCodes `json:"cbe_ifb_product_codes"`
	AboveAmount        uint64       `json:"above_amount"`
	AboveServiceFee    uint64       `json:"above_service_fee"`
	PaymentType        string       `json:"payment_type"`
	Tiers              []Tier       `json:"tiers"`
	CBEGLEntry         GLEntry      `json:"cbe_gl_entry"`
	CBEIFBGLEntry      GLEntry      `json:"cbe_ifb_gl_entry"`
}

type RejectActionRequest struct {
	RejectReason string `json:"reject_json"`
}

func (r *RejectActionRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RejectReason, validation.Required),
	)
}
func (v CreateServiceRequest) Validate() error {

	fmt.Println("filter value", v)
	return validation.ValidateStruct(&v,
		validation.Field(&v.ServiceCode, validation.Required),
		validation.Field(&v.ServiceName, validation.Required),
		validation.Field(&v.ServiceType, validation.Required),
		validation.Field(&v.Key, validation.Required),
		validation.Field(&v.Cap, validation.Required),
		validation.Field(&v.CBEProductCodes, validation.Required),
		validation.Field(&v.CBEIFBProductCodes, validation.Required),
		validation.Field(&v.AboveAmount, validation.Required),
		validation.Field(&v.AboveServiceFee, validation.Required),
		validation.Field(&v.PaymentType, validation.Required),
		validation.Field(&v.Tiers, validation.Required, validation.Length(1, 0)),
	)
}

type Tier struct {
	ID        string
	Min       uint64
	Max       uint64
	FeeAmount uint64
}

type Cap struct {
	KYCLevel  KYCLevel
	SingleCap uint64
	DailyCap  uint64
	MinAmount uint64
	MaxAmount uint64
}

type KYCLevel string

const (
	KYCLevelZero KYCLevel = "ZERO"
	KYCLevelOne  KYCLevel = "ONE"
	KYCLevelTwo  KYCLevel = "TWO"
)

type ProductCodes struct {
	PRD    string
	VATPRD string
	SFPRD  string
	TRXN   string
}

// ... existing code ...
