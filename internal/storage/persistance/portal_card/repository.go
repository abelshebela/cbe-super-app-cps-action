package portal_card

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PortalCardStorage struct {
	dal    dal.MongoDal[model.Card, model.Card]
	client *mongo.Client
	logger utils.Logger
}

func NewPortalCardRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.PortalCardRepository {
	return &PortalCardStorage{
		dal:    dal.NewMongoDal[model.Card, model.Card](client, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (p *PortalCardStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Card], error) {
	filter := bson.M{
		"is_deleted": false,
	}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["card_name"] = searchRegex
	}

	skip := int64((filterParam.Page - 1) * filterParam.PerPage)
	limit := int64(filterParam.PerPage)

	data, err := p.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, err
	}

	total, err := p.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.Card]{
		Data: data,
		Meta: meta,
	}, nil
}
