package productcode

import (
	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/product_code"
)

// ToDomainProductCodeRequest converts an HTTP ProductCodeRequest to a domain-level UpdateProductCodeRequest
func ToDomainProductCodeRequest(httpRequest ProductCodeRequest) entity.UpdateProductCodeRequest {
	return entity.UpdateProductCodeRequest{
		ID:                 httpRequest.ID,
		ProductName:        httpRequest.ProductName,
		CBEProductCodes:    entity.ProductCodes(httpRequest.CBEProductCodes),
		CBEIFBProductCodes: entity.ProductCodes(httpRequest.CBEIFBProductCodes),
	}
}

// ToProductCodeResponse converts an entity.ProductCode to a ProductCodeResponse
func ToProductCodeResponse(productCode entity.ProductCode)ProductCodeResponse {
	return ProductCodeResponse{
		ID:                 productCode.ID,
		ProductName:        productCode.ProductName,
		CBEProductCodes:    ProductCodes(productCode.CBEProductCodes),
		CBEIFBProductCodes: ProductCodes(productCode.CBEIFBProductCodes),
		CreatedAt:          productCode.CreatedAt,
		LastUpdatedAt:      productCode.LastUpdatedAt,
	}
}

// ToProductCodeResponses converts a slice of entity.ProductCode to a slice of ProductCodeResponse
func ToProductCodeResponses(productCodes []*entity.ProductCode) []*ProductCodeResponse {
	responses := make([]*ProductCodeResponse, len(productCodes))
	for i, pc := range productCodes {
		responses[i] = &ProductCodeResponse{
			ID:                 pc.ID,
			ProductName:        pc.ProductName,
			CBEProductCodes:    ProductCodes(pc.CBEProductCodes),
			CBEIFBProductCodes: ProductCodes(pc.CBEIFBProductCodes),
			CreatedAt:          pc.CreatedAt,
			LastUpdatedAt:      pc.LastUpdatedAt,
		}
	}
	return responses
}
