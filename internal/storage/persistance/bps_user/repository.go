package bps_user

import (
	// "cbe-super-app-cps-action/internal/constants/lib"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	// "time"

	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/local_model"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	bpsUserDto "cbe-super-app-cps-action/internal/constants/dto/bps_user"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	bps_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// bpsUserActiveFilter matches non–soft-deleted users. Use $ne: true so documents
// without is_deleted (legacy) or with false are included; only true is excluded.
func bpsUserActiveFilter() bson.M {
	return bson.M{"is_deleted": bson.M{"$ne": true}}
}

type BPSUserStorage struct {
	dal        dal.MongoDal[bps_model.BPSUser, bps_model.BPSUser]
	client     *mongo.Client
	logger     utils.Logger
	collection *mongo.Collection
}

func NewBPSUserRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, logger utils.Logger) storage.BPSUserRepository {
	return &BPSUserStorage{
		dal:        dal.NewMongoDal[bps_model.BPSUser, bps_model.BPSUser](client, cfg, dbName, collection),
		client:     client,
		logger:     logger,
		collection: client.Database(dbName).Collection(collection),
	}
}

func (b *BPSUserStorage) FindByOr(ctx context.Context, phone, email, username string) (*bps_model.BPSUser, error) {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	log.Infof("[BPSUserStorage][FindByOr] searching for BPS user by phone, email, or username")
	// Build conditions dynamically, only for non-empty parameters
	conditions := []bson.M{}

	if phone != "" {
		conditions = append(conditions, bson.M{
			"phone_number": bson.M{"$regex": "^" + regexp.QuoteMeta(phone) + "$", "$options": "i"},
		})
	}

	if username != "" {
		conditions = append(conditions, bson.M{
			"username": bson.M{"$regex": "^" + regexp.QuoteMeta(username) + "$", "$options": "i"},
		})
	}

	if email != "" {
		conditions = append(conditions, bson.M{
			"email": bson.M{"$regex": "^" + regexp.QuoteMeta(email) + "$", "$options": "i"},
		})
	}
	if len(conditions) == 0 {
		log.Infof("[BPSUserStorage][FindByOr] no search parameters provided, returning nil")
		return nil, nil
	}
	filter := bson.M{"$or": conditions}
	data, err := b.dal.FindOne(ctx, filter, nil)
	if err != nil {
		log.Errorf("[BPSUserStorage][FindByOr] failed to find bps user: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return data, nil
}

func (b *BPSUserStorage) GetByUserCode(ctx context.Context, userCode string) (*bpsUserDto.BPSUserResposenDTO, error) {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	log.Infof("[BPSUserStorage][GetByUserCode] fetching BPS user by user code")

	pipeline := mongo.Pipeline{
		// Match by user_code and not deleted
		bson.D{{Key: "$match", Value: bson.M{
			"user_code":  userCode,
			"is_deleted": bson.M{"$ne": true},
		}}},

		// Lookup role from roles collection
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "roles"},
			{Key: "localField", Value: "job_title"},
			{Key: "foreignField", Value: "job_title"},
			{Key: "as", Value: "roles"},
		}}},

		// Extract role code
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "role_code", Value: bson.D{
				{Key: "$arrayElemAt", Value: bson.A{"$roles.role", 0}},
			}},
		}}},

		// Lookup job role details
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "job_roles"},
			{Key: "localField", Value: "role_code"},
			{Key: "foreignField", Value: "code"},
			{Key: "as", Value: "job_roles"},
		}}},

		// Add user_role field
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "role", Value: bson.D{
				{Key: "$arrayElemAt", Value: bson.A{"$job_roles.name", 0}},
			}},
		}}},

		// Clean up temporary fields
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "roles", Value: 0},
			{Key: "job_roles", Value: 0},
			{Key: "role_code", Value: 0},
		}}},
	}

	cursor, err := b.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Errorf("[BPSUserStorage][GetByUserCode] failed aggregation: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer cursor.Close(ctx)

	var results []bpsUserDto.BPSUserResposenDTO
	if err := cursor.All(ctx, &results); err != nil {
		log.Errorf("[BPSUserStorage][GetByUserCode] failed to decode user: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	if len(results) == 0 {
		log.Infof("[BPSUserStorage][GetByUserCode] no BPS user found with user_code: %s", userCode)
		return nil, errors.New(localization.ErrorUserNotFound.Code)
	}

	log.Infof("[BPSUserStorage][GetByUserCode] BPS user retrieved successfully with role")
	return &results[0], nil
}
func (b *BPSUserStorage) GetByUserID(ctx context.Context, userID string) (*bps_model.BPSUser, error) {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	log.Infof("[BPSUserStorage][GetByUserID] fetching BPS user by user ID")
	objID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		log.Errorf("[BPSUserStorage][GetByUserID] invalid user ID format: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": bson.M{"$ne": true}}
	result, err := b.dal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		log.Errorf("[BPSUserStorage][GetByUserID] failed to find BPS user: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	log.Infof("[BPSUserStorage][GetByUserID] BPS user retrieved successfully")
	return result, nil
}

func (s *BPSUserStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]bpsUserDto.BPSUserResposenDTO], error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	// Base match: only active users (include legacy docs without is_deleted)
	match := bpsUserActiveFilter()
	// Search filters
	if enabledVal, ok := filterParam.Filters["enabled"]; ok {
		match["enabled"] = enabledVal
	}
	filterParam.Search = strings.TrimSpace(filterParam.Search)
	if strings.HasPrefix(filterParam.Search, "09") || strings.HasPrefix(filterParam.Search, "07") {
		filterParam.Search = filterParam.Search[1:]
	}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		match["$or"] = []bson.M{
			{"full_name": searchRegex},
			{"username": searchRegex},
			{"user_code": searchRegex},
			{"phone_number": searchRegex},
			{"email": searchRegex},
			{"branch_name": searchRegex},
		}
	}

	// Count total documents (for pagination metadata)
	total, err := s.collection.CountDocuments(ctx, match)
	if err != nil {
		log.Errorf("[BPSUserStorage][FindAllWithPagination] failed to count users: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	offset := filterParam.PerPage * (filterParam.Page - 1)
	if offset >= int(total) {
		// No data for this page, return empty result
		meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
		log.Infof("[BPSUserStorage][FindAllWithPagination] offset >= total, returning empty result")
		return &types.PaginatedResponse[[]bpsUserDto.BPSUserResposenDTO]{
			Data: []bpsUserDto.BPSUserResposenDTO{},
			Meta: meta,
		}, nil
	}
	if offset+filterParam.PerPage > int(total) {
		filterParam.PerPage = int(total) - offset
	}

	pipeline := mongo.Pipeline{

		bson.D{{Key: "$match", Value: match}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "roles"},
			{Key: "localField", Value: "job_title"},
			{Key: "foreignField", Value: "job_title"},
			{Key: "as", Value: "roles"},
		}}},

		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "role_code", Value: bson.D{
				{Key: "$arrayElemAt", Value: bson.A{"$roles.role", 0}}, // Make sure field is "role" not "code"
			}},
		}}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "job_roles"},
			{Key: "localField", Value: "role_code"}, // MUST use the extracted role_code
			{Key: "foreignField", Value: "code"},    // matches job_roles.code
			{Key: "as", Value: "job_roles"},
		}}},

		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "user_role", Value: bson.D{
				{Key: "$arrayElemAt", Value: bson.A{"$job_roles.name", 0}},
			}},
		}}},

		bson.D{{Key: "$project", Value: bson.D{
			{Key: "roles", Value: 0},
			{Key: "job_roles", Value: 0},
			{Key: "role_code", Value: 0},
		}}},

		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "created_at", Value: -1},
		}}},
		bson.D{{Key: "$skip", Value: filterParam.PerPage * (filterParam.Page - 1)}},
		bson.D{{Key: "$limit", Value: filterParam.PerPage}},
	}

	//  Execute aggregation
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Errorf("[BPSUserStorage][FindAllWithPagination] failed aggregation: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	var users []bpsUserDto.BPSUserResposenDTO
	if err := cursor.All(ctx, &users); err != nil {
		log.Errorf("[BPSUserStorage][FindAllWithPagination] failed to decode users: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	//  Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	log.Infof("[BPSUserStorage][FindAllWithPagination] retrieved %d BPS users", len(users))

	//  Return paginated response
	return &types.PaginatedResponse[[]bpsUserDto.BPSUserResposenDTO]{
		Data: users,
		Meta: meta,
	}, nil
}

func (b *BPSUserStorage) Update(ctx context.Context, BpsUser *bps_model.BPSUser) error {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	log.Infof("[BPSUserStorage][Update] updating BPS user")
	filter := bson.M{"_id": BpsUser.ID, "is_deleted": bson.M{"$ne": true}}
	_, err := b.dal.UpdateOne(ctx, filter, BPSUserMapper(*BpsUser))
	if err != nil {
		log.Errorf("[BPSUserStorage][Update] failed to update BPS user: %v", err)
		return local_util.HandleDBError(err)
	}
	log.Infof("[BPSUserStorage][Update] BPS user updated successfully")
	return nil

}

func (b *BPSUserStorage) Create(ctx context.Context, req bps_model.BPSUser) error {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "CreateBPSUser", "BPS User", "CreateBPSUser")
	defer span.End()

	log.Infof("[BPSUserStorage][Create] creating BPS user with user_code: %s", req.UserCode)
	req.CreatedAt = time.Now()
	req.LastModifiedAt = time.Now()
	req.IsFirstTimeLogin = true
	req.IsDeleted = false
	// Save the new user
	_, err := b.dal.InsertOne(ctx, req)
	if err != nil {
		// span.AddEvent("[CreateBPSUser] failed to create BPS user", trace.WithAttributes(
		//     attribute.String("error", err.Error()),
		//     attribute.String("user_code", req.UserCode),
		// ))
		log.Errorf("[BPSUserStorage][Create] failed to create BPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	log.Infof("[BPSUserStorage][Create] BPS user created successfully for user_code: %s", req.UserCode)
	return nil
}

func (b *BPSUserStorage) FindByFilterKey(ctx context.Context, field, value string) (*bps_model.BPSUser, error) {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	log.Infof("[BPSUserStorage][FindByFilterKey] searching BPS user by %s: %s", field, value)
	var filter bson.M

	if field == "id" {
		field = "_id"
		objID, err := bson.ObjectIDFromHex(value)
		if err != nil {
			log.Errorf("[BPSUserStorage][FindByFilterKey] invalid ObjectID: %v", err)
			return nil, errors.New(localization.ErrorInvalidID.Code)
		}
		filter = bson.M{field: objID, "is_deleted": bson.M{"$ne": true}}
	} else {
		filter = bson.M{field: value, "is_deleted": bson.M{"$ne": true}}
	}

	result, err := b.dal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		log.Errorf("[BPSUserStorage][FindByFilterKey] failed to find BPS user: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}

// func BPSUserMapper(user bps_model.BPSUser) bson.M {
