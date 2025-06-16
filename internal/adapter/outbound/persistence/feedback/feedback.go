package customer

import (
	"context"
	"errors"
	"fmt"
	"net/http"


	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/feedback/entity"
	outbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/outbound/feedback"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type FeedbackRepo struct {
	client   *mongo.Client
	mongoDal dal.MongoDal[entity.Feedback, entity.Feedback]
	logger   utils.Logger
}

var _ outbound.FeedbackRepository = (*FeedbackRepo)(nil)

func InitFeedback(client *mongo.Client, database string, collection string, logger utils.Logger) *FeedbackRepo {
	mongoDal := dal.NewMongoDal[entity.Feedback, entity.Feedback](client, database, collection)

	return &FeedbackRepo{
		client:   client,
		mongoDal: mongoDal,
		logger:   logger,
	}
}

func (c *FeedbackRepo) GetFeedbacks(ctx context.Context, filterParams *constant.Filter) (*entity.FeedbackResponse, error) {
	filter := bson.M{}
	projection := bson.M{}

	if filterParams.Search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"name": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			},
		}
	}

	if filterParams.Filters != "" {
		filter["rating"] = filterParams.Filters
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage

	feedbacks, err := c.mongoDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(filterParams.PerPage))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.logger.Errorf("no feedback data found", err)
			err = fmt.Errorf("feedback not found %w", constant.ErrorDefinition{
				Code:    http.StatusNotFound,
				Message: "no feedback data found",
			})
			return nil, err
		}
		c.logger.Errorf("failed to get feedback data", err)
		err = fmt.Errorf("failed to get feedback %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	total, err := c.mongoDal.TotalCount(ctx, bson.M{})
	if err != nil {
		c.logger.Errorf("failed to get total counts", err)
		err := fmt.Errorf("failed to get total counts %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})

		return nil, err
	}

	return &entity.FeedbackResponse{
		Page:      1,
		Feedbacks: feedbacks,
		Limit:     constant.DefaultPerPage,
		Total:     total,
	}, nil
}


func (c *FeedbackRepo) GetFeedbackByID(ctx context.Context, id string) (*entity.Feedback, error) {
	feedbackID, err := bson.ObjectIDFromHex(id)
	if feedbackID.IsZero() {
		c.logger.Errorf("invalid id provided", err)
		err = fmt.Errorf("invalid id provided %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
		})
		return nil, err
	}

	if err != nil {
		c.logger.Errorf("failed to convert id to object id", err)
		err = fmt.Errorf("failed to convert id to object id %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	filter := bson.M{"_id": feedbackID}
	projection := bson.M{}
	feedback, err := c.mongoDal.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.logger.Errorf("feedback not found", err)
			err = fmt.Errorf("feedback not found %w", constant.ErrorDefinition{
				Code:    http.StatusNotFound,
				Message: "feedback not found",
			})
			return nil, err
		}
		c.logger.Errorf("failed to get feedback", err)
		err = fmt.Errorf("failed to get feedback %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}
	return feedback, nil
}
