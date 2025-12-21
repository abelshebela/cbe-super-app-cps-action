package productcode

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

// ProductCodeResponse represents the response format for a product code
type ProductCodeResponse struct {
	ID                 string             `json:"id"`
	ProductName        string             `json:"product_name"`
	CBEProductCodes    model.ProductCodes `json:"cbe_product_codes"`
	CBEIFBProductCodes model.ProductCodes `json:"cbe_ifb_product_codes"`
	CreatedAt          time.Time          `json:"created_at"`
	LastUpdatedAt      time.Time          `json:"last_modified_at"`
}
