package productcode

import "cbe-super-app-cps-action/internal/constants/model"

// Mapper functions
// ToDomainProductCodeRequest converts an HTTP ProductCodeRequest to a domain-level UpdateProductCodeRequest
func ToDomainProductCodeRequest(httpRequest UpdateProductCodeRequest) UpdateProductCodeRequest {
	return UpdateProductCodeRequest{
		ID:                 httpRequest.ID,
		ProductName:        httpRequest.ProductName,
		CBEProductCodes:    httpRequest.CBEProductCodes,
		CBEIFBProductCodes: httpRequest.CBEIFBProductCodes,
	}
}

// ToProductCodeResponse converts a domain ProductCode to a ProductCodeResponse
func ToProductCodeResponse(productCode model.ProductCode) ProductCodeResponse {
	return ProductCodeResponse{
		ID:                 productCode.ID,
		ProductName:        productCode.ProductName,
		CBEProductCodes:    productCode.CBEProductCodes,
		CBEIFBProductCodes: productCode.CBEIFBProductCodes,
		CreatedAt:          productCode.CreatedAt,
		LastUpdatedAt:      productCode.LastUpdatedAt,
	}
}

// ToProductCodeResponses converts a slice of domain ProductCode to a slice of ProductCodeResponse
func ToProductCodeResponses(productCodes []*model.ProductCode) []*ProductCodeResponse {
	responses := make([]*ProductCodeResponse, len(productCodes))
	for i, pc := range productCodes {
		responses[i] = &ProductCodeResponse{
			ID:                 pc.ID,
			ProductName:        pc.ProductName,
			CBEProductCodes:    pc.CBEProductCodes,
			CBEIFBProductCodes: pc.CBEIFBProductCodes,
			CreatedAt:          pc.CreatedAt,
			LastUpdatedAt:      pc.LastUpdatedAt,
		}
	}
	return responses
}

func ToProducCode(service model.ServiceDetails) *model.ProductCode {
	return &model.ProductCode{
		ID:                 service.ID.Hex(),
		ProductName:        service.ServiceName,
		CBEProductCodes:    model.ProductCodes(service.CBEProductCodes),
		CBEIFBProductCodes: model.ProductCodes(service.CBEIFBProductCodes),
		CreatedAt:          service.CreatedAt,
		LastUpdatedAt:      service.LastModifiedAt,
	}
}
