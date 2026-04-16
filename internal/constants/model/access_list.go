package model

import (
	"cbe-super-app-cps-action/internal/constants/types"
)

// APPAccessList is the app/bulk-service shape (Mongo + API). Oracle persistence uses
// AccessListOracle (SERVICE_KEY → Key, NAME → AccessListName, IS_ENABLED → Enabled).
type APPAccessList struct {
	ID             string                `json:"id" bson:"id"`
	Key            string                `json:"key" bson:"key"`
	Enabled        bool                  `json:"enabled" bson:"enabled"`
	AccessListName string                `json:"access_list_name" bson:"access_list_name"`
	SubAccessList  []types.SubAccessList `json:"sub_access_list" bson:"sub_access_list"`
	USSDEnabled    bool                  `json:"ussd_enabled" bson:"ussd_enabled"`
}
