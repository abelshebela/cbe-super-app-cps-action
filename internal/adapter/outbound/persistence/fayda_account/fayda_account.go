package faydaaccount

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/fayda_account"

	model "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type FaydaAccountRepo struct {
	client      *mongo.Client
	logger      utils.Logger
	cpsDal      dal.MongoDal[model.CPSAction, model.CPSAction]
	customerDal dal.MongoDal[member.User, member.User]
}

func InitFaydaAccountPersistence(client *mongo.Client, database string, cpsCollection []string, logger utils.Logger) outbound.FaydaRepository {

	cpsDal := dal.NewMongoDal[model.CPSAction, model.CPSAction](client, database, cpsCollection[0])
	customerDal := dal.NewMongoDal[member.User, member.User](client, database, cpsCollection[1])

	return &FaydaAccountRepo{
		client:      client,
		logger:      logger,
		cpsDal:      cpsDal,
		customerDal: customerDal,
	}
}

func (f *FaydaAccountRepo) AuthorizeFaydaAccountEnableDisable(ctx context.Context, req *entities.CPSAction) (*entities.CPSAction, error) {

	actionData := make(map[string]interface{})
	byte, err := json.Marshal(req.PreviousAction)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(byte, &actionData); err != nil {
		return nil, err
	}

	customerFilter := bson.M{
		"phone_number": actionData["phone_number"],
		"kyc.level":    1,
	}
	customerUpdate := bson.M{
		"is_account_blocked": true,
	}

	updated, err := f.customerDal.UpdateOne(ctx, customerFilter, customerUpdate)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("NOT_FOUND")
		}
		f.logger.Errorf("failed to update  action", err)
		return nil, fmt.Errorf("FAILED_TO_UPDATE_USER_DATA")
	}
	req.CurrentAction = updated
	return req, nil
}
