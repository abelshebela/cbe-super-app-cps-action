package core

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"fmt"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

func CreatArchiveUserDataWithLinkedAccount(ctx context.Context, archivedUserRepo storage.ArchivedUserRepository, archivedLinkedAccountRepo storage.ArchivedLinkedAccountRepository, userOldData *model.User, linkedAccountOldData *model.LinkedAccount, haveAccount bool) (error, error) {
	var archUserErr, archLinkedAccErr error

	lib.GoRoutinBaker(types.BakerOptions{Sequential: false, UseMutex: false},
		func() {
			archUserErr = archivedUserRepo.Create(ctx, userOldData)
		},
		func() {
			if haveAccount {
				archLinkedAccErr = archivedLinkedAccountRepo.Create(ctx, linkedAccountOldData)
			}
		},
	)

	return archUserErr, archLinkedAccErr
}

func DeleteUserDataWithLinkedAccount(ctx context.Context, userRepo storage.UserRepository, LinkedAccountRepo storage.LinkedAccountRepository, user model.User, likedAccount model.LinkedAccount, haveAccount bool) (error, error) {
	var userErr, linkedAccErr error
	lib.GoRoutinBaker(types.BakerOptions{Sequential: false, UseMutex: false},
		func() {
			userErr = userRepo.Delete(ctx, user.ID.Hex())
			fmt.Println("Deletion of user on the unlink error: %v", userErr)
		},
		func() {
			if haveAccount {
				linkedAccErr = LinkedAccountRepo.Delete(ctx, likedAccount.ID.Hex())
				fmt.Println("Deletion account on the linked account error: %v", linkedAccErr)
			}
		},
	)

	// fmt.Println("User Deletion from user collection user error: %v linked error: %v", userErr, likedAccount.ID)
	return userErr, linkedAccErr
}
