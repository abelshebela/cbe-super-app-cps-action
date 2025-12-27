package core

import (
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

func RoleExistenChecker(ctx context.Context, roleId string, update imodel.JobRole, roleRepo storage.JobRoleRepository) error {

	var role *imodel.JobRole
	var err error
	if roleId != "" {
		role, err = roleRepo.FindByID(ctx, roleId)
		if err != nil {
			return err
		}
	}

	res, err := roleRepo.Find(ctx, bson.M{"code": update.Code, "name": update.Name})
	if err != nil {
		return err
	}

	if role != nil && res != nil {
		if role.ID.Hex() != res.ID.Hex() {
			return errors.New(localization.ErrorUsedRoleExisting.Code)
		}
	}

	return nil
}
