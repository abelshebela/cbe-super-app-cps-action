package productcode

import "time"

// ProductCode represents a product code entity
type ProductCode struct {
	ID                 string       `json:"id" bson:"_id"`
	ProductName        string       `json:"product_name" bson:"product_name"`
	CBEProductCodes    ProductCodes `json:"cbe_product_codes" bson:"cbe_product_codes"`
	CBEIFBProductCodes ProductCodes `json:"cbe_ifb_product_codes" bson:"cbe_ifb_product_codes"`
	CreatedAt          time.Time    `json:"created_at" bson:"created_at"`
	LastUpdatedAt      time.Time    `json:"last_updated_at" bson:"last_updated_at"`
}

// ProductCodes represents the codes for a product
type ProductCodes struct {
	PRD    string `json:"prd" bson:"prd"`
	VATPRD string `json:"vat_prd" bson:"vatprd"`
	SFPRD  string `json:"sf_prd" bson:"sfprd"`
	TRXN   string `json:"trxn" bson:"trxn"`
}
