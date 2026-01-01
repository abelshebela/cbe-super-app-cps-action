package role_repo

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
)

type RoleRepository struct {
	client     *mongo.Client
	mongoDal   dal.MongoDal[model.Role, model.Role]
	logger     utils.Logger
	collection *mongo.Collection
}

func NewRoleRepository(client *mongo.Client,cfg *config.VaultConfig, database, collection string, logger utils.Logger) storage.RoleRepository {
	return &RoleRepository{
		client:     client,
		mongoDal:   dal.NewMongoDal[model.Role, model.Role](client, cfg,database, collection),
		logger:     logger,
		collection: client.Database(database).Collection(collection),
	}
}

func (r *RoleRepository) Exists(ctx context.Context, id string) (bool, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}
	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": oid})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RoleRepository) ExistsMany(ctx context.Context, ids []string) (bool, error) {
	if len(ids) == 0 {
		return true, nil
	}
	oids := make([]bson.ObjectID, 0, len(ids))
	for _, id := range ids {
		oid, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return false, err
		}
		oids = append(oids, oid)
	}

	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": bson.M{"$in": oids}})
	if err != nil {
		return false, err
	}
	return count == int64(len(ids)), nil
}

func (r *RoleRepository) Create(ctx context.Context, role *model.Role) error {
	role.ID = bson.NewObjectID()
	_, err := r.mongoDal.InsertOne(ctx, *role)
	if err != nil {
		r.logger.Errorf("[Role Repository] Error while creating Error: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (r *RoleRepository) Update(ctx context.Context, id string, role *model.Role) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("[Role Repository][Update] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	update := bson.M{}
	if role.JobTitle != "" {
		update["job_title"] = role.JobTitle
	}
	if role.Role != "" {
		update["role"] = role.Role
	}
	if !role.UpdatedAt.IsZero() {
		update["updated_at"] = role.UpdatedAt
	}

	filter := bson.M{"_id": objID}
	_, err = r.mongoDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		r.logger.Errorf("[Role Repository][Update] failed to update: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (r *RoleRepository) FindByID(ctx context.Context, id string) (*model.Role, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("[Role Repository][FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID}
	result, err := r.mongoDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		r.logger.Errorf("[Role Repository][FindByID] failed to find: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return result, nil
}

func (r *RoleRepository) FindByName(ctx context.Context, name string) (*model.Role, error) {
	filter := bson.M{"job_title": bson.M{"$regex": "^" + name + "$", "$options": "i"}}
	result, err := r.mongoDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		r.logger.Errorf("[Role Repository][FindByName] failed to find: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return result, nil
}
func (r *RoleRepository) FindByRole(ctx context.Context, name string) (*model.Role, error) {
	filter := bson.M{"role": bson.M{"$regex": "^" + name + "$", "$options": "i"}}
	result, err := r.mongoDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		r.logger.Errorf("[Role Repository][FindByName] failed to find: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return result, nil
}

func (r *RoleRepository) FindByCode(ctx context.Context, code string) (*model.Role, error) {
	filter := bson.M{"code": code}
	result, err := r.mongoDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		r.logger.Errorf("[Role Repository][FindByCode] failed to find: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return result, nil
}

func (r *RoleRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.Role], error) {
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

	data, err := r.mongoDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		r.logger.Errorf("[Role Repository][FindAllWithPagination] fetch error: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := r.mongoDal.TotalCount(ctx, filter)
	if err != nil {
		r.logger.Errorf("[Role Repository][FindAllWithPagination] count error: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]model.Role]{
		Data: data,
		Meta: meta,
	}, nil
}

func (r *RoleRepository) FindByFilterKey(ctx context.Context, field string, value string) (*model.Role, error) {
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
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		r.logger.Errorf("[RoleRepository][FindByFilterKey] failed to find: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return result, nil
}
