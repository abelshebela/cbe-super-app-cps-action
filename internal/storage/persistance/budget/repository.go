package budget

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BudgetStorage struct {
	iconDal  dal.MongoDal[model.Icon, model.Icon]
	colorDal dal.MongoDal[model.Color, model.Color]
	cpsDal   dal.MongoDal[model.CPSAction, model.CPSAction]
	client   *mongo.Client
	logger   utils.Logger
}

var _ storage.BudgetRepository = (*BudgetStorage)(nil)

func NewBudgetRepository(client *mongo.Client, dbName string, collections []string, logger utils.Logger) storage.BudgetRepository {
	return &BudgetStorage{
		iconDal:  dal.NewMongoDal[model.Icon, model.Icon](client, dbName, collections[0]),
		colorDal: dal.NewMongoDal[model.Color, model.Color](client, dbName, collections[1]),
		client:   client,
		logger:   logger,
	}
}

func (b *BudgetStorage) CreateIcon(ctx context.Context, icon *model.Icon) error {
	_, err := b.iconDal.InsertOne(ctx, *icon)
	if err != nil {
		b.logger.Errorf("failed to create icon: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (b *BudgetStorage) FetchIcons(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Icon], error) {
	filter := bson.M{"is_deleted": false}
	allowedKeys := []string{"enabled"}

	filter, skip, limit := lib.FilterBuilder(*filterParams, bson.M{}, allowedKeys)

	if filterParams.Search != "" {
		searchRegex := bson.M{"$regex": filterParams.Search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"icon": searchRegex},
		}
	}

	data, err := b.iconDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		b.logger.Errorf("failed to fetch icons: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := b.iconDal.TotalCount(ctx, filter)
	if err != nil {
		b.logger.Errorf("failed counting icons: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParams.Page, filterParams.PerPage)
	return &types.PaginatedResponse[[]*model.Icon]{Data: data, Meta: meta}, nil
}

func (b *BudgetStorage) UpdateIcon(ctx context.Context, id string, icon *model.Icon) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorPINInvalid.Code)
	}

	filter := bson.M{"_id": objectID, "is_deleted": false}

	updateData := IconMapper(*icon)
	_, err = b.iconDal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		b.logger.Errorf("failed to create update icon: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (b *BudgetStorage) CreateColor(ctx context.Context, color *model.Color) error {
	// makerData := local_util.ExtractUserFromContext(ctx)
	// filterPending := bson.M{
	// 	"maker_id":      makerData.UserCode,
	// 	"department":    makerData.Department,
	// 	"action_status": "PENDING",
	// }

	// _, err := b.cpsDal.FindOne(ctx, filterPending, nil)
	// if err == nil {
	// 	return nil, errors.New(localization.ErrorPendingActionExists.Code)
	// }
	// if err != nil && err != mongo.ErrNoDocuments {
	// 	b.logger.Errorf("failed to check pending actions: %v", err)
	// 	return nil, errors.New(localization.ErrorUnexpectedError.Code)
	// }

	// existing, err := b.colorDal.FindOne(ctx, bson.M{"color": colorName}, nil)
	// if err != nil && err != mongo.ErrNoDocuments {
	// 	b.logger.Errorf("failed checking existing color: %v", err)
	// 	return nil, errors.New(localization.ErrorUnexpectedError.Code)
	// }
	// if existing != nil {
	// 	return nil, errors.New(localization.ErrorDuplicateColorExists.Code)
	// }

	// cpsAction.CreatedAt = time.Now()
	// cpsAction.MakerActionTime = time.Now()
	// cpsAction.ID = bson.NewObjectID()
	// cpsAction.LastModifiedAt = time.Now()

	_, err := b.colorDal.InsertOne(ctx, *color)
	if err != nil {
		b.logger.Errorf("failed to create cps action for color: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (b *BudgetStorage) FetchColors(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Color], error) {
	filter := bson.M{"is_deleted": false}
	allowedKeys := []string{"enabled"}

	filter, skip, limit := lib.FilterBuilder(*filterParams, bson.M{}, allowedKeys)
	if filterParams.Search != "" {
		searchRegex := bson.M{"$regex": filterParams.Search, "$options": "i"}
		filter["$or"] = []bson.M{{"color": searchRegex}}
	}

	data, err := b.colorDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		b.logger.Errorf("failed to fetch colors: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := b.colorDal.TotalCount(ctx, filter)
	if err != nil {
		b.logger.Errorf("failed to count colors: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParams.Page, filterParams.PerPage)
	return &types.PaginatedResponse[[]*model.Color]{Data: data, Meta: meta}, nil
}

func (b *BudgetStorage) UpdateColor(ctx context.Context, id string, color *model.Color) error {
	// var colorData model.Color
	// raw, err := json.Marshal(cpsAction.CurrentAction)
	// if err != nil {
	// 	b.logger.Errorf("failed to marshal current action: %v", err)
	// 	return nil, errors.New(localization.ErrorUnexpectedError.Code)
	// }
	// if err := json.Unmarshal(raw, &colorData); err != nil {
	// 	b.logger.Errorf("failed to unmarshal current action: %v", err)
	// 	return nil, errors.New(localization.ErrorUnexpectedError.Code)
	// }

	// found, err := b.colorDal.FindOne(ctx, bson.M{"color": colorData.Color}, nil)
	// if err != nil && err != mongo.ErrNoDocuments {
	// 	b.logger.Errorf("failed checking existing color: %v", err)
	// 	return nil, errors.New(localization.ErrorUnexpectedError.Code)
	// }
	// if found != nil && found.Color == colorData.Color {
	// 	return nil, errors.New(localization.ErrorDuplicateColorExists.Code)
	// }

	// makerData := local_util.ExtractUserFromContext(ctx)
	// filterPending := bson.M{"maker_id": makerData.UserCode, "department": makerData.Department, "action_status": "PENDING"}
	// _, err = b.cpsDal.FindOne(ctx, filterPending, nil)
	// if err == nil {
	// 	return nil, errors.New(localization.ErrorPendingActionExists.Code)
	// }
	// if err != nil && err != mongo.ErrNoDocuments {
	// 	b.logger.Errorf("failed checking pending: %v", err)
	// 	return nil, errors.New(localization.ErrorUnexpectedError.Code)
	// }

	// rawAction, err := bson.Marshal(cpsAction)
	// if err != nil {
	// 	b.logger.Errorf("failed to marshal cpsAction: %v", err)
	// 	return nil, errors.New(localization.ErrorUnexpectedError.Code)
	// }
	// var updateMap bson.M
	// if err := bson.Unmarshal(rawAction, &updateMap); err != nil {
	// 	b.logger.Errorf("failed to unmarshal cpsAction to map: %v", err)
	// 	return nil, errors.New(localization.ErrorUnexpectedError.Code)
	// }

	// action, err := b.cpsDal.UpdateOne(ctx, bson.M{"_id": cpsAction.ID}, bson.M{"$set": updateMap})
	// if err != nil {
	// 	if errors.Is(err, mongo.ErrNoDocuments) {
	// 		return nil, errors.New(localization.ErrorFileNotFound.Code)
	// 	}
	// 	b.logger.Errorf("failed to update cps action: %v", err)
	// 	return nil, errors.New(localization.ErrorUnexpectedError.Code)
	// }

	// return &action, nil

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorPINInvalid.Code)
	}

	filter := bson.M{"_id": objectID, "is_deleted": false}

	updateData := ColorMapper(*color)
	_, err = b.colorDal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		b.logger.Errorf("failed to create update color: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (b *BudgetStorage) AuthorizeCPSAction(ctx context.Context, cpsAction *model.CPSAction) error {
	if cpsAction == nil {
		return errors.New(localization.ErrorInvalidInputParameter.Code)
	}

	_, err := b.ApproveAction(ctx, *cpsAction)
	if err != nil {
		return err
	}
	return nil
}

func (b *BudgetStorage) GetByIDColor(ctx context.Context, id string) (*model.Color, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	color, err := b.colorDal.FindOne(ctx, bson.M{"_id": objectID}, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		b.logger.Errorf("failed to fetch color: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return color, nil
}

func (b *BudgetStorage) CheckColorExist(ctx context.Context, colorName string) (bool, error) {
	_, err := b.colorDal.FindOne(ctx, bson.M{"color": colorName}, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		b.logger.Errorf("error checking color existence: %v", err)
		return false, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return true, nil
}

func (b *BudgetStorage) CreateAction(ctx context.Context, cpsAction model.CPSAction) (*model.CPSAction, error) {
	if cpsAction.RequestAction == "BUDGET_UPDATE_COLOR" || cpsAction.RequestAction == "BUDGET_CREATE_COLOR" {
		var colorData model.Color
		raw, err := json.Marshal(cpsAction.CurrentAction)
		if err == nil {
			_ = json.Unmarshal(raw, &colorData)
		}
		if colorData.Color != "" {
			found, err := b.colorDal.FindOne(ctx, bson.M{"color": colorData.Color}, nil)
			if err != nil && err != mongo.ErrNoDocuments {
				b.logger.Errorf("failed checking color existence: %v", err)
				return nil, errors.New(localization.ErrorUnexpectedError.Code)
			}
			if found != nil && found.Color == colorData.Color {
				return nil, errors.New(localization.ErrorDuplicateColorExists.Code)
			}
		}
	}

	cpsAction.ID = bson.NewObjectID()
	created, err := b.cpsDal.InsertOne(ctx, cpsAction)
	if err != nil {
		b.logger.Errorf("failed to insert cps action: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return &created, nil
}

func (b *BudgetStorage) ApproveAction(ctx context.Context, approver model.CPSAction) (*model.CPSAction, error) {
	action, err := b.cpsDal.FindOne(ctx, bson.M{"action_code": approver.ActionCode}, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		b.logger.Errorf("failed to query cps action: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	if action.ActionStatus != "PENDING" {
		return nil, errors.New(localization.ErrorActionAlreadyExists.Code)
	}

	now := time.Now()
	action.CheckerActionTime = &now
	action.ActionStatus = string(constants.Verified)

	castToBsonM := func(input interface{}) (bson.M, error) {
		raw, err := bson.Marshal(input)
		if err != nil {
			return nil, err
		}
		var out bson.M
		err = bson.Unmarshal(raw, &out)
		return out, err
	}

	switch action.RequestAction {
	case "BUDGET_ICON", "BUDGET_CREATE_ICON":
		current, err := castToBsonM(action.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidActionFormat.Code)
		}
		iconURL, _ := current["icon_url"].(string)
		if iconURL == "" {
			return nil, errors.New(localization.ErrorInvalidActionFormat.Code)
		}

		switch action.ActionType {
		case constants.CREATE:
			icon := model.Icon{
				Icon:         iconURL,
				Enabled:      true,
				IsDeleted:    false,
				CreatedAt:    time.Now(),
				LastModified: time.Now(),
			}
			if _, err := b.iconDal.InsertOne(ctx, icon); err != nil {
				b.logger.Errorf("failed inserting icon on approve: %v", err)
				return nil, errors.New(localization.ErrorUnexpectedError.Code)
			}
		case constants.UPDATE:
			prev, err := castToBsonM(action.PreviousAction)
			if err != nil {
				return nil, errors.New(localization.ErrorInvalidActionFormat.Code)
			}
			iconID := prev["icon_id"]
			filter := bson.M{"_id": iconID, "is_deleted": false}
			update := bson.M{"$set": bson.M{"icon": iconURL, "last_modified": time.Now()}}
			if _, err := b.iconDal.UpdateOne(ctx, filter, update); err != nil {
				b.logger.Errorf("failed updating icon on approve: %v", err)
				return nil, errors.New(localization.ErrorUnexpectedError.Code)
			}
		default:
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

	case "BUDGET_COLOR", "BUDGET_CREATE_COLOR":
		current, err := castToBsonM(action.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidActionFormat.Code)
		}
		colorName, _ := current["color"].(string)
		if colorName == "" {
			return nil, errors.New(localization.ErrorInvalidActionFormat.Code)
		}

		switch action.ActionType {
		case constants.CREATE:
			color := model.Color{
				Color:     colorName,
				Enabled:   true,
				IsDeleted: false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if _, err := b.colorDal.InsertOne(ctx, color); err != nil {
				b.logger.Errorf("failed inserting color on approve: %v", err)
				return nil, errors.New(localization.ErrorUnexpectedError.Code)
			}
		case constants.UPDATE:
			prev, err := castToBsonM(action.PreviousAction)
			if err != nil {
				return nil, errors.New(localization.ErrorInvalidActionFormat.Code)
			}
			colorID := prev["color_id"]
			filter := bson.M{"_id": colorID, "is_deleted": false}
			update := bson.M{"$set": bson.M{"color": colorName, "updated_at": time.Now()}}
			if _, err := b.colorDal.UpdateOne(ctx, filter, update); err != nil {
				b.logger.Errorf("failed updating color on approve: %v", err)
				return nil, errors.New(localization.ErrorUnexpectedError.Code)
			}
		default:
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

	case "BUDGET_DELETE_ICON", "BUDGET_DELETE_COLOR":
		current, err := castToBsonM(action.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidActionFormat.Code)
		}
		if action.RequestAction == "BUDGET_DELETE_ICON" {
			id, _ := current["icon_id"].(string)
			if id == "" {
				return nil, errors.New(localization.ErrorInvalidActionFormat.Code)
			}
			filter := bson.M{"_id": id, "is_deleted": false}
			update := bson.M{"$set": bson.M{"is_deleted": true, "last_modified": time.Now()}}
			if _, err := b.iconDal.UpdateOne(ctx, filter, update); err != nil {
				b.logger.Errorf("failed to mark icon deleted: %v", err)
				return nil, errors.New(localization.ErrorUnexpectedError.Code)
			}
		} else {
			id, _ := current["color_id"].(string)
			if id == "" {
				return nil, errors.New(localization.ErrorInvalidActionFormat.Code)
			}
			filter := bson.M{"_id": id, "is_deleted": false}
			update := bson.M{"$set": bson.M{"is_deleted": true, "last_modified": time.Now()}}
			if _, err := b.colorDal.UpdateOne(ctx, filter, update); err != nil {
				b.logger.Errorf("failed to mark color deleted: %v", err)
				return nil, errors.New(localization.ErrorUnexpectedError.Code)
			}
		}

	default:
		b.logger.Errorf("unknown request action: %s", action.RequestAction)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": action.ID}
	update := bson.M{
		"$set": bson.M{
			"action_status":        constants.Verified,
			"checker_action_time":  action.CheckerActionTime,
			"checker_id":           approver.CheckerID,
			"checker_name":         approver.CheckerName,
			"checker_phone_number": approver.CheckerPhoneNumber,
			"last_modified_at":     time.Now(),
		},
	}
	res, err := b.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		b.logger.Errorf("failed to update cps action after approve: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return &res, nil
}

func UnmarshalMap(input interface{}, output interface{}) error {
	bytes, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, output)
}
