package newscategory_repo

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	"context"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
)

type NewsCategoryRepository struct {
	client        *mongo.Client
	mongoDal      dal.MongoDal[model.NewsCategory, model.NewsCategory]
	kafkaProducer kafka.ClientOrchestrationProducer
	logger        utils.Logger
}

func NewNewsCategoryRepository(client *mongo.Client,cfg *config.VaultConfig, database string, collection string, kafkaProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.NewsCategoryRepository {
	mongoDal := dal.NewMongoDal[model.NewsCategory, model.NewsCategory](client,cfg, database, collection)
	return &NewsCategoryRepository{
		client:        client,
		mongoDal:      mongoDal,
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

// Create implements storage.NewsCategoryRepository.
func (n *NewsCategoryRepository) Create(ctx context.Context, categoryName []string) error {
	n.logger.Infof("Create called with category names: %v", categoryName)
	session, err := n.client.StartSession()
	if err != nil {
		n.logger.Errorf("Create: failed to start session: %v", err)
		return err
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sc context.Context) error {
		if err := session.StartTransaction(); err != nil {
			n.logger.Errorf("Create: failed to start transaction: %v", err)
			return err
		}

		for _, name := range categoryName {
			cat := model.NewsCategory{
				ID:             bson.NewObjectID(),
				CategoryName:   name,
				CreatedAt:      time.Now(),
				LastModifiedAt: time.Now(),
			}
			if _, err := n.mongoDal.InsertOne(sc, cat); err != nil {
				n.logger.Errorf("Create: failed to insert category '%s': %v", name, err)
				_ = session.AbortTransaction(sc)
				return err
			}
		}

		if err := session.CommitTransaction(sc); err != nil {
			n.logger.Errorf("Create: failed to commit transaction: %v", err)
			return err
		}
		n.logger.Infof("Create: transaction committed successfully for categories: %v", categoryName)
		return nil
	})
	if err != nil {
		n.logger.Errorf("Create: transaction failed: %v", err)
	} else {
		n.kafkaProducer.PublishMessage(ctx, categoryName, string(constants.ClientOrchestrationNewsCategoryTopic), string(constants.ClientOrchestrationNewsCategoryTopic), "new news categories created")
		n.logger.Infof("Create completed for category names: %v", categoryName)
	}
	return err
}

// Delete implements storage.NewsCategoryRepository.
func (n *NewsCategoryRepository) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return localization.ErrorNewsCategoryInvalidID
	}

	err = n.mongoDal.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	return nil
}

// FindAllWithPagination implements storage.NewsCategoryRepository.
func (n *NewsCategoryRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.NewsCategory], error) {
	searchKeys := bson.M{}

	allowedKeys := []string{"category_name"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"category_name": searchRegex},
		}
	}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	if skip < 0 {
		skip = 0
	}
	if limit < 0 {
		limit = 10
	}

	filter["is_deleted"] = false
	catagories, err := n.mongoDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, err
	}

	totalDocs, err := n.mongoDal.TotalCount(ctx, filter)
	if err != nil {
		totalDocs = int64(skip) + int64(len(catagories))
	}

	// Ensure sensible defaults
	if limit <= 0 {
		if len(catagories) > 0 {
			limit = int64(len(catagories))
		} else {
			limit = 1
		}
	}
	page := filterParam.Page
	if page <= 0 {
		page = 1
	}

	// calculate pagination meta
	totalPages := int((totalDocs + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	pagingCounter := (int64(page)-1)*limit + 1
	hasPrev := page > 1
	hasNext := page < totalPages

	var prevPage *int
	if hasPrev {
		p := page - 1
		prevPage = &p
	}
	var nextPage *int
	if hasNext {
		nxt := page + 1
		nextPage = &nxt
	}

	return &types.PaginatedResponse[[]model.NewsCategory]{
		Data: catagories,
		Meta: types.PaginationMeta{
			TotalDocs:     totalDocs,
			Limit:         int(limit),
			TotalPages:    totalPages,
			Page:          page,
			PagingCounter: int(pagingCounter),
			HasPrevPage:   hasPrev,
			HasNextPage:   hasNext,
			PrevPage:      prevPage,
			NextPage:      nextPage,
		},
	}, nil

}

// FindByNames implements storage.NewsCategoryRepository.
func (n *NewsCategoryRepository) FindByNames(ctx context.Context, names []string) (*model.NewsCategory, error) {
	var orQueries []bson.M
	for _, name := range names {
		orQueries = append(orQueries, bson.M{
			"category_name": bson.M{
				"$regex":   "^" + name + "$",
				"$options": "i",
			},
		})
	}
	filter := bson.M{"$or": orQueries}
	news_category, err := n.mongoDal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		return nil, err
	}

	return news_category, nil
}

// Get implements storage.NewsCategoryRepository.
func (n *NewsCategoryRepository) Get(ctx context.Context, id string) (*model.NewsCategory, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, localization.ErrorNewsCategoryInvalidID
	}

	category, err := n.mongoDal.FindOne(ctx, bson.M{"_id": objID, "is_deleted": false}, bson.M{})
	if err != nil {
		return nil, err
	}
	return category, nil
}

// Update implements storage.NewsCategoryRepository.
func (n *NewsCategoryRepository) Update(ctx context.Context, id string, categoryName string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return localization.ErrorNewsCategoryInvalidID
	}
	// if err:= n.mongoDal.FindOne(ctx,bson.M{}, bson.M{}) err==nil{

	// }
	updatedNewsCategory, err := n.mongoDal.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"category_name": categoryName})
	if err != nil {
		return err
	}
	n.kafkaProducer.PublishMessage(ctx, updatedNewsCategory, string(constants.ClientOrchestrationNewsCategoryTopic), string(constants.ClientOrchestrationNewsCategoryTopic), "news category updated")
	return nil
}
