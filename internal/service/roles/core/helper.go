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

	if types == constants.CREATE {
		resByName, err := roleRepo.Find(ctx, bson.M{"name": update.Name, "code": update.Code})
		if err != nil {
			if err.Error() != localization.ErrorResourceNotFound.Code {
				return err
			}
		}
		// On create, both code and name must be unique
		if resByName != nil {
			return errors.New(localization.ErrorUsedRoleExisting.Code)
		}
	} else if types == constants.UPDATE {
		if roleId == "" {
			return errors.New(localization.ErrorRoleIDMissing.Code)
		}
		role, err = roleRepo.FindByID(ctx, roleId)
		if err != nil {
			return err
		}

		// Check name/type uniqueness only when provided.
		if update.Name != "" || update.Type != "" {
			filter := bson.M{}
			if update.Name != "" {
				filter["name"] = update.Name
			}
			if update.Type != "" {
				filter["type"] = update.Type
			}
			if len(filter) > 0 {
				resByName, err := roleRepo.Find(ctx, filter)
				if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
					return err
				}
				if resByName != nil && role != nil && role.ID.Hex() != resByName.ID.Hex() {
					return errors.New(localization.ErrorUsedRoleExisting.Code)
				}
			}
		}

		// Check code uniqueness only when provided.
		if update.Code != "" {
			resByCode, err := roleRepo.Find(ctx, bson.M{"code": update.Code})
			if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
				return err
			}
			if resByCode != nil && role != nil && role.ID.Hex() != resByCode.ID.Hex() {
				return errors.New(localization.ErrorUsedRoleExisting.Code)
			}
		}
	}

	return nil
}
