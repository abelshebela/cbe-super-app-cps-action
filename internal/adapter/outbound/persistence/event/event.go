package event

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	common_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type EventPersistence struct {
	eventDal dal.MongoDal[EventDocument, EventDocument]
	cpsDal   dal.MongoDal[CPSActionDocument, CPSActionDocument]
	logger   utils.Logger
}

func InitEventPersistence(client *mongo.Client, dbName string, logger utils.Logger) *EventPersistence {
	return &EventPersistence{
		eventDal: dal.NewMongoDal[EventDocument, EventDocument](client, dbName, "events"),
		cpsDal:   dal.NewMongoDal[CPSActionDocument, CPSActionDocument](client, dbName, "cps_actions"),
		logger:   logger,
	}
}

var _ event.EventRepository = (*EventPersistence)(nil)

func (e *EventPersistence) CreateCpsAction(ctx context.Context, action action.CPSAction) (*action.CPSAction, error) {
	action.CreatedAt = time.Now()
	action.LastModifiedAt = time.Now()
	actionDoc, err := ToCpsActionDocument(action)
	if err != nil {
		
		return nil, err
	}
	res, err := e.cpsDal.InsertOne(ctx, *actionDoc)
	if err != nil {
		return nil, fmt.Errorf(common_utils.GeneralDBInsertFailed)
	}

	actionDomin := res.toModel()
	return &actionDomin, nil
}

func (e *EventPersistence) UpdateCpsAction(ctx context.Context, action action.CPSAction) error {
	filter := bson.M{"action_code": action.ActionCode}
	actionDoc, err := ToCpsActionDocument(action)
	if err != nil {
		return err
	}
	update := bson.M{"$set": actionDoc}
	_, err = e.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf(common_utils.GeneralDBUpdateFailed)
	}
	return nil
}

func (e *EventPersistence) FetchCpsActionByID(ctx context.Context, action_ID string) (*action.CPSAction, error) {
	filter := bson.M{"action_code": action_ID}
	res, err := e.cpsDal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf(common_utils.ActionNotFound)
		}
		return nil, fmt.Errorf(common_utils.GeneralDBQueryFailed)
	}

	resDomain := res.toModel()
	return &resDomain, nil
}

func (e *EventPersistence) CreateEvent(ctx context.Context, event event.Event) (*event.Event, error) {
	event.CreatedAt = time.Now()
	event.LastModifiedAt = time.Now()
	eventDoc, err := ToEventDocument(event)

	if err != nil {
		return nil, err
	}

	res, err := e.eventDal.InsertOne(ctx, *eventDoc)
	if err != nil {
		return nil, fmt.Errorf(common_utils.GeneralDBInsertFailed)
	}

	eventDomain := res.toModel()
	return &eventDomain, nil
}

func (e *EventPersistence) FetchEventById(ctx context.Context, event_id string) (*event.Event, error) {
	objID, err := common_utils.ParsePrimitiveObjectID(event_id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}

	res, err := e.eventDal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf(common_utils.NotFound)
		}
		return nil, fmt.Errorf(common_utils.GeneralDBQueryFailed)
	}

	eventDomain := res.toModel()

	return &eventDomain, nil
}

func (e *EventPersistence) FetchEvent(ctx context.Context, limit, offset int) ([]*event.Event, error) {
	skip := int64(offset)
	lim := int64(limit)
	res, err := e.eventDal.FindAllWithPagination(ctx, nil, nil, skip, lim)
	if err != nil {
		return nil, fmt.Errorf(common_utils.GeneralDBQueryFailed)
	}
	var events []*event.Event
	for _, ptr := range res {
		if ptr != nil {
			eventDomain := ptr.toModel()
			events = append(events, &eventDomain)
		}
	}
	return events, nil
}
