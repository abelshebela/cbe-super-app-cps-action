package dto

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/service"
)

type ServiceDetailsResponse struct {
	ID                 string               `json:"id"`
	ServiceID          string               `json:"service_id"`
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
	Enabled            bool                 `json:"enabled"`
	IsDeleted          bool                 `json:"is_deleted"`
	CreatedAt          time.Time            `json:"created_at"`
	LastModifiedAt     time.Time            `json:"last_modified_at"`
	DeletedAt          time.Time            `json:"deleted_at,omitempty"`
}

type UpdateServiceDetailsRequest struct {
	ID                 string               `json:"id"`
	ServiceID          string               `json:"service_id"`
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
	Enabled            bool                 `json:"enabled"`
	MakerID            string               `json:"maker_id"`
}

type UpdateServiceDetailsResponse struct {
	ActionID string `json:"action_id"`
}

type ApproveServiceDetailsRequest struct {
	ActionID        string `json:"action_id"`
	Approve         bool   `json:"approve"`
	CheckerID       string `json:"checker_id"`
	RejectionReason string `json:"rejection_reason,omitempty"`
}

type ServiceDetailsErrorResponse struct {
	Error string `json:"error"`
}
