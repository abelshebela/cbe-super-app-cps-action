package model

import (
	"cbe-super-app-cps-action/internal/constants/types"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type APPAccessList struct {
	ID             bson.ObjectID         `json:"id,omitempty" bson:"_id"`
	Key            string                `json:"key" bson:"key"`
	Enabled        bool                  `json:"enabled" bson:"enabled"`
	AccessListName string                `json:"accessListName" bson:"accessListName"`
	SubAccessList  []types.SubAccessList `json:"subAccessList" bson:"subAccessList"`
	USSDEnabled    bool                  `json:"USSDEnabled" bson:"USSDEnabled"`
}
