package core

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
)

func CheckRoleExistent(ctx context.Context, role string, roleRepo storage.RoleRepository) error {

	_, err := roleRepo.FindByRole(ctx, role)
	if err != nil {
		return errors.New(localization.ErrorRoleNotFound.Code)
	}
	return nil
}
func CheckJobTitleExistent(ctx context.Context, jobTitle string, roleRepo storage.RoleRepository) error {

	_, err := roleRepo.FindByName(ctx, jobTitle)
	if err != nil {
		return errors.New(localization.ErrorRoleNotFound.Code)
	}
	return nil
}
