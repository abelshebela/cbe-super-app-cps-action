package miniappmerchant

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/mini_app_merchant"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
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
		Enabled:       domain.Enabled,
		IsDeleted:     domain.IsDeleted,
		CreatedAt:     domain.CreatedAt,
		LastModified:  domain.LastModifiedAt,
	}
}

// Convert DTO to Domain model for service layer
func ToMiniAppMerchantDomainFromUpdateDTO(d *dto.MiniAppMerchantDTO) *model.MiniAppMerchant {
	return &model.MiniAppMerchant{
		ID:                bson.NewObjectID(),
		MerchantType:      d.Type,
		Code:              d.MerchantCode,
		MerchantName:      d.MerchantName,
		BankAccountNumber: d.AccountNumber,
		Email:             d.Email,
		PhoneNumber:       d.PhoneNumber,
		KYC: types.KYC{
			Representative: types.KYCInformation{
				Name:  d.MerchantRepresentativeName,
				Phone: d.PhoneNumber,
				Email: d.Email,
			},
			Status: "",
		},
		Branches:       d.Branches,
		Enabled:        true,
		IsDeleted:      false,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
	}
}
