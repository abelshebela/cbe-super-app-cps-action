package account_product_category_core

import (
	"encoding/json"

	imodel "cbe-super-app-cps-action/internal/constants/model"
)

// MapFromAction deserialises a CPS CurrentAction payload back to AccountProductCategory.
func MapFromAction(raw interface{}) (imodel.AccountProductCategory, error) {
	b, err := json.Marshal(raw)
	if err != nil {
		return imodel.AccountProductCategory{}, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return imodel.AccountProductCategory{}, err
	}

	apc := imodel.AccountProductCategory{}

	if v, ok := m["id"].(string); ok {
		apc.ID = v
	}
	if v, ok := m["account_type"].(string); ok {
		apc.AccountType = v
	}
	if v, ok := m["category_name"].(string); ok {
		apc.CategoryName = v
	}
	if v, ok := m["cbs_category_code"].(string); ok {
		apc.CBSCategoryCode = v
	}
	if v, ok := m["description"].(string); ok {
		apc.Description = v
	}
	if v, ok := m["is_enabled"].(bool); ok {
		apc.IsEnabled = v
	}
	if v, ok := m["is_deleted"].(bool); ok {
		apc.IsDeleted = v
	}

	return apc, nil
}
