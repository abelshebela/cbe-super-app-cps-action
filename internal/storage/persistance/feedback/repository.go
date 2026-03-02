package feedback

import (
	"cbe-super-app-cps-action/internal/constants/dto/feedback"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	local_model "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type FeedbackStorage struct {
	dal                 dal.MongoDal[local_model.Feedback, local_model.Feedback]
	customerFeedbackDal dal.MongoDal[local_model.CustomerFeedback, local_model.CustomerFeedback]
	surveyFeedbackDal   dal.MongoDal[local_model.SurveyFeedback, local_model.SurveyFeedback]
	client              *mongo.Client
	feedbackCollection  *mongo.Collection
	customerCollection  *mongo.Collection
	logger              utils.Logger
}

func NewFeedbackRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, feedbackCollection, customerFeedbackCollection, surveyFeedbackCollection string, logger utils.Logger) storage.FeedbackRepository {
	return &FeedbackStorage{
		dal:                 dal.NewMongoDal[local_model.Feedback, local_model.Feedback](client, cfg, dbName, feedbackCollection),
		customerFeedbackDal: dal.NewMongoDal[local_model.CustomerFeedback, local_model.CustomerFeedback](client, cfg, dbName, customerFeedbackCollection),
		surveyFeedbackDal:   dal.NewMongoDal[local_model.SurveyFeedback, local_model.SurveyFeedback](client, cfg, dbName, surveyFeedbackCollection),
		client:              client,
		feedbackCollection:  client.Database(dbName).Collection(feedbackCollection),
		customerCollection:  client.Database(dbName).Collection(customerFeedbackCollection),
		logger:              logger,
	}
}

func (f *FeedbackStorage) Create(ctx context.Context, feedback *local_model.Feedback) error {
	f.logger.Infof("[FeedbackStorage][Create] creating feedback")
	feed, err := f.dal.InsertOne(ctx, *feedback)
	if err != nil {
		f.logger.Errorf("[FeedbackStorage][Create] failed to create feedback: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	f.logger.Infof("[FeedbackStorage][Create] feedback created successfully", feed, feed.UserCode)

	return nil
}

func (f *FeedbackStorage) CreateSurveyFeedback(ctx context.Context, surveyFeedback *local_model.SurveyFeedback) error {
	f.logger.Infof("[FeedbackStorage][CreateSurveyFeedback] creating survey feedback")
	feed, err := f.surveyFeedbackDal.InsertOne(ctx, *surveyFeedback)
	if err != nil {
		f.logger.Errorf("[FeedbackStorage][CreateSurveyFeedback] failed to create survey feedback: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	f.logger.Infof("[FeedbackStorage][CreateSurveyFeedback] survey feedback created successfully", feed, feed.UserCode)

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
				{Key: "phone_number", Value: "$user.phone_number"},
			}},
		}}},
	}

	f.logger.Infof("[FeedbackStorage][FindByID] fetching feedback by id: %s", id)
	cursor, err := f.feedbackCollection.Aggregate(ctx, pipeline)
	if err != nil {
		f.logger.Errorf("[FeedbackStorage][FindByID] failed to aggregate feedback: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		f.logger.Errorf("[FeedbackStorage][FindByID] feedback not found")
		return nil, errors.New(localization.ErrorFileNotFound.Code)
	}

	var resp feedback.FeedbackResponse
	if err := cursor.Decode(&resp); err != nil {
		f.logger.Errorf("[FeedbackStorage][FindByID] failed to decode feedback response: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	f.logger.Infof("[FeedbackStorage][FindByID] feedback retrieved successfully")
	return &resp, nil
}

// func (f *FeedbackStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponseForFeedback[[]*feedback.FeedbackResponse], error) {
// 	filter := bson.M{}
// 	searchKeys := bson.M{}
// 	allowedKeys := []string{"created_at", "user_id", "responses"}
// 	if filterParam.Search != "" {
// 		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
// 		searchKeys["responses"] = searchRegex
// 	}
// 	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

// 	pipeline := mongo.Pipeline{
// 		{{Key: "$match", Value: filter}},
// 		{{Key: "$facet", Value: bson.D{
// 			{Key: "docs", Value: bson.A{
// 				bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
// 				bson.D{{Key: "$skip", Value: skip}},
// 				bson.D{{Key: "$limit", Value: limit}},
// 				bson.D{{Key: "$addFields", Value: bson.D{
// 					// {Key: "user_id_obj", Value: bson.D{{Key: "$toObjectId", Value: "$user_id"}}},
// 					{Key: "user_id_obj", Value: bson.D{
// 						{Key: "$convert", Value: bson.D{
// 							{Key: "input", Value: "$user_id"},
// 							{Key: "to", Value: "objectId"},
// 							{Key: "onError", Value: nil},
// 							{Key: "onNull", Value: nil},
// 						}},
// 					}},
// 				}}},
// 				bson.D{{Key: "$lookup", Value: bson.D{
// 					{Key: "from", Value: "members"},
// 					{Key: "localField", Value: "user_id_obj"},
// 					{Key: "foreignField", Value: "_id"},
// 					{Key: "as", Value: "user"},
// 				}}},
// 				bson.D{{Key: "$unwind", Value: bson.D{
// 					{Key: "path", Value: "$user"},
// 					{Key: "preserveNullAndEmptyArrays", Value: true},
// 				}}},
// 				bson.D{{Key: "$project", Value: bson.D{
// 					{Key: "_id", Value: 1},
// 					{Key: "responses", Value: 1},
// 					{Key: "created_at", Value: 1},
// 					{Key: "updated_at", Value: 1},
// 					{Key: "rate", Value: "$responses.user_experience.answer"},
// 					{Key: "user", Value: bson.D{
// 						{Key: "_id", Value: bson.M{"$toString": "$user._id"}},
// 						{Key: "user_code", Value: "$user.user_code"},
// 						{Key: "full_name", Value: "$user.full_name"},
// 						{Key: "phone_number", Value: "$user.phone_number"},
// 					}},
// 				}}},
// 			}},
// 			{Key: "totalCount", Value: bson.A{
// 				bson.D{{Key: "$count", Value: "count"}},
// 			}},
// 			{Key: "averageRatings", Value: bson.A{
// 				bson.D{{Key: "$project", Value: bson.D{
// 					{Key: "ratings", Value: bson.D{
// 						{Key: "$filter", Value: bson.D{
// 							{Key: "input", Value: bson.D{{Key: "$objectToArray", Value: "$responses"}}},
// 							{Key: "as", Value: "item"},
// 							{Key: "cond", Value: bson.D{
// 								{Key: "$and", Value: bson.A{
// 									bson.D{{Key: "$eq", Value: bson.A{"$$item.v.type", "rating"}}},
// 									bson.D{{Key: "$ne", Value: bson.A{"$$item.v.answer", nil}}},
// 								}},
// 							}},
// 						}},
// 					}},
// 				}}},
// 				bson.D{{Key: "$unwind", Value: "$ratings"}},
// 				bson.D{{Key: "$group", Value: bson.D{
// 					{Key: "_id", Value: "$ratings.k"},
// 					{Key: "average", Value: bson.D{{Key: "$avg", Value: bson.D{{Key: "$toDouble", Value: "$ratings.v.answer"}}}}},
// 				}}},
// 			}},
// 		}}},
// 	}

// 	f.logger.Infof("[FindAllWithPagination] fetching feedbacks with pagination")
// 	cursor, err := f.feedbackCollection.Aggregate(ctx, pipeline)
// 	if err != nil {
// 		f.logger.Errorf("[FindAllWithPagination] failed to aggregate feedbacks: %v", err)
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	if !cursor.Next(ctx) {
// 		f.logger.Infof("[FindAllWithPagination] no feedbacks found")
// 		return nil, nil
// 	}

// 	var result struct {
// 		Docs       []*feedback.FeedbackResponse `bson:"docs"`
// 		TotalCount []struct {
// 			Count int64 `bson:"count"`
// 		} `bson:"totalCount"`
// 		AverageRatings []struct {
// 			ID      string  `bson:"_id"`
// 			Average float64 `bson:"average"`
// 		} `bson:"averageRatings"`
// 	}

// 	if err := cursor.Decode(&result); err != nil {
// 		f.logger.Errorf("[FindAllWithPagination] failed to decode aggregated result: %v", err)
// 		return nil, err
// 	}

// 	var total int64
// 	if len(result.TotalCount) > 0 {
// 		total = result.TotalCount[0].Count
// 	}

// 	averages := make(map[string]float64)
// 	for _, avg := range result.AverageRatings {
// 		averages[avg.ID] = math.Round(avg.Average*100) / 100
// 	}

// 	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
// 	f.logger.Infof("[FindAllWithPagination] retrieved %d feedbacks", len(result.Docs))

//		return &types.PaginatedResponseForFeedback[[]*feedback.FeedbackResponse]{
//			Data:           result.Docs,
//			Meta:           meta,
//			AverageRatings: averages,
//		}, nil
//	}

func (f *FeedbackStorage) FindAllCustomerFeedbacks(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]local_model.CustomerFeedback], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"search"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"customer_name": searchRegex},
			{"email": searchRegex},
			{"phone_number": searchRegex},
			{"account_number": searchRegex},
			{"message": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	f.logger.Infof("[FeedbackStorage][FindAllCustomerFeedbacks] fetching customer feedbacks with pagination")
	data, err := f.customerFeedbackDal.FindAllWithPaginationE(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		f.logger.Errorf("[FeedbackStorage][FindAllCustomerFeedbacks] failed to fetch customer feedbacks: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := f.customerFeedbackDal.TotalCount(ctx, filter)
	if err != nil {
		f.logger.Errorf("[FeedbackStorage][FindAllCustomerFeedbacks] failed to count customer feedbacks: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	f.logger.Infof("[FeedbackStorage][FindAllCustomerFeedbacks] retrieved %d customer feedbacks", len(data))

	return &types.PaginatedResponse[[]local_model.CustomerFeedback]{
		Data: data,
		Meta: meta,
	}, nil
}

func (f *FeedbackStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]local_model.Feedback], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"search"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"customer_name": searchRegex},
			{"email": searchRegex},
			{"phone_number": searchRegex},
			{"account_number": searchRegex},
			{"message": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	f.logger.Infof("[FindAllWithPagination] fetching feedbacks with pagination")
	data, err := f.dal.FindAllWithPaginationE(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		f.logger.Errorf("[FeedbackStorage][FindAllWithPagination] failed to fetch feedbacks: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := f.dal.TotalCount(ctx, filter)
	if err != nil {
		f.logger.Errorf("[FeedbackStorage][FindAllWithPagination] failed to count feedbacks: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	f.logger.Infof("[FeedbackStorage][FindAllWithPagination] retrieved %d feedbacks", len(data))

	return &types.PaginatedResponse[[]local_model.Feedback]{
		Data: data,
		Meta: meta,
	}, nil
}

// FindFeedbackByID implements [storage.FeedbackRepository].
func (f *FeedbackStorage) FindFeedbackByID(ctx context.Context, id string) (*local_model.Feedback, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		f.logger.Errorf("[FeedbackStorage][FindFeedbackByID] invalid id: %s", id)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	f.logger.Infof("[FeedbackStorage][FindFeedbackByID] fetching feedback by id: %s", id)

	result, err := f.dal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}

	f.logger.Infof("[FeedbackStorage][FindFeedbackByID] feedback retrieved successfully for id: %s", id)
	return result, nil
}

func (f *FeedbackStorage) FindAllSurveyFeedbacks(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]local_model.SurveyFeedback], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"search"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"customer_name": searchRegex},
			{"email": searchRegex},
			{"phone_number": searchRegex},
			{"account_number": searchRegex},
			{"message": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	f.logger.Infof("[FeedbackStorage][FindAllSurveyFeedbacks] fetching survey feedbacks with pagination")
	data, err := f.surveyFeedbackDal.FindAllWithPaginationE(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		f.logger.Errorf("[FeedbackStorage][FindAllSurveyFeedbacks] failed to fetch survey feedbacks: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := f.surveyFeedbackDal.TotalCount(ctx, filter)
	if err != nil {
		f.logger.Errorf("[FeedbackStorage][FindAllSurveyFeedbacks] failed to count survey feedbacks: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	f.logger.Infof("[FeedbackStorage][FindAllSurveyFeedbacks] retrieved %d survey feedbacks", len(data))

	return &types.PaginatedResponse[[]local_model.SurveyFeedback]{
		Data: data,
		Meta: meta,
	}, nil
}

func (f *FeedbackStorage) FindSurveyFeedbackByID(ctx context.Context, id string) (*local_model.SurveyFeedback, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		f.logger.Errorf("[FeedbackStorage][FindSurveyFeedbackByID] invalid id: %s", id)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	f.logger.Infof("[FeedbackStorage][FindSurveyFeedbackByID] fetching survey feedback by id: %s", id)

	result, err := f.surveyFeedbackDal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}

	f.logger.Infof("[FeedbackStorage][FindSurveyFeedbackByID] survey feedback retrieved successfully for id: %s", id)
	return result, nil
}

func (f *FeedbackStorage) FindCustomerFeedbackByID(ctx context.Context, id string) (*local_model.CustomerFeedback, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		f.logger.Errorf("[FeedbackStorage][FindCustomerFeedbackByID] invalid id: %s", id)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	f.logger.Infof("[FeedbackStorage][FindCustomerFeedbackByID] fetching customer feedback by id: %s", id)

	result, err := f.customerFeedbackDal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}

	f.logger.Infof("[FeedbackStorage][FindCustomerFeedbackByID] customer feedback retrieved successfully for id: %s", id)
	return result, nil
}
