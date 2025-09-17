package service_details

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/service_details"
	"cbe-super-app-cps-action/internal/constants/model"
)

// MapToMinimumTransferCapResponse maps ServiceDetails to MinimumTransferCapResponse
func MapToMinimumTransferCapResponse(service *model.ServiceDetails) *dto.MinimumTransferCapResponse {
	return &dto.MinimumTransferCapResponse{
		ID:          service.ID,
		ServiceName: service.ServiceName,
		ServiceCode: service.ServiceCode,
		ServiceType: service.ServiceType,
		MinimumAmount:service.Cap.MinAmount,
		CreatedAt:   service.CreatedAt,
	}
}

// MapToMaximumTransferCapResponse maps ServiceDetails to MaximumTransferCapResponse
func MapToMaximumTransferCapResponse(service *model.ServiceDetails) *dto.MaximumTransferCapResponse {
	return &dto.MaximumTransferCapResponse{
		ID:              service.ID,
		ServiceName:     service.ServiceName,
		ServiceKey:      service.Key,
		ServiceCode:     service.ServiceCode,
		ServiceType:     service.ServiceType,
		Cap:             service.Cap,
		ProductCodes:    service.CBEProductCodes,
		IFBProductCodes: service.CBEIFBProductCodes,
		GLEntry:         service.CBEGLEntry,
		CreatedAt:       service.CreatedAt,
	}
}

// MapToServiceFeeResponse maps ServiceDetails to ServiceFeeResponse
func MapToServiceFeeResponse(service *model.ServiceDetails) *dto.ServiceFeeResponse {
	return &dto.ServiceFeeResponse{
		ID:              service.ID,
		ServiceName:     service.ServiceName,
		ServiceKey:      service.Key,
		ServiceType:     service.ServiceType,
		PaymentType:     service.PaymentType,
		Tiers:           service.Tiers,
		MinAmount:       service.AboveAmount,
		ProductCodes:    service.CBEProductCodes,
		IFBProductCodes: service.CBEIFBProductCodes,
		GLEntry:         service.CBEGLEntry,
		CreatedAt:       service.CreatedAt,
	}
}

// MapToServiceFeeDetailResponse maps ServiceDetails to ServiceFeeDetailResponse
func MapToServiceFeeDetailResponse(service *model.ServiceDetails) *dto.ServiceFeeDetailResponse {
	return &dto.ServiceFeeDetailResponse{
		ID:                 service.ID,
		ServiceCode:        service.ServiceCode,
		ServiceName:        service.ServiceName,
		ServiceType:        service.ServiceType,
		Key:                service.Key,
		Cap:                service.Cap,
		CBEProductCodes:    service.CBEProductCodes,
		CBEIFBProductCodes: service.CBEIFBProductCodes,
		AboveAmount:        service.AboveAmount,
		AboveServiceFee:    service.AboveServiceFee,
		PaymentType:        service.PaymentType,
		Tiers:              service.Tiers,
		CBEGLEntry:         service.CBEGLEntry,
		CBEIFBGLEntry:      service.CBEIFBGLEntry,
		Enabled:            service.Enabled,
		IsDeleted:          service.IsDeleted,
		CreatedAt:          service.CreatedAt,
		LastModifiedAt:     service.LastModifiedAt,
		DeletedAt:          service.DeletedAt,
	}
}

// MapToTotalTransferCapResponse maps HQ to TotalTransferCapResponse
func MapToTotalTransferCapResponse(hq *model.HQ) *dto.TotalTransferCapResponse {
	return &dto.TotalTransferCapResponse{
		ID:                hq.ID,
		TotalCap:          hq.TotalCap,
		CreatedAt:         hq.CreatedAt,
		UpdatedAtTotalCap: hq.UpdatedAtTotalCap,
	}
}

// MapSliceToMinimumTransferCapResponse maps a slice of ServiceDetails to MinimumTransferCapResponse
func MapSliceToMinimumTransferCapResponse(services []*model.ServiceDetails) []*dto.MinimumTransferCapResponse {
	result := make([]*dto.MinimumTransferCapResponse, len(services))
	for i, service := range services {
		result[i] = MapToMinimumTransferCapResponse(service)
	}
	return result
}

// MapSliceToMaximumTransferCapResponse maps a slice of ServiceDetails to MaximumTransferCapResponse
func MapSliceToMaximumTransferCapResponse(services []*model.ServiceDetails) []*dto.MaximumTransferCapResponse {
	result := make([]*dto.MaximumTransferCapResponse, len(services))
	for i, service := range services {
		result[i] = MapToMaximumTransferCapResponse(service)
	}
	return result
}

// MapSliceToServiceFeeResponse maps a slice of ServiceDetails to ServiceFeeResponse
func MapSliceToServiceFeeResponse(services []*model.ServiceDetails) []*dto.ServiceFeeResponse {
	result := make([]*dto.ServiceFeeResponse, len(services))
	for i, service := range services {
		result[i] = MapToServiceFeeResponse(service)
	}
	return result
}
