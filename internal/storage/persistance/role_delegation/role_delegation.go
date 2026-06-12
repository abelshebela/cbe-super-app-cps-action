package role_delegation_repo

import (
	cpsuser "cbe-super-app-cps-action/internal/constants/dto/cps_user"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"time"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type roleDelegationRepository struct {
	repo              dal.MongoDal[imodel.RoleDelegation, imodel.RoleDelegation]
	cpsRepo           dal.MongoDal[imodel.CPSUser, imodel.CPSUser]
	logger            utils.Logger
	collection        *mongo.Collection
	cpsUserCollection *mongo.Collection
	bpsUserCollection *mongo.Collection
}

// Create implements [storage.RoleDelegationRepository].
func (r *roleDelegationRepository) CreateWithExistingUser(ctx context.Context, role *imodel.RoleDelegation) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[RoleDelegationRepository][Create] creating role delegation for user=%s job_title=%s", role.DelegatedUserID, role.DelegatedUserJobTitle)
	session, err := r.collection.Database().Client().StartSession()
	if err != nil {
		log.Errorf("[RoleDelegationRepository][Create] failed to start session: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sc context.Context) (any, error) {
		role.ID = bson.NewObjectID()
		role.Enable = true
		result, err := r.collection.InsertOne(sc, role)
		if err != nil {
			log.Errorf("[RoleDelegationRepository][Create] failed to create role delegation: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		insertedID, ok := result.InsertedID.(bson.ObjectID)
		if !ok {
			log.Errorf("[RoleDelegationRepository][Create] inserted id has unexpected type: %T", result.InsertedID)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		updatePipeline := mongo.Pipeline{
			bson.D{{Key: "$set", Value: bson.M{"is_delegation_active": true, "delegation_id": insertedID}}},
		}

		if role.DelegatedUserUserType == "CPS" {
			_, err = r.cpsUserCollection.UpdateOne(sc, bson.M{"username": role.DelegatedUserID, "is_deleted": false}, updatePipeline)
			if err != nil {
				log.Errorf("[RoleDelegationRepository][Create] failed to update cps user delegate: %v", err)
				return nil, local_util.HandleDBError(err)
			}
		} else {
			_, err = r.bpsUserCollection.UpdateOne(sc, bson.M{"username": role.DelegatedUserID, "is_deleted": false}, updatePipeline)
			if err != nil {
				log.Errorf("[RoleDelegationRepository][Create] failed to update cps user delegate: %v", err)
				return nil, local_util.HandleDBError(err)
			}
		}

		if role.RevokeExistingDelegation {
			_, err = r.collection.UpdateMany(sc, bson.M{"delegated_user_id": role.DelegatedUserID, "_id": bson.M{"$ne": insertedID}}, bson.M{"$set": bson.M{"enable": false}})
			if err != nil {
				log.Errorf("[RoleDelegationRepository][Create] failed to disable existing delegations: %v", err)
				return nil, local_util.HandleDBError(err)
			}
		}

		return nil, nil
	})
	if err != nil {
		r.logger.Errorf("[RoleDelegationRepository][Create] transaction failed: %v", err)
		return localization.ErrorUnexpectedError
	}

	log.Infof("[RoleDelegationRepository][Create] role delegation created successfully")
	return nil
}

// Create implements [storage.RoleDelegationRepository].
func (r *roleDelegationRepository) CreateWithNewUser(ctx context.Context, role *imodel.RoleDelegation) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[RoleDelegationRepository][Create] creating role delegation for user=%s job_title=%s", role.DelegatedUserID, role.DelegatedUserJobTitle)
	session, err := r.collection.Database().Client().StartSession()
	if err != nil {
		log.Errorf("[RoleDelegationRepository][Create] failed to start session: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sc context.Context) (any, error) {
		role.ID = bson.NewObjectID()
		role.Enable = true
		if role.DelegatedUserUserType == "CPS" {
			role.DelegatedUserUserCode = local_util.GenerateCPSUserCode()
		} else {
			role.DelegatedUserUserCode = local_util.GenerateBPSUserCode()
		}
		result, err := r.collection.InsertOne(sc, role)
		if err != nil {
			log.Errorf("[RoleDelegationRepository][Create] failed to create role delegation: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		insertedID, ok := result.InsertedID.(bson.ObjectID)
		if !ok {
			log.Errorf("[RoleDelegationRepository][Create] inserted id has unexpected type: %T", result.InsertedID)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		now := time.Now()
		if role.DelegatedUserUserType == "CPS" {
			assignedDepartment, err := bson.ObjectIDFromHex(role.DelegatedUserDepartmentOrBranch)
			if err != nil {
				log.Errorf("[RoleDelegationRepository][Create] invalid department id: %v", err)
				return nil, localization.ErrorUnexpectedError
			}
			_, err = r.cpsUserCollection.InsertOne(sc, model.CPSUser{
				UserCode:           role.DelegatedUserUserCode,
				FullName:           role.DelegatedUserFullName,
				Role:               role.DelegatedUserExistingRole,
				PhoneNumber:        role.DelegatedUserPhoneNumber,
				Email:              role.DelegatedUserEmail,
				UserName:           role.DelegatedUserID,
				JobTitle:           role.DelegatedUserJobTitle,
				Department:         assignedDepartment,
				Enabled:            true,
				DateJoined:         &now,
				CreatedAt:          now,
				IsDelegationActive: true,
				DelegationID:       insertedID,
			})
			if err != nil {
				log.Errorf("[RoleDelegationRepository][Create] failed to update cps user delegate: %v", err)
				return nil, local_util.HandleDBError(err)
			}
		} else {
			_, err = r.bpsUserCollection.InsertOne(sc, model.BPSUser{
				UserCode:    role.DelegatedUserUserCode,
				FullName:    role.DelegatedUserFullName,
				Role:        role.DelegatedUserExistingRole,
				PhoneNumber: role.DelegatedUserPhoneNumber,
				Email:       role.DelegatedUserEmail,
				UserName:    role.DelegatedUserID,
				JobTitle:    role.DelegatedUserJobTitle,
				BranchCode:  []string{role.DelegatedUserDepartmentOrBranch},
				Enabled:     true,
				CreatedAt:   now,
			})
			if err != nil {
				log.Errorf("[RoleDelegationRepository][Create] failed to update bps user delegate: %v", err)
				return nil, local_util.HandleDBError(err)
			}
		}

		return nil, nil
	})
	if err != nil {
		return err
	}

	log.Infof("[RoleDelegationRepository][Create] role delegation created successfully")
	return nil
}

// EnableOrDisable implements [storage.RoleDelegationRepository].
func (r *roleDelegationRepository) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[RoleDelegationRepository][EnableOrDisable] updating role delegation id=%s enabled=%v", id, enable)
	objID, err := local_util.ParseObjectID(id)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][EnableOrDisable] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"enable": enable, "updated_at": time.Now()}

	_, err = r.repo.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][EnableOrDisable] failed to update role delegation: %v", err)
		return local_util.HandleDBError(err)
	}

	log.Infof("[RoleDelegationRepository][EnableOrDisable] role delegation updated successfully")
	return nil
}

// Find implements [storage.RoleDelegationRepository].
func (r *roleDelegationRepository) Find(ctx context.Context, filter bson.M) (*imodel.RoleDelegation, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	result, err := r.repo.FindOne(ctx, filter, nil)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][Find] failed to find role delegation: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}

// FindAll implements [storage.RoleDelegationRepository].
func (r *roleDelegationRepository) FindAll(ctx context.Context) (*[]imodel.RoleDelegation, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	data, err := r.repo.FindAllWithCursorBasedPagination(ctx, dal.FilterOp{})
	if err != nil {
		log.Errorf("[RoleDelegationRepository][FindAll] failed to fetch role delegations: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	return &data, nil
}

// FindAllWithPagination implements [storage.RoleDelegationRepository].
func (r *roleDelegationRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.RoleDelegation], error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	searchKeys := bson.M{}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"delegated_user_job_title": searchRegex},
			{"delegated_user_id": searchRegex},
			{"delegated_user_email": searchRegex},
			{"delegated_user_phone_number": searchRegex},
			{"delegated_user_full_name": searchRegex},
		}
	}

	allowedKeys := []string{"enable", "delegated_user_email", "delegated_user_phone_number", "delegated_user_job_title", "delegated_user_id", "delegated_user_full_name", "start_at", "end_at"}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	if val, ok := filterParam.Filters["portal"]; ok {
		if val == "BPS" || val == "CPS" {
			filter["delegation_type"] = val
		}
	}

	if val, ok := filterParam.Filters["status"]; ok {
		if val == "Expired" {
			filter["enable"] = true
			filter["end_at"] = bson.M{"$lte": time.Now()}
		}
		if val == "Revoked" {
			filter["enable"] = false
		}
		if val == "Active" {
			filter["enable"] = true
			filter["start_at"] = bson.M{"$lte": time.Now()}
			filter["end_at"] = bson.M{"$gt": time.Now()}
		}
		if val == "Suspended" {
			filter["enable"] = true
			filter["start_at"] = bson.M{"$gt": time.Now()}
			filter["end_at"] = bson.M{"$gt": time.Now()}
		}
	}

	total, err := r.repo.TotalCount(ctx, filter)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][FindByUsername] failed to count role delegations: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	if total <= 0 || skip >= total {
		return &types.PaginatedResponse[[]imodel.RoleDelegation]{}, nil
	}

	data, err := r.repo.FindAllWithPaginationD(ctx, dal.FilterParam{
		Filter:     filter,
		Projection: bson.M{},
		Sort:       bson.D{bson.E{Key: "created_at", Value: -1}},
		Skip:       skip,
		Limit:      limit,
	})
	if err != nil {
		log.Errorf("[RoleDelegationRepository][FindByUsername] failed to fetch role delegations: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]imodel.RoleDelegation]{
		Data: data,
		Meta: meta,
	}, nil
}

func (r *roleDelegationRepository) FindByUsername(ctx context.Context, id string, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.RoleDelegation], error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	allowedKeys := []string{}
	filter, skip, limit := lib.FilterBuilder(filterParam, bson.M{}, allowedKeys)
	filter["delegated_user_id"] = id

	total, err := r.repo.TotalCount(ctx, filter)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][FindAllWithPagination] failed to count role delegations: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	if total <= 0 || skip >= total {
		return &types.PaginatedResponse[[]imodel.RoleDelegation]{}, nil
	}

	data, err := r.repo.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][FindAllWithPagination] failed to fetch role delegations: %v", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &types.PaginatedResponse[[]imodel.RoleDelegation]{}, localization.ErrRoleDelegationForUserNotFound
		}
		return nil, localization.ErrorUnexpectedError
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]imodel.RoleDelegation]{
		Data: data,
		Meta: meta,
	}, nil
}

func (r *roleDelegationRepository) FindByUserCode(ctx context.Context, usercode string) (*cpsuser.CpsUserPopulatedResponse, error) {
	r.logger.Infof("[RoleDelegationRepository][FindByUsername] fetching active delegation for user code: %s", usercode)
	log := local_util.LoggerFromCtx(ctx, r.logger)
	currentTime := time.Now()

	data, err := r.repo.FindOne(ctx, bson.M{
		"delegated_user_user_code": usercode,
		"start_at":                 bson.M{"$lte": currentTime},
		"end_at":                   bson.M{"$gt": currentTime},
		"enable":                   true,
	}, bson.M{})
	if err != nil {
		log.Errorf("[RoleDelegationRepository][FindByUsername] failed to fetch role delegations: %v", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, localization.ErrorUserNotFound
		}
		return nil, localization.ErrorUnexpectedError
	}

	return &cpsuser.CpsUserPopulatedResponse{
		ID:       data.ID,
		UserCode: data.DelegatedUserUserCode,
		FullName: data.DelegatedUserFullName,
		Role: cpsuser.RoleResponse{
			Code: data.NewRoleID,
		},
		Department: &cpsuser.DepartmentResponse{
			Name: data.NewDepartmentOrBranch,
		},
		PhoneNumber:        data.DelegatedUserPhoneNumber,
		Email:              data.DelegatedUserEmail,
		UserName:           data.DelegatedUserID,
		Realm:              data.DelegatedUserUserType,
		IsDelegationActive: true,
		DelegatedRole:      data.NewRoleID,
	}, nil
}

// FindByID implements [storage.RoleDelegationRepository].
func (r *roleDelegationRepository) FindByID(ctx context.Context, id string) (*imodel.RoleDelegation, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	objID, err := local_util.ParseObjectID(id)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	result, err := r.repo.FindOne(ctx, bson.M{"_id": objID}, nil)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][FindByID] failed to find role delegation: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}

// Update implements [storage.RoleDelegationRepository].
func (r *roleDelegationRepository) Update(ctx context.Context, id string, role *imodel.RoleDelegation) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[RoleDelegationRepository][Update] updating role delegation id=%s", id)
	objID, err := local_util.ParseObjectID(id)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][Update] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	update := bson.M{
		"start_at":   role.StartAt,
		"end_at":     role.EndAt,
		"updated_at": time.Now(),
	}

	_, err = r.repo.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][Update] failed to update role delegation: %v", err)
		return local_util.HandleDBError(err)
	}

	return nil
}

func (r *roleDelegationRepository) CheckIfDelegationAlreadyExists(ctx context.Context, userID string, start, end time.Time) (bool, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	filter := bson.M{
		"delegated_user_id": userID,
		"enabled":           true,
		"$and": []bson.M{
			{
				"start_at": bson.M{"$lte": start},
				"end_at":   bson.M{"$gte": start},
			},
		},
	}
	count, err := r.repo.TotalCount(ctx, filter)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][CheckIfDelegationAlreadyExists] failed to check existing delegation: %v", err)
		return false, local_util.HandleDBError(err)
	}
	return count > 0, nil
}

func (r *roleDelegationRepository) Delete(ctx context.Context, userID string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	objID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][Delete] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{
		"_id": objID,
	}

	err = r.repo.DeleteOneH(ctx, filter)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][Delete] failed to delete delegations: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func NewRoleDelegationRepository(client *mongo.Client, cfg *config.VaultConfig, database string, collection []string, logger utils.Logger) storage.RoleDelegationRepository {
	return &roleDelegationRepository{
		repo:              dal.NewMongoDal[imodel.RoleDelegation, imodel.RoleDelegation](client, cfg, database, collection[0]),
		cpsRepo:           dal.NewMongoDal[imodel.CPSUser, imodel.CPSUser](client, cfg, database, collection[1]),
		logger:            logger,
		collection:        client.Database(database).Collection(collection[0]),
		cpsUserCollection: client.Database(database).Collection(collection[1]),
		bpsUserCollection: client.Database(database).Collection(collection[2]),
	}
}
