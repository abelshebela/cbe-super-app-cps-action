package service_details

import (
	"time"

	imodel "cbe-super-app-cps-action/internal/constants/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ServiceDetailsMapper maps a ServiceDetails model to a bson.M for updates
func ServiceDetailsMapper(details imodel.ServiceDetails) bson.M {
	return bson.M{
		"service_code":          details.ServiceCode,
		"service_name":          details.ServiceName,
		"service_type":          details.ServiceType,
		"key":                   details.Key,
		"cap":                   details.Cap,
		"cbe_product_codes":     details.CBEProductCodes,
		"cbe_ifb_product_codes": details.CBEIFBProductCodes,
		"above_amount":          details.AboveAmount,
		"above_service_fee":     details.AboveServiceFee,
		"payment_type":          details.PaymentType,
		"tiers":                 details.Tiers,
		"cbe_gl_entry":          details.CBGLEntry,
		"cbe_ifb_gl_entry":      details.CBIFBGLEntry,
		"enabled":               details.Enabled,
		"last_modified_at":      time.Now(),
	}
}
