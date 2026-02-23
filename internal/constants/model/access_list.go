package model

import (
	"cbe-super-app-cps-action/internal/constants/types"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type APPAccessList struct {
	ID             bson.ObjectID         `json:"_id" bson:"_id"`
	Key            string                `json:"key" bson:"key"`
	Enabled        bool                  `json:"enabled" bson:"enabled"`
	AccessListName string                `json:"access_list_name" bson:"access_list_name"`
	SubAccessList  []types.SubAccessList `json:"sub_access_list" bson:"sub_access_list"`
	USSDEnabled    bool                  `json:"ussd_enabled" bson:"ussd_enabled"`
}
