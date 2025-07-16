// Package event provides services and interfaces for handling event-related business logic.
package event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	common_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
)

type EventService interface {
	CreateEventRequest(ctx context.Context, event Event, ticket Ticket, maker Maker) (string, error)
	ApproveEventRequest(ctx context.Context, action_id string, action_taken bool, checkerID, checkerName, checkerPhone string) error
	FetchEvent(ctx context.Context, limit, offset int) ([]*Event, error)
	FetchEventByID(ctx context.Context, event_id string) (*Event, error)
}

type Service struct {
	Repository EventRepository
}

func NewEventService(repository EventRepository) EventService {
	return &Service{
		Repository: repository,
	}
}

func (s *Service) CreateEventRequest(ctx context.Context, event Event, ticket Ticket, maker Maker) (string, error) {
	actionID := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})

	a := action.CPSAction{
		ActionCode:       actionID,
		MakerID:          maker.ID,
		MakerName:        maker.FullName,
		MakerPhoneNumber: maker.PhoneNumber,
		ActionType:       action.ActionCreate,
		RequestAction:    action.RequestUpdateServiceRule,
		ActionStatus:     action.ActionPending,
		Department:       maker.Department,
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
	// check if pending action exist
	err := s.Repository.CPSActionExists(ctx, action.CreateCPSAction{
		Department:    a.Department,
		RequestAction: a.RequestAction,
		MakerUser: action.User{
			UserID:      a.MakerID,
			FullName:    a.MakerName,
			PhoneNumber: a.MakerPhoneNumber,
			Department:  maker.Department,
		},
	})
	if err != nil {
		return "", err
	}
	data, err := s.Repository.CreateCpsAction(ctx, a)
	if err != nil {
		return "", err
	}
	return data.ActionCode, nil
}
func (s *Service) ApproveEventRequest(ctx context.Context, actionID string, actionTaken bool, checkerID, checkerName, checkerPhone string) error {
	// Input validation
	if actionID == "" || checkerID == "" || checkerName == "" || checkerPhone == "" {
		return errors.New(common_utils.InvalidInput)
	}

	cpsAction, err := s.Repository.FetchCpsActionByID(ctx, actionID)
	if err != nil {
		return errors.New(err.Error())
	}

	// Update CPSAction status and checker details
	if actionTaken {
		cpsAction.ActionStatus = action.ActionApproved
	} else {
		cpsAction.ActionStatus = action.ActionRejected
	}
	cpsAction.CheckerID = checkerID
	cpsAction.CheckerName = checkerName
	cpsAction.CheckerPhoneNumber = checkerPhone

	// Check for nil CurrentAction
	if cpsAction.CurrentAction == nil {
		return fmt.Errorf("CurrentAction is nil for actionID %s", actionID)
	}

	// Type assertion for CurrentAction
	var e Event
	switch v := cpsAction.CurrentAction.(type) {
	case Event:
		e = v
	case map[string]interface{}:
		// Extract the "event" key from the map
		eventData, ok := v["event"]
		if !ok {
			return fmt.Errorf("CurrentAction map does not contain 'event' key for actionID %s; map contents: %+v", actionID, v)
		}
		// Convert eventData to JSON and unmarshal into Event struct
		eventBytes, err := json.Marshal(eventData)
		if err != nil {
			return fmt.Errorf("failed to marshal event data to JSON for actionID %s: %w", actionID, err)
		}
		if err := json.Unmarshal(eventBytes, &e); err != nil {
			return fmt.Errorf("failed to unmarshal event data to Event for actionID %s: %w", actionID, err)
		}
	default:
		return fmt.Errorf("failed to assert CurrentAction to Event type for actionID %s; actual type: %T, value: %+v", actionID, cpsAction.CurrentAction, cpsAction.CurrentAction)
	}

	// Use consistent timestamp
	now := time.Now()
	cpsAction.LastModifiedAt = now
	err = s.Repository.UpdateCpsAction(ctx, *cpsAction)
	if err != nil {
		return err
	}

	e.LastModifiedAt = now
	e.Enabled = actionTaken // Set Enabled based on actionTaken
	_, err = s.Repository.CreateEvent(ctx, e)
	if err != nil {
		return fmt.Errorf("failed to create Event for actionID %s: %w", actionID, err)
	}

	return nil
}
func (s *Service) FetchEvent(ctx context.Context, limit, offset int) ([]*Event, error) {
	data, err := s.Repository.FetchEvent(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) FetchEventByID(ctx context.Context, event_id string) (*Event, error) {
	event, err := s.Repository.FetchEventByID(ctx, event_id)
	if err != nil {
		return nil, nil
	}
	return event, nil
}
