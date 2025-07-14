package event

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
)

func (s *Service) CreateEventRequest(ctx context.Context, event Event, ticket Ticket, makerID, makerName, makerPhone string) (string, error) {
	actionId := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})

	a := action.CPSAction{
		ActionCode:       actionId,
		MakerID:          makerID,
		MakerName:        makerName,
		MakerPhoneNumber: makerPhone,
		ActionType:       action.ActionCreate,
		RequestAction:    action.RequestUpdateServiceRule,
		ActionStatus:     action.ActionPending,
		CurrentAction: struct {
			Ticket Ticket
			Event  Event
		}{
			Ticket: ticket,
			Event:  event,
		},
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
	}
	data, err := s.Repository.CreateCpsAction(ctx, a)
	if err != nil {
		return "", err
	}
	return data.ActionCode, nil
}
func (s *Service) ApproveEventRequest(ctx context.Context, actionID string, actionTaken bool, checkerID, checkerName, checkerPhone string) error {
	cpsAction, err := s.Repository.FetchCpsActionById(ctx, actionID)
	if err != nil {
		return err
	}
	if actionTaken {
		cpsAction.ActionStatus = action.ActionApproved
	} else {
		cpsAction.ActionStatus = action.ActionRejected
	}
	cpsAction.CheckerID = checkerID
	cpsAction.CheckerName = checkerName
	cpsAction.CheckerPhoneNumber = checkerPhone

	var e Event
	switch v := cpsAction.CurrentAction.(type) {
	case Event:
		e = v
	case map[string]interface{}:

		return fmt.Errorf("CurrentAction is a map, manual decoding required")
	default:
		return fmt.Errorf("failed to assert CurrentAction to Event type")
	}

	cpsAction.LastModifiedAt = time.Now()
	err = s.Repository.UpdateCpsAction(ctx, cpsAction)
	if err != nil {
		return err
	}
	e.LastModifiedAt = time.Now()
	e.Enabled = true
	_, err = s.Repository.CreateEvent(ctx, e)
	if err != nil {
		return fmt.Errorf("error: %v", err.Error())
	}
	return nil
}
func (s *Service) FetchEvent(ctx context.Context, limit, offset int) ([]Event, error) {
	data, err := s.Repository.FetchEvent(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) FetchEventById(ctx context.Context, event_id string) (Event, error) {
	event, err := s.Repository.FetchEventById(ctx, event_id)
	if err != nil {
		return Event{}, nil
	}
	return event, nil
}
