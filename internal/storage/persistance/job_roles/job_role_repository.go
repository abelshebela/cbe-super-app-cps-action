package role_repo

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/job_roles/core"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type JobRoleRepository struct {
	client                *mongo.Client
	mongoDal              dal.MongoDal[imodel.JobRole, imodel.JobRole]
	logger                utils.Logger
	dbName                string
	cpsUserCollectionName string
	collection            *mongo.Collection
	jobCollection         *mongo.Collection
}

func NewJobRoleRepository(client *mongo.Client, cfg *config.VaultConfig, database string, collection []string, logger utils.Logger) storage.JobRoleRepository {
	return &JobRoleRepository{
		client:                client,
		mongoDal:              dal.NewMongoDal[imodel.JobRole, imodel.JobRole](client, cfg, database, collection[0]),
		logger:                logger,
		dbName:                database,
		collection:            client.Database(database).Collection(collection[1]),
		jobCollection:         client.Database(database).Collection(collection[0]),
		cpsUserCollectionName: collection[2],
	}
}

func (r *JobRoleRepository) Exists(ctx context.Context, id string) (bool, error) {
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

func (r *JobRoleRepository) ExistsMany(ctx context.Context, ids []string) (bool, error) {
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

func (r *JobRoleRepository) Create(ctx context.Context, role *imodel.JobRole) error {
	role.ID = bson.NewObjectID()
	_, err := r.mongoDal.InsertOne(ctx, *role)
	if err != nil {
		r.logger.Errorf("[RoleRepository][Create] failed to create role: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (r *JobRoleRepository) Update(ctx context.Context, id string, jobRole *imodel.JobRole) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("[RoleRepository][Update] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	update := bson.M{}
	if jobRole.JobTitle != "" {
		update["job_title"] = jobRole.JobTitle
	}
	if jobRole.Role != "" {
		update["role"] = jobRole.Role
	}
	if jobRole.Code != "" {
		update["code"] = jobRole.Code
	}
	if !jobRole.UpdateAt.IsZero() {
		update["updated_at"] = jobRole.UpdateAt
	}

	filter := bson.M{"_id": objID}

	// find by name so that we can update users with the new role name if it changes
	existingRole, err := r.FindByID(ctx, id)
	if err != nil {
		r.logger.Errorf("[RoleRepository][Update] failed to find role by name: %s err: %v", jobRole.JobTitle, err)
	}

	_, err = r.mongoDal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("[RoleRepository][Update] failed to update role: %v", err)
		return local_util.HandleDBError(err)
	}
	if existingRole != nil && existingRole.JobTitle != jobRole.JobTitle {
		r.logger.Infof("[RoleRepository][Update] updating cps users for job title change from %s to %s", existingRole.JobTitle, jobRole.JobTitle)
		if err := core.UpdateCpsUsers(ctx, r.dbName, r.cpsUserCollectionName, r.client, existingRole.JobTitle, jobRole.JobTitle); err != nil {
			r.logger.Errorf("[RoleRepository][Update] failed to update cps users for job title change: %v", err)
			// Not returning error since role update succeeded, and user update failure shouldn't block it
		}
	}
	return nil
}

func (r *JobRoleRepository) EnableOrDisable(ctx context.Context, id string, enable bool) error {
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

func (r *JobRoleRepository) SoftDelete(ctx context.Context, id string) error {
	r.logger.Infof("[RoleRepository][SoftDelete] soft deleting role with id: %s", id)
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

func (r *JobRoleRepository) FindByID(ctx context.Context, id string) (*imodel.JobRole, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("[RoleRepository][FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	// Try to find from roles collection first (without job_roles lookup)
	role, err := r.mongoDal.FindOne(ctx, bson.M{"_id": objID}, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Errorf("[RoleRepository][FindByID] role not found in roles collection")
			return nil, local_util.HandleDBError(mongo.ErrNoDocuments)
		}
		r.logger.Errorf("[RoleRepository][FindByID] find failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	// If we need the job role info (for display purposes), try to get it separately
	// But for delete operations, we only need the basic role info
	pipeline := roleWithJobRolePipeline(bson.M{"_id": objID}, r.collection.Name(), 0, 1)
	cursor, err := r.jobCollection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Warnf("[RoleRepository][FindByID] job_roles lookup failed, returning basic role info: %v", err)
		return role, nil
	}
	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		var resultWithJobRole imodel.JobRole
		if err := cursor.Decode(&resultWithJobRole); err != nil {
			r.logger.Warnf("[RoleRepository][FindByID] decode failed, returning basic role info: %v", err)
			return role, nil
		}
		return &resultWithJobRole, nil
	}

	// No job role found, but role exists - return basic role info
	r.logger.Warnf("[RoleRepository][FindByID] no job role found for role %s, returning basic role info", role.Role)
	return role, nil
}

func (r *JobRoleRepository) CheckIfExists(ctx context.Context, id string) (*imodel.JobRole, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("[RoleRepository][CheckIfExists] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	// Try to find from roles collection first (without job_roles lookup)
	role, err := r.mongoDal.FindOne(ctx, bson.M{"_id": objID}, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Errorf("[RoleRepository][CheckIfExists] role not found in roles collection")
			return nil, local_util.HandleDBError(mongo.ErrNoDocuments)
		}
		r.logger.Errorf("[RoleRepository][CheckIfExists] find failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return role, nil
}

func (r *JobRoleRepository) FindByName(ctx context.Context, name string) (*imodel.JobRole, error) {
	filter := bson.M{"job_title": bson.M{"$regex": "^" + name + "$", "$options": "i"}, "is_deleted": false}
	result, err := r.mongoDal.FindOne(ctx, filter, nil)
	if err != nil {
		r.logger.Errorf("[RoleRepository][FindByName] failed to find role: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}
func (r *JobRoleRepository) FindByRole(ctx context.Context, name string) (*imodel.JobRole, error) {
	filter := bson.M{"role": bson.M{"$regex": "^" + name + "$", "$options": "i"}, "is_deleted": false}
	result, err := r.mongoDal.FindOne(ctx, filter, nil)
	if err != nil {
		r.logger.Errorf("[RoleRepository][FindByRole] failed to find role: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}

func (r *JobRoleRepository) FindByCode(ctx context.Context, code string) (*imodel.JobRole, error) {
	filter := bson.M{"code": code, "is_deleted": false}
	result, err := r.mongoDal.FindOne(ctx, filter, nil)
	if err != nil {
		r.logger.Errorf("[RoleRepository][FindByCode] failed to find role: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}

// func (r *RoleRepository) FindAll(ctx context.Context) (*[]imodel.Role, error) {
// 	data, err := r.mongoDal.FindAll(ctx, bson.M{"enabled": true}, nil)
// 	if err != nil {
// 		r.logger.Errorf("[RoleRepository][FindAll] failed to find roles: %v", err)
// 		return nil, local_util.HandleDBError(err)
// 	}
// 	return &data, nil
// }

func (r *JobRoleRepository) FindAll(ctx context.Context) (*[]imodel.JobRole, error) {
	pipeline := roleWithJobRolePipeline(bson.M{"enabled": true, "is_deleted": false}, r.collection.Name(), 0, 0)
	cursor, err := r.jobCollection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("[RoleRepository][FindAll] failed to aggregate job_roles: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer cursor.Close(ctx)
	var data []imodel.JobRole
	if err := cursor.All(ctx, &data); err != nil {
		r.logger.Errorf("[RoleRepository][FindAll] failed to decode job_roles: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return &data, nil
}

func (r *JobRoleRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.JobRole], error) {
	searchKeys := bson.M{}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"job_title": searchRegex},
			{"role": searchRegex},
		}
	}

	allowedKeys := []string{"enabled", "job_title", "role"}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["is_deleted"] = false

	total, err := r.mongoDal.TotalCount(ctx, filter)
	if err != nil {
		r.logger.Errorf("[RoleRepository][FindAllWithPagination] failed to count roles: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	pipeline := roleWithJobRolePipeline(filter, r.collection.Name(), skip, limit)
	cursor, err := r.jobCollection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("[RoleRepository][FindAllWithPagination] aggregate failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer cursor.Close(ctx)

	var data []imodel.JobRole
	if err := cursor.All(ctx, &data); err != nil {
		r.logger.Errorf("[RoleRepository][FindAllWithPagination] decode failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]imodel.JobRole]{
		Data: data,
		Meta: meta,
	}, nil
}

func (r *JobRoleRepository) FindByFilterKey(ctx context.Context, field string, value string) (*imodel.JobRole, error) {
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

// roleWithJobRolePipeline matches roles, joins job_roles on roles.role == job_roles.code, and projects
// type, role_code, and role_name from job_roles. skip/limit are applied only when > 0 (pagination / single doc).
func roleWithJobRolePipeline(match bson.M, jobRolesCollName string, skip, limit int64) []bson.M {
	p := []bson.M{
		{"$match": match},
		{"$sort": bson.M{"created_at": -1}},
		{"$lookup": bson.M{
			"from":         jobRolesCollName,
			"localField":   "role",
			"foreignField": "code",
			"as":           "role_info",
			"pipeline": []bson.M{
				{"$project": bson.M{"code": 1, "name": 1, "type": 1, "_id": 0}},
			},
		}},
		{"$addFields": bson.M{
			"role_info": bson.M{"$arrayElemAt": []interface{}{"$role_info", 0}},
		}},
		{"$project": bson.M{
			"_id":        1,
			"code":       1,
			"job_title":  1,
			"role":       1,
			"enabled":    1,
			"updated_at": 1,
			"created_at": 1,
			"type":       bson.M{"$ifNull": []interface{}{"$role_info.type", ""}},
			"role_code":  bson.M{"$ifNull": []interface{}{"$role_info.code", ""}},
			"role_name":  bson.M{"$ifNull": []interface{}{"$role_info.name", ""}},
		}},
	}
	if skip > 0 {
		p = append(p, bson.M{"$skip": skip})
	}
	if limit > 0 {
		p = append(p, bson.M{"$limit": limit})
	}

	if skip > 0 || limit > 0 {
		p = append(p, bson.M{"$sort": bson.M{"created_at": -1}})
	}
	return p
}

// HasActiveJobRole implements [storage.JobRoleRepository].
func (s *JobRoleRepository) HasActiveJobRole(ctx context.Context, roleCode string) (bool, error) {
	s.logger.Infof("[RoleRepository][HasActiveJobRole] checking for active job role with code: %s", roleCode)
	filter := bson.M{"role": roleCode, "enabled": true}
	_, err := s.mongoDal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.logger.Infof("[RoleRepository][HasActiveJobRole] no active job role found with code: %s", roleCode)
			return false, nil
		}
		s.logger.Errorf("[RoleRepository][HasActiveJobRole] failed to check for active job role: %v", err)
		return false, local_util.HandleDBError(err)
	}
	s.logger.Infof("[RoleRepository][HasActiveJobRole] active job role found with code: %s", roleCode)
	return true, nil
}
