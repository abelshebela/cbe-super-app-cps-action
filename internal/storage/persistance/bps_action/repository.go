package bps_action

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	bps_action "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type bpsActionRepository struct {
	client    *mongo.Client
	actionDal dal.MongoDal[bps_action.BPSAction, bps_action.BPSAction]
	logger    utils.Logger
}

// GetBPSActionByUserID implements [storage.BPSActionRepository].
func (b bpsActionRepository) GetBPSActionByUserID(ctx context.Context, userID string, filterParam types.Filter) (types.PaginatedResponse[[]bps_action.BPSAction], error) {
	b.logger.Infof("[GetBPSActionByUserID] fetching BPS actions for user ID: %s", userID)
	objID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		b.logger.Errorf("[GetBPSActionByUserID] invalid user ID: %v", err)
		return types.PaginatedResponse[[]bps_action.BPSAction]{}, err
	}
	filter, skip, limit := lib.FilterBuilder(filterParam, bson.M{}, nil)
	filter["user_information.user_id"] = objID

	cus, err := b.actionDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		b.logger.Errorf("[GetBPSActionByUserID] failed to fetch BPS actions: %v", err)
		return types.PaginatedResponse[[]bps_action.BPSAction]{}, errors.New("failed to fetch BPS actions")
	}

	filter = bson.M{"user_information.user_id": objID}
	total, err := b.actionDal.TotalCount(ctx, filter)
	if err != nil {
		b.logger.Errorf("[GetBPSActionByUserID] failed to fetch total count: %v", err)
		return types.PaginatedResponse[[]bps_action.BPSAction]{}, errors.New("error failed to get bps action")
	}
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return types.PaginatedResponse[[]bps_action.BPSAction]{
		Data: cus,
		Meta: meta,
	}, nil
}

func NewBPSActionRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger, cfg *config.VaultConfig) storage.BPSActionRepository {
	return bpsActionRepository{
		client:    client,
		actionDal: dal.NewMongoDal[bps_action.BPSAction, bps_action.BPSAction](client, cfg, dbName, collection),
		logger:    logger,
	}
}
