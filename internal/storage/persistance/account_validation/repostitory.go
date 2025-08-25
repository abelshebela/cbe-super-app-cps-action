package accountvalidation

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AccountValidationStore struct {
	dal    dal.MongoDal[model.ValidationRule, model.ValidationRule]
	client *mongo.Client
	logger utils.Logger
}

// NewAccountValidationStore returns a ValidationRuleRepository
func NewAccountValidationStore(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.ValidationRuleRepository {
	return &AccountValidationStore{
		dal:    dal.NewMongoDal[model.ValidationRule, model.ValidationRule](client, dbName, collection),
		client: client,
		logger: logger,
	}
}

// GetAccountValidationByID implements ValidationRuleRepository
func (l *AccountValidationStore) GetAccountValidationByID(ctx context.Context, id string) (*model.ValidationRule, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}

	result, err := l.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateAccountValidation implements ValidationRuleRepository
func (a *AccountValidationStore) UpdateAccountValidation(ctx context.Context, id string, rule *model.ValidationRule) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := AccountValidationMapper(*rule)

	_, err = a.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

// GetAllAccountValidation implements ValidationRuleRepository
func (l *AccountValidationStore) GetAllAccountValidation(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.ValidationRule], error) {
	filter := bson.M{"is_deleted": false}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["name"] = searchRegex
	}

	skip := int64((filterParam.Page - 1) * filterParam.PerPage)
	limit := int64(filterParam.PerPage)

	data, err := l.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, err
	}

	total, err := l.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.ValidationRule]{
		Data: data,
		Meta: meta,
	}, nil
}
