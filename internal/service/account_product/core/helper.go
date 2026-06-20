package account_product_core

import (
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func MapFromAction(raw interface{}) (imodel.AccountProduct, error) {
	b, err := bson.Marshal(raw)
	if err != nil {
		return imodel.AccountProduct{}, err
	}
	var ap imodel.AccountProduct
	if err := bson.Unmarshal(b, &ap); err != nil {
		return imodel.AccountProduct{}, err
	}
	return ap, nil
}
