package core

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
)

func CreatArchiveUserDataWithLinkedAccount(ctx context.Context, archivedUserRepo storage.ArchivedUserRepository, archivedLinkedAccountRepo storage.ArchivedLinkedAccountRepository, userOldData *model.User, linkedAccountOldData *model.LinkedAccount) (error, error) {
	var archUserErr, archLinkedAccErr error

	lib.GoRoutinBaker(types.BakerOptions{Sequential: false, UseMutex: false},
		func() {
			archUserErr = archivedUserRepo.Create(ctx, userOldData)
		},
		func() {
			archLinkedAccErr = archivedLinkedAccountRepo.Create(ctx, linkedAccountOldData)
		},
	)

	return archUserErr, archLinkedAccErr
}

func DeleteUserDataWithLinkedAccount(ctx context.Context, userRepo storage.UserRepository, LinkedAccountRepo storage.LinkedAccountRepository, userID string, likedAccountId string) (error, error) {
	var archUserErr, archLinkedAccErr error
	// archUserErr = userRepo.Delete(ctx, userID)
	// archLinkedAccErr = LinkedAccountRepo.Delete(ctx, likedAccountId)

	lib.GoRoutinBaker(types.BakerOptions{Sequential: false, UseMutex: false},
		func() {
			archUserErr = userRepo.Delete(ctx, userID)
		},
		func() {
			archLinkedAccErr = LinkedAccountRepo.Delete(ctx, likedAccountId)
		},
	)

	return archUserErr, archLinkedAccErr
}
