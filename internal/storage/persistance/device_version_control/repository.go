package deviceversioncontrol

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DeviceVersionControlRepository struct {
	deviceDal     dal.MongoDal[model.DeviceVersionControl, model.DeviceVersionControl]
	client        *mongo.Client
	kafkaProducer kafka.ClientOrchestrationProducer
	logger        utils.Logger
}

func NewDeviceVersionControlRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, kafkaProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.DeviceVersionControlRepository {
	return &DeviceVersionControlRepository{
		deviceDal:     dal.NewMongoDal[model.DeviceVersionControl, model.DeviceVersionControl](client, cfg, dbName, collection),
		client:        client,
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

func (d *DeviceVersionControlRepository) Save(ctx context.Context, deviceVersionControl model.DeviceVersionControl) error {
	newDeviceVersion, err := d.deviceDal.InsertOne(ctx, deviceVersionControl)
	if err != nil {
		d.logger.Errorf("[DeviceVersionControl][Save] failed to save device version control: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	_ = newDeviceVersion

	return nil
}

func (d *DeviceVersionControlRepository) Update(ctx context.Context, id string, deviceVersionControl bson.M) error {
	d.logger.Infof("[DeviceVersionControl][Update] updating device version control for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		d.logger.Errorf("[DeviceVersionControl][Update] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID}

	deviceVersionControl["updated_at"] = time.Now()
	deviceVersionControl["last_modified_at"] = time.Now()
	updatedDeviceVersion, err := d.deviceDal.UpdateOne(ctx, filter, deviceVersionControl)
	if err != nil {
		d.logger.Errorf("[DeviceVersionControl][Update] failed to update device version control: %v", err)
		return local_util.HandleDBError(err)
	}
	_ = updatedDeviceVersion

	d.logger.Infof("[DeviceVersionControl][Update] device version control updated successfully")
	return nil
}

func (d *DeviceVersionControlRepository) Delete(ctx context.Context, id string) error {
	d.logger.Infof("[DeviceVersionControl][Delete] deleting device version control for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		d.logger.Errorf("[DeviceVersionControl][Delete] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	err = d.deviceDal.DeleteOne(ctx, filter)
	if err != nil {
		d.logger.Errorf("[DeviceVersionControl][Delete] failed to delete device version control: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	d.logger.Infof("[DeviceVersionControl][Delete] device version control deleted successfully")
	return nil
}

func (d *DeviceVersionControlRepository) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	d.logger.Infof("[DeviceVersionControl][EnableOrDisable] processing device version control enable/disable for id: %s, enabled: %v", id, enable)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		d.logger.Errorf("[DeviceVersionControl][EnableOrDisable] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"enabled": enable}
	updatedDeviceVersion, err := d.deviceDal.UpdateOne(ctx, filter, update)
	if err != nil {
		d.logger.Errorf("[DeviceVersionControl][EnableOrDisable] failed to enable/disable device version control: %v", err)
		return local_util.HandleDBError(err)
	}
	_ = updatedDeviceVersion

	d.logger.Infof("[DeviceVersionControl][EnableOrDisable] device version control enable/disable completed successfully")
	return nil
}

func (d *DeviceVersionControlRepository) FindByID(ctx context.Context, id string) (model.DeviceVersionControl, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		d.logger.Errorf("[DeviceVersionControl][FindByID] invalid object id: %v", err)
		return model.DeviceVersionControl{}, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID}

	d.logger.Infof("[DeviceVersionControl][FindByID] fetching device version control by id: %s", id)
	result, err := d.deviceDal.FindOne(ctx, filter, nil)
	if err != nil {
		d.logger.Errorf("[DeviceVersionControl][FindByID] failed to find device version control: %v", err)
		return model.DeviceVersionControl{}, local_util.HandleDBError(err)
	}
	d.logger.Infof("[DeviceVersionControl][FindByID] device version control retrieved successfully")
	return *result, nil
}

func (d *DeviceVersionControlRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (types.PaginatedResponse[[]model.DeviceVersionControl], error) {

	searchKeys := bson.M{}

	// 2. Allowed filterable/searchable fields
	allowedKeys := []string{"latest_version", "platform", "enabled", "force_update", "created_at", "last_modified", "updated_at", "updated_by", "created_by", "enabled"}

	// 3. Add search (if provided)
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"latest_version": searchRegex},
			{"platform": searchRegex},
			{"enabled": searchRegex},
			{"updated_at": searchRegex},
			{"force_update": searchRegex},
			{"updated_by": searchRegex},
			{"created_by": searchRegex},
			{"created_at": searchRegex},
			{"last_modified_at": searchRegex},
		}

	}
	// 4. Build filter, skip, limit
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	// 5. Fetch data
	data, err := d.deviceDal.FindAllWithPaginationE(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		d.logger.Errorf("[DeviceVersionControl][FindAllWithPagination] failed to fetch device version controls: %v", err)
		return types.PaginatedResponse[[]model.DeviceVersionControl]{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// 6. Count total
	total, err := d.deviceDal.TotalCount(ctx, filter)
	if err != nil {
		d.logger.Errorf("[DeviceVersionControl][FindAllWithPagination] failed to count device version controls: %v", err)
		return types.PaginatedResponse[[]model.DeviceVersionControl]{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	d.logger.Infof("[DeviceVersionControl][FindAllWithPagination] retrieved %d device version controls", len(data))

	// 8. Return standard paginated response
	return types.PaginatedResponse[[]model.DeviceVersionControl]{
		Data: data,
		Meta: meta,
	}, nil
}

func (d *DeviceVersionControlRepository) FindOne(ctx context.Context, platform, lastVersion string) (model.DeviceVersionControl, error) {
	filter := bson.M{"platform": platform}
	if lastVersion != "" {
		filter["latest_version"] = lastVersion
	}
	result, err := d.deviceDal.FindOne(ctx, filter, nil)
	if err != nil {
		d.logger.Errorf("[DeviceVersionControl][FindOne] failed to find device version control: %v", err)
		return model.DeviceVersionControl{}, local_util.HandleDBError(err)
	}
	return *result, nil
}
