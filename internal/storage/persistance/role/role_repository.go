package role_repo

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RoleRepository struct {
	client     *mongo.Client
	mongoDal   dal.MongoDal[imodel.Role, imodel.Role]
	logger     utils.Logger
	collection *mongo.Collection
}

func NewRoleRepository(client *mongo.Client, cfg *config.VaultConfig, database string, collection []string, logger utils.Logger) storage.RoleRepository {
	return &RoleRepository{
		client:     client,
		mongoDal:   dal.NewMongoDal[imodel.Role, imodel.Role](client, cfg, database, collection[0]),
		logger:     logger,
		collection: client.Database(database).Collection(collection[1]),
	}
}

func (r *RoleRepository) Exists(ctx context.Context, id string) (bool, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("[RoleRepository][Exists] invalid object id: %v", err)
		return false, err
	}
	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": oid})
	if err != nil {
		r.logger.Errorf("[RoleRepository][Exists] failed to count documents: %v", err)
		return false, err
	}
	return count > 0, nil
}

func (r *RoleRepository) ExistsMany(ctx context.Context, ids []string) (bool, error) {
	if len(ids) == 0 {
		return true, nil
	}
	oids := make([]string, 0, len(ids))
	for _, id := range ids {
		oids = append(oids, id)
	}

	count, err := r.collection.CountDocuments(ctx, bson.M{"code": bson.M{"$in": oids}})
	if err != nil {
		r.logger.Errorf("[RoleRepository][ExistsMany] failed to count documents: %v", err)
		return false, err
	}
	return count == int64(len(ids)), nil
}

func (r *RoleRepository) Create(ctx context.Context, role *imodel.Role) error {
	role.ID = bson.NewObjectID()
	_, err := r.mongoDal.InsertOne(ctx, *role)
	if err != nil {
		r.logger.Errorf("[RoleRepository][Create] failed to create role: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (r *RoleRepository) Update(ctx context.Context, id string, role *imodel.Role) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("[RoleRepository][Update] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	update := bson.M{}
	if role.JobTitle != "" {
		update["job_title"] = role.JobTitle
	}
	if role.Role != "" {
		update["role"] = role.Role
	}
	if !role.UpdateAt.IsZero() {
		update["updated_at"] = role.UpdateAt
	}

	filter := bson.M{"_id": objID}
	_, err = r.mongoDal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("[RoleRepository][Update] failed to update role: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (r *RoleRepository) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("[RoleRepository][EnableOrDisable] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"enabled": enable, "updated_at": time.Now()}

	_, err = r.mongoDal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("[RoleRepository][EnableOrDisable] failed to enable/disable role: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (r *RoleRepository) SoftDelete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("[RoleRepository][SoftDelete] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"is_deleted": true, "deleted_at": time.Now()}

	_, err = r.mongoDal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("[RoleRepository][SoftDelete] failed to soft delete role: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (r *RoleRepository) FindByID(ctx context.Context, id string) (*imodel.Role, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("[RoleRepository][FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID}
	result, err := r.mongoDal.FindOne(ctx, filter, nil)
	if err != nil {
		r.logger.Errorf("[RoleRepository][FindByID] failed to find role: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}

func (r *RoleRepository) FindByName(ctx context.Context, name string) (*imodel.Role, error) {
	filter := bson.M{"job_title": bson.M{"$regex": "^" + name + "$", "$options": "i"}}
	result, err := r.mongoDal.FindOne(ctx, filter, nil)
	if err != nil {
		r.logger.Errorf("[RoleRepository][FindByName] failed to find role: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}
func (r *RoleRepository) FindByRole(ctx context.Context, name string) (*imodel.Role, error) {
	filter := bson.M{"role": bson.M{"$regex": "^" + name + "$", "$options": "i"}}
	result, err := r.mongoDal.FindOne(ctx, filter, nil)
	if err != nil {
		r.logger.Errorf("[RoleRepository][FindByRole] failed to find role: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}

func (r *RoleRepository) FindByCode(ctx context.Context, code string) (*imodel.Role, error) {
	filter := bson.M{"code": code}
	result, err := r.mongoDal.FindOne(ctx, filter, nil)
	if err != nil {
		r.logger.Errorf("[RoleRepository][FindByCode] failed to find role: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}

func (r *RoleRepository) FindAll(ctx context.Context) (*[]imodel.Role, error) {
	data, err := r.mongoDal.FindAll(ctx, bson.M{}, nil)
	if err != nil {
		r.logger.Errorf("[RoleRepository][FindAll] failed to find roles: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return &data, nil
}
func (r *RoleRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.Role], error) {
	searchKeys := bson.M{}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"job_title": searchRegex},
			{"role": searchRegex},
		}
	}

	allowedKeys := []string{"enabled"}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	data, err := r.mongoDal.FindAllWithPaginationE(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		r.logger.Errorf("[RoleRepository][FindAllWithPagination] failed to fetch roles: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	total, err := r.mongoDal.TotalCount(ctx, filter)
	if err != nil {
		r.logger.Errorf("[RoleRepository][FindAllWithPagination] failed to count roles: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]imodel.Role]{
		Data: data,
		Meta: meta,
	}, nil
}

func (r *RoleRepository) FindByFilterKey(ctx context.Context, field string, value string) (*imodel.Role, error) {
	var filter bson.M

	if field == "id" || field == "_id" {
		objID, err := bson.ObjectIDFromHex(value)
		if err != nil {
			r.logger.Errorf("[RoleRepository][FindByFilterKey] invalid ObjectID: %v", err)
			return nil, errors.New(localization.ErrorInvalidID.Code)
		}
		filter = bson.M{"_id": objID}
	} else {
		filter = bson.M{field: value}
	}

	result, err := r.mongoDal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}
