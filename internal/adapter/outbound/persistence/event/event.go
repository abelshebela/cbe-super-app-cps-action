package event

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
	event_outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/event"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type EventPersistence struct {
	eventDal dal.MongoDal[EventDocument, EventDocument]
	logger   utils.Logger
}

func InitEventPersistence(client *mongo.Client, dbName string, collections string, logger utils.Logger) event_outbound.EventRepository {
	return &EventPersistence{
		eventDal: dal.NewMongoDal[EventDocument, EventDocument](client, dbName, collections),
		logger:   logger,
	}
}

func (e *EventPersistence) CreateEvent(ctx context.Context, event event.Event) (*event.Event, error) {
	eventDoc, err := ToEventDocument(event)
	if err != nil {
		return nil, err
	}

	res, err := e.eventDal.InsertOne(ctx, *eventDoc)
	if err != nil {
		return nil, fmt.Errorf(common_util.GeneralDBInsertFailed)
	}
	result := res.toModel()
	return &result, nil
}

func (e *EventPersistence) FetchEventByID(ctx context.Context, id string) (*event.Event, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	eventDoc, err := e.eventDal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf(common_util.NotFound)
		}
		return nil, fmt.Errorf(common_util.GeneralDBQueryFailed)
	}

	result := eventDoc.toModel()
	return &result, nil
}

func (e *EventPersistence) FetchEvent(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*event.Event], error) {
	filter := bson.M{"is_deleted": false}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"event_code": searchRegex},
			{"event_name": searchRegex},
			{"event_city": searchRegex},
			{"account_number": searchRegex},
			{"event_venue": searchRegex},
		}
	}

	skip := (filterParam.Page - 1) * filterParam.PerPage
	limit := filterParam.PerPage

	eventDocs, err := e.eventDal.FindAllWithPagination(ctx, filter, bson.M{}, int64(skip), int64(limit))
	if err != nil {
		return nil, fmt.Errorf(common_util.GeneralDBQueryFailed)
	}

	var events []*event.Event
	for _, doc := range eventDocs {
		converted := doc.toModel()
		events = append(events, &converted)
	}

	total, err := e.eventDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	meta := common_util.BuildPaginationMeta(total, filterParam.Page, limit)
	return &common_util.PaginatedResponse[[]*event.Event]{
		Data: events,
		Meta: meta,
	}, nil
}

func (e *EventPersistence) UpdateEvent(ctx context.Context, event event.Event) (*event.Event, error) {
	objID, err := common_util.ParsePrimitiveObjectID(event.ID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}

	update := bson.M{}
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
	if event.Restriction.Type != "" || event.Restriction.Description != "" {
		update["restriction"] = event.Restriction
	}
	if len(event.Ticket) > 0 {
		update["ticket"] = event.Ticket
	}
	if event.Status != "" {
		update["status"] = event.Status
	}

	if !event.EventInformation.StartDate.IsZero() ||
		!event.EventInformation.DueDate.IsZero() ||
		event.EventInformation.Description != "" ||
		event.EventInformation.Cover != "" ||
		event.EventInformation.VideoLink != "" {
		update["event_information"] = event.EventInformation
	}
	if event.TicketStatistics.Category != "" ||
		event.TicketStatistics.Revenue != 0 ||
		event.TicketStatistics.NumberOfSoldTicket != 0 {
		update["ticket_statistics"] = event.TicketStatistics
	}
	if event.TicketInformation.TotalNumberOfTicket != 0 ||
		event.TicketInformation.TotalNumberOfAvailableTicket != 0 ||
		event.TicketInformation.TotalNumberOFUnsoldTicket != 0 {
		update["ticket_information"] = event.TicketInformation
	}
	if event.MerchantInformation.MerchantID != "" ||
		event.MerchantInformation.MercahntName != "" ||
		event.MerchantInformation.MerchantPhoneNumber != "" ||
		event.MerchantInformation.MerchantEmail != "" {
		update["merchant_information"] = event.MerchantInformation
	}
	update["has_restriction"] = event.HasRestriction
	update["last_modified_at"] = time.Now()

	if len(update) == 1 {
		return nil, fmt.Errorf(common_util.NoDataProvidedForUpdate)
	}

	eventDoc, err := e.eventDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf(common_util.NotFound)
		}
		e.logger.Warnf(err.Error(), "while updating event")
		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	result := eventDoc.toModel()
	return &result, nil
}

func (e *EventPersistence) DeleteEvent(ctx context.Context, id string) (*event.Event, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		return nil, fmt.Errorf(common_util.InvalidID)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"is_deleted": true, "deleted_at": time.Now()}

	eventDoc, err := e.eventDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf(common_util.NotFound)
		}
		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	result := eventDoc.toModel()
	return &result, nil
}

func (e *EventPersistence) EnableDisableEvent(ctx context.Context, id string, enable bool) (*event.Event, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}

	update := bson.M{
		"enabled":          enable,
		"last_modified_at": time.Now(),
	}

	eventDoc, err := e.eventDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			e.logger.Errorf("Event with ID %s not found for enable/disable", id)
			return nil, fmt.Errorf(common_util.NotFound)
		}
		e.logger.Errorf("Failed to update enabled state for Event ID %s: %v", id, err)
		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	result := eventDoc.toModel()
	return &result, nil
}

func (e *EventPersistence) EventNameExists(ctx context.Context, name string, id *string) (bool, error) {
	if name == "" {
		return false, fmt.Errorf(common_util.InvalidInput)
	}

	filter := bson.M{
		"event_name": bson.M{"$regex": fmt.Sprintf("^%s$", name), "$options": "i"},
		"is_deleted": false,
	}

	if id != nil {
		objID, err := common_util.ParsePrimitiveObjectID(*id)
		if err != nil {
			return false, fmt.Errorf(common_util.InvalidID)
		}
		filter["_id"] = bson.M{"$ne": objID}
	}

	count, err := e.eventDal.TotalCount(ctx, filter)
	if err != nil {
		e.logger.Warnf("Failed to check event name existence for name %s: %v", name, err)
		return false, fmt.Errorf(common_util.GeneralDBQueryFailed)
	}

	return count > 0, nil
}
