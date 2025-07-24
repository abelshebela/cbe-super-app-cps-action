package miniappmerchant

import (
	"time"

	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp_merchant"
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
		MiniAppIDs:    domain.MiniAppIDs,
		Enabled:       domain.Enabled,
		IsDeleted:     domain.IsDeleted,
		CreatedAt:     domain.CreatedAt,
		LastModified:  domain.LastModifiedAt,
	}
}

func ToMiniAppMerchantDomainFromCreateDTO(d *CreateMiniAppMerchantDTO) *domain.MiniAppMerchant {
	return &domain.MiniAppMerchant{
		MerchantType: d.Type,
		MerchantName: d.MerchantName,
		KYC: domain.KYC{
			Status: domain.KYCStatusComplete,
			Representative: domain.KYCInformation{
				Name:  d.MerchantRepresentativeName,
				Email: d.Email,
				Phone: d.PhoneNumber,
			},
		},
		PhoneNumber:       d.PhoneNumber,
		Email:             d.Email,
		BankAccountNumber: d.AccountNumber,
		CreatedAt:         time.Now(),
		LastModifiedAt:    time.Now(),
	}
}
func ToMiniAppMerchantDomainFromUpdateDTO(d *UpdateMiniAppMerchantDTO, id string) *domain.MiniAppMerchant {
	now := time.Now()
	merchant := &domain.MiniAppMerchant{
		ID:             id,
		LastModifiedAt: now,
	}

	if d.Type != nil {
		merchant.MerchantType = *d.Type
	}
	if d.MerchantName != nil {
		merchant.MerchantName = *d.MerchantName
	}

	// Safely set KYC fields only if any of them exist
	if d.MerchantRepresentativeName != nil || d.Email != nil || d.PhoneNumber != nil {
		merchant.KYC.Representative = domain.KYCInformation{}

		if d.MerchantRepresentativeName != nil && *d.MerchantRepresentativeName != "" {
			merchant.KYC.Representative.Name = *d.MerchantRepresentativeName
		}
		if d.Email != nil && *d.Email != "" {
			merchant.KYC.Representative.Email = *d.Email
		}
		if d.PhoneNumber != nil && *d.PhoneNumber != "" {
			merchant.KYC.Representative.Phone = *d.PhoneNumber
		}
	}

	if d.PhoneNumber != nil {
		merchant.PhoneNumber = *d.PhoneNumber
	}
	if d.Email != nil {
		merchant.Email = *d.Email
	}
	if d.AccountNumber != nil {
		merchant.BankAccountNumber = *d.AccountNumber
	}
	if d.MiniAppID != nil {
		merchant.MiniAppIDs = []string{*d.MiniAppID}
	}
	if d.Enabled != nil {
		merchant.Enabled = *d.Enabled
	}

	return merchant
}
