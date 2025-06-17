package event

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
)

func (s *Service) CreateEventRequest(ctx context.Context, event Event, ticket Ticket, makerID, makerName, makerPhone string) (string, error) {
	actionId := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})

	a := action.CPSAction{
		ActionCode: actionId,
		Maker: action.User{
			UserID:      makerID,
			FullName:    makerName,
			PhoneNumber: makerPhone,
			Timestamp:   time.Now(),
		},
		ActionType:    action.ActionCreate,
		RequestAction: action.RequestUpdateServiceRule,
		ActionStatus:  action.ActionPending,
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
func (s *Service) ApproveEventRequest(ctx context.Context, action_id string, action_taken bool, checkerID, checkerName, checkerPhone string) error {
	cps_action, err := s.Repository.FetchCpsActionById(ctx, action_id)
	if err != nil {
		return err
	}
	if action_taken {
		cps_action.ActionStatus = action.ActionApproved
	} else {
		cps_action.ActionStatus = action.ActionRejected
	}
	cps_action.Checker = action.User{
		UserID: checkerID,
	}
	var e Event
	e, ok := cps_action.CurrentAction.(Event)
	if !ok {
		return fmt.Errorf("failed to assert CurrentAction to Event type")
	}

	cps_action.LastModifiedAt = time.Now()
	err = s.Repository.UpdateCpsAction(ctx, cps_action)
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
