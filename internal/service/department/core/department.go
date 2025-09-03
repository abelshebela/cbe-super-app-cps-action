package department_core

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"
)

func Department_mapper(a interface{}) model.Department {
	dept := model.Department{}

	action := a.(map[string]interface{})

	if v, ok := action["department_code"].(string); ok {
		dept.DepartmentCode = v
	}
	if v, ok := action["department"].(string); ok {
		dept.Department = v
	}

	if v, ok := action["portal_cards"].([]interface{}); ok {
		for _, item := range v {
			if s, ok := item.(string); ok {
				dept.PortalCards = append(dept.PortalCards, s)
			}
		}
	}
	if v, ok := action["enabled"].(bool); ok {
		dept.Enabled = v
	}
	if v, ok := action["is_deleted"].(bool); ok {
		dept.IsDeleted = v
	}
	if v, ok := action["created_at"].(time.Time); ok {
		dept.CreatedAt = v
	}
	if v, ok := action["last_modified"].(time.Time); ok {
		dept.LastModified = v
	}

	return dept
}
