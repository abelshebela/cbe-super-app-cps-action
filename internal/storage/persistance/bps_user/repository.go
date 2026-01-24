package bps_user

import (
	// "cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"regexp"
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

	filter := bson.M{"$or": conditions}
	data, err := b.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			b.logger.Infof("[FindByOr] no bps user found matching the criteria")
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		b.logger.Errorf("[FindByOr] failed to find bps user: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return data, nil
}

func (b *BPSUserStorage) GetByUserCode(ctx context.Context, userCode string) (*bps_model.BPSUser, error) {
	b.logger.Infof("[GetByUserCode] fetching BPS user by user code")
	filter := bson.M{"user_code": userCode, "is_deleted": false}
	result, err := b.dal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		b.logger.Errorf("[GetByUserCode] failed to find BPS user: %v", err)
		return nil, err
	}
	b.logger.Infof("[GetByUserCode] BPS user retrieved successfully")
	return result, nil
}

func (s *BPSUserStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]bpsUserDto.BPSUserResposenDTO], error) {
	// Base match: only active users
	match := bson.M{"is_deleted": false}

	// Search filters
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		match["$or"] = []bson.M{
			{"full_name": searchRegex},
			{"username": searchRegex},
			{"user_code": searchRegex},
			{"phone_number": searchRegex},
		}
	}

	pipeline := mongo.Pipeline{
		// Match stage
		bson.D{{Key: "$match", Value: match}},

		// Lookup stage
		bson.D{{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "roles"},
				{Key: "localField", Value: "job_title"},
				{Key: "foreignField", Value: "job_title"},
				{Key: "as", Value: "roles"},
			},
		}},

		// Unwind roles
		bson.D{{
			Key: "$unwind",
			Value: bson.D{
				{Key: "path", Value: "$roles"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			},
		}},

		// Add role_name field
		bson.D{{
			Key: "$addFields",
			Value: bson.D{
				{Key: "user_role", Value: "$roles.role"},
			},
		}},

		// Remove roles array
		bson.D{{
			Key: "$project",
			Value: bson.D{
				{Key: "roles", Value: 0},
			},
		}},

		// Pagination
		bson.D{{Key: "$skip", Value: filterParam.PerPage * (filterParam.Page - 1)}},
		bson.D{{Key: "$limit", Value: filterParam.PerPage}},
	}

	//  Execute aggregation
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		s.logger.Errorf("[FindAllWithPagination] failed aggregation: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	defer cursor.Close(ctx)

	var users []bpsUserDto.BPSUserResposenDTO
	if err := cursor.All(ctx, &users); err != nil {
		s.logger.Errorf("[FindAllWithPagination] failed to decode users: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// Count total documents (for pagination metadata)
	total, err := s.collection.CountDocuments(ctx, match)
	if err != nil {
		s.logger.Errorf("[FindAllWithPagination] failed to count users: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	//  Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	s.logger.Infof("[FindAllWithPagination] retrieved %d BPS users", len(users))

	//  Return paginated response
	return &types.PaginatedResponse[[]bpsUserDto.BPSUserResposenDTO]{
		Data: users,
		Meta: meta,
	}, nil
}

// func (s *BPSUserStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]bps_model.BPSUser], error) {
// 	// 1. Base filter (only active records)
// 	filter := bson.M{"is_deleted": false}
// 	searchKeys := bson.M{}

// 	// 2. Allowed filterable/searchable fields
// 	allowedKeys := []string{"branch_code", "branch_name", "enabled", "role", "first_password_set", "full_name", "user_code", "phone_number", "username", "branch_code"}

// 	// 3. Add search (if provided)
// 	if filterParam.Search != "" {
// 		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
// 		searchKeys["$or"] = []bson.M{
// 			{"full_name": searchRegex},
// 			{"username": searchRegex},
// 			{"user_code": searchRegex},
// 			{"phone_number": searchRegex},
// 		}
// 	}
// 	// 4. Build filter, skip, limit
// 	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

// 	// 5. Fetch data
// 	data, err := s.dal.FindAllWithPaginationE(ctx, filter, bson.M{}, skip, limit)
// 	if err != nil {
// 		s.logger.Errorf("[FindAllWithPagination] failed to fetch BPS users: %v", err)
// 		return nil, errors.New(localization.ErrorUnexpectedError.Message)
// 	}

// 	// 6. Count total
// 	total, err := s.dal.TotalCount(ctx, filter)
// 	if err != nil {
// 		s.logger.Errorf("[FindAllWithPagination] failed to count BPS users: %v", err)
// 		return nil, errors.New(localization.ErrorUnexpectedError.Message)
// 	}

// 	// 7. Build pagination metadata
// 	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
// 	s.logger.Infof("[FindAllWithPagination] retrieved %d BPS users", len(data))

// 	// 8. Return standard paginated response
// 	return &types.PaginatedResponse[[]bps_model.BPSUser]{
// 		Data: data,
// 		Meta: meta,
// 	}, nil
// }

func (b *BPSUserStorage) Update(ctx context.Context, BpsUser *bps_model.BPSUser) error {
	b.logger.Infof("[Update] updating BPS user")
	filter := bson.M{"_id": BpsUser.ID, "is_deleted": false}
	_, err := b.dal.UpdateOne(ctx, filter, BPSUserMapper(*BpsUser))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			b.logger.Errorf("[Update] BPS user not found")
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		b.logger.Errorf("[Update] failed to update BPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	b.logger.Infof("[Update] BPS user updated successfully")
	return nil

}

func (b *BPSUserStorage) Create(ctx context.Context, req bps_model.BPSUser) error {

	ctx, span := local_util.TraceLogger(ctx, "service", "CreateBPSUser", "BPS User", "CreateBPSUser")
	defer span.End()

	b.logger.Infof("[CreateBPSUser] creating BPS user with user_code: %s", req.UserCode)
	req.CreatedAt = time.Now()
	req.LastModifiedAt = time.Now()
	req.IsFirstTimeLogin = true
	// Save the new user
	_, err := b.dal.InsertOne(ctx, req)
	if err != nil {
		// span.AddEvent("[CreateBPSUser] failed to create BPS user", trace.WithAttributes(
		//     attribute.String("error", err.Error()),
		//     attribute.String("user_code", req.UserCode),
		// ))
		b.logger.Errorf("[CreateBPSUser] failed to create BPS user: %v", err)
		return err
	}

	b.logger.Infof("[CreateBPSUser] BPS user created successfully for user_code: %s", req.UserCode)
	return nil
}

func (b *BPSUserStorage) FindByFilterKey(ctx context.Context, field, value string) (*bps_model.BPSUser, error) {
	b.logger.Infof("[FindByFilterKey] searching BPS user by %s: %s", field, value)
	var filter bson.M

	if field == "id" {
		field = "_id"
		objID, err := bson.ObjectIDFromHex(value)
		if err != nil {
			b.logger.Errorf("[FindByFilterKey] invalid ObjectID: %v", err)
			return nil, errors.New("invalid id format")
		}
		filter = bson.M{field: objID, "is_deleted": false}
	} else {
		filter = bson.M{field: value, "is_deleted": false}
	}

	result, err := b.dal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		b.logger.Errorf("[FindByFilterKey] failed to find BPS user by %s: %v", field, err)
		return nil, err
	}
	return result, nil
}
