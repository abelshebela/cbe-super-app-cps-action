package budget_category

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func BudgetCategoryMapper(budgetCategory model.BudgetCategory) bson.M {
	result := bson.M{}

	if budgetCategory.Name != "" {
		result["name"] = budgetCategory.Name
	}
	if budgetCategory.Color != "" {
		result["color"] = budgetCategory.Color
	}
	if budgetCategory.Icon != "" {
		result["icon"] = budgetCategory.Icon
	}
	result["enabled"] = budgetCategory.Enabled
	result["is_deleted"] = budgetCategory.IsDeleted

	if !budgetCategory.CreatedAt.IsZero() {
		result["created_at"] = budgetCategory.CreatedAt
	}
	if !budgetCategory.UpdatedAt.IsZero() {
		result["updated_at"] = budgetCategory.UpdatedAt
	}

	return result
}
