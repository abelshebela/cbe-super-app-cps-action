package feedback

import (
	"cbe-super-app-cps-action/internal/constants/dto/feedback"
	"cbe-super-app-cps-action/internal/constants/lib"
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

type FeedbackStorage struct {
	dal        dal.MongoDal[model.Feedback, model.Feedback]
	client     *mongo.Client
	collection *mongo.Collection
	logger     utils.Logger
}

func NewFeedbackRepository(client *mongo.Client, dbName string, collectionName string, logger utils.Logger) storage.FeedbackRepository {
	db := client.Database(dbName)
	collection := db.Collection(collectionName)
	return &FeedbackStorage{
		dal:        dal.NewMongoDal[model.Feedback, model.Feedback](client, dbName, collectionName),
		client:     client,
		collection: collection,
		logger:     logger,
	}
}

func (f *FeedbackStorage) Create(ctx context.Context, feedback *model.Feedback) error {
	_, err := f.dal.InsertOne(ctx, *feedback)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (s *FeedbackStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*feedback.FeedbackResponse], error) {
	filter := bson.M{}
	searchKeys := bson.M{}
	allowedKeys := []string{"responses"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["responses"] = searchRegex
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "user_id_obj", Value: bson.D{
				{Key: "$toObjectId", Value: "$user_id"},
			}},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "members"},
			{Key: "localField", Value: "user_id_obj"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "user"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$user"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "responses", Value: 1},
			{Key: "created_at", Value: 1},
			{Key: "updated_at", Value: 1},
			{Key: "user", Value: bson.D{
				{Key: "_id", Value: bson.M{"$toString": "$user._id"}},
				{Key: "user_code", Value: "$user.user_code"},
				{Key: "full_name", Value: "$user.full_name"},
			}},
		}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
	}

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		s.logger.Errorf("aggregate feedbacks: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	defer cursor.Close(ctx)

	var data []*feedback.FeedbackResponse
	if err := cursor.All(ctx, &data); err != nil {
		s.logger.Errorf("decode feedbacks: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	countPipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$count", Value: "total"}},
	}

	countCursor, err := s.collection.Aggregate(ctx, countPipeline)
	if err != nil {
		s.logger.Errorf("count feedbacks: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	defer countCursor.Close(ctx)

	var countResult []bson.M
	if err := countCursor.All(ctx, &countResult); err != nil {
		s.logger.Errorf("decode count: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	var total int64
	if len(countResult) > 0 {
		if val, ok := countResult[0]["total"]; ok {
			if intVal, ok := val.(int32); ok {
				total = int64(intVal)
			}
		}
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*feedback.FeedbackResponse]{
		Data: data,
		Meta: meta,
	}, nil
}

func (f *FeedbackStorage) FindByID(ctx context.Context, id string) (*feedback.FeedbackResponse, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "_id", Value: objID},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "user_id_obj", Value: bson.D{
				{Key: "$toObjectId", Value: "$user_id"},
			}},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "members"},
			{Key: "localField", Value: "user_id_obj"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "user"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$user"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "responses", Value: 1},
			{Key: "created_at", Value: 1},
			{Key: "updated_at", Value: 1},
			{Key: "user", Value: bson.D{
				{Key: "_id", Value: bson.M{"$toString": "$user._id"}},
				{Key: "user_code", Value: "$user.user_code"},
				{Key: "full_name", Value: "$user.full_name"},
			}},
		}}},
	}

	cursor, err := f.collection.Aggregate(ctx, pipeline)
	if err != nil {
		f.logger.Errorf("aggregate feedback by id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		return nil, errors.New(localization.ErrorFileNotFound.Code)
	}

	var resp feedback.FeedbackResponse
	if err := cursor.Decode(&resp); err != nil {
		f.logger.Errorf("decode feedback response: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	return &resp, nil
}
