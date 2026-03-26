package linked_account

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
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

type LinkedAccountStorage struct {
	dal           dal.MongoDal[model.LinkedAccount, model.LinkedAccount]
	client        *mongo.Client
	kafkaProducer kafka.ClientOrchestrationProducer
	logger        utils.Logger
}

func NewLinkedAccountRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, kafkaProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.LinkedAccountRepository {
	return &LinkedAccountStorage{
		dal:           dal.NewMongoDal[model.LinkedAccount, model.LinkedAccount](client, cfg, dbName, collection),
		client:        client,
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

func (l *LinkedAccountStorage) FindByCustomerNumber(ctx context.Context, customerNumber string) (*model.LinkedAccount, error) {
	filter := bson.M{
		"customer_number": customerNumber,
	}

	result, err := l.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}

func (l *LinkedAccountStorage) Create(ctx context.Context, account *model.LinkedAccount) error {
	newLinkedAccount, err := l.dal.InsertOne(ctx, *account)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	l.kafkaProducer.PublishMessage(ctx, newLinkedAccount, string(constants.ClientOrchestrationLinkedAccountTopic), string(constants.ClientOrchestrationLinkedAccountTopic), "new linked account created")

	return nil
}

func (l *LinkedAccountStorage) Update(ctx context.Context, id string, account *model.LinkedAccount) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := LinkedAccountMapper(*account)

	updatedLinkedAccount, err := l.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		return local_util.HandleDBError(err)
	}

	l.kafkaProducer.PublishMessage(ctx, updatedLinkedAccount, string(constants.ClientOrchestrationLinkedAccountTopic), string(constants.ClientOrchestrationLinkedAccountTopic), "linked account updated")

	return nil
}

func (l *LinkedAccountStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}
	return l.dal.DeleteOneH(ctx, filter)
}

func (l *LinkedAccountStorage) FindByID(ctx context.Context, id string) (*model.LinkedAccount, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}

	result, err := l.dal.FindOne(ctx, filter, nil)

	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}

func (l *LinkedAccountStorage) FindByAccountNumber(ctx context.Context, accountNumber string) (*model.LinkedAccount, error) {
	filter := bson.M{
		"account_number": accountNumber,
	}
	result, err := l.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}

	return result, nil
}

func (l *LinkedAccountStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.LinkedAccount], error) {
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

	return &types.PaginatedResponse[[]model.LinkedAccount]{
		Data: data,
		Meta: meta,
	}, nil
}
