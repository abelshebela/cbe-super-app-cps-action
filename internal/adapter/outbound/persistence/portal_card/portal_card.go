package faydaaccount

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	portalCardDomain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/portal_card"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PortalCardRepo struct {
	client        *mongo.Client
	logger        utils.Logger
	portalCardDal dal.MongoDal[model.Card, model.Card]
}

func InitPortalCardPersistence(client *mongo.Client, database string, collection string, logger utils.Logger) *PortalCardRepo {
	portalCardDal := dal.NewMongoDal[model.Card, model.Card](client, database, collection)
	return &PortalCardRepo{
		client:        client,
		portalCardDal: portalCardDal,
		logger:        logger,
	}
}

func (o *PortalCardRepo) GetAllPortalCard(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*portalCardDomain.Card], error) {
	filter := map[string]interface{}{}
	projection := map[string]interface{}{}

	// Add search functionality
	if filterParams != nil && filterParams.Search != "" {
		filter["$or"] = []map[string]interface{}{
			{"card_name": map[string]interface{}{"$regex": filterParams.Search, "$options": "i"}},
		}
	}

	// Add filter functionality (if needed)
	if filterParams != nil && filterParams.Filters != "" {
		filter["card_name"] = filterParams.Filters
	}

	page := 1
	limit := 10
	skip := 0
	if filterParams != nil {
		page = filterParams.Page
		limit = filterParams.PerPage
		skip = (page - 1) * limit
	}

	data, err := o.portalCardDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		return nil, err
	}
	totalDocs, err := o.portalCardDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	var dataList []*portalCardDomain.Card
	for _, s := range data {
		if s == nil {
			continue
		}
		dataList = append(dataList, &portalCardDomain.Card{
			ID:       s.ID,
			CardName: s.CardName,
			SubCards: s.SubCards,
		})
	}

	meta := common_util.BuildPaginationMeta(totalDocs, page, limit)

	return &common_util.PaginatedResponse[[]*portalCardDomain.Card]{
		Data: dataList,
		Meta: meta,
	}, nil
}
