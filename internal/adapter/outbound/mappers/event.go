package mappers

import (
	"fmt"

	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	event "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func  ToEventModel(e *model.EventDocument) event.Event {
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

func ToEventDocument(event event.Event) (*model.EventDocument, error) {
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

	return &model.EventDocument{
		ID:                  objectID,
		EventCode:           event.EventCode,
		EventName:           event.EventName,
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
