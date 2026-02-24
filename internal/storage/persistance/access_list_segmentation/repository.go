package access_list_segmentation_repository

import (
	"cbe-super-app-cps-action/internal/constants"
	access_list_segmentation_dto "cbe-super-app-cps-action/internal/constants/dto/access_list_segmentation"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	local_model "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
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

type AccessListSegmentation struct {
	repo           dal.MongoDal[local_model.AccessListSegmentation, local_model.AccessListSegmentation]
	client         *mongo.Client
	accBlock       storage.AccountBlockRepository
	dbName         string
	collectionName string
	kafkaProducer  kafka.AccessListSegmentationProducer
	logger         utils.Logger
}

// FindParentChildRelationship implements [storage.AccessListSegmentationRepository].
func (a *AccessListSegmentation) FindParentChildRelationship(ctx context.Context) ([]local_model.AccessItemRelation, error) {
	collection := a.client.Database(a.dbName).Collection("access_items_relation")
	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{})

	if err != nil {
		a.logger.Errorf("[AccessListSegmentation][FindParentChildRelationship] failed to aggregate: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	var relations []local_model.AccessItemRelation
	if err := cursor.All(ctx, &relations); err != nil {
		a.logger.Errorf("[AccessListSegmentation][FindParentChildRelationship] failed to decode: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return relations, nil
}

// FindBySegmentIDAndAccessListKeys implements [storage.AccessListSegmentationRepository].
func (a *AccessListSegmentation) FindBySegmentIDAndAccessListKeys(ctx context.Context, id string, keys []string) (*local_model.AccessListSegmentation, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[AccessListSegmentation][FindBySegmentIDAndAccessListKeys] invalid ObjectID: %s", id)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"segmented_id": objID, "access_list_key": bson.M{"$in": keys}, "enabled": true}
	response, err := a.repo.FindOne(ctx, filter, nil)
	if err != nil {
		a.logger.Errorf("[AccessListSegmentation][FindBySegmentIDAndAccessListKeys] find error: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	return response, nil
}

// FindByAccountSegmentationAndServiceID implements storage.AccessListSegmentationRepository.
func (a *AccessListSegmentation) FindByAccountSegmentationAndAccessListKeys(ctx context.Context, customerSegments string, segmentKeys []string) (*local_model.AccessListSegmentation, error) {
	seg, err := a.repo.FindOne(ctx, bson.M{"segmentation_code": customerSegments, "service_id": bson.M{"$in": segmentKeys}, "enabled": true}, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		a.logger.Errorf("[AccessListSegmentation][FindByAccountSegmentationAndAccessListKeys] failed to find access list segmentation by segmentation id and service id: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return seg, nil
}

// FindBySegmentationAndServiceID implements storage.AccessListSegmentationRepository.
func (a *AccessListSegmentation) FindBySegmentationAndServiceID(ctx context.Context, segmentationID string, serviceID string) (*local_model.AccessListSegmentation, error) {
	seg, err := a.repo.FindOne(ctx, bson.M{"segmented_id": segmentationID, "service_id": serviceID, "enabled": true}, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		a.logger.Errorf("[AccessListSegmentation][FindBySegmentationAndServiceID] failed to find access list segmentation by segmentation id and service id: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return seg, nil
}

// FindByIDS implements storage.AccessListSegmentationRepository.
func (a *AccessListSegmentation) FindByIDS(ctx context.Context, ids []string, t string) (*local_model.AccessListSegmentation, error) {
	var filter bson.M
	var objIDs []bson.ObjectID
	for _, idStr := range ids {
		objID, err := bson.ObjectIDFromHex(idStr)
		if err != nil {
			a.logger.Errorf("[AccessListSegmentation][FindByIDS] invalid ObjectID: %s", idStr)
			return nil, errors.New(localization.ErrorInvalidID.Code)
		}
		objIDs = append(objIDs, objID)
	}
	filter = bson.M{"segmented_id": bson.M{"$in": objIDs}, "type": t, "enabled": true}

	als, err := a.repo.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		a.logger.Errorf("[AccessListSegmentation][FindByIDS] failed to find access list segmentation by ids: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return als, nil
}

func (a *AccessListSegmentation) CreateAccountSegment(ctx context.Context, accessListSegmentation access_list_segmentation_dto.CreateAccessListSegmentationRequest) error {
	docs := []local_model.AccessListSegmentation{}
	for i, idStr := range accessListSegmentation.AccessListKeys {
		doc := local_model.AccessListSegmentation{
			ID:               bson.NewObjectID(),
			Type:             accessListSegmentation.Type,
			AccessListKey:    idStr,
			AccessListName:   accessListSegmentation.AccessListNames[i],
			SegmentationCode: accessListSegmentation.SegmentCode,
			SegmentationName: accessListSegmentation.SegmentName,
			SegmentationType: accessListSegmentation.SegmentType,
			Enabled:          true,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		docs = append(docs, doc)
	}

	collection := a.client.Database(a.dbName).Collection(a.collectionName)
	_, err := collection.InsertMany(ctx, docs)
	if err != nil {
		a.logger.Errorf("[AccessListSegmentation][CreateAccountSegment] failed to insert documents: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	a.kafkaProducer.PublishMessage(ctx, docs, "create", string(constants.AccessListSegmentationTopic), "create account-segment")
	return nil
}

func (a *AccessListSegmentation) CreateBlockSegment(ctx context.Context, accessListSegmentation access_list_segmentation_dto.CreateAccessListSegmentationRequest) error {
	if len(accessListSegmentation.SegmentedID) == 0 {
		return errors.New(localization.ErrorAccessListSegmentationIDSRequired.Code)
	}

	var docs []interface{}
	for i, idStr := range accessListSegmentation.AccessListKeys {
		objID, err := bson.ObjectIDFromHex(accessListSegmentation.SegmentedID)
		if err != nil {
			a.logger.Errorf("[AccessListSegmentation][CreateBlockSegment] invalid ObjectID: %s", accessListSegmentation.SegmentedID)
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
		doc := local_model.AccessListSegmentation{
			ID:               bson.NewObjectID(),
			Type:             accessListSegmentation.Type,
			AccessListKey:    idStr,
			AccessListName:   accessListSegmentation.AccessListNames[i],
			SegmentationType: accessListSegmentation.SegmentType,
			SegmentedID:      objID,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			Enabled:          true,
		}
		docs = append(docs, doc)
	}

	collection := a.client.Database(a.dbName).Collection(a.collectionName)
	_, err := collection.InsertMany(ctx, docs)
	if err != nil {
		a.logger.Errorf("[AccessListSegmentation][CreateBlockSegment] failed to insert documents: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

// EnableOrDisable implements storage.AccessListSegmentationRepository.
func (a *AccessListSegmentation) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[AccessListSegmentation][EnableOrDisable] invalid ObjectID: %s", id)
		return errors.New(localization.ErrorInvalidID.Code)
	}
	update := bson.M{"enabled": enable}
	filter := bson.M{"_id": objID, "enabled": true}

	_, err = a.repo.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("[AccessListSegmentation][EnableOrDisable] failed to enable/disable access list segmentation: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (a *AccessListSegmentation) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]local_model.AccessListSegmentation], error) {
	a.logger.Infof("[AccessListSegmentation][FindAllWithPagination] filter params: %+v", filterParam)

	searchKeys := bson.M{}
	allowedKeys := []string{"type", "segmented_id", "created_at", "updated_at", "enabled"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"type": searchRegex},
			{"enabled": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	data, err := a.repo.FindAllWithPagination(ctx, filter, nil, skip, limit)
	if err != nil {
		a.logger.Errorf("[AccessListSegmentation][FindAllWithPagination] find error: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := a.repo.TotalCount(ctx, filter)
	if err != nil {
		a.logger.Errorf("[AccessListSegmentation][FindAllWithPagination] count error: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	a.logger.Infof("[AccessListSegmentation][FindAllWithPagination] retrieved %d segmentations", len(data))

	return &types.PaginatedResponse[[]local_model.AccessListSegmentation]{
		Data: data,
		Meta: meta,
	}, nil
}

// FindByID implements storage.AccessListSegmentationRepository.
func (a *AccessListSegmentation) FindByID(ctx context.Context, id string) (*local_model.AccessListSegmentation, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[AccessListSegmentation][FindByID] invalid ObjectID: %s", id)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "enabled": true}
	response, err := a.repo.FindOne(ctx, filter, nil)
	if err != nil {
		a.logger.Errorf("[AccessListSegmentation][FindByID] find error: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	return response, nil

}

// FindByIDAndType implements storage.AccessListSegmentationRepository.
func (a *AccessListSegmentation) FindByIDAndType(ctx context.Context, id string, t string) (*local_model.AccessListSegmentation, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[AccessListSegmentation][FindByIDAndType] invalid ObjectID: %s", id)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"segmented_id": objID, "type": t, "enabled": true}
	response, err := a.repo.FindOne(ctx, filter, nil)
	if err != nil {
		a.logger.Errorf("[AccessListSegmentation][FindByIDAndType] find error: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	return response, nil

}

// Update implements storage.AccessListSegmentationRepository.
func (a *AccessListSegmentation) Update(ctx context.Context, id string, accessListSegmentation local_model.AccessListSegmentation) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[AccessListSegmentation][Update] invalid ObjectID: %s", id)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filer := bson.M{"_id": objID}
	update := bson.M{}
	if accessListSegmentation.Type != "" {
		update["type"] = accessListSegmentation.Type
	}
	if accessListSegmentation.SegmentationType != "" {
		update["segmentation_type"] = accessListSegmentation.SegmentationType
	}
	if !accessListSegmentation.SegmentedID.IsZero() {
		update["segmented_id"] = accessListSegmentation.SegmentedID
	}
	if accessListSegmentation.AccessListKey != "" {
		update["access_list_key"] = accessListSegmentation.AccessListKey
	}
	if accessListSegmentation.SegmentationCode != "" {
		update["segmentation_code"] = accessListSegmentation.SegmentationCode
	}
	if accessListSegmentation.AccessListName != "" {
		update["access_list_name"] = accessListSegmentation.AccessListName
	}
	update["updated_at"] = time.Now()

	_, err = a.repo.UpdateOne(ctx, filer, update)
	if err != nil {
		a.logger.Errorf("[AccessListSegmentation][Update] update error: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}
func (a *AccessListSegmentation) FindAllBySegmentIDorSegmentCode(ctx context.Context, segmentIDorCode string) ([]local_model.AccessListSegmentation, error) {
	var filter bson.M
	var objID bson.ObjectID
	objID, err := bson.ObjectIDFromHex(segmentIDorCode)
	if err != nil {
		a.logger.Infof("[AccessListSegmentation][FindAllBySegmentIDorSegmentCode] segmentIDorCode is not a valid ObjectID, treating as segmentation code")
		filter = bson.M{
			"segmentation_code": segmentIDorCode,
		}
	} else {
		a.logger.Infof("[AccessListSegmentation][FindAllBySegmentIDorSegmentCode] segmentIDorCode is a valid ObjectID, treating as segmented ID")
		filter = bson.M{
			"segmented_id": objID,
		}
	}
	filter["enabled"] = true
	als, err := a.repo.FindAll(ctx, filter, nil)
	if err != nil {
		a.logger.Errorf("[AccessListSegmentation][FindAllBySegmentIDorSegmentCode] failed to find access list segmentation by segmentation id or segment code: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return als, nil
}
func (a *AccessListSegmentation) FindAllBySegmentIDorSegmentCodeAndKeys(ctx context.Context, segmentIDorCode string, ac []string) ([]local_model.AccessListSegmentation, error) {
	var filter bson.M
	var objID bson.ObjectID
	objID, err := bson.ObjectIDFromHex(segmentIDorCode)
	if err != nil {
		a.logger.Infof("[AccessListSegmentation][FindAllBySegmentIDorSegmentCodeAndKeys] segmentIDorCode is not a valid ObjectID, treating as segmentation code")
		filter = bson.M{
			"segmentation_code": segmentIDorCode,
		}
	} else {
		a.logger.Infof("[AccessListSegmentation][FindAllBySegmentIDorSegmentCodeAndKeys] segmentIDorCode is a valid ObjectID, treating as segmented ID")
		filter = bson.M{
			"segmented_id": objID,
		}
	}
	filter["access_list_key"] = bson.M{"$in": ac}
	filter["enabled"] = true
	als, err := a.repo.FindAll(ctx, filter, nil)
	if err != nil {
		a.logger.Errorf("[AccessListSegmentation][FindAllBySegmentIDorSegmentCodeAndKeys] failed to find access list segmentation: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return als, nil
}

// BulkDisable implements [storage.AccessListSegmentationRepository].
func (a *AccessListSegmentation) BulkDisable(ctx context.Context, req access_list_segmentation_dto.BulkDisableAccessListSegmentationRequest) error {
	var filter bson.M
	var objID bson.ObjectID
	objID, err := bson.ObjectIDFromHex(req.ID)
	if err != nil {
		a.logger.Infof("[AccessListSegmentation][BulkDisable] segmentIDorCode is not a valid ObjectID, treating as segmentation code")
		filter = bson.M{
			"segmentation_code": req.ID,
		}
	} else {
		a.logger.Infof("[AccessListSegmentation][BulkDisable] segmentIDorCode is a valid ObjectID, treating as segmented ID")
		filter = bson.M{
			"segmented_id": objID,
		}
	}
	filter["access_list_key"] = bson.M{"$in": req.Keys}
	update := bson.M{
		"$set": bson.M{
			"enabled": false,
		},
	}
	collection := a.client.Database(a.dbName).Collection(a.collectionName)
	_, err = collection.UpdateMany(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("[AccessListSegmentation][BulkDisable] failed to bulk disable access list segmentation: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	a.kafkaProducer.PublishMessage(ctx, req, "delete", string(constants.AccessListSegmentationTopic), "bulk disable access-list-segmentation")
	return nil
}

func NewAccessListSegmentationRepository(client *mongo.Client, cfg *config.VaultConfig, dbName, collectionName string, customerSegmentationProducer kafka.AccessListSegmentationProducer, logger utils.Logger) storage.AccessListSegmentationRepository {
	return &AccessListSegmentation{
		repo:           dal.NewMongoDal[local_model.AccessListSegmentation, local_model.AccessListSegmentation](client, cfg, dbName, collectionName),
		client:         client,
		kafkaProducer:  customerSegmentationProducer,
		dbName:         dbName,
		collectionName: collectionName,
		logger:         logger,
	}
}
