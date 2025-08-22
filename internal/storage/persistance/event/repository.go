package event

import (
	"context"
	"errors"
	"time"

	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type EventStorage struct {
	dal    dal.MongoDal[model.EventDocument, model.EventDocument]
	logger utils.Logger
}

func NewEventRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.EventRepository {
	return &EventStorage{
		dal:    dal.NewMongoDal[model.EventDocument, model.EventDocument](client, dbName, collection),
		logger: logger,
	}
}

func (e *EventStorage) Create(ctx context.Context, event *model.Event) error {
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}
	event.LastModifiedAt = time.Now()

	eventDoc := EventDocumentMapper(*event)

	_, err := e.dal.InsertOne(ctx, *eventDoc)
	if err != nil {
		e.logger.Errorf("Create Event failed", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (e *EventStorage) Update(ctx context.Context, id string, event *model.Event) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	event.LastModifiedAt = time.Now()
	updateDoc := EventDocumentToBsonM(model.EventDocument(*event))

	filter := bson.M{"_id": objID, "is_deleted": false}

	_, err = e.dal.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		e.logger.Errorf("Update Event failed", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (e *EventStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateFields := bson.M{
		"is_deleted": true,
		"deleted_at": time.Now(),
	}

	_, err = e.dal.UpdateOne(ctx, filter, updateFields)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		e.logger.Errorf("Delete Event failed", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (e *EventStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{
		"enabled":          enable,
		"last_modified_at": time.Now(),
	}
	_, err = e.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		e.logger.Errorf("EnableOrDisable Event failed", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}
func (e *EventStorage) FindByID(ctx context.Context, id string) (*model.Event, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}

	doc, err := e.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		e.logger.Errorf("FindByID Event failed", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	result := EventMapper(doc)
	return &result, nil
}

func (e *EventStorage) FindByName(ctx context.Context, name string) (*model.Event, error) {
	if name == "" {
		e.logger.Warnf("FindByName called with empty name")
		return nil, errors.New(localization.ErrorEventNameRequired.Code)
	}

	filter := bson.M{
		"event_name": name,
		"is_deleted": false,
	}

	doc, err := e.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			e.logger.Warnf("No event found with name: %s", name)
			return nil, nil
		}
		e.logger.Errorf("FindByName Event failed for name %s: %v", name, err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	result := EventMapper(doc)
	return &result, nil
}

func (e *EventStorage) FindAll(ctx context.Context, filter bson.M, projection bson.M) ([]*model.Event, error) {
	filter["is_deleted"] = false

	docs, err := e.dal.FindAll(ctx, filter, projection)
	if err != nil {
		e.logger.Errorf("FindAll Event failed", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	var events []*model.Event
	for _, doc := range docs {
		event := EventMapper(doc)
		events = append(events, &event)
	}
	return events, nil
}

func (e *EventStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Event], error) {
	e.logger.Infof("FindAllWithPagination called with filter: %+v", filterParam)

	filter := bson.M{"is_deleted": false}
	searchKeys := bson.M{}

	allowedKeys := []string{"event_code", "event_name", "event_city", "account_number", "event_venue"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["event_name"] = searchRegex
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	e.logger.Debugf("Mongo filter: %+v, skip: %d, limit: %d", filter, skip, limit)

	docs, err := e.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		e.logger.Errorf("FindAllWithPagination Event failed", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	var events []*model.Event
	for _, doc := range docs {
		event := EventMapper(doc)
		events = append(events, &event)
	}

	total, err := e.dal.TotalCount(ctx, filter)
	if err != nil {
		e.logger.Errorf("Count Event failed", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	e.logger.Infof("FindAllWithPagination returning %d events, total: %d", len(events), total)

	return &types.PaginatedResponse[[]*model.Event]{
		Data: events,
		Meta: meta,
	}, nil
}
