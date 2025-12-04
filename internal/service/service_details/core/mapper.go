package service_details

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/service_details"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
)

// MapToMinimumTransferCapResponse maps ServiceDetails to MinimumTransferCapResponse
func MapToMinimumTransferCapResponse(service *model.ServiceDetails) *dto.MinimumTransferCapResponse {
	return &dto.MinimumTransferCapResponse{
		ID:            service.ID,
		ServiceName:   service.ServiceName,
		ServiceCode:   service.ServiceCode,
		ServiceType:   service.ServiceType,
		MinimumAmount: service.Cap.MinAmount,
		CreatedAt:     service.CreatedAt,
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
		GLEntry:         service.CBGLEntry,
		CreatedAt:       service.CreatedAt,
	}
}

func ServiceMapper(service *model.ServiceDetails, input dto.ServiceFeeDetailDTO) model.ServiceDetails {

	service.ServiceType = input.ServiceType
	service.PaymentType = string(input.PaymentType)
	service.SingleCapLevelOne = uint64(input.SingleCapLevelOne)
	service.DailyCapLevelOne = uint64(input.DailyCapLevelOne)
	service.MinAmountVirtual = uint64(input.MinAmountVIRTUAL)
	service.AboveAmount = input.AboveAmount
	service.AboveServiceFee = input.AboveServiceFee
	// map tiers from DTO to model types
	tiers := make([]types.Tier, len(input.Tiers))
	for i, t := range input.Tiers {
		tiers[i] = types.Tier{
			Min:       t.Min,
			Max:       t.Max,
			FeeAmount: t.FeeAmount,
		}
	}

	service.Tiers = tiers
	// map GL entries from DTO to model types
	service.CBGLEntry = types.GLEntry{
		ProductAccount:    input.CBglEntry.ProductAccount,
		ProductBranchCode: input.CBglEntry.ProductBranchCode,
		ServiceAccount:    input.CBglEntry.ServiceAccount,
		ServiceBranchCode: input.CBglEntry.ServiceBranchCode,
		VatAccount:        input.CBglEntry.VatAccount,
		VatBranchCode:     input.CBglEntry.VatBranchCode,
	}
	service.CBIFBGLEntry = types.IFBglEntry{
		ProductAccount:    input.IFBglEntry.ProductAccount,
		ProductBranchCode: input.IFBglEntry.ProductBranchCode,
		ServiceAccount:    input.IFBglEntry.ServiceAccount,
		ServiceBranchCode: input.IFBglEntry.ServiceBranchCode,
		VatAccount:        input.IFBglEntry.VatAccount,
		VatBranchCode:     input.IFBglEntry.VatBranchCode,
	}

	return *service
}

// MapToServiceFeeResponse maps ServiceDetails to ServiceFeeResponse
func MapToServiceFeeResponse(service *model.ServiceDetails) *dto.ServiceFeeResponse {
	return &dto.ServiceFeeResponse{
		ID:                service.ID,
		ServiceName:       service.ServiceName,
		DailyCapLevelOne:  service.DailyCapLevelOne,
		MinAmountVirtual:  service.MinAmountVirtual,
		SingleCapLevelOne: service.SingleCapLevelOne,
		ServiceKey:        service.Key,
		ServiceType:       service.ServiceType,
		PaymentType:       service.PaymentType,
		Tiers:             service.Tiers,
		MinAmount:         service.AboveAmount,
		ProductCodes:      service.CBEProductCodes,
		IFBProductCodes:   service.CBEIFBProductCodes,
		GLEntry:           service.CBGLEntry,
		UpdatedAt:         service.LastModifiedAt,
		CreatedAt:         service.CreatedAt,
	}
}

// MapToServiceFeeDetailResponse maps ServiceDetails to ServiceFeeDetailResponse
func MapToServiceFeeDetailResponse(service *model.ServiceDetails) *dto.ServiceFeeDetailResponse {
	return &dto.ServiceFeeDetailResponse{
		ID:                 service.ID,
		ServiceCode:        service.ServiceCode,
		ServiceName:        service.ServiceName,
		ServiceType:        service.ServiceType,
		DailyCapLevelOne:   service.DailyCapLevelOne,
		MinAmountVirtual:   service.MinAmountVirtual,
		SingleCapLevelOne:  service.SingleCapLevelOne,
		Key:                service.Key,
		Cap:                service.Cap,
		CBEProductCodes:    service.CBEProductCodes,
		CBEIFBProductCodes: service.CBEIFBProductCodes,
		AboveAmount:        service.AboveAmount,
		AboveServiceFee:    service.AboveServiceFee,
		PaymentType:        service.PaymentType,
		Tiers:              service.Tiers,
		CBEGLEntry:         service.CBGLEntry,
		CBEIFBGLEntry:      service.CBIFBGLEntry,
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
