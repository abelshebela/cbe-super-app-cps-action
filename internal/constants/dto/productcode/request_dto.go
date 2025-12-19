package productcode

import "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

type UpdateProductCodeRequest struct {
	ID                 string             `json:"id" form:"id"`
	ProductName        string             `json:"product_name" form:"product_name"`
	CBEProductCodes    model.ProductCodes `json:"cbe_product_codes" form:"cbe_product_codes"`
	CBEIFBProductCodes model.ProductCodes `json:"cbe_ifb_product_codes" form:"cbe_ifb_product_codes"`
}
