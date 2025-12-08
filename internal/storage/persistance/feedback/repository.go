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
	"math"

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

func (f *FeedbackStorage) FindByID(ctx context.Context, id string) (*feedback.FeedbackResponse, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "_id", Value: objID}}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "user_id_obj", Value: bson.D{{Key: "$toObjectId", Value: "$user_id"}}},
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

func (f *FeedbackStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponseForFeedback[[]*feedback.FeedbackResponse], error) {
	filter := bson.M{}
	searchKeys := bson.M{}
	allowedKeys := []string{"created_at", "user_id","responses"}
if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["responses"] = searchRegex
	}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$facet", Value: bson.D{
			{Key: "docs", Value: bson.A{
				bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
				bson.D{{Key: "$skip", Value: skip}},
				bson.D{{Key: "$limit", Value: limit}},
				bson.D{{Key: "$addFields", Value: bson.D{
					{Key: "user_id_obj", Value: bson.D{{Key: "$toObjectId", Value: "$user_id"}}},
				}}},
				bson.D{{Key: "$lookup", Value: bson.D{
					{Key: "from", Value: "members"},
					{Key: "localField", Value: "user_id_obj"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "user"},
				}}},
				bson.D{{Key: "$unwind", Value: bson.D{
					{Key: "path", Value: "$user"},
					{Key: "preserveNullAndEmptyArrays", Value: true},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
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
			}},
			{Key: "totalCount", Value: bson.A{
				bson.D{{Key: "$count", Value: "count"}},
			}},
			{Key: "averageRatings", Value: bson.A{
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "ratings", Value: bson.D{
						{Key: "$filter", Value: bson.D{
							{Key: "input", Value: bson.D{{Key: "$objectToArray", Value: "$responses"}}},
							{Key: "as", Value: "item"},
							{Key: "cond", Value: bson.D{
								{Key: "$and", Value: bson.A{
									bson.D{{Key: "$eq", Value: bson.A{"$$item.v.type", "rating"}}},
									bson.D{{Key: "$ne", Value: bson.A{"$$item.v.answer", nil}}},
								}},
							}},
						}},
					}},
				}}},
				bson.D{{Key: "$unwind", Value: "$ratings"}},
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: "$ratings.k"},
					{Key: "average", Value: bson.D{{Key: "$avg", Value: bson.D{{Key: "$toDouble", Value: "$ratings.v.answer"}}}}},
				}}},
			}},
		}}},
	}

	cursor, err := f.collection.Aggregate(ctx, pipeline)
	if err != nil {
		f.logger.Errorf("aggregate feedbacks: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		return nil, nil
	}

	var result struct {
		Docs       []*feedback.FeedbackResponse `bson:"docs"`
		TotalCount []struct {
			Count int64 `bson:"count"`
		} `bson:"totalCount"`
		AverageRatings []struct {
			ID      string  `bson:"_id"`
			Average float64 `bson:"average"`
		} `bson:"averageRatings"`
	}

	if err := cursor.Decode(&result); err != nil {
		f.logger.Errorf("decode aggregated result: %v", err)
		return nil, err
	}

	var total int64
	if len(result.TotalCount) > 0 {
		total = result.TotalCount[0].Count
	}

	averages := make(map[string]float64)
	for _, avg := range result.AverageRatings {
		averages[avg.ID] = math.Round(avg.Average*100) / 100
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponseForFeedback[[]*feedback.FeedbackResponse]{
		Data:           result.Docs,
		Meta:           meta,
		AverageRatings: averages,
	}, nil
}
