package core

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func UpdateDataBuilder(ctx context.Context, avatarRepo storage.AvatarRepository, id, label string, fromEnableDisable bool, enable bool) (*model.Avatar, *model.Avatar, error) {

	var existed *model.Avatar
	var update model.Avatar
	var existedData []*model.Avatar
	var err error

	if !fromEnableDisable {
		existedData, err = avatarRepo.FindAll(ctx, bson.M{"label": label}, nil)
		if err != nil {
			return nil, nil, err
		}

		if len(existedData) > 0 {
			return nil, nil, errors.New(localization.ErrorAvatarAlreadyExist.Code)
		}

	} else {
		existed, err = avatarRepo.FindByID(ctx, id)
		if err != nil {
			return nil, nil, err
		}

		if existed.Enable == enable {
			if enable {
				return nil, nil, errors.New(localization.ErrorAvatarAlreadyEnabled.Code)
			} else {
				return nil, nil, errors.New(localization.ErrorAvatarAlreadyDisabled.Code)
			}
		}
		update = *existed
		update.Enable = enable

		return existed, &update, nil
	}

	existed, err = avatarRepo.FindByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	update = *existed

	return existed, &update, nil
}
