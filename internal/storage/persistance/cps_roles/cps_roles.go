package cpsroles

// import (
// 	"cbe-super-app-cps-action/internal/constants"
// 	"cbe-super-app-cps-action/internal/constants/lib"
// 	"cbe-super-app-cps-action/internal/constants/localization"
// 	"cbe-super-app-cps-action/internal/constants/types"
// 	"cbe-super-app-cps-action/internal/storage"
// 	"cbe-super-app-cps-action/internal/storage/kafka"
// 	local_util "cbe-super-app-cps-action/pkgs/utils"
// 	"context"
// 	"regexp"
// 	"strings"

// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
// 	"go.mongodb.org/mongo-driver/v2/bson"

// 	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
// 	imodel "cbe-super-app-cps-action/internal/constants/model"
// 	"errors"
// 	"time"

// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
// 	"go.mongodb.org/mongo-driver/v2/mongo"
// 	"go.mongodb.org/mongo-driver/v2/mongo/options"
// )

// type cpsRoleStorage struct {
// 	cfg                     *config.VaultConfig
// 	dal                     dal.MongoDal[imodel.CPSRoles, imodel.CPSRoles]
// 	accessDal               dal.MongoDal[model.APPAccessList, model.APPAccessList]
// 	client                  *mongo.Client
// 	dbName                  string
// 	collection              string
// 	accessListCollection    string
// 	accessListSegCollection string
// 	kafkaProducer           kafka.ClientOrchestrationProducer
// 	logger                  utils.Logger
// }

// func NewCPSRolesStorage(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection []string, kafkaProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.CPSRolesRepository {
// 	return &cpsRoleStorage{
// 		cfg:                     cfg,
// 		dal:                     dal.NewMongoDal[imodel.CPSRoles, imodel.CPSRoles](client, cfg, dbName, collection[0]),
// 		client:                  client,
// 		dbName:                  dbName,
// 		collection:              collection[0],
// 		accessListCollection:    collection[1],
// 		accessListSegCollection: collection[2],
// 		kafkaProducer:           kafkaProducer,
// 		logger:                  logger,
// 	}
// }

// func (m *cpsRoleStorage) Create(ctx context.Context, req imodel.CPSRoles) error {
// 	req.Name = strings.ToUpper(req.Name)
// 	_, err := m.dal.InsertOne(ctx, req)
// 	if err != nil {
// 		m.logger.Errorf("[CPSRolesStorage][Create] failed to create cps role: %v", err)
// 		return errors.New(localization.ErrorUnexpectedError.Code)
// 	}
// 	return nil
// }

// func (m *cpsRoleStorage) Update(ctx context.Context, id string, req imodel.CPSRoles) error {
// 	objID, err := bson.ObjectIDFromHex(id)
// 	if err != nil {
// 		m.logger.Errorf("[CPSRolesStorage][Update] invalid id format: %s, error: %v", id, err)
// 		return errors.New(localization.ErrorInvalidID.Code)
// 	}

// 	filter := bson.M{"_id": objID}
// 	update := bson.M{
// 		"updated_at": time.Now(),
// 	}

// 	if req.Name != "" {
// 		update["name"] = strings.ToUpper(req.Name)
// 	}
// 	if req.RoleCode != "" {
// 		update["role_code"] = req.RoleCode
// 	}
// 	if req.Description != "" {
// 		update["description"] = req.Description
// 	}

// 	updatedCPSRole, err := m.dal.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		return local_util.HandleDBError(err)
// 	}

// 	m.kafkaProducer.PublishMessage(
// 		ctx,
// 		updatedCPSRole,
// 		string(constants.ClientOrchestrationServicesTopic),
// 		string(constants.CustomerRoleUpdatedTopic),
// 		"cps customer role updated",
// 	)

// 	return nil
// }

// func (m *cpsRoleStorage) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]imodel.CPSRoles], error) {
// 	searchKeys := bson.M{}
// 	allowedKeys := []string{"search", "enabled"}

// 	if filterParam.Search != "" {
// 		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
// 		searchKeys["name"] = searchRegex
// 	}

// 	filter, skip, limit := lib.FilterBuilder(*filterParam, searchKeys, allowedKeys)

// 	filter["is_deleted"] = false

// 	// Use aggregation to apply sort before skip/limit so pagination reflects ordering
// 	col := m.client.Database(m.dbName).Collection(m.collection)
// 	pipeline := mongo.Pipeline{
// 		bson.D{{Key: "$match", Value: filter}},
// 		bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
// 		bson.D{{Key: "$facet", Value: bson.M{
// 			"data":  []bson.D{{{Key: "$skip", Value: skip}}, {{Key: "$limit", Value: limit}}},
// 			"total": []bson.D{{{Key: "$count", Value: "count"}}},
// 		}}},
// 	}

// 	cursor, err := col.Aggregate(ctx, pipeline)
// 	if err != nil {
// 		m.logger.Errorf("[CPSRolesStorage][FindAllWithPagination] aggregation error: %v", err)
// 		return nil, errors.New(localization.ErrorUnexpectedError.Code)
// 	}
// 	defer cursor.Close(ctx)

// 	var aggResult []struct {
// 		Data  []imodel.CPSRoles `bson:"data"`
// 		Total []struct {
// 			Count int64 `bson:"count"`
// 		} `bson:"total"`
// 	}
// 	if err := cursor.All(ctx, &aggResult); err != nil {
// 		m.logger.Errorf("[CPSRolesStorage][FindAllWithPagination] failed to decode aggregation result: %v", err)
// 		return nil, errors.New(localization.ErrorUnexpectedError.Code)
// 	}

// 	var data []imodel.CPSRoles
// 	var total int64
// 	if len(aggResult) > 0 {
// 		data = aggResult[0].Data
// 		if len(aggResult[0].Total) > 0 {
// 			total = aggResult[0].Total[0].Count
// 		}
// 	}

// 	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

// 	return &types.PaginatedResponse[[]imodel.CPSRoles]{
// 		Data: data,
// 		Meta: meta,
// 	}, nil
// }

// func (m *cpsRoleStorage) FindById(ctx context.Context, id string) (*imodel.CPSRoles, error) {
// 	m.logger.Infof("[CPSRolesStorage][FindById] fetching cps role by id: %s", id)
// 	objID, err := bson.ObjectIDFromHex(id)
// 	if err != nil {
// 		m.logger.Errorf("[CPSRolesStorage][FindById] invalid id format: %s, error: %v", id, err)
// 		return nil, errors.New(localization.ErrorInvalidID.Code)
// 	}

// 	filter := bson.M{"_id": objID, "is_deleted": false}
// 	result, err := m.dal.FindOne(ctx, filter, bson.M{})
// 	if err != nil {
// 		m.logger.Errorf("[CPSRolesStorage][FindById] failed to find cps role, id: %s, error: %v", id, err)
// 		return nil, local_util.HandleDBError(err)
// 	}
// 	m.logger.Infof("[CPSRolesStorage][FindById] cps role retrieved successfully, id: %s, result: %v", id, result)

// 	// Use MongoDB pipeline to group action_names by maker/checker/auditor index
// 	approverCol := m.client.Database(m.dbName).Collection("cps_action_approver_index")
// 	pipeline := mongo.Pipeline{
// 		bson.D{{Key: "$match", Value: bson.D{{"role_id", result.RoleCode}}}},
// 		bson.D{{Key: "$facet", Value: bson.M{
// 			"maker": mongo.Pipeline{
// 				bson.D{{Key: "$match", Value: bson.D{{"maker_index", bson.D{{"$ne", nil}}}}}},
// 				bson.D{{Key: "$group", Value: bson.D{{"_id", nil}, {"actions", bson.D{{"$addToSet", "$action_name"}}}}}},
// 			},
// 			"checker": mongo.Pipeline{
// 				bson.D{{Key: "$match", Value: bson.D{{"checker_index", bson.D{{"$ne", nil}}}}}},
// 				bson.D{{Key: "$group", Value: bson.D{{"_id", nil}, {"actions", bson.D{{"$addToSet", "$action_name"}}}}}},
// 			},
// 			"auditor": mongo.Pipeline{
// 				bson.D{{Key: "$match", Value: bson.D{{"auditor_index", bson.D{{"$ne", nil}}}}}},
// 				bson.D{{Key: "$group", Value: bson.D{{"_id", nil}, {"actions", bson.D{{"$addToSet", "$action_name"}}}}}},
// 			},
// 		}}},
// 	}
// 	cursor, err := approverCol.Aggregate(ctx, pipeline)
// 	if err != nil {
// 		m.logger.Errorf("[CPSRolesStorage][FindById] error aggregating cps_action_approver_index: %v", err)
// 		return nil, local_util.HandleDBError(err)
// 	}
// 	defer cursor.Close(ctx)

// 	type actionGroup struct {
// 		Actions []string `bson:"actions"`
// 	}

// 	type facetActions struct {
// 		Maker   []actionGroup `bson:"maker"`
// 		Checker []actionGroup `bson:"checker"`
// 		Auditor []actionGroup `bson:"auditor"`
// 	}
// 	var aggResult []facetActions

// 	if err := cursor.All(ctx, &aggResult); err != nil {
// 		m.logger.Errorf("[CPSRolesStorage][FindById] error decoding facet result: %v", err)
// 		return nil, local_util.HandleDBError(err)
// 	}
// 	m.logger.Infof("[CPSRolesStorage][FindById] facet aggregation result: %v", aggResult)

// 	var maker, checker, auditor []string
// 	if len(aggResult) > 0 {
// 		if len(aggResult[0].Maker) > 0 {
// 			maker = aggResult[0].Maker[0].Actions
// 		}
// 		if len(aggResult[0].Checker) > 0 {
// 			checker = aggResult[0].Checker[0].Actions
// 		}
// 		if len(aggResult[0].Auditor) > 0 {
// 			auditor = aggResult[0].Auditor[0].Actions
// 		}
// 	} else {
// 		m.logger.Warnf("[CPSRolesStorage][FindById] no approver index found for role_id: %s", result.RoleCode)
// 	}

// 	result.MakerActions = maker
// 	result.CheckerActions = checker
// 	result.AuditorActions = auditor

// 	// Fetch enabled/disabled service lists for this role
// 	enabledServices, disabledServices, err := m.fetchServiceAccessLists(ctx, result.ID)
// 	if err != nil {
// 		m.logger.Warnf("[CPSRolesStorage][FindById] failed to fetch service access lists for role %s: %v", id, err)
// 	}
// 	result.EnabledServices = enabledServices
// 	result.DisabledServices = disabledServices

// 	segEnabled, segErr := m.fetchServiceAccessListsByRoleNameSegmentationCode(ctx, result.Name)
// 	if segErr != nil {
// 		m.logger.Warnf("[CPSRolesStorage][FindById] failed to fetch service access lists by role name segmentation_code for role %s: %v", id, segErr)
// 	} else {
// 		result.EnabledServices = segEnabled
// 		// result.DisabledServices = segDisabled
// 	}

// 	return result, nil
// }

// func (m *cpsRoleStorage) fetchServiceAccessListsByRoleNameSegmentationCode(ctx context.Context, cpsRoleName string) ([]imodel.ServiceAccessInfo, error) {
// 	roleName := strings.TrimSpace(cpsRoleName)
// 	if roleName == "" {
// 		return nil, nil
// 	}

// 	// Optimized:
// 	// - Get all excluded keys in one query (distinct)
// 	// - Fetch all enabled access_list items excluding those keys in one query ($nin)
// 	segCol := m.client.Database(m.dbName).Collection(m.accessListSegCollection)
// 	accessListCol := m.client.Database(m.dbName).Collection(m.accessListCollection)

// 	segCodeRegex := "^\\s*" + regexp.QuoteMeta(roleName) + "\\s*$"
// 	segFilter := bson.M{
// 		"enabled":           true,
// 		"segmentation_code": bson.M{"$regex": segCodeRegex, "$options": "i"},
// 	}

// 	// Get excluded keys (enabled segmentation rows) in one query.
// 	segCursor, err := segCol.Find(ctx, segFilter, options.Find().SetProjection(bson.M{"access_list_key": 1}))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer segCursor.Close(ctx)

// 	type segKeyDoc struct {
// 		AccessListKey string `bson:"access_list_key"`
// 	}
// 	var segKeys []segKeyDoc
// 	if err := segCursor.All(ctx, &segKeys); err != nil {
// 		return nil, err
// 	}

// 	excludedSet := make(map[string]struct{}, len(segKeys))
// 	for _, d := range segKeys {
// 		if k := strings.TrimSpace(d.AccessListKey); k != "" {
// 			excludedSet[k] = struct{}{}
// 		}
// 	}

// 	excludedKeys := make([]string, 0, len(excludedSet))
// 	for k := range excludedSet {
// 		excludedKeys = append(excludedKeys, k)
// 	}

// 	accessFilter := bson.M{"enabled": true}
// 	if len(excludedKeys) > 0 {
// 		accessFilter["key"] = bson.M{"$nin": excludedKeys}
// 	}

// 	cursor, err := accessListCol.Find(ctx, accessFilter, options.Find().SetProjection(bson.M{
// 		"key":              1,
// 		"access_list_name": 1,
// 	}))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	type accessListDoc struct {
// 		Key            string `bson:"key"`
// 		AccessListName string `bson:"access_list_name"`
// 	}
// 	var docs []accessListDoc
// 	if err := cursor.All(ctx, &docs); err != nil {
// 		return nil, err
// 	}

// 	out := make([]imodel.ServiceAccessInfo, 0, len(docs))
// 	for _, d := range docs {
// 		out = append(out, imodel.ServiceAccessInfo{Key: d.Key, AccessListName: d.AccessListName})
// 	}
// 	return out, nil
// }

// // fetchServiceAccessLists retrieves enabled and disabled service lists for a CPS role.
// //   - Enabled: services in access_list (enabled=true) NOT blocked for this role
// //     (i.e., not in access_list_segmentation, or in segmentation with enabled=false)
// //   - Disabled: services in access_list (enabled=true) that ARE blocked for this role
// //     (i.e., in access_list_segmentation with enabled=true)
// func (m *cpsRoleStorage) fetchServiceAccessLists(ctx context.Context, roleID bson.ObjectID) ([]imodel.ServiceAccessInfo, []imodel.ServiceAccessInfo, error) {
// 	// 1. Fetch all globally enabled services from access_list
// 	accessListCol := m.client.Database(m.dbName).Collection(m.accessListCollection)
// 	alCursor, err := accessListCol.Find(ctx, bson.M{"enabled": true})
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	defer alCursor.Close(ctx)

// 	type accessListDoc struct {
// 		Key            string `bson:"key"`
// 		AccessListName string `bson:"access_list_name"`
// 	}
// 	var allServices []accessListDoc
// 	if err := alCursor.All(ctx, &allServices); err != nil {
// 		return nil, nil, err
// 	}

// 	// 2. Fetch segmentation entries for this role from access_list_segmentation
// 	segCol := m.client.Database(m.dbName).Collection(m.accessListSegCollection)
// 	segCursor, err := segCol.Find(ctx, bson.M{"segmented_id": roleID})
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	defer segCursor.Close(ctx)

// 	type segDoc struct {
// 		AccessListKey string `bson:"access_list_key"`
// 		Enabled       bool   `bson:"enabled"`
// 	}
// 	var segEntries []segDoc
// 	if err := segCursor.All(ctx, &segEntries); err != nil {
// 		return nil, nil, err
// 	}

// 	// 3. Build a map of blocked keys (segmentation enabled=true means disabled for this role)
// 	blockedKeys := make(map[string]bool, len(segEntries))
// 	for _, seg := range segEntries {
// 		if seg.Enabled {
// 			blockedKeys[seg.AccessListKey] = true
// 		}
// 	}

// 	// 4. Categorize services
// 	var enabledServices, disabledServices []imodel.ServiceAccessInfo
// 	for _, svc := range allServices {
// 		info := imodel.ServiceAccessInfo{
// 			Key:            svc.Key,
// 			AccessListName: svc.AccessListName,
// 		}
// 		if blockedKeys[svc.Key] {
// 			disabledServices = append(disabledServices, info)
// 		} else {
// 			enabledServices = append(enabledServices, info)
// 		}
// 	}

// 	return enabledServices, disabledServices, nil
// }

// func (m *cpsRoleStorage) FindByNameOrRoleCode(ctx context.Context, name, roleCode string) (*imodel.CPSRoles, error) {
// 	filter := []bson.M{}
// 	if roleCode != "" {
// 		filter = append(filter, bson.M{"name": bson.M{"$regex": "^" + name, "$options": "i"}})
// 	}
// 	if name != "" {
// 		filter = append(filter, bson.M{"role_code": bson.M{"$regex": "^" + roleCode, "$options": "i"}})
// 	}
// 	orFilter := bson.M{"$or": filter, "is_deleted": false}
// 	result, err := m.dal.FindOne(ctx, orFilter, bson.M{})
// 	if err != nil {
// 		if errors.Is(err, mongo.ErrNoDocuments) {
// 			return nil, errors.New(localization.ErrorResourceNotFound.Code)
// 		}
// 		m.logger.Errorf("[CPSRolesStorage][FindByNameOrRoleCode] failed to find cps role, name: %s, roleCode: %s, error: %v", name, roleCode, err)
// 		return nil, local_util.HandleDBError(err)
// 	}
// 	return result, nil
// }

// func (m *cpsRoleStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
// 	objID, err := bson.ObjectIDFromHex(id)
// 	if err != nil {
// 		m.logger.Errorf("[CPSRolesStorage][EnableOrDisable] invalid id format: %s, error: %v", id, err)
// 		return errors.New(localization.ErrorInvalidID.Code)
// 	}

// 	filter := bson.M{"_id": objID}
// 	update := bson.M{"enabled": enable, "updated_at": time.Now()}

// 	_, err = m.dal.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		m.logger.Errorf("[CPSRolesStorage][EnableOrDisable] failed to enable/disable cps role, id: %s, error: %v", id, err)
// 		return local_util.HandleDBError(err)
// 	}
// 	return nil
// }

// // EnableServiceAccess removes the block for the given access list keys on this role.
// // It deletes segmentation entries with enabled=true, or sets them to enabled=false.
// func (m *cpsRoleStorage) EnableServiceAccess(ctx context.Context, roleID string, accessListKeys []string) error {
// 	objID, err := bson.ObjectIDFromHex(roleID)
// 	if err != nil {
// 		m.logger.Errorf("[CPSRolesStorage][EnableServiceAccess] invalid role id: %s, error: %v", roleID, err)
// 		return errors.New(localization.ErrorInvalidID.Code)
// 	}

// 	segCol := m.client.Database(m.dbName).Collection("access_list_segmentation")
// 	filter := bson.M{
// 		"segmented_id":    objID,
// 		"access_list_key": bson.M{"$in": accessListKeys},
// 		"enabled":         true,
// 	}
// 	update := bson.M{"$set": bson.M{"enabled": false, "updated_at": time.Now()}}

// 	_, err = segCol.UpdateMany(ctx, filter, update)
// 	if err != nil {
// 		m.logger.Errorf("[CPSRolesStorage][EnableServiceAccess] failed to enable services for role %s: %v", roleID, err)
// 		return local_util.HandleDBError(err)
// 	}

// 	return nil
// }

// // DisableServiceAccess blocks the given access list keys for this role.
// // It upserts segmentation entries with enabled=true.
// func (m *cpsRoleStorage) DisableServiceAccess(ctx context.Context, roleID string, accessListKeys []string) error {
// 	objID, err := bson.ObjectIDFromHex(roleID)
// 	if err != nil {
// 		m.logger.Errorf("[CPSRolesStorage][DisableServiceAccess] invalid role id: %s, error: %v", roleID, err)
// 		return errors.New(localization.ErrorInvalidID.Code)
// 	}

// 	// Fetch the role to get name for segmentation_name
// 	role, err := m.dal.FindOne(ctx, bson.M{"_id": objID}, bson.M{})
// 	if err != nil {
// 		m.logger.Errorf("[CPSRolesStorage][DisableServiceAccess] role not found: %s, error: %v", roleID, err)
// 		return errors.New(localization.ErrorResourceNotFound.Code)
// 	}

// 	// Fetch access list names for the given keys
// 	accessListCol := m.client.Database(m.dbName).Collection("access_list")
// 	alCursor, err := accessListCol.Find(ctx, bson.M{"key": bson.M{"$in": accessListKeys}})
// 	if err != nil {
// 		m.logger.Errorf("[CPSRolesStorage][DisableServiceAccess] failed to fetch access list: %v", err)
// 		return errors.New(localization.ErrorUnexpectedError.Code)
// 	}
// 	defer alCursor.Close(ctx)

// 	type alDoc struct {
// 		Key            string `bson:"key"`
// 		AccessListName string `bson:"access_list_name"`
// 	}
// 	var alDocs []alDoc
// 	if err := alCursor.All(ctx, &alDocs); err != nil {
// 		m.logger.Errorf("[CPSRolesStorage][DisableServiceAccess] failed to decode access list: %v", err)
// 		return errors.New(localization.ErrorUnexpectedError.Code)
// 	}
// 	alNameMap := make(map[string]string, len(alDocs))
// 	for _, doc := range alDocs {
// 		alNameMap[doc.Key] = doc.AccessListName
// 	}

// 	segCol := m.client.Database(m.dbName).Collection("access_list_segmentation")
// 	for _, key := range accessListKeys {
// 		filter := bson.M{
// 			"segmented_id":    objID,
// 			"access_list_key": key,
// 		}
// 		update := bson.M{
// 			"$set": bson.M{
// 				"enabled":    true,
// 				"updated_at": time.Now(),
// 			},
// 			"$setOnInsert": bson.M{
// 				"_id":               bson.NewObjectID(),
// 				"type":              "block",
// 				"segmentation_type": "cps_role",
// 				"access_list_key":   key,
// 				"access_list_name":  alNameMap[key],
// 				"segmented_id":      objID,
// 				"segmentation_code": role.RoleCode,
// 				"segmentation_name": role.Name,
// 				"created_at":        time.Now(),
// 			},
// 		}
// 		opts := options.UpdateOne().SetUpsert(true)
// 		_, err := segCol.UpdateOne(ctx, filter, update, opts)
// 		if err != nil {
// 			m.logger.Errorf("[CPSRolesStorage][DisableServiceAccess] failed to upsert segmentation for key %s: %v", key, err)
// 			return errors.New(localization.ErrorUnexpectedError.Code)
// 		}
// 	}

// 	return nil
// }

// func (m *cpsRoleStorage) Delete(ctx context.Context, id string) error {
// 	objID, err := bson.ObjectIDFromHex(id)
// 	if err != nil {
// 		m.logger.Errorf("[CPSRolesStorage][Delete] invalid id format: %s, error: %v", id, err)
// 		return errors.New(localization.ErrorInvalidID.Code)
// 	}

// 	filter := bson.M{"_id": objID, "is_deleted": false}
// 	update := bson.M{"is_deleted": true, "deleted_at": time.Now()}

// 	_, err = m.dal.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		return local_util.HandleDBError(err)
// 	}
// 	return nil
// }

// func (m *cpsRoleStorage) FindByCustomerSegmentation(ctx context.Context, customerSegment string) (*imodel.CPSRoles, error) {
// 	filter := bson.M{"name": customerSegment, "enabled": true}
// 	result, err := m.dal.FindOne(ctx, filter, bson.M{})
// 	if err != nil {
// 		return nil, local_util.HandleDBError(err)
// 	}
// 	return result, nil
// }
