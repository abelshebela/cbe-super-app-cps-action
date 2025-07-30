package productcode

type UpdateProductCodeRequest struct {
	ID                 string       `json:"id" form:"id"`
	ProductName        string       `json:"product_name" form:"product_name"`
	CBEProductCodes    ProductCodes `json:"cbe_product_codes" form:"cbe_product_codes"`
	CBEIFBProductCodes ProductCodes `json:"cbe_ifb_product_codes" form:"cbe_ifb_product_codes"`
}
