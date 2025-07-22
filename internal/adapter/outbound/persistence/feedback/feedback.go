package customer

import (
	"context"
	"errors"
	"fmt"

	// "net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/feedback/entity"
	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/feedback"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
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

func (c *FeedbackRepo) GetFeedbacks(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.Feedback], error) {
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
	limit := filterParams.PerPage

	feedbacks, err := c.mongoDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.logger.Errorf("no feedback data found", err)
			return &common_util.PaginatedResponse[[]*entity.Feedback]{
				Data: []*entity.Feedback{},
				Meta: common_util.BuildPaginationMeta(0, filterParams.Page, filterParams.PerPage),
			}, nil
		}
		c.logger.Errorf("failed to get feedback data", err)
		return nil, fmt.Errorf("FAILED_TO_GET_FEEDBACK")
	}

	total, err := c.mongoDal.TotalCount(ctx, filter)
	if err != nil {
		c.logger.Errorf("failed to get total counts", err)
		return nil, fmt.Errorf("FAILED_TO_GET_FEEDBACK_COUNTS")
	}

	meta := common_util.BuildPaginationMeta(total, filterParams.Page, filterParams.PerPage)

	if feedbacks == nil {
		feedbacks = []*entity.Feedback{}
	}

	return &common_util.PaginatedResponse[[]*entity.Feedback]{
		Data: feedbacks,
		Meta: meta,
	}, nil
}

func (c *FeedbackRepo) GetFeedbackByID(ctx context.Context, id string) (*entity.Feedback, error) {
	feedbackID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.logger.Errorf("invalid id provided: %v", err)
		return nil, fmt.Errorf("FAILED_TO_CONVERT_ID")
	}

	if feedbackID.IsZero() {
		c.logger.Errorf("invalid id provided", err)

		return nil, fmt.Errorf("INVALID_ID")
	}

	filter := bson.M{"_id": feedbackID}
	projection := bson.M{}
	// fmt.Println(filter, projection,"dataaaaaaa")
	feedback, err := c.mongoDal.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.logger.Errorf("feedback not found", err)

			return nil, fmt.Errorf("FAILED_TO_GET_FEEDBACK_DATA")
		}
		c.logger.Errorf("failed to get feedback", err)
		return nil, fmt.Errorf("FAILED_TO_GET_FEEDBACK")
	}
	return feedback, nil
}
