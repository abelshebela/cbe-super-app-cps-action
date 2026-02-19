package archived_linked_account

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ArchivedLinkedAccountStorage struct {
	dal    dal.MongoDal[model.LinkedAccount, model.ArchivedLinkedAccount]
	client *mongo.Client
	logger utils.Logger
}

type ArchivedLinkedAccountRepository interface {
	Create(ctx context.Context, user *model.LinkedAccount) error
	FindByID(ctx context.Context, id string, isUserId bool) (*model.ArchivedLinkedAccount, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.ArchivedLinkedAccount], error)
}

func NewArchivedLinkedAccountRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, logger utils.Logger) storage.ArchivedLinkedAccountRepository {
	return &ArchivedLinkedAccountStorage{
		dal:    dal.NewMongoDal[model.LinkedAccount, model.ArchivedLinkedAccount](client, cfg, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (l *ArchivedLinkedAccountStorage) FindByCustomerNumber(ctx context.Context, customerNumber string) (*model.ArchivedLinkedAccount, error) {
	filter := bson.M{
		"customer_number": customerNumber,
	}

	result, err := l.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}

func (l *ArchivedLinkedAccountStorage) Create(ctx context.Context, account *model.LinkedAccount) error {
	_, err := l.dal.InsertOne(ctx, *account)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (l *ArchivedLinkedAccountStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	return l.dal.DeleteOne(ctx, filter)
}

func (l *ArchivedLinkedAccountStorage) FindByID(ctx context.Context, id string, isUserId bool) (*model.ArchivedLinkedAccount, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	var filter bson.M
	if isUserId {
		filter = bson.M{"user_id": objID}
	} else {
		filter = bson.M{"_id": objID}
	}

	result, err := l.dal.FindOne(ctx, filter, nil)

	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}

func (l *ArchivedLinkedAccountStorage) FindByAccountNumber(ctx context.Context, accountNumber string) (*model.ArchivedLinkedAccount, error) {
	filter := bson.M{
		"account_number": accountNumber,
	}
	result, err := l.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}

	return result, nil
}

func (l *ArchivedLinkedAccountStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.ArchivedLinkedAccount], error) {
	filter := bson.M{
		"is_deleted": false,
	}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"account_holder_name": searchRegex},
			{"account_number": searchRegex},
		}
	}

	skip := int64((filterParam.Page - 1) * filterParam.PerPage)
	limit := int64(filterParam.PerPage)

	data, err := l.dal.FindAllWithPaginationE(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := l.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]model.ArchivedLinkedAccount]{
		Data: data,
		Meta: meta,
	}, nil
}
