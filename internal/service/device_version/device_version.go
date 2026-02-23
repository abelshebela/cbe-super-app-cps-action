package deviceversion

import (
	"cbe-super-app-cps-action/internal/constants"
	deviceversion "cbe-super-app-cps-action/internal/constants/dto/device_version"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/device_version/core"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	"context"
	"errors"
	"fmt"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type DeviceVersionService struct {
	deviceVersionRepo storage.DeviceVersionControlRepository
	cpsService        service.CPSActionService
	kafkaProducer     *kafka.ClientOrchestrationProducer
	logger            utils.Logger
}

func NewDeviceVersionService(deviceVersionRepo storage.DeviceVersionControlRepository, cpsService service.CPSActionService, kafkaProducer *kafka.ClientOrchestrationProducer, logger utils.Logger) service.DeviceVersionServiceSrv {
	return &DeviceVersionService{
		deviceVersionRepo: deviceVersionRepo,
		cpsService:        cpsService,
		kafkaProducer:     kafkaProducer,
		logger:            logger,
	}

}

// Authorize applies the approved CPS action for Device Version operations.
func (d *DeviceVersionService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "DeviceVersion", "Authorize")
	defer span.End()

	d.logger.Infof("[DevVerSvc][Authorize] action: %s", cpsAction.RequestAction)
	actionData, err := local_util.JsonUnmarshal[model.DeviceVersionControl](cpsAction.CurrentAction)
	if err != nil {
		d.logger.Errorf("[DevVerSvc][Authorize] unmarshal err: %v", err)
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
			d.logger.Errorf("[DevVerSvc][Authorize] disable existing err: %v", err)
			span.AddEvent("Failed to disable existing device version", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("platform", actionData.Platform),
			))
			return nil, err
		}
		if err = d.deviceVersionRepo.Save(ctx, *actionData); err != nil {
			d.logger.Errorf("[DevVerSvc][Authorize] create err: %v", err)
			span.AddEvent("Device version create action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("platform", actionData.Platform),
			))
			return nil, err
		}

		createdDeviceVersion, err := d.deviceVersionRepo.FindOne(ctx, actionData.Platform, actionData.LatestVersion)
		if err != nil {
			d.logger.Errorf("[DevVerSvc][Authorize] fetch created for kafka err: %v", err)
			span.AddEvent("Failed to fetch created device version for kafka publish", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("platform", actionData.Platform),
				attribute.String("version", actionData.LatestVersion),
			))
		} else {
			if d.kafkaProducer != nil {
				if err := d.kafkaProducer.PublishMessage(ctx, createdDeviceVersion, string(constants.DeviceVersionControlTopic), string(constants.DeviceVersionControlTopic), "new device version control created"); err != nil {
					d.logger.Errorf("[DevVerSvc][Authorize] kafka publish err: %v", err)
					span.AddEvent("Kafka publish failed for created device version", trace.WithAttributes(
						attribute.String("error", err.Error()),
						attribute.String("platform", actionData.Platform),
						attribute.String("version", actionData.LatestVersion),
					))
				}
			}
		}
		d.logger.Infof("[DevVerSvc][Authorize] created platform: %s", actionData.Platform)
	case string(constants.RequestUpdateDeviceVersion), string(constants.RequestEnableDisableDeviceVersion):
		updateData, err := core.UpdateDeviceVersionBsonForDb(*actionData, cpsAction.MakerName)
		if err != nil {
			d.logger.Errorf("[DevVerSvc][Authorize] prepare update err: %v", err)
			span.AddEvent("Failed to prepare update data", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		if err := d.deviceVersionRepo.Update(ctx, cpsAction.UniqueId, updateData); err != nil {
			d.logger.Errorf("[DevVerSvc][Authorize] update err: %v", err)
			span.AddEvent("Device version update/enable-disable action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}

		updatedDeviceVersion, err := d.deviceVersionRepo.FindByID(ctx, cpsAction.UniqueId)
		if err != nil {
			d.logger.Errorf("[DevVerSvc][Authorize] fetch updated for kafka err: %v", err)
			span.AddEvent("Failed to fetch updated device version for kafka publish", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
		} else {
			if d.kafkaProducer != nil {
				if err := d.kafkaProducer.PublishMessage(ctx, updatedDeviceVersion, string(constants.DeviceVersionControlTopic), string(constants.DeviceVersionControlTopic), "update device version control"); err != nil {
					d.logger.Errorf("[DevVerSvc][Authorize] kafka publish updated err: %v", err)
					span.AddEvent("Kafka publish failed for updated device version", trace.WithAttributes(
						attribute.String("error", err.Error()),
						attribute.String("unique_id", cpsAction.UniqueId),
					))
				}
			}
		}
		d.logger.Infof("[DevVerSvc][Authorize] updated id: %s", cpsAction.UniqueId)

	case string(constants.RequestDeleteDeviceVersion):
		if err := d.deviceVersionRepo.Delete(ctx, cpsAction.UniqueId); err != nil {
			d.logger.Errorf("[DevVerSvc][Authorize] delete err: %v", err)
			span.AddEvent("Device version delete action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		d.logger.Infof("[DevVerSvc][Authorize] deleted id: %s", cpsAction.UniqueId)
	default:
		d.logger.Errorf("[DevVerSvc][Authorize] unsupported: %s", cpsAction.RequestAction)
		span.AddEvent("Unsupported action", trace.WithAttributes(
			attribute.String("error", localization.ErrorUnsupportedAction.Code),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	d.logger.Infof("[DevVerSvc][Authorize] done: %s", cpsAction.RequestAction)
	return cpsAction, nil
}

// CreateDeviceVersion implements service.DeviceVersionService.
func (d *DeviceVersionService) CreateDeviceVersion(ctx context.Context, deviceVersion deviceversion.CreateDeviceVersionRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateDeviceVersion", "DeviceVersion", "CreateDeviceVersion")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		d.logger.Errorf("[DevVerSvc][Create] incomplete user")
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.IncompleteUserInfo),
			attribute.String("platform", deviceVersion.Platform),
		))
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	device, err := d.deviceVersionRepo.FindOne(ctx, deviceVersion.Platform, deviceVersion.LatestVersion)
	if err != nil {
		if err.Error() != localization.ErrorResourceNotFound.Code {
			span.AddEvent("Failed to check existing device version", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("platform", deviceVersion.Platform),
			))
			return err
		}
	} else if !device.ID.IsZero() {
		return fmt.Errorf("device version control with version %s and platform %s already exists", deviceVersion.LatestVersion, deviceVersion.Platform)
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
		d.logger.Errorf("[DevVerSvc][Create] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("platform", deviceVersion.Platform),
		))
		return err
	}
	d.logger.Infof("[DevVerSvc][Create] request created platform: %s", deviceVersion.Platform)
	return nil
}

// EnableDisableDeviceVersion implements service.DeviceVersionService.
func (d *DeviceVersionService) EnableDisableDeviceVersion(ctx context.Context, id string, enableDisable bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableDisableDeviceVersion", "DeviceVersion", "EnableDisableDeviceVersion")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		d.logger.Errorf("[DevVerSvc][EnableDisable] incomplete user")
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.IncompleteUserInfo),
			attribute.String("id", id),
		))
		return errors.New(constants.IncompleteUserInfo)
	}
	deviceVersion, err := d.deviceVersionRepo.FindByID(ctx, id)
	if err != nil {
		d.logger.Errorf("[DevVerSvc][EnableDisable] find err: %v", err)
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
		d.logger.Errorf("[DevVerSvc][EnableDisable] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	d.logger.Infof("[DevVerSvc][EnableDisable] request created id: %s enabled: %v", id, enableDisable)
	return nil
}

// GetAllDeviceVersions implements service.DeviceVersionService.
func (d *DeviceVersionService) GetAllDeviceVersions(ctx context.Context, filterParams *types.Filter) (types.PaginatedResponse[[]model.DeviceVersionControl], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllDeviceVersions", "DeviceVersion", "GetAllDeviceVersions")
	defer span.End()

	// repository returns value slice; map to pointer slice to match signature
	res, err := d.deviceVersionRepo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		d.logger.Errorf("[DevVerSvc][GetAll] fetch err: %v", err)
		span.AddEvent("Failed to fetch device versions", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return types.PaginatedResponse[[]model.DeviceVersionControl]{}, err
	}
	d.logger.Infof("[DevVerSvc][GetAll] count: %d", len(res.Data))
	return res, nil
}

// GetDeviceVersionByID implements service.DeviceVersionService.
func (d *DeviceVersionService) GetDeviceVersionByID(ctx context.Context, id string) (model.DeviceVersionControl, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetDeviceVersionByID", "DeviceVersion", "GetDeviceVersionByID")
	defer span.End()

	dv, err := d.deviceVersionRepo.FindByID(ctx, id)
	if err != nil {
		d.logger.Errorf("[DevVerSvc][GetByID] find err: %v", err)
		span.AddEvent("Failed to find device version", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return model.DeviceVersionControl{}, err
	}
	d.logger.Infof("[DevVerSvc][GetByID] found id: %s", id)
	return dv, nil

}

// UpdateDeviceVersion implements service.DeviceVersionService.
func (d *DeviceVersionService) UpdateDeviceVersion(ctx context.Context, id string, req deviceversion.UpdateDeviceVersionRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateDeviceVersion", "DeviceVersion", "UpdateDeviceVersion")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		d.logger.Errorf("[DevVerSvc][Update] incomplete user")
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.IncompleteUserInfo),
			attribute.String("id", id),
		))
		return errors.New(constants.IncompleteUserInfo)
	}
	existing, err := d.deviceVersionRepo.FindByID(ctx, id)
	if err != nil {
		d.logger.Errorf("[DevVerSvc][Update] find err: %v", err)
		span.AddEvent("Failed to find device version", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	if existing == (model.DeviceVersionControl{}) {
		d.logger.Errorf("[DevVerSvc][Update] not found: %s", id)
		span.AddEvent("Device version not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorResourceNotFound.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	// apply updates
	update, err := core.UpdateDeviceVersionBson(req, makerData.FullName)
	if err != nil {
		d.logger.Errorf("[DevVerSvc][Update] prepare data err: %v", err)
		span.AddEvent("Failed to prepare update data", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	action := lib.CpsModelBuilder(existing.ID.Hex(), makerData, &existing, update, string(constants.RequestUpdateDeviceVersion), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &action); err != nil {
		d.logger.Errorf("[DevVerSvc][Update] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	d.logger.Infof("[DevVerSvc][Update] request created id: %s", id)
	return nil
}

func (d *DeviceVersionService) DisableExistingDeviceVersion(ctx context.Context, platform string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DisableExistingDeviceVersion", "DeviceVersion", "DisableExistingDeviceVersion")
	defer span.End()

	d.logger.Infof("[DevVerSvc][DisableExisting] platform: %s", platform)
	existing, err := d.deviceVersionRepo.FindOne(ctx, platform, "")
	if err != nil {
		if err.Error() == localization.ErrorResourceNotFound.Code {
			d.logger.Infof("[DevVerSvc][DisableExisting] none found platform: %s", platform)
			return nil
		}
		d.logger.Errorf("[DevVerSvc][DisableExisting] find err: %v", err)
		span.AddEvent("Failed to find existing device version", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("platform", platform),
		))
		return err
	}
	if existing == (model.DeviceVersionControl{}) {
		d.logger.Infof("[DevVerSvc][DisableExisting] none found platform: %s", platform)
		return nil
	}
	err = d.deviceVersionRepo.EnableOrDisable(ctx, existing.ID.Hex(), false)
	if err != nil {
		d.logger.Errorf("[DevVerSvc][DisableExisting] disable err: %v", err)
		span.AddEvent("Failed to disable existing device version", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("platform", platform),
			attribute.String("id", existing.ID.Hex()),
		))
		return errors.New(localization.ErrorOnDisablingExistingDeviceControl.Code)
	}
	d.logger.Infof("[DevVerSvc][DisableExisting] disabled platform: %s", platform)
	return nil
}
