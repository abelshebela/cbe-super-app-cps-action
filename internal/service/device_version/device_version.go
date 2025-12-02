package deviceversion

import (
	"cbe-super-app-cps-action/internal/constants"
	deviceversion "cbe-super-app-cps-action/internal/constants/dto/device_version"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/device_version/core"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"fmt"
	"time"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type DeviceVersionService struct {
	deviceVersionRepo storage.DeviceVersionControlRepository
	cpsService        service.CPSActionService
	logger            utils.Logger
}

func NewDeviceVersionService(deviceVersionRepo storage.DeviceVersionControlRepository, cpsService service.CPSActionService, logger utils.Logger) service.DeviceVersionServiceSrv {
	return &DeviceVersionService{
		deviceVersionRepo: deviceVersionRepo,
		cpsService:        cpsService,
		logger:            logger,
	}

}

// Authorize applies the approved CPS action for Device Version operations.
func (d *DeviceVersionService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	actionData, err := local_util.JsonUnmarshal[model.DeviceVersionControl](cpsAction.CurrentAction)
	if err != nil {
		return nil, fmt.Errorf("%s", localization.ErrorUnexpectedError.Code)
	}

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestCreateDeviceVersion):
		actionData.CreatedAt = time.Now()
		err = d.DisableExistingDeviceVersion(ctx,actionData.Platform)
		if err!=nil{
			return nil, err
		}
		if err = d.deviceVersionRepo.Save(ctx, *actionData); err != nil {
			d.logger.Errorf("DeviceVersion create action failed: %v", err)
			return nil, err
		}
	case string(constants.RequestUpdateDeviceVersion), string(constants.RequestEnableDisableDeviceVersion):
		updateData, err := core.UpdateDeviceVersionBsonForDb(*actionData, cpsAction.MakerName)
		if err != nil {
			return nil, err
		}
		if err := d.deviceVersionRepo.Update(ctx, cpsAction.UniqueId, updateData); err != nil {
			d.logger.Errorf("DeviceVersion update/enable-disable action failed: %v", err)
			return nil, err
		}

	case string(constants.RequestDeleteDeviceVersion):
		if err := d.deviceVersionRepo.Delete(ctx, cpsAction.UniqueId); err != nil {
			d.logger.Errorf("DeviceVersion delete action failed: %v", err)
			return nil, err
		}
	default:
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	return cpsAction, nil
}

// CreateDeviceVersion implements service.DeviceVersionService.
func (d *DeviceVersionService) CreateDeviceVersion(ctx context.Context, deviceVersion deviceversion.CreateDeviceVersionRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		d.logger.Errorf("Create Device Version failed incomplete user data")
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	_, err := d.deviceVersionRepo.FindOne(ctx, bson.M{"platform": deviceVersion.Platform, "latest_version": deviceVersion.LatestVersion})
	if err != nil {
		if err.Error() != localization.ErrorResourceNotFound.Code {
			return err
		}
	}

	new_device_version := model.DeviceVersionControl{
		LatestVersion: deviceVersion.LatestVersion,
		Platform:      deviceVersion.Platform,
		CreatedBy:     makerData.FullName,
		CreatedAt:     time.Now(),
		ForceUpdate:   deviceVersion.ForceUpdate,
		ReleaseNotes:  deviceVersion.ReleaseNotes,
		Enabled:       true,
	}

	action := lib.CpsModelBuilder("", makerData, nil, new_device_version, string(constants.RequestCreateDeviceVersion), constants.CREATE)
	if err := d.cpsService.CreateCPSAction(ctx, &action); err != nil {
		return err
	}
	return nil
}

// EnableDisableDeviceVersion implements service.DeviceVersionService.
func (d *DeviceVersionService) EnableDisableDeviceVersion(ctx context.Context, id string, enableDisable bool) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		return errors.New(constants.IncompleteUserInfo)
	}
	objID, err := core.IdProvider(ctx, id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	deviceVersion, err := d.deviceVersionRepo.FindOne(ctx, filter)
	if err != nil {
		return err
	}

	// short-circuit if already in desired state
	if deviceVersion.Enabled && enableDisable {
		return errors.New(localization.ErrorDeviceVersionAlreadyEnabled.Code)
	}
	if !deviceVersion.Enabled && !enableDisable {
		return errors.New(localization.ErrorDeviceVersionAlreadyDisabled.Code)
	}

	updated := deviceVersion
	updated.Enabled = enableDisable
	updated.UpdatedBy = makerData.FullName
	updated.UpdatedAt = time.Now()
	// create CPS action for enable/disable
	action := lib.CpsModelBuilder(id, makerData, &deviceVersion, updated, string(constants.RequestEnableDisableDeviceVersion), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &action); err != nil {
		return err
	}
	return nil
}

// GetAllDeviceVersions implements service.DeviceVersionService.
func (d *DeviceVersionService) GetAllDeviceVersions(ctx context.Context, filterParams *types.Filter) (types.PaginatedResponse[[]model.DeviceVersionControl], error) {
	// repository returns value slice; map to pointer slice to match signature
	res, err := d.deviceVersionRepo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		return types.PaginatedResponse[[]model.DeviceVersionControl]{}, err
	}

	return res, nil
}

// GetDeviceVersionByID implements service.DeviceVersionService.
func (d *DeviceVersionService) GetDeviceVersionByID(ctx context.Context, id string) (model.DeviceVersionControl, error) {
	objID, err := core.IdProvider(ctx, id)
	if err != nil {
		return model.DeviceVersionControl{}, err
	}
	filter := bson.M{"_id": objID}
	dv, err := d.deviceVersionRepo.FindOne(ctx, filter)
	if err != nil {
		return model.DeviceVersionControl{}, err
	}
	return dv, nil

}

// UpdateDeviceVersion implements service.DeviceVersionService.
func (d *DeviceVersionService) UpdateDeviceVersion(ctx context.Context, id string, req deviceversion.UpdateDeviceVersionRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		return errors.New(constants.IncompleteUserInfo)
	}
	// fetch existing
	objID, err := core.IdProvider(ctx, id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	existing, err := d.deviceVersionRepo.FindOne(ctx, filter)
	if err != nil {
		return err
	}

	if existing == (model.DeviceVersionControl{}) {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	// apply updates
	update, err := core.UpdateDeviceVersionBson(req, makerData.FullName)
	if err != nil {
		return err
	}

	action := lib.CpsModelBuilder(existing.ID.Hex(), makerData, &existing, update, string(constants.RequestUpdateDeviceVersion), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &action); err != nil {
		return err
	}

	return nil
}

func (d *DeviceVersionService) DisableExistingDeviceVersion(ctx context.Context,platform string) error {
	filter := bson.M{"enabled": true,"platform":platform}
	existing, err := d.deviceVersionRepo.FindOne(ctx, filter)
	if err != nil {
		return err
	}
	if existing == (model.DeviceVersionControl{}) {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	err = d.deviceVersionRepo.EnableOrDisable(ctx,existing.ID.Hex(),false)
	if err != nil {
		return errors.New(localization.ErrorOnDisablingExistingDeviceControl.Code)
	}

	return nil
}