package cpsroles

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type cpsRoleStorage struct {
	cfg           *config.VaultConfig
	dal           dal.MongoDal[imodel.CPSRoles, imodel.CPSRoles]
	client        *mongo.Client
	dbName        string
	collection    string
	kafkaProducer kafka.ClientOrchestrationProducer
	logger        utils.Logger
}

func NewCPSRolesStorage(client *mongo.Client, cfg *config.VaultConfig, dbName, collection string, kafkaProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.CPSRolesRepository {
	return &cpsRoleStorage{
		cfg:           cfg,
		dal:           dal.NewMongoDal[imodel.CPSRoles, imodel.CPSRoles](client, cfg, dbName, collection),
		client:        client,
		dbName:        dbName,
		collection:    collection,
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

func (m *cpsRoleStorage) Create(ctx context.Context, req imodel.CPSRoles) error {
	_, err := m.dal.InsertOne(ctx, req)
	if err != nil {
		m.logger.Errorf("[Create] failed to create cps role: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (m *cpsRoleStorage) Update(ctx context.Context, id string, req imodel.CPSRoles) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("[Update] invalid id format: %s, error: %v", id, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"name":       req.Name,
		"enabled":    req.Enabled,
		"updated_at": time.Now(),
	}

	updatedCPSRole, err := m.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Warnf("[Update] cps role not found, id: %s", id)
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		m.logger.Errorf("[Update] failed to update cps role, id: %s, error: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	m.kafkaProducer.PublishMessage(
		ctx,
		updatedCPSRole,
		string(constants.ClientOrchestrationServicesTopic),
		"cps-customer-role-updated",
		"cps customer role updated",
	)

	return nil
}

func (m *cpsRoleStorage) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]imodel.CPSRoles], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"search", "enabled"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"name": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(*filterParam, searchKeys, allowedKeys)

	data, err := m.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		m.logger.Errorf("[FindAllWithPagination] failed to fetch paginated cps roles: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter["is_deleted"] = false
	total, err := m.dal.TotalCount(ctx, filter)
	if err != nil {
		m.logger.Errorf("[FindAllWithPagination] failed to count total cps roles: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]imodel.CPSRoles]{
		Data: data,
		Meta: meta,
	}, nil
}

func (m *cpsRoleStorage) FindById(ctx context.Context, id string) (*imodel.CPSRoles, error) {
	m.logger.Infof("[FindById] fetching cps role by id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("[FindById] invalid id format: %s, error: %v", id, err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	result, err := m.dal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Warnf("[FindById] cps role not found, id: %s", id)
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		m.logger.Errorf("[FindById] failed to find cps role, id: %s, error: %v", id, err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	m.logger.Infof("[FindById] cps role retrieved successfully, id: %s, result: %v", id, result)

	// Use MongoDB pipeline to group action_names by maker/checker/auditor index
	approverCol := m.client.Database(m.dbName).Collection("cps_action_approver_index")
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{"role_id", result.RoleCode}}}},
		bson.D{{Key: "$facet", Value: bson.M{
			"maker": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.D{{"maker_index", bson.D{{"$ne", nil}}}}}},
				bson.D{{Key: "$group", Value: bson.D{{"_id", nil}, {"actions", bson.D{{"$addToSet", "$action_name"}}}}}},
			},
			"checker": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.D{{"checker_index", bson.D{{"$ne", nil}}}}}},
				bson.D{{Key: "$group", Value: bson.D{{"_id", nil}, {"actions", bson.D{{"$addToSet", "$action_name"}}}}}},
			},
			"auditor": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.D{{"auditor_index", bson.D{{"$ne", nil}}}}}},
				bson.D{{Key: "$group", Value: bson.D{{"_id", nil}, {"actions", bson.D{{"$addToSet", "$action_name"}}}}}},
			},
		}}},
	}
	cursor, err := approverCol.Aggregate(ctx, pipeline)
	if err != nil {
		m.logger.Errorf("Error aggregating cps_action_approver_index: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	type actionGroup struct {
		Actions []string `bson:"actions"`
	}

	type facetActions struct {
		Maker   []actionGroup `bson:"maker"`
		Checker []actionGroup `bson:"checker"`
		Auditor []actionGroup `bson:"auditor"`
	}
	var aggResult []facetActions

	if err := cursor.All(ctx, &aggResult); err != nil {
		m.logger.Errorf("Error decoding facet result: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	m.logger.Infof("Facet aggregation result: %v", aggResult)

	var maker, checker, auditor []string
	if len(aggResult) > 0 {
		if len(aggResult[0].Maker) > 0 {
			maker = aggResult[0].Maker[0].Actions
		}
		if len(aggResult[0].Checker) > 0 {
			checker = aggResult[0].Checker[0].Actions
		}
		if len(aggResult[0].Auditor) > 0 {
			auditor = aggResult[0].Auditor[0].Actions
		}
	} else {
		m.logger.Warnf("No approver index found for role_id: %s", result.RoleCode)
	}

	result.MakerActions = maker
	result.CheckerActions = checker
	result.AuditorActions = auditor

	// Fetch enabled/disabled service lists for this role
	enabledServices, disabledServices, err := m.fetchServiceAccessLists(ctx, result.ID)
	if err != nil {
		m.logger.Warnf("[FindById] failed to fetch service access lists for role %s: %v", id, err)
	}
	result.EnabledServices = enabledServices
	result.DisabledServices = disabledServices

	return result, nil
}

// fetchServiceAccessLists retrieves enabled and disabled service lists for a CPS role.
//   - Enabled: services in access_list (enabled=true) NOT blocked for this role
//     (i.e., not in access_list_segmentation, or in segmentation with enabled=false)
//   - Disabled: services in access_list (enabled=true) that ARE blocked for this role
//     (i.e., in access_list_segmentation with enabled=true)
func (m *cpsRoleStorage) fetchServiceAccessLists(ctx context.Context, roleID bson.ObjectID) ([]imodel.ServiceAccessInfo, []imodel.ServiceAccessInfo, error) {
	// 1. Fetch all globally enabled services from access_list
	accessListCol := m.client.Database(m.dbName).Collection("access_list")
	alCursor, err := accessListCol.Find(ctx, bson.M{"enabled": true})
	if err != nil {
		return nil, nil, err
	}
	defer alCursor.Close(ctx)

	type accessListDoc struct {
		Key            string `bson:"key"`
		AccessListName string `bson:"access_list_name"`
	}
	var allServices []accessListDoc
	if err := alCursor.All(ctx, &allServices); err != nil {
		return nil, nil, err
	}

	// 2. Fetch segmentation entries for this role from access_list_segmentation
	segCol := m.client.Database(m.dbName).Collection("access_list_segmentation")
	segCursor, err := segCol.Find(ctx, bson.M{"segmented_id": roleID})
	if err != nil {
		return nil, nil, err
	}
	defer segCursor.Close(ctx)

	type segDoc struct {
		AccessListKey string `bson:"access_list_key"`
		Enabled       bool   `bson:"enabled"`
	}
	var segEntries []segDoc
	if err := segCursor.All(ctx, &segEntries); err != nil {
		return nil, nil, err
	}

	// 3. Build a map of blocked keys (segmentation enabled=true means disabled for this role)
	blockedKeys := make(map[string]bool, len(segEntries))
	for _, seg := range segEntries {
		if seg.Enabled {
			blockedKeys[seg.AccessListKey] = true
		}
	}

	// 4. Categorize services
	var enabledServices, disabledServices []imodel.ServiceAccessInfo
	for _, svc := range allServices {
		info := imodel.ServiceAccessInfo{
			Key:            svc.Key,
			AccessListName: svc.AccessListName,
		}
		if blockedKeys[svc.Key] {
			disabledServices = append(disabledServices, info)
		} else {
			enabledServices = append(enabledServices, info)
		}
	}

	return enabledServices, disabledServices, nil
}

func (m *cpsRoleStorage) FindByNameOrRoleCode(ctx context.Context, name, roleCode string) (*imodel.CPSRoles, error) {
	filter := []bson.M{}
	if roleCode == "" {
		filter = append(filter, bson.M{"name": bson.M{"$regex": "^" + name, "$options": "i"}})
	}
	if name == "" {
		filter = append(filter, bson.M{"role_code": bson.M{"$regex": "^" + roleCode, "$options": "i"}})
	}
	orFilter := bson.M{"$or": filter}
	result, err := m.dal.FindOne(ctx, orFilter, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Warnf("[FindByNameOrRoleCode] cps role not found, name: %s, roleCode: %s", name, roleCode)
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		m.logger.Errorf("[FindByNameOrRoleCode] failed to find cps role, name: %s, roleCode: %s, error: %v", name, roleCode, err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return result, nil
}

func (m *cpsRoleStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("[EnableOrDisable] invalid id format: %s, error: %v", id, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"enabled": enable, "updated_at": time.Now()}

	_, err = m.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Warnf("[EnableOrDisable] cps role not found, id: %s", id)
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		m.logger.Errorf("[EnableOrDisable] failed to enable/disable cps role, id: %s, error: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

// EnableServiceAccess removes the block for the given access list keys on this role.
// It deletes segmentation entries with enabled=true, or sets them to enabled=false.
func (m *cpsRoleStorage) EnableServiceAccess(ctx context.Context, roleID string, accessListKeys []string) error {
	objID, err := bson.ObjectIDFromHex(roleID)
	if err != nil {
		m.logger.Errorf("[EnableServiceAccess] invalid role id: %s, error: %v", roleID, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	segCol := m.client.Database(m.dbName).Collection("access_list_segmentation")
	filter := bson.M{
		"segmented_id":    objID,
		"access_list_key": bson.M{"$in": accessListKeys},
		"enabled":         true,
	}
	update := bson.M{"$set": bson.M{"enabled": false, "updated_at": time.Now()}}

	_, err = segCol.UpdateMany(ctx, filter, update)
	if err != nil {
		m.logger.Errorf("[EnableServiceAccess] failed to enable services for role %s: %v", roleID, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

// DisableServiceAccess blocks the given access list keys for this role.
// It upserts segmentation entries with enabled=true.
func (m *cpsRoleStorage) DisableServiceAccess(ctx context.Context, roleID string, accessListKeys []string) error {
	objID, err := bson.ObjectIDFromHex(roleID)
	if err != nil {
		m.logger.Errorf("[DisableServiceAccess] invalid role id: %s, error: %v", roleID, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	// Fetch the role to get name for segmentation_name
	role, err := m.dal.FindOne(ctx, bson.M{"_id": objID}, bson.M{})
	if err != nil {
		m.logger.Errorf("[DisableServiceAccess] role not found: %s, error: %v", roleID, err)
		return errors.New(localization.ErrorResourceNotFound.Code)
	}

	// Fetch access list names for the given keys
	accessListCol := m.client.Database(m.dbName).Collection("access_list")
	alCursor, err := accessListCol.Find(ctx, bson.M{"key": bson.M{"$in": accessListKeys}})
	if err != nil {
		m.logger.Errorf("[DisableServiceAccess] failed to fetch access list: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer alCursor.Close(ctx)

	type alDoc struct {
		Key            string `bson:"key"`
		AccessListName string `bson:"access_list_name"`
	}
	var alDocs []alDoc
	if err := alCursor.All(ctx, &alDocs); err != nil {
		m.logger.Errorf("[DisableServiceAccess] failed to decode access list: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	alNameMap := make(map[string]string, len(alDocs))
	for _, doc := range alDocs {
		alNameMap[doc.Key] = doc.AccessListName
	}

	segCol := m.client.Database(m.dbName).Collection("access_list_segmentation")
	for _, key := range accessListKeys {
		filter := bson.M{
			"segmented_id":    objID,
			"access_list_key": key,
		}
		update := bson.M{
			"$set": bson.M{
				"enabled":    true,
				"updated_at": time.Now(),
			},
			"$setOnInsert": bson.M{
				"_id":               bson.NewObjectID(),
				"type":              "block",
				"segmentation_type": "cps_role",
				"access_list_key":   key,
				"access_list_name":  alNameMap[key],
				"segmented_id":      objID,
				"segmentation_code": role.RoleCode,
				"segmentation_name": role.Name,
				"created_at":        time.Now(),
			},
		}
		opts := options.UpdateOne().SetUpsert(true)
		_, err := segCol.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			m.logger.Errorf("[DisableServiceAccess] failed to upsert segmentation for key %s: %v", key, err)
			return errors.New(localization.ErrorUnexpectedError.Code)
		}
	}

	return nil
}

func (m *cpsRoleStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("[Delete] invalid id format: %s, error: %v", id, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"is_deleted": true, "deleted_at": time.Now()}

	_, err = m.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Warnf("[Delete] cps role not found, id: %s", id)
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		m.logger.Errorf("[Delete] failed to soft delete cps role, id: %s, error: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (m *cpsRoleStorage) FindByCustomerSegmentation(ctx context.Context, customerSegment string) (*imodel.CPSRoles, error) {
	filter := bson.M{"name": customerSegment, "enabled": true}
	result, err := m.dal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Warnf("[FindByCustomerSegmentation] cps role not found for customer segment: %s", customerSegment)
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		m.logger.Errorf("[FindByCustomerSegmentation] failed to find cps role for customer segment: %s, error: %v", customerSegment, err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return result, nil
}
