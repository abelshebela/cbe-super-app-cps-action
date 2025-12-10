package bps_user

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BPSUserStorage struct {
	dal    dal.MongoDal[model.BPSUser, model.BPSUser]
	client *mongo.Client
	logger utils.Logger
}

func NewBPSUserRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.BPSUserRepository {
	return &BPSUserStorage{
		dal:    dal.NewMongoDal[model.BPSUser, model.BPSUser](client, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (b *BPSUserStorage) GetByUserCode(ctx context.Context, userCode string) (*model.BPSUser, error) {
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

func (s *BPSUserStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.BPSUser], error) {
	// 1. Base filter (only active records)
	filter := bson.M{"is_deleted": false}
	searchKeys := bson.M{}

	// 2. Allowed filterable/searchable fields
	allowedKeys := []string{"branch_code", "branch_name", "enabled", "role", "first_password_set", "full_name", "user_code", "phone_number", "username", "branch_code"}

	// 3. Add search (if provided)
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"full_name": searchRegex},
			{"username": searchRegex},
			{"user_code": searchRegex},
			{"phone_number": searchRegex},
		}
	}
	// 4. Build filter, skip, limit
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	// 5. Fetch data
	data, err := s.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		s.logger.Errorf("[FindAllWithPagination] failed to fetch BPS users: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 6. Count total
	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		s.logger.Errorf("[FindAllWithPagination] failed to count BPS users: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	s.logger.Infof("[FindAllWithPagination] retrieved %d BPS users", len(data))

	// 8. Return standard paginated response
	return &types.PaginatedResponse[[]*model.BPSUser]{
		Data: data,
		Meta: meta,
	}, nil
}

func (b *BPSUserStorage) Update(ctx context.Context, BpsUser *model.BPSUser) error {
	b.logger.Infof("[Update] updating BPS user")
	filter := bson.M{"user_code": BpsUser.UserCode, "is_deleted": false}
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
