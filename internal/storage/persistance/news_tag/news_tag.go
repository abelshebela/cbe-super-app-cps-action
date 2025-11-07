package newstag_repo

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type NewsTagRepository struct {
	client   *mongo.Client
	mongoDal dal.MongoDal[model.NewsTag, model.NewsTag]
	logger   utils.Logger
}

func NewNewsTagRepository(client *mongo.Client, database string, collection string, logger utils.Logger) storage.NewsTagRepository {
	mongoDal := dal.NewMongoDal[model.NewsTag, model.NewsTag](client, database, collection)
	return &NewsTagRepository{
		client:   client,
		mongoDal: mongoDal,
		logger:   logger,
	}
}

// Create implements storage.NewsTagRepository.
func (n *NewsTagRepository) Create(ctx context.Context, tagName []string) error {
	session, err := n.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sc context.Context) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		for _, name := range tagName {
			cat := model.NewsTag{
				ID:             bson.NewObjectID(),
				TagName:        name,
				CreatedAt:      time.Now(),
				LastModifiedAt: time.Now(),
			}
			if _, err := n.mongoDal.InsertOne(sc, cat); err != nil {
				_ = session.AbortTransaction(sc)
				return err
			}
		}

		if err := session.CommitTransaction(sc); err != nil {
			return err
		}
		return nil
	})
	return err
}

// Delete implements storage.NewsTagRepository.
func (n *NewsTagRepository) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return localization.ErrorNewsTagInvalidID
	}

	err = n.mongoDal.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	return nil
}

// FindAllWithPagination implements storage.NewsTagRepository.
func (n *NewsTagRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.NewsTag], error) {

	searchKeys := bson.M{}

	allowedKeys := []string{"tag_name"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"tag_name": searchRegex},
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
	categories, err := n.mongoDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, err
	}

	totalDocs, err := n.mongoDal.TotalCount(ctx, filter)
	if err != nil {
		totalDocs = int64(skip) + int64(len(categories))
	}

	// Ensure sensible defaults
	if limit <= 0 {
		if len(categories) > 0 {
			limit = int64(len(categories))
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

	return &types.PaginatedResponse[[]*model.NewsTag]{
		Data: categories,
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

// FindByNames implements storage.NewsTagRepository.
func (n *NewsTagRepository) FindByNames(ctx context.Context, names []string) (*model.NewsTag, error) {
	var orQueries []bson.M
	for _, name := range names {
		orQueries = append(orQueries, bson.M{
			"tag_name": bson.M{
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

// Get implements storage.NewsTagRepository.
func (n *NewsTagRepository) Get(ctx context.Context, id string) (*model.NewsTag, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, localization.ErrorNewsTagInvalidID
	}

	tag, err := n.mongoDal.FindOne(ctx, bson.M{"_id": objID, "is_deleted": false}, bson.M{})
	if err != nil {
		return nil, err
	}
	return tag, nil
}

// Update implements storage.NewsTagRepository.
func (n *NewsTagRepository) Update(ctx context.Context, id string, tagName string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return localization.ErrorNewsCategoryInvalidID
	}

	_, err = n.mongoDal.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"tag_name": tagName})
	if err != nil {
		return err
	}
	return nil
}
