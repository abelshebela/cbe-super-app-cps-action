package linked_account

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

type LinkedAccountStorage struct {
	dal    dal.MongoDal[model.LinkedAccount, model.LinkedAccount]
	client *mongo.Client
	logger utils.Logger
}

func NewLinkedAccountRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.LinkedAccountRepository {
	return &LinkedAccountStorage{
		dal:    dal.NewMongoDal[model.LinkedAccount, model.LinkedAccount](client, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (l *LinkedAccountStorage) FindByCustomerNumber(ctx context.Context, customerNumber string) (*model.LinkedAccount, error) {
	filter := bson.M{
		"customer_number": customerNumber,
		"is_deleted":      false,
	}

	result, err := l.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return result, nil
}

func (l *LinkedAccountStorage) Create(ctx context.Context, account *model.LinkedAccount) error {
	_, err := l.dal.InsertOne(ctx, *account)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (l *LinkedAccountStorage) Update(ctx context.Context, id string, account *model.LinkedAccount) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := LinkedAccountMapper(*account)

	_, err = l.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (l *LinkedAccountStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	return l.dal.DeleteOne(ctx, filter)
}

func (l *LinkedAccountStorage) FindByID(ctx context.Context, id string) (*model.LinkedAccount, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}

	result, err := l.dal.FindOne(ctx, filter, nil)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (l *LinkedAccountStorage) FindByAccountNumber(ctx context.Context, accountNumber string) (*model.LinkedAccount, error) {
	filter := bson.M{
		"account_number": accountNumber,
	}

	result, err := l.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return result, nil
}

func (l *LinkedAccountStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.LinkedAccount], error) {
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

	data, err := l.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, err
	}

	total, err := l.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.LinkedAccount]{
		Data: data,
		Meta: meta,
	}, nil
}
