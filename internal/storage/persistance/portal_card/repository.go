package portal_card

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/localization"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

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

func (p *PortalCardStorage) Create(ctx context.Context, card *model.Card) error {
	_, err := p.dal.InsertOne(ctx, *card)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (p *PortalCardStorage) Update(ctx context.Context, id string, card *model.Card) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{
		"$set": bson.M{
			"card_name": card.CardName,
			"sub_cards": card.SubCards,
		},
	}

	_, err = p.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (p *PortalCardStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	return p.dal.DeleteOne(ctx, filter)
}

func (p *PortalCardStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"enabled": enable}}
	_, err = p.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
	}
	return nil
}

func (p *PortalCardStorage) FindByID(ctx context.Context, id string) (*model.Card, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}

	result, err := p.dal.FindOne(ctx, filter, nil)

	if err != nil {
		return nil, err
	}
	return result, nil
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