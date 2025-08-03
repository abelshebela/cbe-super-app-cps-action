package external_call

import (
	"context"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/errors"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/model"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ExternalCallRepository struct {
	externalCallDal dal.MongoDal[model.ExternalCall, model.ExternalCall]
	logger          utils.Logger
}

func NewExternalCallRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.ExternalCallRepository {
	return &ExternalCallRepository{
		externalCallDal: dal.NewMongoDal[model.ExternalCall, model.ExternalCall](client, dbName, collection),
		logger:          logger,
	}
}

func (e *ExternalCallRepository) Save(ctx context.Context, externalCall *model.ExternalCall) error {
	if externalCall == nil {
		e.logger.Errorf("Save external call failed: external call is nil")
		return errors.ErrUnexpected
	}

	// Set timestamps
	now := time.Now()
	if externalCall.CreatedAt.IsZero() {
		externalCall.CreatedAt = now
	}
	externalCall.UpdatedAt = now

	_, err := e.externalCallDal.InsertOne(ctx, *externalCall)
	if err != nil {
		e.logger.Errorf("Unexpected error while saving external call. error=%v, externalCall=%+v", err, externalCall)
		return errors.ErrUnexpected
	}
	e.logger.Infof("External call saved successfully. externalCall=%+v", externalCall)
	return nil
}

func (e *ExternalCallRepository) FindById(ctx context.Context, id string) (*model.ExternalCall, error) {
	if id == "" {
		e.logger.Errorf("Find external call failed: id is empty")
		return nil, errors.ErrIdEmpty
	}

	filter, err := FilterIdForExternalCall(id)
	if err != nil {
		e.logger.Errorf("Find external call failed: error creating filter. error=%v, id=%s", err, id)
		return nil, err
	}

	projection := ExternalCallProjection()
	externalCall, err := e.externalCallDal.FindOne(ctx, filter, projection)
	if err != nil {
		if err == errors.ErrUnexpected {
			e.logger.Errorf("External call not found. id=%s", id)
			return nil, errors.ErrUnexpected
		}
		e.logger.Errorf("Unexpected error while finding external call. error=%v, id=%s", err, id)
		return nil, errors.ErrUnexpected
	}
	e.logger.Infof("External call found successfully. externalCall=%+v", externalCall)
	return externalCall, nil
}

func (e *ExternalCallRepository) FindByURL(ctx context.Context, url string) ([]*model.ExternalCall, error) {
	if url == "" {
		e.logger.Errorf("Find external call failed: url is empty")
		return nil, errors.ErrUnexpected
	}

	filter := bson.M{"url": url}
	projection := ExternalCallProjection()
	externalCalls, err := e.externalCallDal.Find(ctx, filter, projection)
	if err != nil {
		e.logger.Errorf("Unexpected error while finding external calls by URL. error=%v, url=%s", err, url)
		return nil, errors.ErrUnexpected
	}
	e.logger.Infof("External calls found successfully. count=%d, url=%s", len(externalCalls), url)
	return externalCalls, nil
}

func (e *ExternalCallRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*model.ExternalCall, error) {
	filter := bson.M{
		"created_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}
	projection := ExternalCallProjection()
	externalCalls, err := e.externalCallDal.Find(ctx, filter, projection)
	if err != nil {
		e.logger.Errorf("Unexpected error while finding external calls by date range. error=%v, startDate=%v, endDate=%v", err, startDate, endDate)
		return nil, errors.ErrUnexpected
	}
	e.logger.Infof("External calls found successfully. count=%d, startDate=%v, endDate=%v", len(externalCalls), startDate, endDate)
	return externalCalls, nil
}

func (e *ExternalCallRepository) Update(ctx context.Context, id string, update *model.ExternalCall) error {
	if id == "" {
		e.logger.Errorf("Update external call failed: id is empty")
		return errors.ErrIdEmpty
	}

	if update == nil {
		e.logger.Errorf("Update external call failed: update is nil")
		return errors.ErrUnexpected
	}

	filter, err := FilterIdForExternalCall(id)
	if err != nil {
		e.logger.Errorf("Update external call failed: error creating filter. error=%v, id=%s", err, id)
		return err
	}

	update.UpdatedAt = time.Now()
	_, err = e.externalCallDal.UpdateOne(ctx, filter, update)
	if err != nil {
		e.logger.Errorf("Update external call failed: error updating external call. error=%v, id=%s", err, id)
		return err
	}
	e.logger.Infof("External call updated successfully. id=%s", id)
	return nil
}

func (e *ExternalCallRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		e.logger.Errorf("Delete external call failed: id is empty")
		return errors.ErrIdEmpty
	}

	filter, err := FilterIdForExternalCall(id)
	if err != nil {
		e.logger.Errorf("Delete external call failed: error creating filter. error=%v, id=%s", err, id)
		return err
	}

	if err := e.externalCallDal.DeleteOne(ctx, filter); err != nil {
		e.logger.Errorf("Delete external call failed: error deleting external call. error=%v, id=%s", err, id)
		return err
	}
	e.logger.Infof("External call deleted successfully. id=%s", id)
	return nil
}
