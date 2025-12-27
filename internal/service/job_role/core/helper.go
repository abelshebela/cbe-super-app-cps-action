package core

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
)

func CheckRoleExistent(ctx context.Context, roleId string, roleRepo storage.RoleRepository) error {

	_, err := roleRepo.FindByName(ctx, roleId)
	if err != nil {
		return errors.New(localization.ErrorRoleNotFound.Code)
	}
	return nil
}
