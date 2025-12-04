package action_role_repo

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

type ActionRoleRepository struct {
	client   *mongo.Client
	mongoDal dal.MongoDal[model.ActionRole, model.ActionRole]
	logger   utils.Logger
}

func NewActionRoleRepository(client *mongo.Client, database, collection string, logger utils.Logger) storage.ActionRoleRepository {
	return &ActionRoleRepository{
		client:   client,
		mongoDal: dal.NewMongoDal[model.ActionRole, model.ActionRole](client, database, collection),
		logger:   logger,
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

	data, err := r.mongoDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := r.mongoDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.ActionRole]{
		Data: data,
		Meta: meta,
	}, nil
}
