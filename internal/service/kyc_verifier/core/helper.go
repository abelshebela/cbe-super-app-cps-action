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

func AccountCreateAndLink(ctx context.Context, actionData model.CPSAction, id string, userData model.User, accountLookupService accountLookup.Account, userRepo storage.UserRepository, linkedAccountRepo storage.LinkedAccountRepository, logger utils.Logger) error {

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

	if err := AccountLinker(ctx, actionData, id, userData, accountResponse, userRepo, linkedAccountRepo, logger); err != nil {
		logger.Errorf("failed to link account: %v", err)
		return err
	}

	return nil
}

func AccountLinker(ctx context.Context, actionData model.CPSAction, id string, userData model.User, account types.Account, userRepo storage.UserRepository, linkedAccountRepo storage.LinkedAccountRepository, logger utils.Logger) error {

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
				LastLinkedStatus:  userData.LastAccountLinked,
				LinkedAt:          time.Now(),
				LinkedBranch:      userData.BranchCode,
				RegistrationType:  constants.RegistrationTypeNew,
				IsAccountActive:   true,
				AndOrStatus:       false,
				AccountBranchCode: account.AccountBranchCode,
				CurrencyCode:      account.AccountCurrency,
				IsMain:            true,
				MakerAndChecker: types.MakerChecker{
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
			userData.LastAccountLinked = true
			userData.IsActivated = true
			userData.MainAccount = account.AccountNumber
			userData.CustomerNumber = account.CustomerNumber
			userData.AccountType = constants.AccountType(account.AccountType)
			userData.MemberType = constants.MemberType(account.AccountBranchType)
			userData.IsVerified = true
			userData.IsSelfRegister = true
			userData.Gender = constants.Gender(account.Gender)
			userData.Address.Zone = account.CustomerAddress
			userData.Address.Region = account.CustomerAddress
			userData.Address.Woreda = account.CustomerAddress
			userData.Address.Kebele = account.CustomerAddress
			userData.MotherName = account.CustomerMotherName
			userData.Avatar = account.Picture
			userData.BranchCode = account.AccountBranchCode

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
		user.Gender = constants.Gender(updated.KYCData.Gender)
	}
	if user.Avatar == "" && updated.KYCData.Picture != "" {
		user.Avatar = updated.KYCData.Picture
	}
	if user.Nationality == "" && updated.KYCData.Nationality != "" {
		user.Nationality = updated.KYCData.Nationality
	}
	if user.BirthDate.IsZero() && !updated.KYCData.BirthDate.IsZero() {
		user.BirthDate = updated.KYCData.BirthDate
	}
	if user.DocumentFront == "" && updated.KYCData.DocumentFront != "" {
		user.DocumentFront = updated.KYCData.DocumentFront
	}
	if user.DocumentBack == "" && updated.KYCData.DocumentBack != "" {
		user.DocumentBack = updated.KYCData.DocumentBack
	}
	if string(user.AccountType) == "" && updated.KYCData.AccountType != "" {
		user.AccountType = constants.AccountType(updated.KYCData.AccountType)
	}
	// Assign Address fields individually if not set
	if user.Address.Zone == "" && updated.KYCData.Address.Zone != "" {
		user.Address.Zone = updated.KYCData.Address.Zone
	}
	if user.Address.Kebele == "" && updated.KYCData.Address.Kebele != "" {
		user.Address.Kebele = updated.KYCData.Address.Kebele
	}
	if user.Address.Woreda == "" && updated.KYCData.Address.Woreda != "" {
		user.Address.Woreda = updated.KYCData.Address.Woreda
	}
	if user.Address.Region == "" && updated.KYCData.Address.Region != "" {
		user.Address.Region = updated.KYCData.Address.Region
	}
	if user.MotherName == "" && updated.KYCData.MothersName != "" {
		user.MotherName = updated.KYCData.MothersName
	}
	// Optionally map other fields if needed
	if user.KYCLevel == 0 && updated.KYCLevel != 0 {
		user.KYCLevel = updated.KYCLevel
	}

	if err := userRepo.Update(ctx, user.ID.Hex(), user); err != nil {
		logger.Errorf("failed to update user from KYC: %v", err)
		return err
	}
	return nil
}
