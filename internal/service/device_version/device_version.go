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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "DeviceVersion", "Authorize")
	defer span.End()

	d.logger.Infof("[Authorize] authorizing device version action: %s", cpsAction.RequestAction)
	actionData, err := local_util.JsonUnmarshal[model.DeviceVersionControl](cpsAction.CurrentAction)
	if err != nil {
		d.logger.Errorf("[Authorize] failed to unmarshal CurrentAction: %v", err)
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		return nil, fmt.Errorf("%s", localization.ErrorUnexpectedError.Code)
	}

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestCreateDeviceVersion):
		actionData.CreatedAt = time.Now()
		err = d.DisableExistingDeviceVersion(ctx, actionData.Platform)
		if err != nil {
			d.logger.Errorf("[Authorize] failed to disable existing device version: %v", err)
			span.AddEvent("Failed to disable existing device version", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("platform", actionData.Platform),
			))
			return nil, err
		}
		if err = d.deviceVersionRepo.Save(ctx, *actionData); err != nil {
			d.logger.Errorf("[Authorize] device version create action failed: %v", err)
			span.AddEvent("Device version create action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("platform", actionData.Platform),
			))
			return nil, err
		}
		d.logger.Infof("[Authorize] device version created successfully for platform: %s", actionData.Platform)
	case string(constants.RequestUpdateDeviceVersion), string(constants.RequestEnableDisableDeviceVersion):
		updateData, err := core.UpdateDeviceVersionBsonForDb(*actionData, cpsAction.MakerName)
		if err != nil {
			d.logger.Errorf("[Authorize] failed to prepare update data: %v", err)
			span.AddEvent("Failed to prepare update data", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		if err := d.deviceVersionRepo.Update(ctx, cpsAction.UniqueId, updateData); err != nil {
			d.logger.Errorf("[Authorize] device version update/enable-disable action failed: %v", err)
			span.AddEvent("Device version update/enable-disable action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		d.logger.Infof("[Authorize] device version updated successfully for id: %s", cpsAction.UniqueId)

	case string(constants.RequestDeleteDeviceVersion):
		if err := d.deviceVersionRepo.Delete(ctx, cpsAction.UniqueId); err != nil {
			d.logger.Errorf("[Authorize] device version delete action failed: %v", err)
			span.AddEvent("Device version delete action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		d.logger.Infof("[Authorize] device version deleted successfully for id: %s", cpsAction.UniqueId)
	default:
		d.logger.Errorf("[Authorize] unsupported action: %s", cpsAction.RequestAction)
		span.AddEvent("Unsupported action", trace.WithAttributes(
			attribute.String("error", localization.ErrorUnsupportedAction.Code),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	d.logger.Infof("[Authorize] device version action authorized successfully: %s", cpsAction.RequestAction)
	return cpsAction, nil
}

// CreateDeviceVersion implements service.DeviceVersionService.
func (d *DeviceVersionService) CreateDeviceVersion(ctx context.Context, deviceVersion deviceversion.CreateDeviceVersionRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateDeviceVersion", "DeviceVersion", "CreateDeviceVersion")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		d.logger.Errorf("Create Device Version failed incomplete user data")
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.IncompleteUserInfo),
			attribute.String("platform", deviceVersion.Platform),
		))
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	_, err := d.deviceVersionRepo.FindOne(ctx, bson.M{"platform": deviceVersion.Platform, "latest_version": deviceVersion.LatestVersion})
	if err != nil {
		if err.Error() != localization.ErrorResourceNotFound.Code {
			span.AddEvent("Failed to check existing device version", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("platform", deviceVersion.Platform),
			))
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
		d.logger.Errorf("[CreateDeviceVersion] failed to create CPS action: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("platform", deviceVersion.Platform),
		))
		return err
	}
	d.logger.Infof("[CreateDeviceVersion] device version creation request created successfully for platform: %s", deviceVersion.Platform)
	return nil
}

// EnableDisableDeviceVersion implements service.DeviceVersionService.
func (d *DeviceVersionService) EnableDisableDeviceVersion(ctx context.Context, id string, enableDisable bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableDisableDeviceVersion", "DeviceVersion", "EnableDisableDeviceVersion")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		d.logger.Errorf("[EnableDisableDeviceVersion] incomplete user data")
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.IncompleteUserInfo),
			attribute.String("id", id),
		))
		return errors.New(constants.IncompleteUserInfo)
	}
	objID, err := core.IdProvider(ctx, id)
	if err != nil {
		d.logger.Errorf("[EnableDisableDeviceVersion] failed to parse id: %v", err)
		span.AddEvent("Failed to parse id", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	filter := bson.M{"_id": objID}
	deviceVersion, err := d.deviceVersionRepo.FindOne(ctx, filter)
	if err != nil {
		d.logger.Errorf("[EnableDisableDeviceVersion] failed to find device version: %v", err)
		span.AddEvent("Failed to find device version", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	// short-circuit if already in desired state
	if deviceVersion.Enabled && enableDisable {
		span.AddEvent("Device version already enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorDeviceVersionAlreadyEnabled.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorDeviceVersionAlreadyEnabled.Code)
	}
	if !deviceVersion.Enabled && !enableDisable {
		span.AddEvent("Device version already disabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorDeviceVersionAlreadyDisabled.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorDeviceVersionAlreadyDisabled.Code)
	}

	updated := deviceVersion
	updated.Enabled = enableDisable
	updated.UpdatedBy = makerData.FullName
	updated.UpdatedAt = time.Now()
	// create CPS action for enable/disable
	action := lib.CpsModelBuilder(id, makerData, &deviceVersion, updated, string(constants.RequestEnableDisableDeviceVersion), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &action); err != nil {
		d.logger.Errorf("[EnableDisableDeviceVersion] failed to create CPS action: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	d.logger.Infof("[EnableDisableDeviceVersion] enable/disable request created successfully for id: %s, enabled: %v", id, enableDisable)
	return nil
}

// GetAllDeviceVersions implements service.DeviceVersionService.
func (d *DeviceVersionService) GetAllDeviceVersions(ctx context.Context, filterParams *types.Filter) (types.PaginatedResponse[[]model.DeviceVersionControl], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllDeviceVersions", "DeviceVersion", "GetAllDeviceVersions")
	defer span.End()

	// repository returns value slice; map to pointer slice to match signature
	res, err := d.deviceVersionRepo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		d.logger.Errorf("[GetAllDeviceVersions] failed to fetch device versions: %v", err)
		span.AddEvent("Failed to fetch device versions", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return types.PaginatedResponse[[]model.DeviceVersionControl]{}, err
	}
	d.logger.Infof("[GetAllDeviceVersions] retrieved %d device versions", len(res.Data))
	return res, nil
}

// GetDeviceVersionByID implements service.DeviceVersionService.
func (d *DeviceVersionService) GetDeviceVersionByID(ctx context.Context, id string) (model.DeviceVersionControl, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetDeviceVersionByID", "DeviceVersion", "GetDeviceVersionByID")
	defer span.End()

	objID, err := core.IdProvider(ctx, id)
	if err != nil {
		d.logger.Errorf("[GetDeviceVersionByID] failed to parse id: %v", err)
		span.AddEvent("Failed to parse id", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return model.DeviceVersionControl{}, err
	}
	filter := bson.M{"_id": objID}
	dv, err := d.deviceVersionRepo.FindOne(ctx, filter)
	if err != nil {
		d.logger.Errorf("[GetDeviceVersionByID] failed to find device version: %v", err)
		span.AddEvent("Failed to find device version", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return model.DeviceVersionControl{}, err
	}
	d.logger.Infof("[GetDeviceVersionByID] device version retrieved successfully for id: %s", id)
	return dv, nil

}

// UpdateDeviceVersion implements service.DeviceVersionService.
func (d *DeviceVersionService) UpdateDeviceVersion(ctx context.Context, id string, req deviceversion.UpdateDeviceVersionRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateDeviceVersion", "DeviceVersion", "UpdateDeviceVersion")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		d.logger.Errorf("[UpdateDeviceVersion] incomplete user data")
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.IncompleteUserInfo),
			attribute.String("id", id),
		))
		return errors.New(constants.IncompleteUserInfo)
	}
	// fetch existing
	objID, err := core.IdProvider(ctx, id)
	if err != nil {
		d.logger.Errorf("[UpdateDeviceVersion] failed to parse id: %v", err)
		span.AddEvent("Failed to parse id", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	filter := bson.M{"_id": objID}
	existing, err := d.deviceVersionRepo.FindOne(ctx, filter)
	if err != nil {
		d.logger.Errorf("[UpdateDeviceVersion] failed to find device version: %v", err)
		span.AddEvent("Failed to find device version", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	if existing == (model.DeviceVersionControl{}) {
		d.logger.Errorf("[UpdateDeviceVersion] device version not found: %s", id)
		span.AddEvent("Device version not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorResourceNotFound.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	// apply updates
	update, err := core.UpdateDeviceVersionBson(req, makerData.FullName)
	if err != nil {
		d.logger.Errorf("[UpdateDeviceVersion] failed to prepare update data: %v", err)
		span.AddEvent("Failed to prepare update data", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	action := lib.CpsModelBuilder(existing.ID.Hex(), makerData, &existing, update, string(constants.RequestUpdateDeviceVersion), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &action); err != nil {
		d.logger.Errorf("[UpdateDeviceVersion] failed to create CPS action: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	d.logger.Infof("[UpdateDeviceVersion] device version update request created successfully for id: %s", id)
	return nil
}

func (d *DeviceVersionService) DisableExistingDeviceVersion(ctx context.Context, platform string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DisableExistingDeviceVersion", "DeviceVersion", "DisableExistingDeviceVersion")
	defer span.End()

	d.logger.Infof("[DisableExistingDeviceVersion] disabling existing device version for platform: %s", platform)
	filter := bson.M{"enabled": true, "platform": platform}
	existing, err := d.deviceVersionRepo.FindOne(ctx, filter)
	if err != nil {
		if err.Error() == localization.ErrorResourceNotFound.Code {
			d.logger.Infof("[DisableExistingDeviceVersion] no existing enabled device version found for platform: %s", platform)
			return nil
		}
		d.logger.Errorf("[DisableExistingDeviceVersion] failed to find existing device version: %v", err)
		span.AddEvent("Failed to find existing device version", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("platform", platform),
		))
		return err
	}
	if existing == (model.DeviceVersionControl{}) {
		d.logger.Infof("[DisableExistingDeviceVersion] no existing enabled device version found for platform: %s", platform)
		return nil
	}
	err = d.deviceVersionRepo.EnableOrDisable(ctx, existing.ID.Hex(), false)
	if err != nil {
		d.logger.Errorf("[DisableExistingDeviceVersion] failed to disable existing device version: %v", err)
		span.AddEvent("Failed to disable existing device version", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("platform", platform),
			attribute.String("id", existing.ID.Hex()),
		))
		return errors.New(localization.ErrorOnDisablingExistingDeviceControl.Code)
	}
	d.logger.Infof("[DisableExistingDeviceVersion] existing device version disabled successfully for platform: %s", platform)
	return nil
}
