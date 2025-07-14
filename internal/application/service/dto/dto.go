package dto

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"
)

type ServiceResponse struct{}
type UpdateServiceRequest struct{}

type Tier struct {
	Min       float64 `json:"min" validate:"required,gte=0"`
	Max       float64 `json:"max" validate:"required,gtfield=Min"`
	FeeAmount float64 `json:"fee_amount" validate:"required,gte=0"`
}

type CreateServiceRequest struct {
	ServiceCode        string               `json:"service_code"`
	ServiceName        string               `json:"service_name"`
	ServiceType        string               `json:"service_type"`
	Key                string               `json:"key"`
	Cap                service.Cap          `json:"cap"`
	CBEProductCodes    service.ProductCodes `json:"cbe_product_codes"`
	CBEIFBProductCodes service.ProductCodes `json:"cbe_ifb_product_codes"`
	AboveAmount        uint64               `json:"above_amount"`
	AboveServiceFee    uint64               `json:"above_service_fee"`
	PaymentType        string               `json:"payment_type"`
	Tiers              []service.Tier       `json:"tiers"`
	CBEGLEntry         service.GLEntry      `json:"cbe_gl_entry"`
	CBEIFBGLEntry      service.GLEntry      `json:"cbe_ifb_gl_entry"`
}
