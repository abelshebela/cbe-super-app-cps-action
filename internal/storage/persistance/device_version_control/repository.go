package deviceversioncontrol

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/kafka"
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

type DeviceVersionControlRepository struct {
	deviceDal     dal.MongoDal[imodel.DeviceVersionControl, imodel.DeviceVersionControl]
	client        *mongo.Client
	kafkaProducer kafka.ClientOrchestrationProducer
	logger        utils.Logger
}

func NewDeviceVersionControlRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, kafkaProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.DeviceVersionControlRepository {
	return &DeviceVersionControlRepository{
		deviceDal:     dal.NewMongoDal[imodel.DeviceVersionControl, imodel.DeviceVersionControl](client, cfg, dbName, collection),
		client:        client,
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

func (d *DeviceVersionControlRepository) Save(ctx context.Context, deviceVersionControl imodel.DeviceVersionControl) error {
	log := local_util.LoggerFromCtx(ctx, d.logger)

	newDeviceVersion, err := d.deviceDal.InsertOne(ctx, deviceVersionControl)
	if err != nil {
		log.Errorf("[DeviceVersionControl][Save] failed to save device version control: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	_ = newDeviceVersion

	return nil
}

func (d *DeviceVersionControlRepository) Update(ctx context.Context, id string, deviceVersionControl bson.M) error {
	log := local_util.LoggerFromCtx(ctx, d.logger)

	log.Infof("[DeviceVersionControl][Update] updating device version control for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[DeviceVersionControl][Update] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID}

	deviceVersionControl["updated_at"] = time.Now()
	deviceVersionControl["last_modified_at"] = time.Now()
	updatedDeviceVersion, err := d.deviceDal.UpdateOne(ctx, filter, deviceVersionControl)
	if err != nil {
		log.Errorf("[DeviceVersionControl][Update] failed to update device version control: %v", err)
		return local_util.HandleDBError(err)
	}
	_ = updatedDeviceVersion

	log.Infof("[DeviceVersionControl][Update] device version control updated successfully")
	return nil
}

func (d *DeviceVersionControlRepository) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, d.logger)

	log.Infof("[DeviceVersionControl][Delete] deleting device version control for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[DeviceVersionControl][Delete] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	err = d.deviceDal.DeleteOne(ctx, filter)
	if err != nil {
		log.Errorf("[DeviceVersionControl][Delete] failed to delete device version control: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	log.Infof("[DeviceVersionControl][Delete] device version control deleted successfully")
	return nil
}

func (d *DeviceVersionControlRepository) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, d.logger)

	log.Infof("[DeviceVersionControl][EnableOrDisable] processing device version control enable/disable for id: %s, enabled: %v", id, enable)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[DeviceVersionControl][EnableOrDisable] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"enabled": enable}
	updatedDeviceVersion, err := d.deviceDal.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Errorf("[DeviceVersionControl][EnableOrDisable] failed to enable/disable device version control: %v", err)
		return local_util.HandleDBError(err)
	}
	_ = updatedDeviceVersion

	log.Infof("[DeviceVersionControl][EnableOrDisable] device version control enable/disable completed successfully")
	return nil
}

func (d *DeviceVersionControlRepository) FindByID(ctx context.Context, id string) (imodel.DeviceVersionControl, error) {
	log := local_util.LoggerFromCtx(ctx, d.logger)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[DeviceVersionControl][FindByID] invalid object id: %v", err)
		return imodel.DeviceVersionControl{}, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID}

	log.Infof("[DeviceVersionControl][FindByID] fetching device version control by id: %s", id)
	result, err := d.deviceDal.FindOne(ctx, filter, nil)
	if err != nil {
		log.Errorf("[DeviceVersionControl][FindByID] failed to find device version control: %v", err)
		return imodel.DeviceVersionControl{}, local_util.HandleDBError(err)
	}
	log.Infof("[DeviceVersionControl][FindByID] device version control retrieved successfully")
	return *result, nil
}

func (d *DeviceVersionControlRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (types.PaginatedResponse[[]imodel.DeviceVersionControl], error) {
	log := local_util.LoggerFromCtx(ctx, d.logger)

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

	// filter["is_deleted"] = false
	log.Infof("[DeviceVersionControl][FindAllWithPagination] applying filter: %v, skip: %d, limit: %d", filter, skip, limit)

	// 5. Fetch data
	data, err := d.deviceDal.FindAllWithPaginationE(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		log.Errorf("[DeviceVersionControl][FindAllWithPagination] failed to fetch device version controls: %v", err)
		return types.PaginatedResponse[[]imodel.DeviceVersionControl]{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// 6. Count total
	total, err := d.deviceDal.TotalCount(ctx, filter)
	if err != nil {
		log.Errorf("[DeviceVersionControl][FindAllWithPagination] failed to count device version controls: %v", err)
		return types.PaginatedResponse[[]imodel.DeviceVersionControl]{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	log.Infof("[DeviceVersionControl][FindAllWithPagination] retrieved %d device version controls", len(data))

	// 8. Return standard paginated response
	return types.PaginatedResponse[[]imodel.DeviceVersionControl]{
		Data: data,
		Meta: meta,
	}, nil
}

func (d *DeviceVersionControlRepository) FindOne(ctx context.Context, platform, lastVersion string) (imodel.DeviceVersionControl, error) {
	log := local_util.LoggerFromCtx(ctx, d.logger)

	filter := bson.M{"platform": platform}
	if lastVersion != "" {
		filter["latest_version"] = lastVersion
	}
	result, err := d.deviceDal.FindOne(ctx, filter, nil)
	if err != nil {
		log.Errorf("[DeviceVersionControl][FindOne] failed to find device version control: %v", err)
		return imodel.DeviceVersionControl{}, local_util.HandleDBError(err)
	}
	return *result, nil
}
