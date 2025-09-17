package productcode

import (
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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

func (pc ProductCodes) validate() error {
	return validation.ValidateStruct(&pc,
		validation.Field(&pc.PRD,
			validation.By(utils.NotBlank),
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&pc.VATPRD,
			validation.By(utils.NotBlank),
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&pc.SFPRD,
			validation.By(utils.NotBlank),
			validation.By(utils.NoSpecialChars),
		),
	)
}

func (u ProductCodeRequest) validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.ProductName,
			validation.By(utils.NotBlank),
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&u.CBEProductCodes,
			validation.By(func(value interface{}) error {
				if pc, ok := value.(ProductCodes); ok {
					return pc.validate()
				}
				return nil
			}),
		),
		validation.Field(&u.CBEIFBProductCodes,
			validation.By(func(value interface{}) error {
				if pc, ok := value.(ProductCodes); ok {
					return pc.validate()
				}
				return nil
			}),
		),
	)
}
