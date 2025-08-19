package feedback

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/localization"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type FeedbackStorage struct {
	dal    dal.MongoDal[model.Feedback, model.Feedback]
	client *mongo.Client
	logger utils.Logger
}

func NewFeedbackRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.FeedbackRepository {
	return &FeedbackStorage{
		dal:    dal.NewMongoDal[model.Feedback, model.Feedback](client, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (f *FeedbackStorage) Create(ctx context.Context, feedback *model.Feedback) error {
	_, err := f.dal.InsertOne(ctx, *feedback)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (f *FeedbackStorage) Update(ctx context.Context, id string, feedback *model.Feedback) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := FeedbackMapper(*feedback)

	_, err = f.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (f *FeedbackStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	return f.dal.DeleteOne(ctx, filter)
}

func (f *FeedbackStorage) FindByID(ctx context.Context, id string) (*model.Feedback, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}

	result, err := f.dal.FindOne(ctx, filter, nil)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (f *FeedbackStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Feedback], error) {
	filter := bson.M{
		"is_deleted": false,
	}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["user_id"] = searchRegex
	}

	skip := int64((filterParam.Page - 1) * filterParam.PerPage)
	limit := int64(filterParam.PerPage)

	data, err := f.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, err
	}

	total, err := f.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.Feedback]{
		Data: data,
		Meta: meta,
	}, nil
}
