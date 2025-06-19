package persistence

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
	amount_based_auth_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/amount_based_auth"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

type AmountBasedAuthRepo struct {
	client   *mongo.Client
	mongoDal dal.MongoDal[amount_based_auth_domain.AuthTier, amount_based_auth_domain.AuthTier]
	timeOut  time.Duration
	logger   utils.Logger
}

func (repo AmountBasedAuthRepo) ApproveAmountBasedAuth(ctx context.Context, id string) (string, error) {
	// Load configuration
	load, loadErr := config.Load()
	if loadErr != nil {
		log.Printf("Config load error: %v", loadErr)
		return "Failed to load config", loadErr
	}

	// Connect to MongoDB
	client, mongoErr := config.ConnectToMongoDB(load.MongoDBURI)
	if mongoErr != nil {
		log.Printf("Mongo connection error: %v", mongoErr)
		return "Failed to connect to database", mongoErr
	}

	// DAL for cpsaction collection
	actionDal := dal.NewMongoDal[action.CPSAction, action.CPSAction](client, load.MongoDBDatabase, "cpsaction")
	actionFilter := bson.M{"uniqueid": id}
	actionProjection := bson.M{}

	// Fetch the action document
	savedAction, fetchErr := actionDal.FindOne(ctx, actionFilter, actionProjection)
	if fetchErr != nil {
		log.Printf("Find action error: %v", fetchErr)
		return "Unable to find Action", fetchErr
	}
	if savedAction == nil {
		return "Unable to find Action", nil
	}

	// Parse string ID to ObjectID
	objectID, idParser := bson.ObjectIDFromHex(id)
	if idParser != nil {
		return "Failed to parse Object ID", nil
	}
	filter := bson.M{"_id": objectID}

	// Struct to unmarshal CurrentAction into
	type TierAction struct {
		MinAmount int64 `bson:"minAmount"`
		MaxAmount int64 `bson:"maxAmount"`
		Enabled   bool  `bson:"enabled"`
	}

	// Decode CurrentAction
	rawDoc, err := bson.Marshal(savedAction.CurrentAction)
	if err != nil {
		log.Printf("Marshal CurrentAction error: %v", err)
		return "Failed to process current action", err
	}

	var tier TierAction
	err = bson.Unmarshal(rawDoc, &tier)
	if err != nil {
		log.Printf("Unmarshal CurrentAction error: %v", err)
		return "Failed to parse current action", err
	}

	// Update the document in the DB
	update := bson.M{
		"minAmount": tier.MinAmount,
		"maxAmount": tier.MaxAmount,
		"enabled":   tier.Enabled,
	}

	_, updateErr := repo.mongoDal.UpdateOne(ctx, filter, update)
	if updateErr != nil {
		log.Printf("Update error: %v", updateErr)
		return "Error while updating tier", updateErr
	}

	// Delete the action document after successful approval
	deleteErr := actionDal.DeleteOne(ctx, actionFilter)
	if deleteErr != nil {
		log.Printf("Delete action error: %v", deleteErr)
		return "Error while deleting action", deleteErr
	}

	return "Update Success", nil
}

func (repo AmountBasedAuthRepo) UpdateAuthTier(request amount_based_auth_domain.AmountBasedAuthRequest, ctx context.Context, logger utils.Logger) (string, error) {
	// 1. Use filter by ID (or unique identifier)
	objectID, idParser := bson.ObjectIDFromHex(request.Id)
	if idParser != nil {
		return "Failed to parse object id", nil
	}

	filter := bson.M{"_id": objectID}
	projection := bson.M{}

	log.Printf("filter %s", filter)
	// 2. Fetch existing document
	tierById, err := repo.mongoDal.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			//log.Printf(err.Error())
			log.Printf("No action found — maybe the collection is empty.")
			// Optionally insert a default record here
		} else {
			log.Printf("Unexpected error: %v", err)
			return "Failed to fetch Auth tier", nil
		}

	}
	if tierById == nil {
		//logger.Errorf("auth tier not found")
		return "There is no such auth tier with the given id", nil
	}

	tier := bson.M{
		"Id":        tierById.ID,
		"minAmount": request.MinAmount,
		"maxAmount": request.MaxAmount,
		"enabled":   true,
	}

	load, loadErr := config.Load()
	if loadErr != nil {
		logger.Infof(loadErr.Error())
	}

	client, mongoErr := config.ConnectToMongoDB(load.MongoDBURI)
	if mongoErr != nil {
		logger.Infof(mongoErr.Error())
	}

	actionDal := dal.NewMongoDal[action.CPSAction, action.CPSAction](client, load.MongoDBDatabase, "cpsaction")
	actionFilter := bson.M{"uniqueid": request.Id}
	actionProjection := bson.M{}

	savedAction, fetchErr := actionDal.FindOne(ctx, actionFilter, actionProjection)
	if fetchErr != nil {
		log.Printf(fetchErr.Error())
		//logger.Infof(fetchErr.Error())
	}

	if savedAction != nil {
		return "There is a pending action. please Approve it first", nil
	}
	cpsAction := action.CPSAction{
		ActionCode:      "",
		UniqueID:        request.Id,
		Maker:           action.User{},
		Checker:         action.User{},
		Unique_ID:       request.Id,
		Department:      "",
		RejectionReason: nil,
		PreviosAction:   nil,
		CurrentAction:   tier,
		ActionStatus:    "",
		ActionType:      "",
		RequestAction:   "",
		CreatedAt:       time.Now(),
		LastModifiedAt:  time.Now(),
	}

	actionInsert, actionInsertErr := actionDal.InsertOne(ctx, cpsAction)
	if actionInsertErr != nil {
		logger.Infof(actionInsertErr.Error())
		return "Failed to insert cps action", actionInsertErr
	}

	log.Printf("actionInsert %v", actionInsert)

	return "Maker Action Success", nil

}

func InitAmountBasedAuth(client *mongo.Client, database string, collection string) *AmountBasedAuthRepo {
	mongoDal := dal.NewMongoDal[amount_based_auth_domain.AuthTier, amount_based_auth_domain.AuthTier](client, database, collection)

	return &AmountBasedAuthRepo{
		client:   client,
		mongoDal: mongoDal,
		timeOut:  5 * time.Second,
		logger:   nil,
	}
}
