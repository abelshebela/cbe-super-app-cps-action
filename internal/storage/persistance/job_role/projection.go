package job_role

import (
	imodel "cbe-super-app-cps-action/internal/constants/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func JobRoleMapper(data imodel.JobRole) bson.M {
	result := bson.M{}
	if data.Code != "" {
		result["code"] = data.Code
	}
	if data.Name != "" {
		result["name"] = data.Name
	}
	if data.PortalCards != nil {
		result["portal_cards"] = data.PortalCards
	}
	if !data.UpdatedAt.IsZero() {
		result["updated_at"] = data.UpdatedAt
	}
	if !data.CreatedAt.IsZero() {
		result["created_at"] = data.CreatedAt
	}
	return result
}
