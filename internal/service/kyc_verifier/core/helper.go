package core

import (
	"cbe-super-app-cps-action/internal/constants"
	accountLookupDto "cbe-super-app-cps-action/internal/constants/dto/account_lookup"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	accountLookup "cbe-super-app-cps-action/internal/storage/external_call/account_lookup"
	"context"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	member "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
shared_constant"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/constants"
shared_types  "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func AccountCreateAndLink(ctx context.Context, actionData model.CPSAction, id string, userData member.User, accountLookupService accountLookup.Account, userRepo storage.UserRepository, linkedAccountRepo storage.LinkedAccountRepository, logger utils.Logger) error {

	data := accountLookupDto.CreateAccountRequest{
		CustomerName:       userData.FullName,
		Gender:             constants.Gender(userData.Gender),
		PhoneNumber:        userData.PhoneNumber,
		AccountType:        string(userData.MemberType),
		AccountBranchType:  constants.AccountType(userData.MemberType),
		Picture:            userData.Avatar,
	}

	accountResponse, err := accountLookupService.CreateAccountWithFayda(ctx, data)
	if err != nil {
		logger.Errorf("failed to create account with fayda: %v", err)
		return err
	}

	if err := AccountLinker(ctx, actionData, id, userData, accountResponse, userRepo, linkedAccountRepo, logger); err != nil {
		logger.Errorf("failed to link account: %v", err)
		return err
	}

	return nil
}

func AccountLinker(ctx context.Context, actionData model.CPSAction, id string, userData member.User, account types.Account, userRepo storage.UserRepository, linkedAccountRepo storage.LinkedAccountRepository, logger utils.Logger) error {

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
				
				LinkedAt:          time.Now(),
				IsAccountActive:   true,
				AndOrStatus:       false,
				AccountBranchCode: account.AccountBranchCode,
				CurrencyCode:      account.AccountCurrency,
				IsMain:            true,
				MakerAndChecker: shared_types.MakerChecker{
					Linkers: struct {
						Maker   string `json:"maker" bson:"maker"`
						Checker string `json:"checker" bson:"checker"`
					}{Maker: actionData.MakerName, Checker: actionData.CheckerName},
					Unlinkers: struct {
						Maker   string `json:"maker" bson:"maker"`
						Checker string `json:"checker" bson:"checker"`
					}{Maker: "", Checker: ""},
				},
				CreatedAt: time.Now(),
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

func MapandUpdateuserFromKYC(ctx context.Context, userRepo storage.UserRepository, updated model.CustomerKYC, logger utils.Logger) error {
	user, err := userRepo.FindById(ctx, updated.UserID)
	if err != nil {
		logger.Errorf("failed to find user by id: %v", err)
		return err
	}
	// Map fields from updated CustomerKYC to user
	if user.FullName == "" && updated.KYCData.FullName != "" {
		user.FullName = updated.KYCData.FullName
	}
	if user.PhoneNumber == "" && updated.KYCData.PhoneNumber != "" {
		user.PhoneNumber = updated.KYCData.PhoneNumber
	}
	if string(user.Gender) == "" && updated.KYCData.Gender != "" {
		user.Gender = shared_constant.Gender(updated.KYCData.Gender)
	}
	if user.Avatar == "" && updated.KYCData.Picture != "" {
		user.Avatar = updated.KYCData.Picture
	}

	if string(user.MemberType) == "" && updated.KYCData.AccountType != "" {
		user.MemberType = shared_constant.MemberType(updated.KYCData.AccountType)
	}


	if err := userRepo.Update(ctx, user.ID.Hex(), user); err != nil {
		logger.Errorf("failed to update user from KYC: %v", err)
		return err
	}
	return nil
}
