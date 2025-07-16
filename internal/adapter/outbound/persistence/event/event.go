package event

import (
	"context"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type EventPersistence struct {
	eventDal dal.MongoDal[event.Event, event.Event]
	cpsDal   dal.MongoDal[action.CPSAction, action.CPSAction]
	logger   utils.Logger
}

func InitEventPersistence(client *mongo.Client, dbName string, logger utils.Logger) *EventPersistence {
	return &EventPersistence{
		eventDal: dal.NewMongoDal[event.Event, event.Event](client, dbName, "events"),
		cpsDal:   dal.NewMongoDal[action.CPSAction, action.CPSAction](client, dbName, "cps_actions"),
		logger:   logger,
	}
}

var _ event.Repository = (*EventPersistence)(nil)

func (e *EventPersistence) CreateCpsAction(ctx context.Context, Action action.CPSAction) (action.CPSAction, error) {
	Action.CreatedAt = time.Now()
	Action.LastModifiedAt = time.Now()
	res, err := e.cpsDal.InsertOne(ctx, Action)
	if err != nil {
		return action.CPSAction{}, fmt.Errorf("FAILED_TO_INSERT_CPS_ACTION")
	}
	return res, nil
}

func (e *EventPersistence) UpdateCpsAction(ctx context.Context, Action action.CPSAction) error {
	filter := bson.M{"action_code": Action.ActionCode}
	update := bson.M{"$set": Action}
	_, err := e.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("FAILED_TO_UPDATE_CPS_ACTION")
	}
	return nil
}

func (e *EventPersistence) FetchCpsActionById(ctx context.Context, Action_Id string) (action.CPSAction, error) {
	filter := bson.M{"action_code": Action_Id}
	res, err := e.cpsDal.FindOne(ctx, filter, nil)
	if err != nil {
		return action.CPSAction{}, fmt.Errorf("FAILED_TO_FETCH_CPS_ACTION")
	}
	if res == nil {
		return action.CPSAction{}, fmt.Errorf("CPS_ACTION_NOT_FOUND")
	}
	return *res, nil
}

func (e *EventPersistence) CreateEvent(ctx context.Context, ev event.Event) (event.Event, error) {
	ev.CreatedAt = time.Now()
	ev.LastModifiedAt = time.Now()
	res, err := e.eventDal.InsertOne(ctx, ev)
	if err != nil {
		return event.Event{}, fmt.Errorf("FAILED_TO_INSERT_EVENT")
	}
	return res, nil
}

func (e *EventPersistence) FetchEventById(ctx context.Context, event_id string) (event.Event, error) {
	filter := bson.M{"_id": event_id}
	res, err := e.eventDal.FindOne(ctx, filter, nil)
	if err != nil {
		return event.Event{}, fmt.Errorf("FAILED_TO_FETCH_EVENT")
	}
	if res == nil {
		return event.Event{}, fmt.Errorf("EVENT_NOT_FOUND")
	}
	return *res, nil
}

func (e *EventPersistence) FetchEvent(ctx context.Context, limit, offset int) ([]event.Event, error) {
	skip := int64(offset)
	lim := int64(limit)
	res, err := e.eventDal.FindAllWithPagination(ctx, nil, nil, skip, lim)
	if err != nil {
		return nil, fmt.Errorf("FAILED_TO_FETCH_EVENT")
	}
	var events []event.Event
	for _, ptr := range res {
		if ptr != nil {
			events = append(events, *ptr)
		}
	}
	return events, nil
}
