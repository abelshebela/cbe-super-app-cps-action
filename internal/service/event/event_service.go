package event

import (
	"cbe-super-app-cps-action/internal/constants"
	eventdto "cbe-super-app-cps-action/internal/constants/dto/event"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/event/core"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type eventService struct {
	repo            storage.EventRepository
	cpsService      service.CPSActionService
	merchantService service.MiniAppMerchantService
	userRepo        storage.UserRepository
	bucketName      string
	logger          utils.Logger
	minio           config.MinioClientInterface
	cfg             *config.VaultConfig
}

func NewEventService(repo storage.EventRepository, cpsActionService service.CPSActionService, merchantService service.MiniAppMerchantService, userRepo storage.UserRepository, minio config.MinioClientInterface, bucketName string, cfg *config.VaultConfig, logger utils.Logger) service.EventService {
	return &eventService{
		repo:            repo,
		cpsService:      cpsActionService,
		merchantService: merchantService,
		userRepo:        userRepo,
		minio:           minio,
		cfg:             cfg,
		bucketName:      bucketName,
		logger:          logger,
	}
}

func (e *eventService) CreateEvent(ctx context.Context, event eventdto.EventRequest) error {
	e.logger.Infof("CreateEvent called", "event_name", event.EventName)

	if err := core.SetMerchantDetails(ctx, e.merchantService, &event); err != nil {
		e.logger.Errorf("SetMerchantDetails failed", "error", err)
		return err
	}

	exist, err := e.repo.FindByName(ctx, event.EventName)
	if err != nil {
		e.logger.Errorf("FindByName failed: %v", err)
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	if exist != nil {
		e.logger.Warnf("Event already exists: %s", event.EventName)
		return errors.New(localization.ErrorEventAlreadyExists.Code)
	}

	code, err := core.GeneratePrefixedName("EVE", event.EventName, e.logger)
	if err != nil {
		e.logger.Errorf("GeneratePrefixedName failed", "error", err)
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	URL, err := lib.UploadFileToMinio(ctx, e.minio, e.bucketName, event.CoverImage, "cover_image", e.cfg.MinioEndPoint, e.logger)
	if err != nil {
		e.logger.Errorf("UploadFileToMinio failed", "error", err)
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	result := core.CreateEventMapper(event, code, URL)

	if err := e.repo.Create(ctx, result); err != nil {
		e.logger.Errorf("repo.Create failed", "error", err)
		return err
	}

	e.logger.Infof("Event created successfully", "event_code", code)
	if err := core.HandleCPSAction(ctx, e.cpsService, "", constants.RequestCreateEvent, result, nil, constants.ActionCreate); err != nil {
		e.logger.Errorf("CPS action failed for event %s: %v", code, err)
	}
	return nil
}

func (e *eventService) UpdateEvent(ctx context.Context, id string, event eventdto.EventRequest) error {
	e.logger.Infof("UpdateEvent called", "event_id", id)

	if err := core.SetMerchantDetails(ctx, e.merchantService, &event); err != nil {
		e.logger.Errorf("SetMerchantDetails failed", "event_id", id, "error", err)
		return err
	}

	prevEvent, err := e.repo.FindByID(ctx, id)
	if err != nil {
		e.logger.Errorf("FindByID failed", "event_id", id, "error", err)
		return errors.New(localization.ErrorEventNotFound.Code)
	}

	var URL string
	if event.CoverImage != nil {
		URL, err = lib.UploadFileToMinio(ctx, e.minio, e.bucketName, event.CoverImage, "cover_image", e.cfg.MinioEndPoint, e.logger)
		if err != nil {
			e.logger.Errorf("UploadFileToMinio failed", "event_id", id, "error", err)
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
	} else {
		URL = prevEvent.EventInformation.Cover
	}

	curAction := core.EventMapperForUpdate(prevEvent, event, URL)

	if err := e.repo.Update(ctx, id, &curAction); err != nil {
		e.logger.Errorf("repo.Update failed", "event_id", id, "error", err)
		return err
	}

	e.logger.Infof("Event updated successfully", "event_id", id)
	err = core.HandleCPSAction(ctx, e.cpsService, id, constants.RequestUpdateEvent, curAction, *prevEvent, constants.ActionUpdate)
	if err != nil {
		e.logger.Errorf("CPS action failed for event %s: %v", curAction.EventName, err)
		return err
	}
	return nil
}

func (e *eventService) DeleteEvent(ctx context.Context, id string) error {
	prevEvent, err := e.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New(localization.ErrorEventNotFound.Code)
	}

	if err := e.repo.Delete(ctx, id); err != nil {
		return err
	}

	curData := *prevEvent
	curData.IsDeleted = true
	curData.DeletedAt = time.Now()

	err = core.HandleCPSAction(ctx, e.cpsService, id, constants.RequestDeleteEvent, curData, *prevEvent, constants.ActionDelete)
	if err != nil {
		e.logger.Errorf("CPS action failed for event %s: %v", curData.EventName, err)
		return err
	}
	return nil
}

func (e *eventService) EnableDisableEvent(ctx context.Context, id string, enable bool) error {
	e.logger.Infof("EnableDisableEvent called", "event_id", id, "enable", enable)

	prevEvent, err := e.repo.FindByID(ctx, id)
	if err != nil {
		e.logger.Errorf("FindByID failed", "event_id", id, "error", err)
		return errors.New(localization.ErrorEventNotFound.Code)
	}

	if enable && prevEvent.Enabled {
		e.logger.Warnf("Event already enabled", "event_id", id)
		return errors.New(localization.ErrorEventAlreadyEnabled.Code)
	}
	if !enable && !prevEvent.Enabled {
		e.logger.Warnf("Event already disabled", "event_id", id)
		return errors.New(localization.ErrorEventAlreadyDisabled.Code)
	}

	if err := e.repo.EnableOrDisable(ctx, id, enable); err != nil {
		e.logger.Errorf("EnableOrDisable failed", "event_id", id, "enable", enable, "error", err)
		return err
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
	err=core.HandleCPSAction(ctx, e.cpsService, id, action, curData, *prevEvent, constants.ActionUpdate)
	if err != nil {
		e.logger.Errorf("CPS action failed for event %s: %v", curData.EventName, err)
		return err
	}
	return nil
}

func (e *eventService) FetchEventByID(ctx context.Context, id string) (*model.Event, error) {
	return e.repo.FindByID(ctx, id)
}

func (e *eventService) FetchEvent(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Event], error) {
	return e.repo.FindAllWithPagination(ctx, filterParam)
}
func (e *eventService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	requestedAction := action.RequestAction

	event, ok := action.CurrentAction.(*model.Event)
	if !ok || event == nil {
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	var err error

	switch requestedAction {
	case string(constants.RequestCreateEvent):
		err = e.repo.Create(ctx, event)

	case string(constants.RequestUpdateEvent):
		if event.ID.IsZero() {
			return nil, errors.New(localization.ErrorEventIDRequired.Code)
		}
		if event.EventName == "" && event.EventVenue == "" && event.EventCity == "" {
			return nil, errors.New(localization.ErrorNoDataProvidedForUpdate.Code)
		}
		err = e.repo.Update(ctx, event.ID.Hex(), event)

	case string(constants.RequestDeleteEvent):
		if event.ID.IsZero() {
			return nil, errors.New(localization.ErrorEventIDRequired.Code)
		}
		err = e.repo.Delete(ctx, event.ID.Hex())

	case string(constants.RequestEnableEvent):
		if event.ID.IsZero() {
			return nil, errors.New(localization.ErrorEventIDRequired.Code)
		}
		err = e.repo.EnableOrDisable(ctx, event.ID.Hex(), true)

	case string(constants.RequestDisableEvent):
		if event.ID.IsZero() {
			return nil, errors.New(localization.ErrorEventIDRequired.Code)
		}
		err = e.repo.EnableOrDisable(ctx, event.ID.Hex(), false)

	default:
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	if err != nil {
		return nil, err
	}

	action.CurrentAction = event
	return action, nil
}
