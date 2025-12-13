package core

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"fmt"
)

func CreatArchiveUserDataWithLinkedAccount(ctx context.Context, archivedUserRepo storage.ArchivedUserRepository, archivedLinkedAccountRepo storage.ArchivedLinkedAccountRepository, userOldData *model.User, linkedAccountOldData *model.LinkedAccount, haveAccount bool) (error, error) {
	ctx, span := local_util.TraceLogger(ctx, "core", "CreatArchiveUserDataWithLinkedAccount", "core", "core")
	defer span.End()
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

func DeleteUserDataWithLinkedAccount(ctx context.Context, userRepo storage.UserRepository, LinkedAccountRepo storage.LinkedAccountRepository, user model.User, likedAccount model.LinkedAccount, haveAccount bool) (error, error) {
	ctx, span := local_util.TraceLogger(ctx, "core", "DeleteUserDataWithLinkedAccount", "core", "core")
	defer span.End()
	var userErr, linkedAccErr error
	lib.GoRoutinBaker(types.BakerOptions{Sequential: false, UseMutex: false},
		func() {
			ctx, spanUser := local_util.TraceLogger(ctx, "core", "DeleteUser", "core", "core")
			defer spanUser.End()
			userErr = userRepo.Delete(ctx, user.ID.Hex())
			fmt.Println("Deletion of user on the unlink error: %v", userErr)
		},
		func() {
			if haveAccount {
				ctx, spanLinkedAcc := local_util.TraceLogger(ctx, "core", "DeleteLinkedAccount", "core", "core")
				defer spanLinkedAcc.End()
				linkedAccErr = LinkedAccountRepo.Delete(ctx, likedAccount.ID.Hex())
				fmt.Println("Deletion account on the linked account error: %v", linkedAccErr)
			}
		},
	)

	// fmt.Println("User Deletion from user collection user error: %v linked error: %v", userErr, likedAccount.ID)
	return userErr, linkedAccErr
}
