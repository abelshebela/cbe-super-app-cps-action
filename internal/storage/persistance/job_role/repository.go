package job_role

import (
	"context"
	"errors"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type JobRoleStorage struct {
	client     *mongo.Client
	collection *mongo.Collection
	logger     utils.Logger
	dal        dal.MongoDal[imodel.JobRole, imodel.JobRole]
}

func NewJobRoleRepository(client *mongo.Client, cfg *config.VaultConfig, database, collection string, logger utils.Logger) storage.JobRoleRepository {
	return &JobRoleStorage{
		client:     client,
		collection: client.Database(database).Collection(collection),
		logger:     logger,
		dal:        dal.NewMongoDal[imodel.JobRole, imodel.JobRole](client, cfg, database, collection),
	}
}

func (s *JobRoleStorage) Create(ctx context.Context, role *imodel.JobRole) error {
	s.logger.Infof("[JobRole/Create] creating job role: code=%s name=%s", role.Code, role.Name)
	role.ID = bson.NewObjectID()
	_, err := s.dal.InsertOne(ctx, *role)
	if err != nil {
		s.logger.Errorf("[JobRole/Create] failed to insert: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (s *JobRoleStorage) Update(ctx context.Context, id string, role *imodel.JobRole) error {
	s.logger.Infof("[JobRole/Update] id=%s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Errorf("[JobRole/Update] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := JobRoleMapper(*role)
	_, err = s.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		return local_util.HandleDBError(err)
	}
	return nil
}

func (s *JobRoleStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
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

func (s *JobRoleStorage) SoftDelete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
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

func (s *JobRoleStorage) FindByID(ctx context.Context, id string) (*imodel.JobRole, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID}
	res, err := s.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	return res, nil
}

func (s *JobRoleStorage) FindAll(ctx context.Context) (*[]imodel.JobRole, error) {

	data, err := s.dal.FindAll(ctx, bson.M{}, bson.M{})
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	return &data, nil
}

func (r *JobRoleStorage) ExistsMany(ctx context.Context, codes []string) (bool, error) {
	if len(codes) == 0 {
		return true, nil
	}

	count, err := r.collection.CountDocuments(ctx, bson.M{"code": bson.M{"$in": codes}})
	if err != nil {
		return false, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return count == int64(len(codes)), nil
}

func (s *JobRoleStorage) FindByCode(ctx context.Context, code string) (*imodel.JobRole, error) {
	filter := bson.M{
		"code": code,
	}
	res, err := s.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	return res, nil
}

func (s *JobRoleStorage) Find(ctx context.Context, filter bson.M) (*imodel.JobRole, error) {

	res, err := s.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	return res, nil
}

func (s *JobRoleStorage) FindByName(ctx context.Context, name string) (*imodel.JobRole, error) {
	filter := bson.M{
		"name": bson.M{
			"$regex":   "^" + strings.ToLower(name) + "$",
			"$options": "i",
		},
	}
	res, err := s.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	return res, nil
}

func (r *JobRoleStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.JobRole], error) {
	searchKeys := bson.M{}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"name": searchRegex},
			{"code": searchRegex},
		}
	}

	allowedKeys := []string{"enabled", "name", "code"}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["is_deleted"] = bson.M{"$ne": true}

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
	return &types.PaginatedResponse[[]model.JobRole]{
		Data: data,
		Meta: meta,
	}, nil
}
