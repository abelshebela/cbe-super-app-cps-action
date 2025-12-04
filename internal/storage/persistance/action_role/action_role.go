package action_role_repo

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ActionRoleRepository struct {
	client  *mongo.Client
	mongoDal dal.MongoDal[model.ActionRole, model.ActionRole]
	logger  utils.Logger
}

func NewActionRoleRepository(client *mongo.Client, database, collection string, logger utils.Logger) storage.ActionRoleRepository {
	return &ActionRoleRepository{
		client:  client,
		mongoDal: dal.NewMongoDal[model.ActionRole, model.ActionRole](client, database, collection),
		logger:  logger,
	}
}

func (r *ActionRoleRepository) Create(ctx context.Context, actionRole *model.ActionRole) error {
	_, err := r.mongoDal.InsertOne(ctx, *actionRole)
	return err
}

func (r *ActionRoleRepository) UpdateByActionCode(ctx context.Context, actionCode string, actionRole *model.ActionRole) error {
	update := bson.M{
		"action_name":       actionRole.ActionName,
		"assigned_makers":   actionRole.AssignedMakers,
		"assigned_checkers": actionRole.AssignedChecker,
		"enabled":           actionRole.Enabled,
		"updated_at":        actionRole.UpdatedAt,
	}
	_, err := r.mongoDal.UpdateOne(ctx, bson.M{"action_code": actionCode}, update)
	return err
}

func (r *ActionRoleRepository) EnableOrDisableByActionCode(ctx context.Context, actionCode string, enable bool) error {
	_, err := r.mongoDal.UpdateOne(ctx, bson.M{"action_code": actionCode}, bson.M{"enabled": enable})
	return err
}

func (r *ActionRoleRepository) FindByActionCode(ctx context.Context, actionCode string) (*model.ActionRole, error) {
	return r.mongoDal.FindOne(ctx, bson.M{"action_code": actionCode}, bson.M{})
}

func (r *ActionRoleRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.ActionRole], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"action_code", "action_name", "enabled"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"action_code": searchRegex},
			{"action_name": searchRegex},
		}
	}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	items, err := r.mongoDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, err
	}
	totalDocs, err := r.mongoDal.TotalCount(ctx, filter)
	if err != nil {
		totalDocs = int64(skip) + int64(len(items))
	}
	if limit <= 0 {
		if len(items) > 0 { limit = int64(len(items)) } else { limit = 1 }
	}
	page := filterParam.Page
	if page <= 0 { page = 1 }
	totalPages := int((totalDocs + int64(limit) - 1) / int64(limit))
	if totalPages == 0 { totalPages = 1 }
	pagingCounter := (int64(page)-1)*limit + 1
	hasPrev := page > 1
	hasNext := page < totalPages
	var prevPage *int
	if hasPrev { p := page - 1; prevPage = &p }
	var nextPage *int
	if hasNext { n := page + 1; nextPage = &n }

	return &types.PaginatedResponse[[]*model.ActionRole]{
		Data: items,
		Meta: types.PaginationMeta{
			TotalDocs:     totalDocs,
			Limit:         int(limit),
			TotalPages:    totalPages,
			Page:          page,
			PagingCounter: int(pagingCounter),
			HasPrevPage:   hasPrev,
			HasNextPage:   hasNext,
			PrevPage:      prevPage,
			NextPage:      nextPage,
		},
	}, nil
}
