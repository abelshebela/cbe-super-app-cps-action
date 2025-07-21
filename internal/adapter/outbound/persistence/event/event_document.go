package event

import (
	"fmt"
	"time"

	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	event "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CPSActionDocument struct {
	ID                 bson.ObjectID          `json:"id" bson:"_id,omitempty"`
	ActionCode         string                 `json:"action_code" bson:"action_code"`
	UniqueId           string                 `json:"unique_id" bson:"unique_id"`
	MakerID            string                 `json:"maker_id" bson:"maker_id"`
	MakerName          string                 `json:"maker_name" bson:"maker_name"`
	MakerPhoneNumber   string                 `json:"maker_phone_number" bson:"maker_phone_number"`
	CheckerID          string                 `json:"checker_id" bson:"checker_id"`
	CheckerName        string                 `json:"checker_name" bson:"checker_name"`
	CheckerPhoneNumber string                 `json:"checker_phone_number" bson:"checker_phone_number"`
	Department         string                 `json:"department" bson:"department"`
	RejectionReason    *string                `json:"rejection_reason" bson:"rejection_reason,omitempty"`
	PreviousAction     any                    `json:"previos_action" bson:"previos_action"`
	CurrentAction      any                    `json:"current_action" bson:"current_action"`
	ActionStatus       entities.ActionStatus  `json:"action_status" bson:"action_status"`
	ActionType         entities.ActionType    `json:"action_type" bson:"action_type"`
	RequestAction      entities.RequestAction `json:"request_action" bson:"request_action"`
	CreatedAt          time.Time              `json:"created_at" bson:"created_at"`
	LastModifiedAt     time.Time              `json:"last_modified_at" bson:"last_modified_at"`
	MakerActionTime    time.Time              `json:"maker_action_time" bson:"maker_action_time"`
	CheckerActionTime  time.Time              `json:"checker_action_time" bson:"checker_action_time"`
}

func (action *CPSActionDocument) toModel() entities.CPSAction {
	var currentAction any
	if action.CurrentAction != nil {
		if bsonD, ok := action.CurrentAction.(bson.D); ok {
			m := make(map[string]any, len(bsonD))
			for _, elem := range bsonD {
				m[elem.Key] = elem.Value
			}
			currentAction = m
		} else {
			currentAction = action.CurrentAction
		}
	}

	return entities.CPSAction{
		ID:                 action.ID.Hex(),
		ActionCode:         action.ActionCode,
		UniqueId:           action.UniqueId,
		MakerID:            action.MakerID,
		MakerName:          action.MakerName,
		MakerPhoneNumber:   action.MakerPhoneNumber,
		CheckerID:          action.CheckerID,
		CheckerName:        action.CheckerName,
		CheckerPhoneNumber: action.CheckerPhoneNumber,
		Department:         action.Department,
		RejectionReason:    action.RejectionReason,
		PreviousAction:     action.PreviousAction,
		CurrentAction:      currentAction,
		ActionStatus:       action.ActionStatus,
		ActionType:         action.ActionType,
		RequestAction:      action.RequestAction,
		CreatedAt:          action.CreatedAt,
		LastModifiedAt:     action.LastModifiedAt,
		MakerActionTime:    action.MakerActionTime,
		CheckerActionTime:  &action.CheckerActionTime,
	}
}

func ToCpsActionDocument(cpsAction entities.CPSAction) (*CPSActionDocument, error) {
	var objectID bson.ObjectID

	if cpsAction.ID == "" {
		objectID = bson.NewObjectID()
	} else {
		id, err := bson.ObjectIDFromHex(cpsAction.ID)
		if err != nil {
			return nil, fmt.Errorf(error_codes.InvalidID)
		}

		objectID = id
	}

	return &CPSActionDocument{
		ID:                 objectID,
		ActionCode:         cpsAction.ActionCode,
		UniqueId:           cpsAction.UniqueId,
		MakerID:            cpsAction.MakerID,
		MakerName:          cpsAction.MakerName,
		MakerPhoneNumber:   cpsAction.MakerPhoneNumber,
		CheckerID:          cpsAction.CheckerID,
		CheckerName:        cpsAction.CheckerName,
		CheckerPhoneNumber: cpsAction.CheckerPhoneNumber,
		Department:         cpsAction.Department,
		RejectionReason:    cpsAction.RejectionReason,
		PreviousAction:     cpsAction.PreviousAction,
		CurrentAction:      cpsAction.CurrentAction,
		ActionStatus:       cpsAction.ActionStatus,
		ActionType:         cpsAction.ActionType,
		RequestAction:      cpsAction.RequestAction,
		CreatedAt:          cpsAction.CreatedAt,
		LastModifiedAt:     cpsAction.LastModifiedAt,
		MakerActionTime:    cpsAction.MakerActionTime,
		CheckerActionTime:  *cpsAction.CheckerActionTime,
	}, nil
}

type EventDocument struct {
	ID                  bson.ObjectID
	Code                string
	Name                string
	Address             event.Address
	AccountNumber       string
	EventVenue          string
	RefundPolicy        []string
	MICSInfo            []string
	Restriction         event.Restriction
	Ticket              event.Ticket
	Status              event.EventStatus
	EventInformation    event.EventInformation
	TicketStatistics    event.TicketStatistics
	TicketInformation   event.TicketInformation
	MerchantInformation event.MerchantInformation
	Enabled             bool
	IsDeleted           bool
	HasRestriction      bool
	CreatedAt           time.Time
	DeletedAt           time.Time
	LastModifiedAt      time.Time
}

func (e *EventDocument) toModel() event.Event {
	return event.Event{
		ID:                  e.ID.Hex(),
		Code:                e.Code,
		Name:                e.Name,
		Address:             e.Address,
		AccountNumber:       e.AccountNumber,
		EventVenue:          e.EventVenue,
		RefundPolicy:        e.RefundPolicy,
		MICSInfo:            e.MICSInfo,
		Restriction:         e.Restriction,
		Ticket:              e.Ticket,
		Status:              e.Status,
		EventInformation:    e.EventInformation,
		TicketStatistics:    e.TicketStatistics,
		TicketInformation:   e.TicketInformation,
		MerchantInformation: e.MerchantInformation,
		Enabled:             e.Enabled,
		IsDeleted:           e.IsDeleted,
		HasRestriction:      e.HasRestriction,
		CreatedAt:           e.CreatedAt,
		DeletedAt:           e.DeletedAt,
		LastModifiedAt:      e.LastModifiedAt,
	}
}

func ToEventDocument(event event.Event) (*EventDocument, error) {
	var objectID bson.ObjectID

	if event.ID == "" {
		objectID = bson.NewObjectID()
	} else {
		id, err := bson.ObjectIDFromHex(event.ID)
		if err != nil {
			return nil, fmt.Errorf(error_codes.InvalidID)
		}
		objectID = id
	}

	return &EventDocument{
		ID:                  objectID,
		Code:                event.Code,
		Name:                event.Name,
		Address:             event.Address,
		AccountNumber:       event.AccountNumber,
		EventVenue:          event.EventVenue,
		RefundPolicy:        event.RefundPolicy,
		MICSInfo:            event.MICSInfo,
		Restriction:         event.Restriction,
		Ticket:              event.Ticket,
		Status:              event.Status,
		EventInformation:    event.EventInformation,
		TicketStatistics:    event.TicketStatistics,
		TicketInformation:   event.TicketInformation,
		MerchantInformation: event.MerchantInformation,
		Enabled:             event.Enabled,
		IsDeleted:           event.IsDeleted,
		HasRestriction:      event.HasRestriction,
		CreatedAt:           event.CreatedAt,
		DeletedAt:           event.DeletedAt,
		LastModifiedAt:      event.LastModifiedAt,
	}, nil
}
