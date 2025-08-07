// Package event provides services and interfaces for handling event-related business logic.
package event

import (
	"context"
	"fmt"
	"time"

	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type EventService interface {
	CreateEvent(ctx context.Context, event EventRequest) (*Event, error)
	UpdateEvent(ctx context.Context, id string, event EventRequest) (*Event, *Event, error)
	DeleteEvent(ctx context.Context, id string) (*Event, *Event, error)
	EnableDisableEvent(ctx context.Context, id string, enable bool) (*Event, *Event, error)

	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)

	FetchEventByID(ctx context.Context, id string) (*Event, error)
	FetchEvent(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*Event], error)
}

type Service struct {
	Repository EventRepository
	logger     shared_utils.Logger
	minio      config.MinioClientInterface
	bucketName string
	cfg        *config.VaultConfig
}

func NewEventService(repository EventRepository,
	minio config.MinioClientInterface,
	bucketName string,
	cfg *config.VaultConfig,
	logger shared_utils.Logger,
) EventService {
	return &Service{
		Repository: repository,
		logger:     logger,
		minio:      minio,
		cfg: cfg,
		bucketName: bucketName,
	}
}

func (e *Service) CreateEvent(ctx context.Context, event EventRequest) (*Event, error) {
	exist, err := e.Repository.EventNameExists(ctx, event.EventName, nil)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, fmt.Errorf("EVENT_NAME_ALREADY_EXISTS")
	}

	code, err := common_util.GeneratePrefixedName("EVE", event.EventName, e.logger)
	if err != nil {
		return nil, fmt.Errorf(common_util.UnhandledServerError)
	}

	if event.CoverImage == nil || e.cfg == nil {
		return nil, fmt.Errorf("Nill")
	}

	URL, err := common_util.UploadFileToMinio(ctx, e.minio, e.bucketName, event.CoverImage, "cover_image", e.cfg.MinioEndPoint, e.logger)
	if err != nil {
		return nil, err
	}
	result := Event{
		EventCode:  code,
		EventName:  event.EventName,
		EventCity:  event.EventCity,
		EventVenue: event.EventVenue,
		Status:     EventUpcomming,
		AccountNumber: event.AccountNumber,
		MerchantInformation: MerchantInformation{
			MerchantID:          event.MerchantID,
			MercahntName:        event.MercahntName,
			MerchantPhoneNumber: event.MerchantPhoneNumber,
			MerchantEmail:       event.MerchantEmail,
		},
		EventInformation: EventInformation{
			StartDate:   event.StartDate,
			DueDate:     event.DueDate,
			Description: event.EventDescription,
			Cover:       URL,
		},
		TicketInformation: TicketInformation{
			TotalNumberOfTicket: uint64(event.TotalTicketCount),
			TotalNumberOfAvailableTicket: uint64(event.TotalTicketCount),
		},
		Ticket: event.Tickets,
		CreatedAt: time.Now(),
		LastModifiedAt: time.Now(),
	}

	return &result, nil
}

// UpdateEvent updates an event, using prevEvent values for empty fields in event
func (e *Service) UpdateEvent(ctx context.Context, id string, event EventRequest) (*Event, *Event, error) {
	e.logger.Infof("Updating event", "id", id)

	// Fetch previous event
	prevEvent, err := e.Repository.FetchEventByID(ctx, id)
	if err != nil {
		e.logger.Errorf("Failed to fetch event", "id", id, "error", err)
		return nil, nil, err
	}

	var URL string
	if event.CoverImage != nil {
		URL, err = common_util.UploadFileToMinio(ctx, e.minio, e.bucketName, event.CoverImage, "cover_image", e.cfg.MinioEndPoint, e.logger)
		if err != nil {
			e.logger.Errorf("Failed to upload cover image to MinIO", "error", err)
			return nil, nil, err
		}
	} else {
		URL = prevEvent.EventInformation.Cover
		e.logger.Infof("Using previous cover image", "url", URL)
	}

	if event.EventName != "" {
		exist, err := e.Repository.EventNameExists(ctx, event.EventName, nil)
		if err != nil {
			return nil, nil, err
		}

		if exist {
			return nil, nil, fmt.Errorf(common_util.EventNameAlreadyExists)
		}
	}

	curAction := Event{
		ID:         prevEvent.ID,
		EventCode:  prevEvent.EventCode,
		EventName:  nonEmptyString(event.EventName, prevEvent.EventName),
		EventCity:  nonEmptyString(event.EventCity, prevEvent.EventCity),
		EventVenue: nonEmptyString(event.EventVenue, prevEvent.EventVenue),
		Status:     EventUpcomming,
		MerchantInformation: MerchantInformation{
			MerchantID:          nonEmptyString(event.MerchantID, prevEvent.MerchantInformation.MerchantID),
			MercahntName:        nonEmptyString(event.MercahntName, prevEvent.MerchantInformation.MercahntName),
			MerchantPhoneNumber: nonEmptyString(event.MerchantPhoneNumber, prevEvent.MerchantInformation.MerchantPhoneNumber),
			MerchantEmail:       nonEmptyString(event.MerchantEmail, prevEvent.MerchantInformation.MerchantEmail),
		},
		EventInformation: EventInformation{
			StartDate:   nonZeroTime(event.StartDate, prevEvent.EventInformation.StartDate),
			DueDate:     nonZeroTime(event.DueDate, prevEvent.EventInformation.DueDate),
			Description: nonEmptyString(event.EventDescription, prevEvent.EventInformation.Description),
			Cover:       URL,
		},
		TicketInformation: TicketInformation{
			TotalNumberOfTicket: nonZeroUint64(uint64(event.TotalTicketCount), prevEvent.TicketInformation.TotalNumberOfTicket),
		},
		Ticket: nonEmptyTickets(event.Tickets, prevEvent.Ticket),
		CreatedAt: time.Now(),
		LastModifiedAt: time.Now(),
		AccountNumber: nonEmptyString(event.AccountNumber, prevEvent.AccountNumber),
	}

	e.logger.Infof("Updated event", "id", id, "merchantAppID", curAction.MerchantInformation.MerchantID)
	return &curAction, generateEvent(*prevEvent), nil
}

func (e *Service) DeleteEvent(ctx context.Context, id string) (*Event, *Event, error) {
	e.logger.Infof("Deleting event", "id", id)
	// Fetch previous event
	prevEvent, err := e.Repository.FetchEventByID(ctx, id)
	if err != nil {
		e.logger.Errorf("Failed to fetch event", "id", id, "error", err)
		return nil, nil, err
	}

	curData := generateEvent(*prevEvent)
	curData.IsDeleted = true
	curData.DeletedAt = time.Now()

	return curData, generateEvent(*prevEvent), nil
}
func (e *Service) EnableDisableEvent(ctx context.Context, id string, enable bool) (*Event, *Event, error) {
	e.logger.Infof("EnableDisable event", "id", id)

	// Fetch previous event
	prevEvent, err := e.Repository.FetchEventByID(ctx, id)
	if err != nil {
		e.logger.Errorf("Failed to fetch event", "id", id, "error", err)
		return nil, nil, err
	}

	if enable && prevEvent.Enabled {
		return nil, nil, fmt.Errorf(common_util.ErrAlreadyEnabled)
	}

	if !enable && !prevEvent.Enabled {
		return nil, nil, fmt.Errorf(common_util.ErrAlreadyDisabled)
	}
	curData := generateEvent(*prevEvent)
	curData.Enabled = enable
	curData.LastModifiedAt = time.Now()

	return curData, generateEvent(*prevEvent), nil
}

func (e *Service) FetchEventByID(ctx context.Context, id string) (*Event, error) {
	e.logger.Infof("FetchEventByID", "id", id)
	return e.Repository.FetchEventByID(ctx, id)
}
func (e *Service) FetchEvent(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*Event], error) {
	return e.Repository.FetchEvent(ctx, filterParam)
}

func generateEvent(event Event) *Event {
	return &Event{
		ID:         event.ID,
		EventCode:  event.EventCode,
		EventName:  event.EventName,
		EventCity:  event.EventCity,
		EventVenue: event.EventVenue,
		Status:     EventUpcomming,
		AccountNumber: event.AccountNumber,
		MerchantInformation: MerchantInformation{
			MerchantID:          event.MerchantInformation.MerchantID,
			MercahntName:        event.MerchantInformation.MercahntName,
			MerchantPhoneNumber: event.MerchantInformation.MerchantPhoneNumber,
			MerchantEmail:       event.MerchantInformation.MerchantEmail,
		},
		EventInformation: EventInformation{
			StartDate:   event.EventInformation.StartDate,
			DueDate:     event.EventInformation.DueDate,
			Description: event.EventInformation.Description,
			Cover:       event.EventInformation.Cover,
		},
		TicketInformation: TicketInformation{
			TotalNumberOfTicket: event.TicketInformation.TotalNumberOfTicket,
		},
		Ticket: event.Ticket,
		CreatedAt: event.CreatedAt,
		LastModifiedAt: event.LastModifiedAt,
	}
}

func (e *Service) Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	requestedAction := action.RequestAction

	var event *Event
	var err error

	bindErr := common_util.BindAction(action.CurrentAction, &event)
	if bindErr != nil {
		e.logger.Errorf("failed to bind current action to event: %v", bindErr)
		return nil, fmt.Errorf(common_util.InvalidActionData)
	}

	switch requestedAction {
	case cps_const.RequestCreateEvent:

		event, err = e.Repository.CreateEvent(ctx, *event)
		if err != nil {
			return nil, err
		}

	case cps_const.RequestUpdateEvent:
		event, err = e.Repository.UpdateEvent(ctx, *event)
		if err != nil {
			return nil, err
		}
	case cps_const.RequestDeleteEvent:
		event, err = e.Repository.DeleteEvent(ctx, event.ID)
		if err != nil {
			return nil, err
		}
	case cps_const.RequestEnableEvent:
		event, err = e.Repository.EnableDisableEvent(ctx, event.ID, true)
		if err != nil {
			return nil, err
		}
	case cps_const.RequestDisableEvent:
		event, err = e.Repository.EnableDisableEvent(ctx, event.ID, false)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf(common_util.ErrUnsupported)
	}

	action.CurrentAction = event
	return action, nil
}





































