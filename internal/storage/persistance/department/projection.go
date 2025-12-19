package department

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func DepartmentMapper(data model.Department) bson.M {
	result := bson.M{}
	if data.Department != "" {
		result["department"] = data.Department
	}
	if data.DepartmentCode != "" {
		result["department_code"] = data.DepartmentCode
	}
	if data.PortalCards != nil {
		result["portal_cards"] = data.PortalCards
	}

	result["enabled"] = data.Enabled
	result["is_deleted"] = data.IsDeleted
	if !data.CreatedAt.IsZero() {
		result["created_at"] = data.CreatedAt
	}
	if !data.LastModified.IsZero() {
		result["last_modified"] = data.LastModified
	}
	return result
}
