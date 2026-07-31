package user_action_log

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

// Upsert inserts a new log entry for (action_code, username) if one does not yet exist,
// or updates the status fields on the existing entry. Immutable fields (IDs, responsibility,
// level, created_at) are written only on first insert via $setOnInsert.
func (r *userActionLogRepository) Upsert(ctx context.Context, entry *imodel.UserActionLog) error {
	reqLog := local_util.LoggerFromCtx(ctx, r.logger)

	reqLog.Infof("[UserActionLog][Upsert] action_code=%s user=%s responsibility=%s", entry.ActionCode, entry.Username, entry.UserActionResponsibilities)

	filter := bson.M{
		"action_code": entry.ActionCode,
		"username":    entry.Username,
		"is_deleted":  false,
	}
	update := bson.D{
		{Key: "$setOnInsert", Value: bson.M{
			"_id":                            entry.ID,
			"action_id":                      entry.ActionID,
			"action_code":                    entry.ActionCode,
			"request_action":                 entry.RequestAction,
			"action_taken_service_name":      entry.ActionTakenServiceName,
			"action_taken_service_unique_id": entry.ActionTakenServiceUniqueID,
			"user_id":                        entry.UserID,
			"username":                       entry.Username,
			"user_phone":                     entry.UserPhone,
			"user_role_code":                 entry.UserRoleCode,
			"user_action_responsibilities":   entry.UserActionResponsibilities,
			"action_type":                    entry.ActionType,
			"checker_level":                  entry.CheckerLevel,
			"auditor_level":                  entry.AuditorLevel,
			"is_deleted":                     false,
			"created_at":                     entry.CreatedAt,
			// action_auditor_status starts empty; bulk methods advance it as the auditor workflow progresses.
			"action_auditor_status": "",
		}},
		{Key: "$set", Value: bson.M{
			"given_action_status":  entry.GivenActionStatus,
			"given_auditor_status": entry.GivenAuditorStatus,
			"last_modified_at":     time.Now(),
		}},
	}

	opts := options.UpdateOne().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		reqLog.Errorf("[UserActionLog][Upsert] failed for action_code=%s user=%s: %v", entry.ActionCode, entry.Username, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (r *userActionLogRepository) GetActionCodesByActionLogFilter(ctx context.Context, filter imodel.UserActionLogActionCodeFilter) ([]string, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[UserActionLog][GetActionCodesByActionLogFilter] filter=%+v", filter)

	pipeline, err := r.buildActionCodeFilterPipeline(filter)
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

func (r *userActionLogRepository) buildActionCodeFilterPipeline(filter imodel.UserActionLogActionCodeFilter) (mongo.Pipeline, error) {
	groupStage := bson.D{{Key: "_id", Value: "$action_code"}}
	matchStage := bson.D{}
	log := local_util.LoggerFromCtx(context.Background(), r.logger) // no need to pass real ctx since we won't log after this point

	if len(filter.ActionStatuses) > 0 {
		log.Infof("[CPSAction][buildActionCodeFilterPipeline] applying action status filter: %v", filter.ActionStatuses)
		groupStage = append(groupStage, bson.E{
			Key:   "action_status_match_count",
			Value: sumWhen(bson.D{{Key: "$in", Value: bson.A{"$given_action_status", filter.ActionStatuses}}}),
		})
		matchStage = append(matchStage, bson.E{Key: "action_status_match_count", Value: bson.M{"$gt": 0}})
	}

	if len(filter.AuditorStatuses) > 0 {
		log.Infof("[CPSAction][buildActionCodeFilterPipeline] applying auditor status filter: %v", filter.AuditorStatuses)

		groupStage = append(groupStage, bson.E{
			Key:   "auditor_status_match_count",
			Value: sumWhen(bson.D{{Key: "$in", Value: bson.A{"$given_auditor_status", filter.AuditorStatuses}}}),
		})
		matchStage = append(matchStage, bson.E{Key: "auditor_status_match_count", Value: bson.M{"$gt": 0}})
	}

	if len(filter.ActionAuditorStatuses) > 0 {
		log.Infof("[CPSAction][buildActionCodeFilterPipeline] applying action_auditor_status filter: %v", filter.ActionAuditorStatuses)
		groupStage = append(groupStage, bson.E{
			Key:   "action_auditor_status_match_count",
			Value: sumWhen(bson.D{{Key: "$in", Value: bson.A{"$action_auditor_status", filter.ActionAuditorStatuses}}}),
		})
		matchStage = append(matchStage, bson.E{Key: "action_auditor_status_match_count", Value: bson.M{"$gt": 0}})
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
		log.Infof("[CPSAction][buildActionCodeFilterPipeline] applying level filter: %v (responsibilities=%v)", filter.Levels, filter.Responsibilities)
		escapedLevels := make([]string, len(filter.Levels))
		for i, lvl := range filter.Levels {
			escapedLevels[i] = regexp.QuoteMeta(lvl)
		}
		levelPattern := fmt.Sprintf("(?i)(?:%s)", strings.Join(escapedLevels, "|"))

		hasCheckerResponsibility := slices.Contains(filter.Responsibilities, string(imodel.CHECKER))
		hasAuditorResponsibility := slices.Contains(filter.Responsibilities, string(imodel.AUDITOR))

		var levelConditions bson.A
		if hasCheckerResponsibility && !hasAuditorResponsibility {
			levelConditions = bson.A{
				bson.D{{Key: "$regexMatch", Value: bson.M{"input": "$checker_level", "regex": levelPattern, "options": "i"}}},
				bson.D{{Key: "$eq", Value: bson.A{"$checker_level", ""}}}, // Include maker-only actions
			}
		} else if hasAuditorResponsibility && !hasCheckerResponsibility {
			levelConditions = bson.A{
				bson.D{{Key: "$regexMatch", Value: bson.M{"input": "$auditor_level", "regex": levelPattern, "options": "i"}}},
				bson.D{{Key: "$eq", Value: bson.A{"$auditor_level", ""}}}, // Include maker-only actions
			}
		} else {
			levelConditions = bson.A{
				bson.D{{Key: "$regexMatch", Value: bson.M{"input": "$checker_level", "regex": levelPattern, "options": "i"}}},
				bson.D{{Key: "$regexMatch", Value: bson.M{"input": "$auditor_level", "regex": levelPattern, "options": "i"}}},
				bson.D{{Key: "$eq", Value: bson.A{"$checker_level", ""}}}, // Include maker-only actions
				bson.D{{Key: "$eq", Value: bson.A{"$auditor_level", ""}}}, // Include maker-only actions
			}
		}

		groupStage = append(groupStage, bson.E{
			Key:   "level_match_count",
			Value: sumWhen(bson.D{{Key: "$or", Value: levelConditions}}),
		})
		matchStage = append(matchStage, bson.E{Key: "level_match_count", Value: bson.M{"$gt": 0}})
	}

	if len(filter.CheckerLevels) > 0 {
		log.Infof("[CPSAction][buildActionCodeFilterPipeline] applying checker level filter: %v", filter.CheckerLevels)
		escapedLevels := make([]string, len(filter.CheckerLevels))
		for i, lvl := range filter.CheckerLevels {
			escapedLevels[i] = regexp.QuoteMeta(lvl)
		}
		levelPattern := fmt.Sprintf("(?i)(?:%s)", strings.Join(escapedLevels, "|"))

		groupStage = append(groupStage, bson.E{
			Key: "checker_level_match_count",
			Value: sumWhen(bson.D{{Key: "$or", Value: bson.A{
				bson.D{{Key: "$regexMatch", Value: bson.M{"input": "$checker_level", "regex": levelPattern, "options": "i"}}},
				bson.D{{Key: "$eq", Value: bson.A{"$checker_level", ""}}}, // Include maker-only actions
			}}}),
		})
		matchStage = append(matchStage, bson.E{Key: "checker_level_match_count", Value: bson.M{"$gt": 0}})
	}

	// Separate auditor level filter - matches only auditor_level field
	// Also includes maker-only actions where auditor_level is empty
	if len(filter.AuditorLevels) > 0 {
		log.Infof("[CPSAction][buildActionCodeFilterPipeline] applying auditor level filter: %v", filter.AuditorLevels)
		escapedLevels := make([]string, len(filter.AuditorLevels))
		for i, lvl := range filter.AuditorLevels {
			escapedLevels[i] = regexp.QuoteMeta(lvl)
		}
		levelPattern := fmt.Sprintf("(?i)(?:%s)", strings.Join(escapedLevels, "|"))

		groupStage = append(groupStage, bson.E{
			Key: "auditor_level_match_count",
			Value: sumWhen(bson.D{{Key: "$or", Value: bson.A{
				bson.D{{Key: "$regexMatch", Value: bson.M{"input": "$auditor_level", "regex": levelPattern, "options": "i"}}},
				bson.D{{Key: "$eq", Value: bson.A{"$auditor_level", ""}}}, // Include maker-only actions
			}}}),
		})
		matchStage = append(matchStage, bson.E{Key: "auditor_level_match_count", Value: bson.M{"$gt": 0}})
	}

	// Checker level+status pairs - each pair requires matching checker_level AND given_action_status
	// Note: given_action_status holds the action status (PENDING, APPROVED, REJECTED) for all user types
	for i, pair := range filter.CheckerLevelStatuses {
		fieldName := fmt.Sprintf("checker_level_claim_%d_count", i)
		groupStage = append(groupStage, bson.E{
			Key: fieldName,
			Value: sumWhen(bson.D{{Key: "$and", Value: bson.A{
				bson.D{{Key: "$regexMatch", Value: bson.M{"input": "$checker_level", "regex": regexp.QuoteMeta(pair.Level), "options": "i"}}},
				bson.D{{Key: "$eq", Value: bson.A{"$given_action_status", pair.Claim}}},
			}}}),
		})
		matchStage = append(matchStage, bson.E{Key: fieldName, Value: bson.M{"$gt": 0}})
	}

	if len(filter.Services) > 0 {
		log.Infof("[CPSAction][buildActionCodeFilterPipeline] applying service filter: %v", filter.Services)
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
		log.Infof("[CPSAction][buildActionCodeFilterPipeline] applying checker user filter: %v", filter.CheckerUserIDs)
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
		log.Infof("[CPSAction][buildActionCodeFilterPipeline] applying auditor user filter: %v", filter.AuditorUserIDs)

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
		log.Infof("[CPSAction][buildActionCodeFilterPipeline] applying maker user filter: %v", filter.MakerUserIDs)

		groupStage = append(groupStage, bson.E{
			Key: "maker_match_count",
			Value: sumWhen(bson.D{{Key: "$and", Value: bson.A{
				bson.D{{Key: "$eq", Value: bson.A{fieldResponsibility, string(imodel.MAKER)}}},
				bson.D{{Key: "$in", Value: bson.A{"$user_id", makerUserIDs}}},
			}}}),
		})
		matchStage = append(matchStage, bson.E{Key: "maker_match_count", Value: bson.M{"$gt": 0}})
	}

	// Maker usernames filter - matches username with MAKER responsibility (case-insensitive, partial search)
	if len(filter.MakerUsernames) > 0 {
		log.Infof("[CPSAction][buildActionCodeFilterPipeline] applying maker username filter: %v", filter.MakerUsernames)
		// Build partial match pattern (e.g., "kid" matches "kidusm", "akid", "skids")
		searchPatterns := make([]string, len(filter.MakerUsernames))
		for i, uname := range filter.MakerUsernames {
			searchPatterns[i] = fmt.Sprintf(".*%s.*", regexp.QuoteMeta(uname))
		}
		usernamePattern := fmt.Sprintf("(?i)(?:%s)", strings.Join(searchPatterns, "|"))

		groupStage = append(groupStage, bson.E{
			Key: "maker_username_match_count",
			Value: sumWhen(bson.D{{Key: "$and", Value: bson.A{
				bson.D{{Key: "$eq", Value: bson.A{fieldResponsibility, string(imodel.MAKER)}}},
				bson.D{{Key: "$regexMatch", Value: bson.M{"input": "$username", "regex": usernamePattern, "options": "i"}}},
			}}}),
		})
		matchStage = append(matchStage, bson.E{Key: "maker_username_match_count", Value: bson.M{"$gt": 0}})
	}

	// Checker usernames filter - matches username with CHECKER responsibility (case-insensitive, partial search)
	if len(filter.CheckerUsernames) > 0 {
		log.Infof("[CPSAction][buildActionCodeFilterPipeline] applying checker username filter: %v", filter.CheckerUsernames)
		searchPatterns := make([]string, len(filter.CheckerUsernames))
		for i, uname := range filter.CheckerUsernames {
			searchPatterns[i] = fmt.Sprintf(".*%s.*", regexp.QuoteMeta(uname))
		}
		usernamePattern := fmt.Sprintf("(?i)(?:%s)", strings.Join(searchPatterns, "|"))

		groupStage = append(groupStage, bson.E{
			Key: "checker_username_match_count",
			Value: sumWhen(bson.D{{Key: "$and", Value: bson.A{
				bson.D{{Key: "$eq", Value: bson.A{fieldResponsibility, string(imodel.CHECKER)}}},
				bson.D{{Key: "$regexMatch", Value: bson.M{"input": "$username", "regex": usernamePattern, "options": "i"}}},
			}}}),
		})
		matchStage = append(matchStage, bson.E{Key: "checker_username_match_count", Value: bson.M{"$gt": 0}})
	}

	// Auditor usernames filter - matches username with AUDITOR responsibility (case-insensitive, partial search)
	if len(filter.AuditorUsernames) > 0 {
		log.Infof("[CPSAction][buildActionCodeFilterPipeline] applying auditor username filter: %v", filter.AuditorUsernames)
		searchPatterns := make([]string, len(filter.AuditorUsernames))
		for i, uname := range filter.AuditorUsernames {
			searchPatterns[i] = fmt.Sprintf(".*%s.*", regexp.QuoteMeta(uname))
		}
		usernamePattern := fmt.Sprintf("(?i)(?:%s)", strings.Join(searchPatterns, "|"))

		groupStage = append(groupStage, bson.E{
			Key: "auditor_username_match_count",
			Value: sumWhen(bson.D{{Key: "$and", Value: bson.A{
				bson.D{{Key: "$eq", Value: bson.A{fieldResponsibility, string(imodel.AUDITOR)}}},
				bson.D{{Key: "$regexMatch", Value: bson.M{"input": "$username", "regex": usernamePattern, "options": "i"}}},
			}}}),
		})
		matchStage = append(matchStage, bson.E{Key: "auditor_username_match_count", Value: bson.M{"$gt": 0}})
	}

	// General search across level, service, and username fields
	if filter.Search != "" {
		log.Infof("[CPSAction][buildActionCodeFilterPipeline] applying general search: %s", filter.Search)
		searchPattern := fmt.Sprintf("(?i).*%s.*", regexp.QuoteMeta(filter.Search))

		groupStage = append(groupStage, bson.E{
			Key: "general_search_match_count",
			Value: sumWhen(bson.D{{Key: "$or", Value: bson.A{
				bson.D{{Key: "$regexMatch", Value: bson.M{"input": "$action_code", "regex": searchPattern, "options": "i"}}},
				bson.D{{Key: "$regexMatch", Value: bson.M{"input": "$checker_level", "regex": searchPattern, "options": "i"}}},
				bson.D{{Key: "$regexMatch", Value: bson.M{"input": "$auditor_level", "regex": searchPattern, "options": "i"}}},
				bson.D{{Key: "$regexMatch", Value: bson.M{"input": "$action_taken_service_name", "regex": searchPattern, "options": "i"}}},
				bson.D{{Key: "$regexMatch", Value: bson.M{"input": "$username", "regex": searchPattern, "options": "i"}}},
				bson.D{{Key: "$regexMatch", Value: bson.M{"input": "$user_phone", "regex": searchPattern, "options": "i"}}},
			}}}),
		})
		matchStage = append(matchStage, bson.E{Key: "general_search_match_count", Value: bson.M{"$gt": 0}})
	}

	// AuditorCustomerBared filter - filter actions where customer was barred by auditor
	if filter.AuditorCustomerBared != nil {
		log.Infof("[CPSAction][buildActionCodeFilterPipeline] applying auditor_customer_bared filter: %v", *filter.AuditorCustomerBared)
		groupStage = append(groupStage, bson.E{
			Key:   "auditor_customer_bared_match_count",
			Value: sumWhen(bson.D{{Key: "$eq", Value: bson.A{"$auditor_customer_bared", *filter.AuditorCustomerBared}}}),
		})
		matchStage = append(matchStage, bson.E{Key: "auditor_customer_bared_match_count", Value: bson.M{"$gt": 0}})
	}

	if slices.Contains(filter.ActionAuditorStatuses, string(constants.AUDITORNOTCHECKED)) {
		filter.Responsibilities = []string{}
	}

	if len(filter.Responsibilities) > 0 {
		log.Infof("[CPSAction][buildActionCodeFilterPipeline] applying responsibility filter: %v", filter.Responsibilities)

		groupStage = append(groupStage, bson.E{
			Key:   "responsibility_match_count",
			Value: sumWhen(bson.D{{Key: "$in", Value: bson.A{fieldResponsibility, filter.Responsibilities}}}),
		})
		matchStage = append(matchStage, bson.E{Key: "responsibility_match_count", Value: bson.M{"$gt": 0}})
	}

	// Each LevelClaimPair requires that the same action_code has at least one log
	// entry where auditor_level == Level AND given_auditor_status == Claim.
	// All pairs are ANDed: every pair must be satisfied within the same action_code.
	for i, pair := range filter.LevelClaimPairs {
		fieldName := fmt.Sprintf("level_claim_%d_count", i)
		groupStage = append(groupStage, bson.E{
			Key: fieldName,
			Value: sumWhen(bson.D{{Key: "$and", Value: bson.A{
				bson.D{{Key: "$eq", Value: bson.A{"$auditor_level", pair.Level}}},
				bson.D{{Key: "$eq", Value: bson.A{"$given_auditor_status", pair.Claim}}},
			}}}),
		})
		matchStage = append(matchStage, bson.E{Key: fieldName, Value: bson.M{"$gt": 0}})
	}

	preMatch := bson.M{"is_deleted": false}
	if len(filter.RequestActions) > 0 {
		log.Infof("[CPSAction][buildActionCodeFilterPipeline] applying request_action filter: %v", filter.RequestActions)
		preMatch["request_action"] = bson.M{"$in": filter.RequestActions}
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: preMatch}},
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

// applyGivenActionStatusFilter normalises "action_status" or "given_action_status"
// into a $in query on given_action_status, accepting string, []string or []interface{}.
func applyGivenActionStatusFilter(filter map[string]interface{}) {
	for _, key := range []string{"action_status", "given_action_status"} {
		raw, ok := filter[key]
		if !ok {
			continue
		}
		delete(filter, key)
		switch v := raw.(type) {
		case string:
			if v != "" {
				filter["given_action_status"] = bson.M{"$in": []string{v}}
			}
		case []string:
			if len(v) > 0 {
				filter["given_action_status"] = bson.M{"$in": v}
			}
		case []interface{}:
			statuses := make([]string, 0, len(v))
			for _, item := range v {
				if s, ok := item.(string); ok && s != "" {
					statuses = append(statuses, s)
				}
			}
			if len(statuses) > 0 {
				filter["given_action_status"] = bson.M{"$in": statuses}
			}
		}
		break
	}
}

// applyGivenAuditorStatusFilter normalises "auditor_status" or "given_auditor_status"
// into a $in query on given_auditor_status, accepting string, []string or []interface{}.
func applyGivenAuditorStatusFilter(filter map[string]interface{}) {
	for _, key := range []string{"auditor_status", "given_auditor_status"} {
		raw, ok := filter[key]
		if !ok {
			continue
		}
		delete(filter, key)
		switch v := raw.(type) {
		case string:
			if v != "" {
				filter["given_auditor_status"] = bson.M{"$in": []string{v}}
			}
		case []string:
			if len(v) > 0 {
				filter["given_auditor_status"] = bson.M{"$in": v}
			}
		case []interface{}:
			statuses := make([]string, 0, len(v))
			for _, item := range v {
				if s, ok := item.(string); ok && s != "" {
					statuses = append(statuses, s)
				}
			}
			if len(statuses) > 0 {
				filter["given_auditor_status"] = bson.M{"$in": statuses}
			}
		}
		break
	}
}

// applyActionAuditorStatusFilter normalises "action_auditor_status" into a $in query,
// accepting string, []string or []interface{}. Values: NOTCHECKED, INPROGRESS, CHECKED.
func applyActionAuditorStatusFilter(filter map[string]interface{}) {
	raw, ok := filter["action_auditor_status"]
	if !ok {
		return
	}
	delete(filter, "action_auditor_status")
	switch v := raw.(type) {
	case string:
		if v != "" {
			filter["action_auditor_status"] = bson.M{"$in": []string{v}}
		}
	case []string:
		if len(v) > 0 {
			filter["action_auditor_status"] = bson.M{"$in": v}
		}
	case []interface{}:
		statuses := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok && s != "" {
				statuses = append(statuses, s)
			}
		}
		if len(statuses) > 0 {
			filter["action_auditor_status"] = bson.M{"$in": statuses}
		}
	}
}

func (r *userActionLogRepository) GetActionCodesByFilter(ctx context.Context, filter map[string]interface{}) ([]string, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[UserActionLog][GetActionCodesByFilter] filter=%v", filter)

	filter["is_deleted"] = false
	applyGivenActionStatusFilter(filter)
	applyGivenAuditorStatusFilter(filter)
	applyActionAuditorStatusFilter(filter)

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
	update := bson.M{"$set": bson.M{
		"given_action_status":   "APPROVED",
		"action_auditor_status": "NOTCHECKED",
		"last_modified_at":      time.Now(),
	}}

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
	update := bson.M{"$set": bson.M{
		"given_action_status":   "REJECTED",
		"action_auditor_status": "NOTCHECKED",
		"last_modified_at":      time.Now(),
	}}

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
	update := bson.M{"$set": bson.M{
		"given_action_status": "CANCELED",
		"last_modified_at":    time.Now(),
	}}

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

// AuditorMarkLogsByActionCode propagates the auditor's mark verdict (givenAuditorStatus: e.g. MARKASRIGHT/MARKASWRONG)
// and the auditor process state (actionAuditorStatus: INPROGRESS/CHECKED) to all logs for that action_code.
func (r *userActionLogRepository) AuditorMarkLogsByActionCode(ctx context.Context, actionCode string, givenAuditorStatus string, actionAuditorStatus string, customerBared bool) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[UserActionLog][AuditorMarkLogsByActionCode] action_code=%s given_auditor_status=%s action_auditor_status=%s customer_bared=%v", actionCode, givenAuditorStatus, actionAuditorStatus, customerBared)

	filter := bson.M{"action_code": actionCode}
	update := bson.M{"$set": bson.M{
		"given_auditor_status":   givenAuditorStatus,
		"action_auditor_status":  actionAuditorStatus,
		"auditor_customer_bared": customerBared,
		"last_modified_at":       time.Now(),
	}}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Errorf("[UserActionLog][AuditorMarkLogsByActionCode] update failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

// UpdateAuditorActionStatusByActionCode bulk-updates action_auditor_status for all logs with the given action_code.
// Called by AuditorClaim (INPROGRESS) and can be used for any auditor process state transition.
func (r *userActionLogRepository) UpdateAuditorActionStatusByActionCode(ctx context.Context, actionCode string, actionAuditorStatus string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[UserActionLog][UpdateAuditorActionStatusByActionCode] action_code=%s action_auditor_status=%s", actionCode, actionAuditorStatus)

	filter := bson.M{"action_code": actionCode}
	update := bson.M{"$set": bson.M{
		"action_auditor_status": actionAuditorStatus,
		"last_modified_at":      time.Now(),
	}}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Errorf("[UserActionLog][UpdateAuditorActionStatusByActionCode] update failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

// UpdateReinstateStatus updates user_action_log entries for a given action_code to mark customer as reinstated.
// This sets IsCustomerReinstated = true, ReinstateReason, and AuditorCustomerBared = false.
func (r *userActionLogRepository) UpdateReinstateStatus(ctx context.Context, actionCode string, reason string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[UserActionLog][UpdateReinstateStatus] action_code=%s reason=%s", actionCode, reason)

	filter := bson.M{"action_code": actionCode}
	update := bson.M{"$set": bson.M{
		"is_customer_reinstated": true,
		"reinstate_reason":       reason,
		"auditor_customer_bared": false,
		"last_modified_at":       time.Now(),
	}}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Errorf("[UserActionLog][UpdateReinstateStatus] update failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}
