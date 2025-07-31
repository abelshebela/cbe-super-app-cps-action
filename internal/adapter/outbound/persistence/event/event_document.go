package event

import (
	"fmt"
	"time"

	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	event "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type EventDocument struct {
	ID                  bson.ObjectID         `bson:"_id,omitempty"`
	EventCode           string                     `bson:"event_code"`
	EventName           string                     `bson:"event_name"`
	EventCity           string                     `bson:"event_city"`
	AccountNumber       string                     `bson:"account_number"`
	EventVenue          string                     `bson:"event_venue"`
	RefundPolicy        []string                   `bson:"refund_policy"`
	MICSInfo            []string                   `bson:"mics_info"`
	Restriction         event.Restriction          `bson:"restriction"`
	Ticket              []event.Ticket             `bson:"ticket"`
	Status              event.EventStatus          `bson:"status"`
	EventInformation    event.EventInformation     `bson:"event_information"`
	TicketStatistics    event.TicketStatistics     `bson:"ticket_statistics"`
	TicketInformation   event.TicketInformation    `bson:"ticket_information"`
	MerchantInformation event.MerchantInformation  `bson:"merchant_information"`
	Enabled             bool                       `bson:"enabled"`
	IsDeleted           bool                       `bson:"is_deleted"`
	HasRestriction      bool                       `bson:"has_restriction"`
	CreatedAt           time.Time                  `bson:"created_at"`
	DeletedAt           time.Time                  `bson:"deleted_at"`
	LastModifiedAt      time.Time                  `bson:"last_modified_at"`
}

func (e *EventDocument) toModel() event.Event {
	return event.Event{
		ID:                  e.ID.Hex(),
		EventCode:           e.EventCode,
		EventName:           e.EventName,
		EventCity:           e.EventCity,
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
		ID:        objectID,
		EventCode: event.EventCode,
		EventName: event.EventName,
		EventCity:           event.EventCity,
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
