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
	if event.IsDeleted {
		event.IsDeleted = false
	}
	if event.Enabled {
		event.Enabled = false
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

func EventDocumentToUpdateBsonM(event model.EventDocument) bson.M {
	update := bson.M{
		"last_modified_at": time.Now(),
	}

	if event.EventCode != "" {
		update["event_code"] = event.EventCode
	}
	if event.EventName != "" {
		update["event_name"] = event.EventName
	}
	if event.EventCity != "" {
		update["event_city"] = event.EventCity
	}
	if event.AccountNumber != "" {
		update["account_number"] = event.AccountNumber
	}
	if event.EventVenue != "" {
		update["event_venue"] = event.EventVenue
	}

	if len(event.RefundPolicy) > 0 {
		update["refund_policy"] = event.RefundPolicy
	}
	if len(event.MICSInfo) > 0 {
		update["mics_info"] = event.MICSInfo
	}
	if len(event.Ticket) > 0 {
		update["ticket"] = event.Ticket
	}

	if event.Restriction.Type != "" || event.Restriction.Description != "" {
		restriction := bson.M{}
		if event.Restriction.Type != "" {
			restriction["type"] = event.Restriction.Type
		}
		if event.Restriction.Description != "" {
			restriction["description"] = event.Restriction.Description
		}
		update["restriction"] = restriction
	}

	// Status
	if string(event.Status) != "" {
		update["status"] = event.Status
	}

	info := bson.M{}
	if !event.EventInformation.StartDate.IsZero() {
		info["start_date"] = event.EventInformation.StartDate
	}
	if !event.EventInformation.DueDate.IsZero() {
		info["due_date"] = event.EventInformation.DueDate
	}
	if event.EventInformation.Description != "" {
		info["description"] = event.EventInformation.Description
	}
	if event.EventInformation.Cover != "" {
		info["cover"] = event.EventInformation.Cover
	}
	if event.EventInformation.VideoLink != "" {
		info["video_link"] = event.EventInformation.VideoLink
	}
	if len(info) > 0 {
		update["event_information"] = info
	}

	stats := bson.M{}
	if event.TicketStatistics.Category != "" {
		stats["category"] = event.TicketStatistics.Category
	}
	if event.TicketStatistics.Revenue != 0 {
		stats["revenue"] = event.TicketStatistics.Revenue
	}
	if event.TicketStatistics.NumberOfSoldTicket != 0 {
		stats["number_of_sold_ticket"] = event.TicketStatistics.NumberOfSoldTicket
	}
	if len(stats) > 0 {
		update["ticket_statistics"] = stats
	}

	tInfo := bson.M{}
	if event.TicketInformation.TotalNumberOfTicket != 0 {
		tInfo["total_number_of_ticket"] = event.TicketInformation.TotalNumberOfTicket
	}
	if event.TicketInformation.TotalNumberOfAvailableTicket != 0 {
		tInfo["total_number_of_available_ticket"] = event.TicketInformation.TotalNumberOfAvailableTicket
	}
	if event.TicketInformation.TotalNumberOFUnsoldTicket != 0 {
		tInfo["total_number_of_unsold_ticket"] = event.TicketInformation.TotalNumberOFUnsoldTicket
	}
	if len(tInfo) > 0 {
		update["ticket_information"] = tInfo
	}

	// MerchantInformation struct
	mInfo := bson.M{}
	if event.MerchantInformation.MerchantID != "" {
		mInfo["merchant_id"] = event.MerchantInformation.MerchantID
	}
	if event.MerchantInformation.MercahntName != "" {
		mInfo["merchant_name"] = event.MerchantInformation.MercahntName
	}
	if event.MerchantInformation.MerchantPhoneNumber != "" {
		mInfo["merchant_phone_number"] = event.MerchantInformation.MerchantPhoneNumber
	}
	if event.MerchantInformation.MerchantEmail != "" {
		mInfo["merchant_email"] = event.MerchantInformation.MerchantEmail
	}
	if len(mInfo) > 0 {
		update["merchant_information"] = mInfo
	}

	update["has_restriction"] = event.HasRestriction

	return update
}
