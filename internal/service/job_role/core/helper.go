package core

import (
	"cbe-super-app-cps-action/internal/constants"
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

func JobTitleExistentChecker(ctx context.Context, types, id, jobTitle string, jobRoleRepo storage.RoleRepository) error {

	res, err := jobRoleRepo.FindByName(ctx, jobTitle)
	if err != nil {
		return err
	}

	if types == constants.CREATE && res != nil {

		return errors.New(localization.ErrorUsedJobTitleExisting.Code)
	} else if types == constants.UPDATE {
		if res != nil {
			if res.JobTitle == jobTitle {
				return errors.New(localization.ErrorNoUpdatedJobTitle.Code)
			} else {
				return errors.New(localization.ErrorUsedJobTitleExisting.Code)
			}
		}
	}

	data, err := jobRoleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if data == nil {
		return errors.New(localization.ErrorRoleNotFound.Code)
	}

	if data.JobTitle == jobTitle {
		return nil
	}

	return nil
}
