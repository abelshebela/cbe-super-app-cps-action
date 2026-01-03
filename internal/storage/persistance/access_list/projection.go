package access_list

import (
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	shared_types "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// AccessListProjection is a struct that implements the storage.AppAccessListRepository interface
func AccessListMapper(data model.APPAccessList) bson.M {
	result := bson.M{}
	if data.Key != "" {
		result["key"] = data.Key
	}
	if data.AccessListName != "" {
		result["accessListName"] = data.AccessListName
	}
	if data.SubAccessList != nil || len(data.SubAccessList) == 0 {
		result["subAccessList"] = local_util.MapSlice(data.SubAccessList, func(item shared_types.SubAccessList) bson.M {
			return bson.M{
				"key":              item.Key,
				"access_list_name": item.AccessListName,
				"enabled":          item.Enabled,
			}
		})
	}

	return result
}
