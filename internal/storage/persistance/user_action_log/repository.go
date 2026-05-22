package user_action_log

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"fmt"

	local_util "cbe-super-app-cps-action/pkgs/utils"

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

func (r *userActionLogRepository) Save(ctx context.Context, entry *imodel.UserActionLog) error {
	reqLog := local_util.LoggerFromCtx(ctx, r.logger)

	reqLog.Infof("[UserActionLog][Save] saving action log for action_code=%s user=%s role=%s", entry.ActionCode, entry.Username, entry.UserActionResponsibilities)
	_, err := r.dal.InsertOne(ctx, *entry)
	if err != nil {
		reqLog.Errorf("[UserActionLog][Save] failed to save action log: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (r *userActionLogRepository) GetActionCodesByActionLogFilter(ctx context.Context, filter imodel.UserActionLogActionCodeFilter) ([]string, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[UserActionLog][GetActionCodesByActionLogFilter] filter=%+v", filter)

	pipeline, err := buildActionCodeFilterPipeline(filter)
	if err != nil {
		log.Errorf("[UserActionLog][GetActionCodesByActionLogFilter] invalid filter: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	cur, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Errorf("[UserActionLog][GetActionCodesByActionLogFilter] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	var results []struct {
		ActionCode string `bson:"action_code"`
	}
	if err := cur.All(ctx, &results); err != nil {
		log.Errorf("[UserActionLog][GetActionCodesByActionLogFilter] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	actionCodes := make([]string, 0, len(results))
	for _, result := range results {
		actionCodes = append(actionCodes, result.ActionCode)
	}

	return actionCodes, nil
}

const fieldResponsibility = "$user_action_responsibilities"

func buildActionCodeFilterPipeline(filter imodel.UserActionLogActionCodeFilter) (mongo.Pipeline, error) {
	groupStage := bson.D{{Key: "_id", Value: "$action_code"}}
	matchStage := bson.D{}

	if len(filter.ActionStatuses) > 0 {
		groupStage = append(groupStage, bson.E{
			Key:   "action_status_match_count",
			Value: sumWhen(bson.D{{Key: "$in", Value: bson.A{"$given_action_status", filter.ActionStatuses}}}),
		})
		matchStage = append(matchStage, bson.E{Key: "action_status_match_count", Value: bson.M{"$gt": 0}})
	}

	if len(filter.AuditorStatuses) > 0 {
		groupStage = append(groupStage, bson.E{
			Key:   "auditor_status_match_count",
			Value: sumWhen(bson.D{{Key: "$in", Value: bson.A{"$given_auditor_status", filter.AuditorStatuses}}}),
		})
		matchStage = append(matchStage, bson.E{Key: "auditor_status_match_count", Value: bson.M{"$gt": 0}})
	}

	privateUserIDs, err := toObjectIDs(filter.PrivateUserIDs)
	if err != nil {
		return nil, fmt.Errorf("private users: %w", err)
	}
	if len(privateUserIDs) > 0 {
		groupStage = append(groupStage, bson.E{
			Key:   "private_user_match_count",
			Value: sumWhen(bson.D{{Key: "$in", Value: bson.A{"$user_id", privateUserIDs}}}),
		})
		matchStage = append(matchStage, bson.E{Key: "private_user_match_count", Value: bson.M{"$gt": 0}})
	}

	if len(filter.Levels) > 0 {
		groupStage = append(groupStage, bson.E{
			Key: "level_match_count",
			Value: sumWhen(bson.D{{Key: "$or", Value: bson.A{
				bson.D{{Key: "$in", Value: bson.A{"$checker_level", filter.Levels}}},
				bson.D{{Key: "$in", Value: bson.A{"$auditor_level", filter.Levels}}},
			}}}),
		})
		matchStage = append(matchStage, bson.E{Key: "level_match_count", Value: bson.M{"$gt": 0}})
	}

	if len(filter.Services) > 0 {
		groupStage = append(groupStage, bson.E{
			Key:   "service_match_count",
			Value: sumWhen(bson.D{{Key: "$in", Value: bson.A{"$action_taken_service_name", filter.Services}}}),
		})
		matchStage = append(matchStage, bson.E{Key: "service_match_count", Value: bson.M{"$gt": 0}})
	}

	checkerUserIDs, err := toObjectIDs(filter.CheckerUserIDs)
	if err != nil {
		return nil, fmt.Errorf("checker users: %w", err)
	}
	if len(checkerUserIDs) > 0 {
		groupStage = append(groupStage, bson.E{
			Key: "checker_match_count",
			Value: sumWhen(bson.D{{Key: "$and", Value: bson.A{
				bson.D{{Key: "$eq", Value: bson.A{fieldResponsibility, string(imodel.CHECKER)}}},
				bson.D{{Key: "$in", Value: bson.A{"$user_id", checkerUserIDs}}},
			}}}),
		})
		matchStage = append(matchStage, bson.E{Key: "checker_match_count", Value: bson.M{"$gt": 0}})
	}

	auditorUserIDs, err := toObjectIDs(filter.AuditorUserIDs)
	if err != nil {
		return nil, fmt.Errorf("auditor users: %w", err)
	}
	if len(auditorUserIDs) > 0 {
		groupStage = append(groupStage, bson.E{
			Key: "auditor_match_count",
			Value: sumWhen(bson.D{{Key: "$and", Value: bson.A{
				bson.D{{Key: "$eq", Value: bson.A{fieldResponsibility, string(imodel.AUDITOR)}}},
				bson.D{{Key: "$in", Value: bson.A{"$user_id", auditorUserIDs}}},
			}}}),
		})
		matchStage = append(matchStage, bson.E{Key: "auditor_match_count", Value: bson.M{"$gt": 0}})
	}

	makerUserIDs, err := toObjectIDs(filter.MakerUserIDs)
	if err != nil {
		return nil, fmt.Errorf("maker users: %w", err)
	}
	if len(makerUserIDs) > 0 {
		groupStage = append(groupStage, bson.E{
			Key: "maker_match_count",
			Value: sumWhen(bson.D{{Key: "$and", Value: bson.A{
				bson.D{{Key: "$eq", Value: bson.A{fieldResponsibility, string(imodel.MAKER)}}},
				bson.D{{Key: "$in", Value: bson.A{"$user_id", makerUserIDs}}},
			}}}),
		})
		matchStage = append(matchStage, bson.E{Key: "maker_match_count", Value: bson.M{"$gt": 0}})
	}

	if len(filter.Responsibilities) > 0 {
		groupStage = append(groupStage, bson.E{
			Key:   "responsibility_match_count",
			Value: sumWhen(bson.D{{Key: "$in", Value: bson.A{fieldResponsibility, filter.Responsibilities}}}),
		})
		matchStage = append(matchStage, bson.E{Key: "responsibility_match_count", Value: bson.M{"$gt": 0}})
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"is_deleted": false}}},
		{{Key: "$group", Value: groupStage}},
	}

	if len(matchStage) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: matchStage}})
	}

	pipeline = append(pipeline, bson.D{{Key: "$project", Value: bson.M{
		"_id":         0,
		"action_code": "$_id",
	}}})

	return pipeline, nil
}

func sumWhen(condition interface{}) bson.D {
	return bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{condition, 1, 0}}}}}
}

func toObjectIDs(userIDs []string) ([]bson.ObjectID, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}

	objectIDs := make([]bson.ObjectID, 0, len(userIDs))
	for _, userID := range userIDs {
		if userID == "" {
			continue
		}

		objectID, err := bson.ObjectIDFromHex(userID)
		if err != nil {
			return nil, fmt.Errorf("invalid object id %q: %w", userID, err)
		}
		objectIDs = append(objectIDs, objectID)
	}

	return objectIDs, nil
}

func (r *userActionLogRepository) GetActionCodesBySearch(ctx context.Context, search string) ([]string, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[UserActionLog][GetActionCodesBySearch] search=%s", search)

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
		log.Errorf("[UserActionLog][GetActionCodesBySearch] failed to find action codes: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	var actionCodes []string
	for cursor.Next(ctx) {
		var entry imodel.UserActionLog
		if err := cursor.Decode(&entry); err != nil {
			log.Errorf("[UserActionLog][GetActionCodesBySearch] failed to decode log: %v", err)
			continue
		}
		actionCodes = append(actionCodes, entry.ActionCode)
	}

	if err := cursor.Err(); err != nil {
		log.Errorf("[UserActionLog][GetActionCodesBySearch] cursor error: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return actionCodes, nil
}

func (r *userActionLogRepository) GetActionCodesByFilter(ctx context.Context, filter map[string]interface{}) ([]string, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[UserActionLog][GetActionCodesByFilter] filter=%v", filter)

	filter["is_deleted"] = false

	pipeline := []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": "$action_code"}},
		{"$project": bson.M{"_id": 0, "action_code": "$_id"}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Errorf("[UserActionLog][GetActionCodesByFilter] failed to find action codes: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	var actionCodes []string
	for cursor.Next(ctx) {
		var entry imodel.UserActionLog
		if err := cursor.Decode(&entry); err != nil {
			log.Errorf("[UserActionLog][GetActionCodesBySearch] failed to decode log: %v", err)
			continue
		}
		actionCodes = append(actionCodes, entry.ActionCode)
	}

	if err := cursor.Err(); err != nil {
		log.Errorf("[UserActionLog][GetActionCodesBySearch] cursor error: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return actionCodes, nil
}

func (r *userActionLogRepository) GetActionCodesByUser(ctx context.Context, userID string, responsibility imodel.UserActionResponsibility) ([]string, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[UserActionLog][GetActionCodesByUser] user_id=%s responsibility=%s", userID, responsibility)

	objectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		log.Errorf("[UserActionLog][GetActionCodesByUser] invalid user_id: %v", err)
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
		log.Errorf("[UserActionLog][GetActionCodesByUser] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	var results []struct {
		ActionCode string `bson:"_id"`
	}
	if err := cur.All(ctx, &results); err != nil {
		log.Errorf("[UserActionLog][GetActionCodesByUser] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	codes := make([]string, 0, len(results))
	for _, r := range results {
		codes = append(codes, r.ActionCode)
	}
	return codes, nil
}

func (r *userActionLogRepository) GetActionCodesByUserAndAuditorStatus(ctx context.Context, userID string, responsibility imodel.UserActionResponsibility, status string) ([]string, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[UserActionLog][GetActionCodesByUserAndStatus] user_id=%s responsibility=%s status=%s", userID, responsibility, status)

	objectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		log.Errorf("[UserActionLog][GetActionCodesByUserAndStatus] invalid user_id: %v", err)
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
		log.Errorf("[UserActionLog][GetActionCodesByUser] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	var results []struct {
		ActionCode string `bson:"_id"`
	}
	if err := cur.All(ctx, &results); err != nil {
		log.Errorf("[UserActionLog][GetActionCodesByUser] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	codes := make([]string, 0, len(results))
	for _, r := range results {
		codes = append(codes, r.ActionCode)
	}
	return codes, nil
}

func (r *userActionLogRepository) GetLogsByActionCode(ctx context.Context, actionCode string) ([]string, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[UserActionLog][GetLogsByActionCode] action_code=%s", actionCode)

	filter := bson.M{"action_code": actionCode}
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.Errorf("[UserActionLog][GetLogsByActionCode] failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	var logs []imodel.UserActionLog
	if err := cur.All(ctx, &logs); err != nil {
		log.Errorf("[UserActionLog][GetLogsByActionCode] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	var actionCodeList []string
	for _, logEntry := range logs {
		actionCodeList = append(actionCodeList, logEntry.ActionCode)
	}
	return actionCodeList, nil
}

func (r *userActionLogRepository) ApproveUserActionsByActionCode(ctx context.Context, actionCode string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[UserActionLog][ApproveUserActionsByActionCode] action_code=%s", actionCode)

	filter := bson.M{"action_code": actionCode}
	update := bson.M{"$set": bson.M{"given_auditor_status": "APPROVED"}}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Errorf("[UserActionLog][ApproveUserActionsByActionCode] update failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (r *userActionLogRepository) RejectUserActionsByActionCode(ctx context.Context, actionCode string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[UserActionLog][RejectUserActionsByActionCode] action_code=%s", actionCode)

	filter := bson.M{"action_code": actionCode}
	update := bson.M{"$set": bson.M{"given_auditor_status": "REJECTED"}}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Errorf("[UserActionLog][RejectUserActionsByActionCode] update failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (r *userActionLogRepository) GetLogsByResponsibility(ctx context.Context, responsibility imodel.UserActionResponsibility) ([]string, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[UserActionLog][GetLogsByResponsibility] responsibility=%s", responsibility)

	filter := bson.M{"user_action_responsibilities": string(responsibility)}
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.Errorf("[UserActionLog][GetLogsByResponsibility] failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	var logs []imodel.UserActionLog
	if err := cur.All(ctx, &logs); err != nil {
		log.Errorf("[UserActionLog][GetLogsByResponsibility] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	var actionCodeList []string
	for _, logEntry := range logs {
		actionCodeList = append(actionCodeList, logEntry.ActionCode)
	}
	return actionCodeList, nil
}

func (r *userActionLogRepository) CancelUserActionsByActionCode(ctx context.Context, actionCode string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[UserActionLog][CancelUserActionsByActionCode] action_code=%s", actionCode)

	filter := bson.M{"action_code": actionCode}
	update := bson.M{"$set": bson.M{"given_auditor_status": "CANCELLED"}}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Errorf("[UserActionLog][CancelUserActionsByActionCode] update failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (r *userActionLogRepository) GetLogsByUserID(ctx context.Context, userID string) ([]string, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[UserActionLog][GetLogsByUserID] user_id=%s", userID)

	objectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		log.Errorf("[UserActionLog][GetLogsByUserID] invalid user_id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"user_id": objectID}
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.Errorf("[UserActionLog][GetLogsByUserID] failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	var logs []imodel.UserActionLog
	if err := cur.All(ctx, &logs); err != nil {
		log.Errorf("[UserActionLog][GetLogsByUserID] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	var actionCodeList []string
	for _, logEntry := range logs {
		actionCodeList = append(actionCodeList, logEntry.ActionCode)
	}
	return actionCodeList, nil
}

func (r *userActionLogRepository) GetLogsByUserIDAndResponsibility(ctx context.Context, userID string, responsibility imodel.UserActionResponsibility) ([]string, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[UserActionLog][GetLogsByUserIDAndResponsibility] user_id=%s responsibility=%s", userID, responsibility)

	objectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		log.Errorf("[UserActionLog][GetLogsByUserIDAndResponsibility] invalid user_id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{
		"user_id":                      objectID,
		"user_action_responsibilities": string(responsibility),
	}
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.Errorf("[UserActionLog][GetLogsByUserIDAndResponsibility] failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	var logs []imodel.UserActionLog
	if err := cur.All(ctx, &logs); err != nil {
		log.Errorf("[UserActionLog][GetLogsByUserIDAndResponsibility] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	var actionCodeList []string

	for _, logEntry := range logs {
		actionCodeList = append(actionCodeList, logEntry.ActionCode)
	}
	return actionCodeList, nil
}

func (r *userActionLogRepository) AuditorMarkLogsByActionCode(ctx context.Context, actionCode string, auditorStatus string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[UserActionLog][AuditorMarkLogsByActionCode] action_code=%s auditor_status=%s", actionCode, auditorStatus)

	filter := bson.M{"action_code": actionCode}
	update := bson.M{"$set": bson.M{"given_auditor_status": auditorStatus}}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Errorf("[UserActionLog][AuditorMarkLogsByActionCode] update failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}
