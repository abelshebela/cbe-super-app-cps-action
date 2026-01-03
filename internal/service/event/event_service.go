package event

import (
	"cbe-super-app-cps-action/internal/constants"
	eventdto "cbe-super-app-cps-action/internal/constants/dto/event"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/event/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"path"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type eventService struct {
	repo            storage.EventRepository
	cpsService      service.CPSActionService
	merchantService service.EcommerceMerchantService
	userRepo        storage.UserRepository
	bucketName      string
	logger          utils.Logger
	minio           *s3.Client
	minioPubUrl     string
	cfg             *config.VaultConfig
}

func NewEventService(repo storage.EventRepository, cpsActionService service.CPSActionService, merchantService service.EcommerceMerchantService, userRepo storage.UserRepository, minio *s3.Client, minioPubUrl string, bucketName string, cfg *config.VaultConfig, logger utils.Logger) service.EventService {
	return &eventService{
		repo:            repo,
		cpsService:      cpsActionService,
		merchantService: merchantService,
		userRepo:        userRepo,
		minio:           minio,
		cfg:             cfg,
		bucketName:      bucketName,
		minioPubUrl:     minioPubUrl,
		logger:          logger,
	}
}

func (e *eventService) CreateEvent(ctx context.Context, event eventdto.EventRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateEvent", "Event", "CreateEvent")
	defer span.End()

	e.logger.Infof("CreateEvent called", "event_name", event.EventName)

	if err := core.SetMerchantDetails(ctx, e.merchantService, &event); err != nil {
		e.logger.Errorf("SetMerchantDetails failed", "error", err)
		span.AddEvent("SetMerchantDetails failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("event_name", event.EventName),
		))
		return err
	}

	exist, err := e.repo.Find(ctx, event.EventName)
	if err != nil {
		e.logger.Errorf("FindByName failed: %v", err)
		span.AddEvent("Failed to find event", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("event_name", event.EventName),
		))
		return err
	}

	if exist != nil {
		e.logger.Warnf("Event already exists: %s", event.EventName)
		span.AddEvent("Event already exists", trace.WithAttributes(
			attribute.String("error", localization.ErrorEventAlreadyExists.Code),
			attribute.String("event_name", event.EventName),
		))
		return errors.New(localization.ErrorEventAlreadyExists.Code)
	}

	code, err := core.GeneratePrefixedName("EVE", event.EventName, e.logger)
	if err != nil {
		e.logger.Errorf("GeneratePrefixedName failed", "error", err)
		span.AddEvent("GeneratePrefixedName failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("event_name", event.EventName),
		))
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	URL, err := lib.UploadFileToMinio(ctx, e.minio, e.bucketName, event.CoverImage, string(constants.EventFolderName), *e.cfg, "", e.logger)
	if err != nil {
		e.logger.Errorf("UploadFileToMinio failed", "error", err)
		span.AddEvent("UploadFileToMinio failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("event_name", event.EventName),
		))
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	result := core.CreateEventMapper(event, code, URL)

	e.logger.Infof("Event created successfully", "event_code", code)
	if err := core.HandleCPSAction(ctx, e.cpsService, "", constants.RequestCreateEvent, result, nil, constants.ActionCreate); err != nil {
		e.logger.Errorf("CPS action failed for event %s: %v", code, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("event_code", code),
		))
	}
	return nil
}

func (e *eventService) UpdateEvent(ctx context.Context, id string, event eventdto.EventRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateEvent", "Event", "UpdateEvent")
	defer span.End()

	e.logger.Infof("UpdateEvent called", "event_id", id)

	if err := core.SetMerchantDetails(ctx, e.merchantService, &event); err != nil {
		e.logger.Errorf("SetMerchantDetails failed", "event_id", id, "error", err)
		span.AddEvent("SetMerchantDetails failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	prevEvent, err := e.repo.FindByID(ctx, id)
	if err != nil {
		e.logger.Errorf("FindByID failed", "event_id", id, "error", err)
		span.AddEvent("Failed to find event", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	var URL string
	if event.CoverImage != nil {
		var objectkey string
		if prevEvent.EventInformation.Cover != "" {
			objectkey = path.Base(prevEvent.EventInformation.Cover)
		}

		URL, err = lib.UploadFileToMinio(ctx, e.minio, e.bucketName, event.CoverImage, string(constants.EventFolderName), *e.cfg, objectkey, e.logger)
		if err != nil {
			e.logger.Errorf("UploadFileToMinio failed", "event_id", id, "error", err)
			span.AddEvent("UploadFileToMinio failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
	} else {
		URL = prevEvent.EventInformation.Cover
	}

	curAction := core.EventMapperForUpdate(prevEvent, event, URL)

	e.logger.Infof("Event updated successfully", "event_id", id)
	err = core.HandleCPSAction(ctx, e.cpsService, id, constants.RequestUpdateEvent, curAction, *prevEvent, constants.ActionUpdate)
	if err != nil {
		e.logger.Errorf("CPS action failed for event %s: %v", curAction.EventName, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	return nil
}

func (e *eventService) DeleteEvent(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteEvent", "Event", "DeleteEvent")
	defer span.End()

	prevEvent, err := e.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Event not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorEventNotFound.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorEventNotFound.Code)
	}

	curData := *prevEvent
	curData.IsDeleted = true
	curData.DeletedAt = time.Now()

	err = core.HandleCPSAction(ctx, e.cpsService, id, constants.RequestDeleteEvent, curData, *prevEvent, constants.ActionDelete)
	if err != nil {
		e.logger.Errorf("CPS action failed for event %s: %v", curData.EventName, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	return nil
}

func (e *eventService) EnableDisableEvent(ctx context.Context, id string, enable bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableDisableEvent", "Event", "EnableDisableEvent")
	defer span.End()

	e.logger.Infof("EnableDisableEvent called", "event_id", id, "enable", enable)

	prevEvent, err := e.repo.FindByID(ctx, id)
	if err != nil {
		e.logger.Errorf("FindByID failed", "event_id", id, "error", err)
		span.AddEvent("Event not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorEventNotFound.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorEventNotFound.Code)
	}

	if enable && prevEvent.Enabled {
		e.logger.Warnf("Event already enabled", "event_id", id)
		span.AddEvent("Event already enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorEventAlreadyEnabled.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorEventAlreadyEnabled.Code)
	}
	if !enable && !prevEvent.Enabled {
		e.logger.Warnf("Event already disabled", "event_id", id)
		span.AddEvent("Event already disabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorEventAlreadyDisabled.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorEventAlreadyDisabled.Code)
	}

	curData := *prevEvent
	curData.Enabled = enable
	curData.LastModifiedAt = time.Now()

	var action constants.RequestAction
	if enable {
		action = constants.RequestEnableEvent
	} else {
		action = constants.RequestDisableEvent
	}

	e.logger.Infof("Event enable/disable action handled", "event_id", id, "enable", enable)
	err = core.HandleCPSAction(ctx, e.cpsService, id, action, curData, *prevEvent, constants.ActionUpdate)
	if err != nil {
		e.logger.Errorf("CPS action failed for event %s: %v", curData.EventName, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	return nil
}

func (e *eventService) FetchEventByID(ctx context.Context, id string) (*model.Event, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchEventByID", "Event", "FetchEventByID")
	defer span.End()

	result, err := e.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to fetch event", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}
	return result, nil
}

func (e *eventService) FetchEvent(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.Event], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchEvent", "Event", "FetchEvent")
	defer span.End()

	result, err := e.repo.FindAllWithPagination(ctx, filterParam)
	if err != nil {
		span.AddEvent("Failed to fetch events", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	return result, nil
}
func (e *eventService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Event", "Authorize")
	defer span.End()

	requestedAction := action.RequestAction

	event, err := local_util.JsonUnmarshal[model.Event](action.CurrentAction)
	if err != nil {
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", action.UniqueId),
		))
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	switch requestedAction {
	case string(constants.RequestCreateEvent):
		err = e.repo.Create(ctx, event)
		if err != nil {
			span.AddEvent("Failed to create event", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestUpdateEvent):
		err = e.repo.Update(ctx, event.ID.String(), event)
		if err != nil {
			span.AddEvent("Failed to update event", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestDeleteEvent):
		err = e.repo.Delete(ctx, action.UniqueId)
		if err != nil {
			span.AddEvent("Failed to delete event", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestEnableEvent):
		err = e.repo.EnableOrDisable(ctx, action.UniqueId, true)
		if err != nil {
			span.AddEvent("Failed to enable event", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestDisableEvent):
		err = e.repo.EnableOrDisable(ctx, action.UniqueId, false)
		if err != nil {
			span.AddEvent("Failed to disable event", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	default:
		span.AddEvent("Invalid request", trace.WithAttributes(
			attribute.String("error", localization.ErrorInvalidRequest.Code),
			attribute.String("request_action", string(requestedAction)),
		))
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	action.CurrentAction = event
	return action, nil
}
