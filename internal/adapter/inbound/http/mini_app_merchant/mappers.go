package miniappmerchant

import (
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp_merchant"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp_merchant"
)

func ToMiniAppMerchantResponseDTO(domain *domain.MiniAppMerchant) *MiniAppMerchantResponseDTO {
	return &MiniAppMerchantResponseDTO{
		ID:           domain.ID,
		Code:         domain.Code,
		Type:         domain.MerchantType,
		MerchantName: domain.MerchantName,
		KYC: KYCDTO{
			Status: string(domain.KYC.Status),
			Representative: RepresentativeDTO{
				Name:  domain.KYC.Representative.Name,
				Phone: domain.KYC.Representative.Phone,
				Email: domain.KYC.Representative.Email,
			},
		},
		AccountNumber: domain.BankAccountNumber,
		MiniApps:      domain.MiniApps,
		Enabled:       domain.Enabled,
		IsDeleted:     domain.IsDeleted,
		CreatedAt:     domain.CreatedAt,
		LastModified:  domain.LastModifiedAt,
	}
}

func ToMiniAppMerchantDomainFromUpdateDTO(d *MiniAppMerchantDTO) *dto.MiniAppMerchantRequest {
	return &dto.MiniAppMerchantRequest{
		Type:                       d.Type,
		MerchantName:               d.MerchantName,
		MerchantRepresentativeName: d.MerchantRepresentativeName,
		PhoneNumber:                d.PhoneNumber,
		Email:                      d.Email,
		AccountNumber:              d.AccountNumber,
	}
}
