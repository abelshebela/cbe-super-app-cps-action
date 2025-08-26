package productcode

import (
	"cbe-super-app-cps-action/internal/constants/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func ProductCodeMapper(productCode model.ProductCode) bson.M {
	return bson.M{
		"$set": bson.M{
			"service_name":          productCode.ProductName,
			"cbe_product_codes":     productCode.CBEProductCodes,
			"cbe_ifb_product_codes": productCode.CBEIFBProductCodes,
			"last_modified_at":      productCode.LastUpdatedAt,
			"created_at":            productCode.CreatedAt,
		},
	}
}
