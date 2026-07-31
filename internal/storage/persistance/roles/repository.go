package job_role

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RoleStorage struct {
	client     *mongo.Client
	collection *mongo.Collection
	logger     utils.Logger
	dal        dal.MongoDal[imodel.Role, imodel.Role]
}

func NewRoleRepository(client *mongo.Client, cfg *config.VaultConfig, database, collection string, logger utils.Logger) storage.RoleRepository {
	return &RoleStorage{
		client:     client,
		collection: client.Database(database).Collection(collection),
		logger:     logger,
		dal:        dal.NewMongoDal[imodel.Role, imodel.Role](client, cfg, database, collection),
	}
}

func (s *RoleStorage) Create(ctx context.Context, role *imodel.Role) error {
	s.logger.Infof("[JobRole/Create] creating job role: code=%s name=%s", role.Code, role.Name)
	role.ID = bson.NewObjectID()
	_, err := s.dal.InsertOne(ctx, *role)
	if err != nil {
		s.logger.Errorf("[JobRole/Create] failed to insert: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (s *RoleStorage) Update(ctx context.Context, id string, role *imodel.Role) error {
	s.logger.Infof("[JobRole/Update] id=%s", id)
	objID, err := local_util.ParseObjectID(id)
	if err != nil {
		s.logger.Errorf("[JobRole/Update] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := RoleMapper(*role)
	_, err = s.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		return local_util.HandleDBError(err)
	}
	return nil
}

func (s *RoleStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := local_util.ParseObjectID(id)
	if err != nil {
		s.logger.Errorf("[JobRole/EnableOrDisable] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"enabled": enable, "updated_at": time.Now()}

	_, err = s.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		return local_util.HandleDBError(err)
	}
	return nil
}

func (s *RoleStorage) SoftDelete(ctx context.Context, id string) error {
	objID, err := local_util.ParseObjectID(id)
	if err != nil {
		s.logger.Errorf("[JobRole/SoftDelete] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"is_deleted": true, "deleted_at": time.Now()}

	_, err = s.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		return local_util.HandleDBError(err)
	}
	return nil
}

func (s *RoleStorage) FindByID(ctx context.Context, id string) (*imodel.Role, error) {
	s.logger.Infof("[JobRole/FindByID] id=%s", id)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Errorf("[JobRole/FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}

	res, err := s.dal.FindOne(ctx, filter, nil)
	if err != nil {
		s.logger.Errorf("[JobRole/FindByID] failed to find by id: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return res, nil
}

func (s *RoleStorage) FindAll(ctx context.Context) (*[]imodel.Role, error) {

	filter := dal.FilterOp{
		Filter: bson.M{"enabled": true, "is_deleted": false},
	}
	data, err := s.dal.FindAllWithCursorBasedPagination(ctx, filter)

	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	return &data, nil
}

func (r *RoleStorage) ExistsMany(ctx context.Context, codes []string) (bool, error) {
	if len(codes) == 0 {
		return true, nil
	}

	count, err := r.collection.CountDocuments(ctx, bson.M{"code": bson.M{"$in": codes}, "is_deleted": false})
	if err != nil {
		return false, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return count == int64(len(codes)), nil
}

func (s *RoleStorage) FindByCode(ctx context.Context, code string) (*imodel.Role, error) {
	filter := bson.M{
		"code":       code,
		"is_deleted": false,
	}
	res, err := s.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	return res, nil
}

func (s *RoleStorage) Find(ctx context.Context, filter bson.M) (*imodel.Role, error) {

	filter["is_deleted"] = false
	res, err := s.dal.FindOne(ctx, filter, nil)
	if err != nil {
		s.logger.Errorf("[JobRole/Find] failed to find: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return res, nil
}

func (s *RoleStorage) FindByName(ctx context.Context, name string) (*imodel.Role, error) {
	filter := bson.M{
		"name": bson.M{
			"$regex":   "^" + strings.ToLower(name) + "$",
			"$options": "i",
		},
		"is_deleted": false,
	}
	res, err := s.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	return res, nil
}

func (r *RoleStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.Role], error) {
	searchKeys := bson.M{}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"name": searchRegex},
			{"code": searchRegex},
		}
	}

	allowedKeys := []string{"enabled", "name", "code", "type"}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["is_deleted"] = false

	data, err := r.dal.FindAllWithPaginationE(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}

	total, err := r.dal.TotalCount(ctx, filter)
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
