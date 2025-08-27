package model

import (
// "cbe-super-app-cps-action/internal/constants/types"

// "go.mongodb.org/mongo-driver/v2/bson"
)

// type APPAccessList struct {
// 	ID             bson.ObjectID   `json:"id" bson:"_id"`
// 	Key            string          `json:"key" bson:"key"`
// 	Enabled        bool            `json:"enabled" bson:"enabled"`
// 	AccessListName string          `json:"access_list_name" bson:"access_list_name"`
// 	SubAccessList  []SubAccessList `json:"sub_access_list" bson:"sub_access_list"`
// 	USSDEnabled    bool            `json:"ussd_enabled" bson:"ussd_enabled"`
// }

// type SubAccessList struct {
// 	Key            string `json:"key" bson:"key"`
// 	Enabled        bool   `json:"enabled" bson:"enabled"`
// 	AccessListName string `json:"accessListName" bson:"accessListName"`
// }

type SubAccessList struct {
	Key            string `json:"key" bson:"key"`
	Enabled        bool   `json:"enabled" bson:"enabled"`
	AccessListName string `json:"accessListName" bson:"accessListName"`
}

type APPAccessList struct {
	Key            string          `json:"key" bson:"key"`
	Enabled        bool            `json:"enabled" bson:"enabled"`
	AccessListName string          `json:"accessListName" bson:"accessListName"`
	SubAccessList  []SubAccessList `json:"subAccessList" bson:"subAccessList"`
	USSDEnabled    bool            `json:"USSDEnabled" bson:"USSDEnabled"`
}
