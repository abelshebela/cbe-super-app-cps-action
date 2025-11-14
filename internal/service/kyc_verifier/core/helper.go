package core

import (
	"cbe-super-app-cps-action/internal/constants"
	accountLookupDto "cbe-super-app-cps-action/internal/constants/dto/account_lookup"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	accountLookup "cbe-super-app-cps-action/internal/storage/external_call/account_lookup"
	"context"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func AccountCreator(ctx context.Context, id string, userData model.User, accountLookupService accountLookup.Account, userRepo storage.UserRepository, linkedAccountRepo storage.LinkedAccountRepository, logger utils.Logger) error {

	data := accountLookupDto.CreateAccountRequest{
		CustomerAddress:    userData.Address.Zone + userData.Address.Region + userData.Address.Woreda + userData.Address.Kebele,
		CustomerName:       userData.FullName,
		Gender:             constants.Gender(userData.Gender),
		PhoneNumber:        userData.PhoneNumber,
		CustomerMotherName: userData.MotherName,
		AccountType:        string(userData.AccountType),
		AccountBranchType:  constants.AccountType(userData.MemberType),
		Picture:            userData.Avatar,
	}

	accountResponse, err := accountLookupService.CreateAccountWithFayda(ctx, data)
	if err != nil {
		logger.Errorf("failed to create account with fayda: %v", err)
		return err
	}

	if err := AccountLinker(ctx, id, userData, accountResponse, userRepo, linkedAccountRepo, logger); err != nil {
		logger.Errorf("failed to link account: %v", err)
		return err
	}

	return nil
}

func AccountLinker(ctx context.Context, id string, userData model.User, account types.Account, userRepo storage.UserRepository, linkedAccountRepo storage.LinkedAccountRepository, logger utils.Logger) error {

	lib.GoRoutinBaker(types.BakerOptions{UseMutex: true},
		func() {
			if err := linkedAccountRepo.Create(ctx, &model.LinkedAccount{
				UserID:            userData.ID,
				CustomerNumber:    account.CustomerNumber,
				AccountNumber:     account.AccountNumber,
				AccountHolderName: account.CustomerName,
				AccountType:       account.AccountType,
				BranchCode:        account.AccountBranchCode,
				LinkedStatus:      true,
				LastLinkedStatus:  false,
				LinkedAt:          time.Now(),
				LinkedBranch:      userData.BranchCode,
				RegistrationType:  constants.RegistrationTypeNew,
				IsAccountActive:   true,
				AndOrStatus:       false,
				AccountBranchCode: account.AccountBranchCode,
				CurrencyCode:      account.AccountCurrency,
				IsMain:            true,
				CreatedAt:         time.Now(),
			}); err != nil {
				logger.Errorf("failed to create linked account: %v", err)
			}
		},
		func() {
			if err := userRepo.Update(ctx, id, &userData); err != nil {
				logger.Errorf("failed to update user: %v", err)
			}
		})
	return nil
}
