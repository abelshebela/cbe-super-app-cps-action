package faydaaccount

import (
	"context"

	portalCardDomain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/portal_card"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
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

func (o *PortalCardRepo) GetAllPortalCard(ctx context.Context) ([]*portalCardDomain.Card, error) {
	data, err := o.portalCardDal.FindAll(ctx, nil, nil)
	if err != nil {
		return nil, err
	}

	result := make([]*portalCardDomain.Card, len(data))

	for i, s := range data {
		if s == nil {
			continue
		}
		result[i] = &portalCardDomain.Card{
			ID:       s.ID,
			CardName: s.CardName,
			SubCards: s.SubCards,
		}
	}
	return result, nil
}
