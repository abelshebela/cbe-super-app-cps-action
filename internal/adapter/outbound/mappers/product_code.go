package mappers

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	domian "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/product_code"
)

func ToProducCode(service model.Service) *domian.ProductCode {
	return &domian.ProductCode{
		ID:                 service.ID.Hex(),
		ProductName:        service.ServiceName,
		CBEProductCodes:    domian.ProductCodes(service.CBEProductCodes),
		CBEIFBProductCodes: domian.ProductCodes(service.CBEIFBProductCodes),
		CreatedAt:          service.CreatedAt,
		LastUpdatedAt:      service.LastModifiedAt,
	}
}
