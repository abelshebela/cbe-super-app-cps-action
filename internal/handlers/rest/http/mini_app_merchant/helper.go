package miniappmerchant

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/mini_app_merchant"
	"cbe-super-app-cps-action/internal/constants/model"
	"time"
)

// Convert domain model to Response DTO
func ToMiniAppMerchantResponseDTO(domain *model.MiniAppMerchant) *dto.MiniAppMerchantResponseDTO {
	return &dto.MiniAppMerchantResponseDTO{
		ID:           domain.ID.Hex(),
		Code:         domain.Code,
		Type:         domain.MerchantType,
		MerchantName: domain.MerchantName,
		KYC: dto.KYCDTO{
			Status: string(domain.KYC.Status),
			Representative: dto.RepresentativeDTO{
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

// Convert DTO to Domain model for service layer
func ToMiniAppMerchantDomainFromUpdateDTO(d *dto.MiniAppMerchantDTO) *model.MiniAppMerchant {
	return &model.MiniAppMerchant{
		ID:                d.ID,
		MerchantType:      d.Type,
		MerchantName:      d.MerchantName,
		BankAccountNumber: d.AccountNumber,
		KYC: model.KYC{
			Representative: model.KYCInformation{
				Name:  d.MerchantRepresentativeName,
				Phone: d.PhoneNumber,
				Email: d.Email,
			},
			Status: "",
		},
		Enabled:        true,
		IsDeleted:      false,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
	}
}
