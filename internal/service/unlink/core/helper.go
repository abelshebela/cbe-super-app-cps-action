package core

import (
	"cbe-super-app-cps-action/internal/constants/dto/customer"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"fmt"

	member "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

//	func CreatArchiveUserDataWithLinkedAccount(ctx context.Context, archivedUserRepo storage.ArchivedUserRepository, archivedLinkedAccountRepo storage.ArchivedLinkedAccountRepository, userOldData *model.User, linkedAccountOldData *model.LinkedAccount, haveAccount bool) (error, error) {
//		ctx, span := local_util.TraceLogger(ctx, "core", "CreatArchiveUserDataWithLinkedAccount", "core", "core")
//		defer span.End()
func CreatArchiveUserDataWithLinkedAccount(ctx context.Context, archivedUserRepo storage.ArchivedUserRepository, archivedLinkedAccountRepo storage.ArchivedLinkedAccountRepository, userOldData *member.User, linkedAccountOldData *model.LinkedAccount, haveAccount bool) (error, error) {
	var archUserErr, archLinkedAccErr error

	lib.GoRoutinBaker(types.BakerOptions{Sequential: false, UseMutex: false},
		func() {
			ctx, spanUser := local_util.TraceLogger(ctx, "core", "CreateArchivedUser", "core", "core")
			defer spanUser.End()
			archUserErr = archivedUserRepo.Create(ctx, userOldData)
		},
		func() {
			if haveAccount {
				ctx, spanLinkedAcc := local_util.TraceLogger(ctx, "core", "CreateArchivedLinkedAccount", "core", "core")
				defer spanLinkedAcc.End()
				archLinkedAccErr = archivedLinkedAccountRepo.Create(ctx, linkedAccountOldData)
			}
		},
	)

	return archUserErr, archLinkedAccErr
}

//	func DeleteUserDataWithLinkedAccount(ctx context.Context, userRepo storage.UserRepository, LinkedAccountRepo storage.LinkedAccountRepository, user model.User, likedAccount model.LinkedAccount, haveAccount bool) (error, error) {
//		ctx, span := local_util.TraceLogger(ctx, "core", "DeleteUserDataWithLinkedAccount", "core", "core")
//		defer span.End()
func DeleteUserDataWithLinkedAccount(ctx context.Context, userRepo storage.UserRepository, LinkedAccountRepo storage.LinkedAccountRepository, user member.User, likedAccount model.LinkedAccount, haveAccount bool) (error, error) {
	var userErr, linkedAccErr error
	lib.GoRoutinBaker(types.BakerOptions{Sequential: false, UseMutex: false},
		func() {
			ctx, spanUser := local_util.TraceLogger(ctx, "core", "DeleteUser", "core", "core")
			defer spanUser.End()
			userErr = userRepo.Delete(ctx, user.ID.Hex())
			fmt.Printf("Deletion of user on the unlink error: %v", userErr)
		},
		func() {
			if haveAccount {
				ctx, spanLinkedAcc := local_util.TraceLogger(ctx, "core", "DeleteLinkedAccount", "core", "core")
				defer spanLinkedAcc.End()
				linkedAccErr = LinkedAccountRepo.Delete(ctx, likedAccount.ID.Hex())
				fmt.Printf("Deletion account on the linked account error: %v", linkedAccErr)
			}
		},
	)

	// fmt.Println("User Deletion from user collection user error: %v linked error: %v", userErr, likedAccount.ID)
	return userErr, linkedAccErr
}

// func buildArchivedUserResponse(user *model.ArchivedUser, archivedLinkedUser *model.ArchivedLinkedAccount) *unlink_dto.ArchivedUserResponse {
// 	return &unlink_dto.ArchivedUserResponse{
// 		ID:           user.ID,
// 		CustomerCode: user.UserCode,
// 		CustomerName: user.FullName,
// 		BranchCode:   archivedLinkedUser.BranchCode,
// 		Branch: unlink_dto.AccountBlockResponse{
// 			ID:   archivedLinkedUser.ID.Hex(),
// 			Name: branch.Name,
// 		},
// 		District: unlink_dto.AccountBlockResponse{
// 			ID:   district.ID.Hex(),
// 			Name: district.Name,
// 		},
// 		PhoneNumber:         user.PhoneNumber,
// 		Language:            user.Language,
// 		Avatar:              user.Avatar,
// 		Email:               user.Email,
// 		PushToken:           user.PushToken,
// 		CustomerNumber:      user.CustomerNumber,
// 		UserCategory:        user.UserCategory,
// 		Industry:            user.Industry,
// 		Sector:              user.Sector,
// 		Ownership:           user.Ownership,
// 		CustomerSegment:     user.CustomerSegment,
// 		BlockedReason:       user.BlockedReason,
// 		DeviceUUID:          user.DeviceUUID,
// 		AppVersion:          user.AppVersion,
// 		Gender:              user.Gender,
// 		MemberType:          user.MemberType,
// 		Platform:            user.Platform,
// 		DeviceStatus:        user.DeviceStatus,
// 		OnboardingMethod:    user.OnboardingMethod,
// 		BlockedOn:           user.BlockedOn,
// 		EnabledChannels:     user.EnabledChannels,
// 		LoginAttemptCount:   user.LoginAttemptCount,
// 		LastLoginAttempt:    user.LastLoginAttempt,
// 		LastLogin:           user.LastLogin,
// 		APPInstallationDate: user.APPInstallationDate,
// 		CreatedAt:           user.CreatedAt,
// 		ExpiryAt:            user.ExpiryAt,
// 		LastModifiedAt:      user.LastModifiedAt,
// 		IsBlocked:           user.IsBlocked,
// 		Enabled:             user.Enabled,
// 		FirstPinSet:         user.FirstPinSet,
// 		IsActivated:         user.IsActivated,
// 	}
// }

func MapToDto(cus *member.User) customer.FindCustomerByIDResponse {
	return customer.FindCustomerByIDResponse{
		ID:                   cus.ID,
		UserCode:             cus.UserCode,
		FullName:             cus.FullName,
		BranchCode:           cus.BranchCode,
		ActivationBranchCode: cus.ActivationBranchCode,
		PhoneNumber:          cus.PhoneNumber,
		Language:             cus.Language,
		Avatar:               cus.Avatar,
		Email:                cus.Email,
		PushToken:            cus.PushToken,
		CustomerNumber:       cus.CustomerNumber,
		UserCategory:         cus.UserCategory,
		Industry:             cus.Industry,
		Sector:               cus.Sector,
		Username:             cus.Username,
		Ownership:            cus.Ownership,
		CustomerSegment:      cus.CustomerSegment,
		BlockedReason:        cus.BlockedReason,
		DeviceUUID:           cus.DeviceUUID,
		AppVersion:           cus.AppVersion,
		Gender:               cus.Gender,
		MemberType:           cus.MemberType,
		Platform:             cus.Platform,
		DeviceStatus:         cus.DeviceStatus,
		OnboardingMethod:     cus.OnboardingMethod,
		KYCLevel:             cus.KYCLevel,
		BlockedOn:            cus.BlockedOn,
		IsUSSDEnabled:        cus.IsUSSDEnabled,
		ISuperappEnabled:     cus.ISuperappEnabled,
		LastLogin:            cus.LastLogin,
		APPInstallationDate:  cus.APPInstallationDate,
		Config:               cus.Config,
		CreatedAt:            cus.CreatedAt,
		ExpiryAt:             cus.ExpiryAt,
		LastModifiedAt:       cus.LastModifiedAt,
		IsBlocked:            cus.IsBlocked,
		IsLocked:             cus.IsLocked,
		Enabled:              cus.Enabled,
		FirstPinSet:          cus.FirstPinSet,
		IsActivated:          cus.IsActivated,
	}
}
