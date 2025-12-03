package deviceversioncontrol

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DeviceVersionControlRepository struct {
	deviceDal dal.MongoDal[model.DeviceVersionControl, model.DeviceVersionControl]
	logger    utils.Logger
	client    *mongo.Client
}

func NewDeviceVersionControlRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.DeviceVersionControlRepository {
	return &DeviceVersionControlRepository{
		deviceDal: dal.NewMongoDal[model.DeviceVersionControl, model.DeviceVersionControl](client, dbName, collection),
		logger:    logger,
		client:    client,
	}
}

func (d *DeviceVersionControlRepository) Save(ctx context.Context, deviceVersionControl model.DeviceVersionControl) error {
	_, err := d.deviceDal.InsertOne(ctx, deviceVersionControl)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (d *DeviceVersionControlRepository) Update(ctx context.Context, id string, deviceVersionControl bson.M) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID}

	deviceVersionControl["updated_at"] = time.Now()
	deviceVersionControl["last_modified_at"] = time.Now()
	_, err = d.deviceDal.UpdateOne(ctx, filter, deviceVersionControl)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (d *DeviceVersionControlRepository) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	return d.deviceDal.DeleteOne(ctx, filter)
}

func (d *DeviceVersionControlRepository) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"enabled": enable}
	_, err = d.deviceDal.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (d *DeviceVersionControlRepository) FindByID(ctx context.Context, id string) (*model.DeviceVersionControl, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID}

	result, err := d.deviceDal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			d.logger.Errorf("Device control not found for ID: %s,error ", id, err)
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		d.logger.Errorf("FindByID Device control failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return result, nil
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
	data, err := d.deviceDal.FindAllWithPaginationN(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return types.PaginatedResponse[[]model.DeviceVersionControl]{}, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 6. Count total
	total, err := d.deviceDal.TotalCount(ctx, filter)
	if err != nil {
		return types.PaginatedResponse[[]model.DeviceVersionControl]{}, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Return standard paginated response
	return types.PaginatedResponse[[]model.DeviceVersionControl]{
		Data: data,
		Meta: meta,
	}, nil
}

func (d *DeviceVersionControlRepository) FindOne(ctx context.Context, filter bson.M) (model.DeviceVersionControl, error) {
	result, err := d.deviceDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.DeviceVersionControl{}, errors.New(localization.ErrorResourceNotFound.Code)
		}
		return model.DeviceVersionControl{}, errors.New(localization.ErrorUnexpectedError.Message)
	}
	return *result, nil
}
