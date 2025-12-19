package productcode

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func ProductCodeMapper(productCode model.ProductCode) bson.M {
	update := bson.M{
		"last_modified_at": time.Now(),
	}

	if productCode.ProductName != "" {
		update["service_name"] = productCode.ProductName
	}
	empty := model.ProductCodes{}
	if productCode.CBEProductCodes != empty {
		update["cbe_product_codes"] = productCode.CBEProductCodes
	}
	if productCode.CBEIFBProductCodes != empty {
		update["cbe_ifb_product_codes"] = productCode.CBEIFBProductCodes
	}

	return update
}
