package core

import (
	"encoding/json"
	"time"

	vaultgroup_category "cbe-super-app-cps-action/internal/constants/dto/vaultgroup_category"
	"cbe-super-app-cps-action/internal/constants/model"
)

func CategoryMapper(action map[string]interface{}) model.VaultGroupCategory {
	category := model.VaultGroupCategory{}

	if v, ok := action["name"]; ok {
		if name, ok := v.(string); ok {
			category.Name = name
		}
	}
	if v, ok := action["cover_image"]; ok {
		if img, ok := v.(string); ok {
			category.CoverImage = img
		}
	}
	return category
}

// ConvertVaultGroupCategoryToMongoSafe converts VaultGroupCategory to a MongoDB-safe format
// by converting decimal.Decimal fields to float64 for proper serialization
// func ConvertVaultGroupCategoryToMongoSafe(vaultGroupCategory *model.VaultGroupCategory) map[string]interface{} {
// 	result := map[string]interface{}{
// 		"id":        vaultGroupCategory.ID,
// 		"name":      vaultGroupCategory.Name,
// 		"isactive":  vaultGroupCategory.IsActive,
// 		"createdat": vaultGroupCategory.CreatedAt,
// 		"updatedat": vaultGroupCategory.UpdatedAt,
// 		"deletedat": vaultGroupCategory.DeletedAt,
// 		"createdby": vaultGroupCategory.CreatedBy,
// 		"updatedby": vaultGroupCategory.UpdatedBy,
// 		"isdeleted": vaultGroupCategory.IsDeleted,
// 	}
// 	return result
// }

func MapVaultGroupCategoryToResponse(vaultGroupCategory *model.VaultGroupCategory) *vaultgroup_category.VaultGroupCategoryResponse {
	return &vaultgroup_category.VaultGroupCategoryResponse{
		ID:         vaultGroupCategory.ID,
		Name:       vaultGroupCategory.Name,
		CoverImage: vaultGroupCategory.CoverImage,
		IsActive:   vaultGroupCategory.IsActive,
		IsDeleted:  vaultGroupCategory.IsDeleted,
		CreatedAt:  vaultGroupCategory.CreatedAt,
		UpdatedAt:  vaultGroupCategory.UpdatedAt,
		DeletedAt:  vaultGroupCategory.DeletedAt,
	}
}

// func BuildUpdateVaultGroupCategory(prev *model.VaultGroupCategory, req *model.VaultGroupCategory) *model.VaultGroupCategory {
// 	if req.Name != "" && req.Name != prev.Name {
// 		req.Name = req.Name
// 	}
// 	if req.CoverImage != "" && req.CoverImage != prev.CoverImage {
// 		req.CoverImage = req.CoverImage
// 	}

// 	return req
// }

func BindVaultGroupCategoryFromCPSAction(current interface{}) (model.VaultGroupCategory, error) {
	var VaultGroupCategory model.VaultGroupCategory
	if v, ok := current.(model.VaultGroupCategory); ok {
		return v, nil
	}
	if s, ok := current.(string); ok {
		if err := json.Unmarshal([]byte(s), &VaultGroupCategory); err == nil {
			return VaultGroupCategory, nil
		}
	}
	bytes, err := json.Marshal(current)
	if err != nil {
		return VaultGroupCategory, err
	}
	return MapCamelCaseToVaultGroupCategory(bytes)
}

func MapCamelCaseToVaultGroupCategory(jsonBytes []byte) (model.VaultGroupCategory, error) {
	var result model.VaultGroupCategory

	var data map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return result, err
	}

	// map camelcase fields to snake_case fields

	result.ID = getString(data, "id")
	result.Name = getString(data, "name")
	result.IsActive = getBool(data, "isactive")
	result.CreatedBy = getString(data, "createdby")
	result.UpdatedBy = getString(data, "updatedby")
	result.IsDeleted = getBool(data, "isdeleted")
	result.CreatedAt = getTime(data, "createdat")
	result.UpdatedAt = getTime(data, "updatedat")
	result.DeletedAt = getTimePtr(data, "deletedat")

	return result, nil
}

func BindVaultGroupCategoryUpdateFromCPSAction(current interface{}) (model.VaultGroupCategory, error) {
	var BV model.VaultGroupCategory
	if v, ok := current.(model.VaultGroupCategory); ok {
		return v, nil
	}

	// If stored as JSON string
	if s, ok := current.(string); ok {
		if err := json.Unmarshal([]byte(s), &BV); err == nil {
			return BV, nil
		}
	}

	// Generic path: marshal then unmarshal
	bytes, err := json.Marshal(current)
	if err != nil {
		return BV, err
	}

	// Use custom mapping function to handle camelCase -> UpdateBankVault conversion
	return MapCamelCaseToUpdateaultGroupCategory(bytes)
}
func MapCamelCaseToUpdateaultGroupCategory(jsonBytes []byte) (model.VaultGroupCategory, error) {
	var result model.VaultGroupCategory
	var data map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return result, err
	}

	// map camelcase fields to snake_case fields
	if desc, ok := data["name"].(string); ok && desc != "" {
		result.Name = desc
	}
	return result, nil
}

func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

func getBool(data map[string]interface{}, key string) bool {
	if val, ok := data[key].(bool); ok {
		return val
	}
	return false
}

func getTime(data map[string]interface{}, key string) time.Time {
	if val, ok := data[key].(string); ok {
		if t, err := time.Parse(time.RFC3339, val); err == nil {
			return t
		}
	}
	return time.Time{}
}

func getTimePtr(data map[string]interface{}, key string) *time.Time {
	if val, ok := data[key].(string); ok {
		if t, err := time.Parse(time.RFC3339, val); err == nil {
			return &t
		}
	}
	return nil
}

// func VaultGroupCategoryUpdate(req *model.VaultGroupCategory) model.VaultGroupCategory {
// 	result := model.VaultGroupCategory{}
// 	if req.Name != "" {
// 		result.Name = req.Name
// 	}
// 	return result
// }
