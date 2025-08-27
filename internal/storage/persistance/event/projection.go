package event

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func EventMapper(e *model.EventDocument) model.Event {
	return model.Event{
		ID:                  e.ID,
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

func EventDocumentToBsonM(event model.EventDocument) bson.M {
	return bson.M{
		"event_code":           event.EventCode,
		"event_name":           event.EventName,
		"event_city":           event.EventCity,
		"account_number":       event.AccountNumber,
		"event_venue":          event.EventVenue,
		"refund_policy":        event.RefundPolicy,
		"mics_info":            event.MICSInfo,
		"restriction":          event.Restriction,
		"ticket":               event.Ticket,
		"status":               event.Status,
		"event_information":    event.EventInformation,
		"ticket_statistics":    event.TicketStatistics,
		"ticket_information":   event.TicketInformation,
		"merchant_information": event.MerchantInformation,
		"enabled":              event.Enabled,
		"is_deleted":           event.IsDeleted,
		"has_restriction":      event.HasRestriction,
		"created_at":           event.CreatedAt,
		"deleted_at":           event.DeletedAt,
		"last_modified_at":     event.LastModifiedAt,
	}
}
func EventDocumentMapper(event model.Event) *model.EventDocument {
	objectID := event.ID
	if objectID == bson.NilObjectID {
		objectID = bson.NewObjectID()
	}

	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}
	if event.LastModifiedAt.IsZero() {
		event.LastModifiedAt = time.Now()
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
	}
}
