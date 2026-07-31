package event

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	eventdto "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/event"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/service"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/service/event/core"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"
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
	log := local_util.LoggerFromCtx(ctx, e.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "CreateEvent", "Event", "CreateEvent")
	defer span.End()

	log.Infof("[EventSvc][Create] name: %s", event.EventName)

	if err := core.SetMerchantDetails(ctx, e.merchantService, &event); err != nil {
		log.Errorf("[EventSvc][Create] merchant details err: %v", err)
		span.AddEvent("SetMerchantDetails failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("event_name", event.EventName),
		))
		return err
	}

	exist, err := e.repo.Find(ctx, event.EventName)
	if err != nil {
		log.Errorf("[EventSvc][Create] find err: %v", err)
		span.AddEvent("Failed to find event", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("event_name", event.EventName),
		))
		return err
	}

	if exist != nil {
		log.Warnf("[EventSvc][Create] already exists: %s", event.EventName)
		span.AddEvent("Event already exists", trace.WithAttributes(
			attribute.String("error", localization.ErrorEventAlreadyExists.Code),
			attribute.String("event_name", event.EventName),
		))
		return errors.New(localization.ErrorEventAlreadyExists.Code)
	}

	code, err := core.GeneratePrefixedName("EVE", event.EventName, e.logger)
	if err != nil {
		log.Errorf("[EventSvc][Create] gen name err: %v", err)
		span.AddEvent("GeneratePrefixedName failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("event_name", event.EventName),
		))
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	URL, err := lib.UploadFileToMinio(ctx, e.minio, e.bucketName, event.CoverImage, string(constants.EventFolderName), *e.cfg, "", e.logger)
	if err != nil {
		log.Errorf("[EventSvc][Create] upload err: %v", err)
		span.AddEvent("UploadFileToMinio failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("event_name", event.EventName),
		))
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	result := core.CreateEventMapper(event, code, URL)

	log.Infof("[EventSvc][Create] created code: %s", code)
	if err := core.HandleCPSAction(ctx, e.cpsService, "", constants.RequestCreateEvent, result, nil, constants.ActionCreate); err != nil {
		log.Errorf("[EventSvc][Create] cps action err for %s: %v", code, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("event_code", code),
		))
	}
	return nil
}

func (e *eventService) UpdateEvent(ctx context.Context, id string, event eventdto.EventRequest) error {
	log := local_util.LoggerFromCtx(ctx, e.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateEvent", "Event", "UpdateEvent")
	defer span.End()

	log.Infof("[EventSvc][Update] id: %s", id)

	if err := core.SetMerchantDetails(ctx, e.merchantService, &event); err != nil {
		log.Errorf("[EventSvc][Update] merchant details err: %v", err)
		span.AddEvent("SetMerchantDetails failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	prevEvent, err := e.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[EventSvc][Update] find err: %v", err)
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
			log.Errorf("[EventSvc][Update] upload err: %v", err)
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

	log.Infof("[EventSvc][Update] updated id: %s", id)
	err = core.HandleCPSAction(ctx, e.cpsService, id, constants.RequestUpdateEvent, curAction, *prevEvent, constants.ActionUpdate)
	if err != nil {
		log.Errorf("[EventSvc][Update] cps action err: %v", err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	return nil
}

func (e *eventService) DeleteEvent(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, e.logger)

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
		log.Errorf("[EventSvc][Delete] cps action err: %v", err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	return nil
}

func (e *eventService) EnableDisableEvent(ctx context.Context, id string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, e.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "EnableDisableEvent", "Event", "EnableDisableEvent")
	defer span.End()

	log.Infof("[EventSvc][EnableDisable] id: %s, enable: %v", id, enable)

	prevEvent, err := e.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[EventSvc][EnableDisable] find err: %v", err)
		span.AddEvent("Event not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorEventNotFound.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorEventNotFound.Code)
	}

	if enable && prevEvent.Enabled {
		log.Warnf("[EventSvc][EnableDisable] already enabled: %s", id)
		span.AddEvent("Event already enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorEventAlreadyEnabled.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorEventAlreadyEnabled.Code)
	}
	if !enable && !prevEvent.Enabled {
		log.Warnf("[EventSvc][EnableDisable] already disabled: %s", id)
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

	log.Infof("[EventSvc][EnableDisable] handled id: %s, enable: %v", id, enable)
	err = core.HandleCPSAction(ctx, e.cpsService, id, action, curData, *prevEvent, constants.ActionUpdate)
	if err != nil {
		log.Errorf("[EventSvc][EnableDisable] cps action err: %v", err)
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
