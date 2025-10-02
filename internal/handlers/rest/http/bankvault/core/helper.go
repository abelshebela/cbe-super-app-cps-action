package core

import (
	helper "cbe-super-app-cps-action/internal/constants/dto/bankvault"

	"cbe-super-app-cps-action/internal/constants/model"
	"time"
)

// map dto to model for create

func ToDomainCreateBankVaultRequest(req helper.CreateBankVaultProductRequest, lockPeriodDays time.Duration) *model.BankVaultProduct {
	return &model.BankVaultProduct{
		Name:              req.Name,
		Description:       req.Description,
		Currency:          req.Currency,
		RateBps:           req.RateBps,
		Method:            helper.ToDomainMethod(req.Method),
		Frequency:         helper.ToDomainFrequency(req.Frequency),
		LockPeriod:        lockPeriodDays,
		MinAmount:         req.MinAmount,
		MaxAmount:         req.MaxAmount,
		EarlyUnlockFeeBps: req.EarlyUnlockFeeBps,
		IsActive:          false,
	}
}

// map updaterequest dto to model

func ToDomainUpdateBankVaultRequest(req helper.UpdateBankVaultProductRequest, lockPeriodDays *time.Duration) *model.UpdateBankVault {
	return &model.UpdateBankVault{
		Description: req.Description,
		MinAmount:   req.MinAmount,
		MaxAmount:   req.MaxAmount,
		UpdatedAt:   time.Now().UTC(),
	}
}
