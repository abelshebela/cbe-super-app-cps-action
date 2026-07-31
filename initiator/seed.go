package initiator

import (
	"context"

	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// SeedCPSActionList upserts the UTILITY service action into cps_action_list on startup.
// Uses $setOnInsert so existing documents are never overwritten.
func SeedCPSActionList(ctx context.Context, client *mongo.Client, dbName string, logger utils.Logger) {
	coll := client.Database(dbName).Collection(CPSActionListCollection)

	filter := bson.M{"action_code": "UTILITY"}
	update := bson.M{
		"$setOnInsert": imodel.CPSActionList{
			ActionName:     "Utility Service",
			ActionCode:     "UTILITY",
			Description:    "Management of utility services",
			PortalCardName: "utility",
			IsConfigured:   false,
			IsViewOnly:     false,
		},
	}

	data := coll.FindOne(ctx, filter)
	if data == nil {
		_, err := coll.InsertOne(ctx, update["$setOnInsert"])
		if err != nil {
			logger.Errorf("SeedCPSActionList: failed to upsert UTILITY action: %v", err)
			return
		}

	}

}
