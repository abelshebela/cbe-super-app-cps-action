package core

import (
	helper "cbe-super-app-cps-action/internal/constants/dto/bankvault"

	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

func ToDomainCreateBankVaultRequest(req helper.CreateBankVaultProductRequest, lockPeriodDays float64) *model.BankVaultProduct {
	return &model.BankVaultProduct{
		Name:                       req.Name,
		Currency:                   "ETB",
		Interest:                   req.Interest,
		Method:                     "COMPOUND",
		Frequency:                  req.Frequency,
		LockPeriod:                 lockPeriodDays,
		MinAmount:                  req.MinAmount,
		MaxAmount:                  req.MaxAmount,
		ApplyInterestOnEarlyUnlock: req.ApplyInterestOnEarlyUnlock,
		IsActive:                   false,
	}
}

func ToDomainUpdateBankVaultRequest(req helper.UpdateBankVaultProductRequest, lockPeriodDays *time.Duration) *model.UpdateBankVault {
	return &model.UpdateBankVault{
		MinAmount: req.MinAmount,
		MaxAmount: req.MaxAmount,
		UpdatedAt: time.Now().UTC(),
	}
}
