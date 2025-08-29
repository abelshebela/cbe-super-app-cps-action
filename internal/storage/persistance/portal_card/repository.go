package portal_card

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"fmt"
	"strings"

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

func (s *PortalCardStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Card], error) {
	// 1. Base filter (only active records)
	filter := bson.M{"is_deleted": false}
	searchKeys := bson.M{}

	// 2. Allowed filterable/searchable fields
	allowedKeys := []string{}

	// 3. Add search (if provided)
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["card_name"] = searchRegex // choose your searchable field(s)
	}

	// 4. Build filter, skip, limit
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	// 5. Fetch data
	data, err := s.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 6. Count total
	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Return standard paginated response
	return &types.PaginatedResponse[[]*model.Card]{
		Data: data,
		Meta: meta,
	}, nil
}

func (o *PortalCardStorage) ValidatePortalCard(ctx context.Context, names []string) (bool, error) {
	if len(names) == 0 {
		return false, fmt.Errorf("PORTAL_CARD_ARRAY_EMPTY")
	}

	var cleaned []string
	for _, name := range names {
		if trimmed := strings.TrimSpace(name); trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}

	if len(cleaned) == 0 {
		return false, fmt.Errorf("NO_VALID_PORTAL_CARD_NAME")
	}

	// Query DB
	cards, err := o.dal.FindAll(ctx,
		bson.M{"card_name": bson.M{"$in": cleaned}},
		bson.M{"card_name": 1},
	)
	if err != nil {
		return false, fmt.Errorf("DB_ERROR: %w", err)
	}

	// Build lookup
	found := make(map[string]struct{})
	for _, card := range cards {
		if card.CardName != "" {
			found[card.CardName] = struct{}{}
		}
	}

	// Detect missing names
	var missing []string
	for _, name := range cleaned {
		if _, ok := found[name]; !ok {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		return false, fmt.Errorf("PORTAL_CARD_NOT_FOUND")
	}

	return true, nil
}
func (o *PortalCardStorage) ValidatePortalCardByID(ctx context.Context, ids []string) (bool, error) {
	if len(ids) == 0 {
		return false, fmt.Errorf("PORTAL_CARD_ARRAY_EMPTY")
	}

	var cleaned []bson.ObjectID
	for _, id := range ids {
		if trimmed := strings.TrimSpace(id); trimmed != "" {
			objID, err := bson.ObjectIDFromHex(trimmed)
			if err != nil {
				return false, errors.New(localization.ErrorInvalidID.Code)
			}
			cleaned = append(cleaned, objID)
		}
	}

	if len(cleaned) == 0 {
		return false, fmt.Errorf("NO_VALID_PORTAL_CARD_ID")
	}
	// Query DB for all given ids
	filter := bson.M{"_id": bson.M{"$in": cleaned}}
	projection := bson.M{"_id": 1}

	cards, err := o.dal.FindAll(ctx, filter, projection)
	if err != nil {
		return false, fmt.Errorf("DB_ERROR: %w", err)
	}

	for _, card := range cards {
		fmt.Println(card)
	}

	if len(ids) != len(cards) {
		return false, fmt.Errorf("PORTAL_CARD_NOT_FOUND")
	}

	return true, nil
}
