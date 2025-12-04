package archived_user

import (
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

// 	"cbe-super-app-cps-action/internal/constants/localization"
// 	local_util "cbe-super-app-cps-action/pkgs/utils"

// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
// 	"go.mongodb.org/mongo-driver/v2/bson"
// 	"go.mongodb.org/mongo-driver/v2/mongo"
// )

type archivedUserStorage struct {
	dal    dal.MongoDal[model.ArchivedUser, model.ArchivedUser]
	client *mongo.Client
	logger utils.Logger
}

func NewArchivedUserRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.ArchivedUserRepository {
	return &archivedUserStorage{
		dal:    dal.NewMongoDal[model.ArchivedUser, model.ArchivedUser](client, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (a *archivedUserStorage) Create(ctx context.Context, user *model.User) error {

	archivedUser := UserToArchivedUser(user)
	_, err := a.dal.InsertOne(ctx, *archivedUser)
	if err != nil {
		a.logger.Errorf("Error while creating archive user error: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (a *archivedUserStorage) FindByID(ctx context.Context, id string) (*model.ArchivedUser, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}
	result, err := a.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *archivedUserStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.ArchivedUser], error) {
	// 1. Base filter (only active records)
	filter := bson.M{"is_deleted": false}
	searchKeys := bson.M{}
	// 2. Allowed filterable/searchable fields
	allowedKeys := []string{"enabled", "kyc_level", "is_blocked", "is_verified"}

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
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 6. Count total
	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.ArchivedUser]{
		Data: data,
		Meta: meta,
	}, nil
}
