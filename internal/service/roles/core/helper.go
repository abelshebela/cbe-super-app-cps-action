package core

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func CheckPortalCardsExistent(ctx context.Context, portalCards []string, portalCardRepo storage.PortalCardRepository) error {

	_, err := portalCardRepo.ValidatePortalCard(ctx, portalCards)
	if err != nil {
		return err
	}

	return nil
}

func RoleExistenChecker(ctx context.Context, types, roleId string, update imodel.JobRole, roleRepo storage.JobRoleRepository) error {

	var role *imodel.JobRole
	var err error

	resByCode, err := roleRepo.Find(ctx, bson.M{"code": update.Code})
	if err != nil {
		return err
	}

	resByName, err := roleRepo.Find(ctx, bson.M{"name": update.Name})
	if err != nil {
		return err
	}

	if types == constants.CREATE {
		if resByCode != nil {
			return errors.New(localization.ErrorUsedRoleExisting.Code)
		}
	} else if types == constants.UPDATE {
		if roleId == "" {
			return errors.New(localization.ErrorRoleIDMissing.Code)
		}
	}

	if roleId != "" {
		role, err = roleRepo.FindByID(ctx, roleId)
		if err != nil {
			return err
		}
	}

	if role != nil && (resByCode != nil || resByName != nil) {
		if role.ID.Hex() != resByCode.ID.Hex() {
			return errors.New(localization.ErrorUsedRoleExisting.Code)
		}

		if role.ID.Hex() != resByName.ID.Hex() {
			return errors.New(localization.ErrorUsedRoleExisting.Code)
		}
	}

	return nil
}
