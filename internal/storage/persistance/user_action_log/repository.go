package user_action_log

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type userActionLogRepository struct {
	dal        dal.MongoDal[imodel.UserActionLog, imodel.UserActionLog]
	collection *mongo.Collection
	logger     utils.Logger
}

func NewUserActionLogRepository(client *mongo.Client, cfg *config.VaultConfig, dbName, collectionName string, logger utils.Logger) storage.UserActionLogRepository {
	return &userActionLogRepository{
		dal:        dal.NewMongoDal[imodel.UserActionLog, imodel.UserActionLog](client, cfg, dbName, collectionName),
		collection: client.Database(dbName).Collection(collectionName),
		logger:     logger,
	}
}

func (r *userActionLogRepository) Save(ctx context.Context, log *imodel.UserActionLog) error {
	r.logger.Infof("[UserActionLog][Save] saving action log for action_code=%s user=%s role=%s", log.ActionCode, log.Username, log.UserActionResponsibilities)
	_, err := r.dal.InsertOne(ctx, *log)
	if err != nil {
		r.logger.Errorf("[UserActionLog][Save] failed to save action log: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (r *userActionLogRepository) GetActionCodesBySearch(ctx context.Context, search string) ([]string, error) {
	r.logger.Infof("[UserActionLog][GetActionCodesBySearch] search=%s", search)

	filter := bson.M{
		"$or": []bson.M{
			{"action_code": bson.M{"$regex": search, "$options": "i"}},
			{"username": bson.M{"$regex": search, "$options": "i"}},
			{"given_action_status": bson.M{"$regex": search, "$options": "i"}},
			{"request_action": bson.M{"$regex": search, "$options": "i"}},
			{"action_taken_service_name": bson.M{"$regex": search, "$options": "i"}},
			{"checker_level": bson.M{"$regex": search, "$options": "i"}},
			{"auditor_level": bson.M{"$regex": search, "$options": "i"}},
			{"user_phone": bson.M{"$regex": search, "$options": "i"}},
			{"user_action_responsibilities": bson.M{"$regex": search, "$options": "i"}},
			{"action_taken_service_unique_id": bson.M{"$regex": search, "$options": "i"}},
			{"user_role_code": bson.M{"$regex": search, "$options": "i"}},
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Errorf("[UserActionLog][GetActionCodesBySearch] failed to find action codes: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	var actionCodes []string
	for cursor.Next(ctx) {
		var log imodel.UserActionLog
		if err := cursor.Decode(&log); err != nil {
			r.logger.Errorf("[UserActionLog][GetActionCodesBySearch] failed to decode log: %v", err)
			continue
		}
		actionCodes = append(actionCodes, log.ActionCode)
	}

	if err := cursor.Err(); err != nil {
		r.logger.Errorf("[UserActionLog][GetActionCodesBySearch] cursor error: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return actionCodes, nil
}

func (r *userActionLogRepository) GetActionCodesByFilter(ctx context.Context, filter map[string]interface{}) ([]string, error) {
	r.logger.Infof("[UserActionLog][GetActionCodesByFilter] filter=%v", filter)

	filter["is_deleted"] = false

	pipeline := []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": "$action_code"}},
		{"$project": bson.M{"_id": 0, "action_code": "$_id"}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("[UserActionLog][GetActionCodesByFilter] failed to find action codes: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	var actionCodes []string
	for cursor.Next(ctx) {
		var log imodel.UserActionLog
		if err := cursor.Decode(&log); err != nil {
			r.logger.Errorf("[UserActionLog][GetActionCodesBySearch] failed to decode log: %v", err)
			continue
		}
		actionCodes = append(actionCodes, log.ActionCode)
	}

	if err := cursor.Err(); err != nil {
		r.logger.Errorf("[UserActionLog][GetActionCodesBySearch] cursor error: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return actionCodes, nil
}

func (r *userActionLogRepository) GetActionCodesByUser(ctx context.Context, userID string, responsibility imodel.UserActionResponsibility) ([]string, error) {
	r.logger.Infof("[UserActionLog][GetActionCodesByUser] user_id=%s responsibility=%s", userID, responsibility)

	objectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		r.logger.Errorf("[UserActionLog][GetActionCodesByUser] invalid user_id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{
		"user_id":                      objectID,
		"user_action_responsibilities": string(responsibility),
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.M{
			"_id": "$action_code",
		}}},
	}

	cur, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("[UserActionLog][GetActionCodesByUser] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	var results []struct {
		ActionCode string `bson:"_id"`
	}
	if err := cur.All(ctx, &results); err != nil {
		r.logger.Errorf("[UserActionLog][GetActionCodesByUser] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	codes := make([]string, 0, len(results))
	for _, r := range results {
		codes = append(codes, r.ActionCode)
	}
	return codes, nil
}

func (r *userActionLogRepository) GetActionCodesByUserAndAuditorStatus(ctx context.Context, userID string, responsibility imodel.UserActionResponsibility, status string) ([]string, error) {
	r.logger.Infof("[UserActionLog][GetActionCodesByUserAndStatus] user_id=%s responsibility=%s status=%s", userID, responsibility, status)

	objectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		r.logger.Errorf("[UserActionLog][GetActionCodesByUserAndStatus] invalid user_id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{
		"user_id":                      objectID,
		"user_action_responsibilities": string(responsibility),
		"given_auditor_status":         status,
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.M{
			"_id": "$action_code",
		}}},
	}

	cur, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("[UserActionLog][GetActionCodesByUser] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	var results []struct {
		ActionCode string `bson:"_id"`
	}
	if err := cur.All(ctx, &results); err != nil {
		r.logger.Errorf("[UserActionLog][GetActionCodesByUser] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	codes := make([]string, 0, len(results))
	for _, r := range results {
		codes = append(codes, r.ActionCode)
	}
	return codes, nil
}

func (r *userActionLogRepository) GetLogsByActionCode(ctx context.Context, actionCode string) ([]imodel.UserActionLog, error) {
	r.logger.Infof("[UserActionLog][GetLogsByActionCode] action_code=%s", actionCode)

	filter := bson.M{"action_code": actionCode}
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Errorf("[UserActionLog][GetLogsByActionCode] failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	var logs []imodel.UserActionLog
	if err := cur.All(ctx, &logs); err != nil {
		r.logger.Errorf("[UserActionLog][GetLogsByActionCode] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return logs, nil
}
