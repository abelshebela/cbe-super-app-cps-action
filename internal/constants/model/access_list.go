package model

import (
	"cbe-super-app-cps-action/internal/constants/types"
)

type APPAccessList struct {
	Key            string                `json:"key" bson:"key"`
	Enabled        bool                  `json:"enabled" bson:"enabled"`
	AccessListName string                `json:"accessListName" bson:"accessListName"`
	SubAccessList  []types.SubAccessList `json:"subAccessList" bson:"subAccessList"`
	USSDEnabled    bool                  `json:"USSDEnabled" bson:"USSDEnabled"`
}
