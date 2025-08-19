package event

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type eventService struct {
	repo   storage.EventRepository
	logger utils.Logger
}

func NewEventService(repo storage.EventRepository, logger utils.Logger) service.EventService {
	return &eventService{
		repo:   repo,
		logger: logger,
	}
}

func (e *eventService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {

	e.logger.Infof("Event service authorizing action: %s", cpsAction.ActionCode)

	// For now, return the action as approved
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}
