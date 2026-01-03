package core

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

func CheckRoleExistent(ctx context.Context, role string, roleRepo storage.JobRoleRepository) error {

	_, err := roleRepo.FindByCode(ctx, role)
	if err != nil {
		return errors.New(localization.ErrorRoleNotFound.Code)
	}
	return nil
}
func CheckJobTitleExistent(ctx context.Context, prev model.Role, jobTitle string, roleRepo storage.RoleRepository) error {

	data, err := roleRepo.FindByName(ctx, jobTitle)
	if err != nil {
		if err.Error() == localization.ErrorResourceNotFound.Code {
			return nil
		}
		return err
	}

	if data != nil {
		if data.ID != prev.ID {
			return errors.New(localization.ErrorUsedJobTitleExisting.Code)
		}
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
