package productcode

import (
	"time"
)

// ProductCodes represents the codes for a product in DTOs
type ProductCodes struct {
	PRD    string `json:"prd"`
	VATPRD string `json:"vat_prd"`
	SFPRD  string `json:"sf_prd"`
	TRXN   string `json:"-"`
}

// ProductCodeResponse represents the response format for a product code
type ProductCodeResponse struct {
	ID                 string       `json:"id"`
	ProductName        string       `json:"product_name"`
	CBEProductCodes    ProductCodes `json:"cbe_product_codes"`
	CBEIFBProductCodes ProductCodes `json:"cbe_ifb_product_codes"`
	CreatedAt          time.Time    `json:"created_at"`
	LastUpdatedAt      time.Time    `json:"last_updated_at"`
}

// UpdateProductCodeRequest represents the request payload for updating a product code
type ProductCodeRequest struct {
	ID                 string       `json:"id" form:"id"`
	ProductName        string       `json:"product_name" form:"product_name"`
	CBEProductCodes    ProductCodes `json:"cbe_product_codes" form:"cbe_product_codes"`
	CBEIFBProductCodes ProductCodes `json:"cbe_ifb_product_codes" form:"cbe_ifb_product_codes"`
}

// isEmpty checks if the ProductCodeRequest is empty, excluding TRXN fields
func (u ProductCodeRequest) isEmpty() bool {
	return u.ProductName == "" &&
		u.CBEProductCodes.PRD == "" && u.CBEProductCodes.VATPRD == "" && u.CBEProductCodes.SFPRD == "" &&
		u.CBEIFBProductCodes.PRD == "" && u.CBEIFBProductCodes.VATPRD == "" && u.CBEIFBProductCodes.SFPRD == ""
}
