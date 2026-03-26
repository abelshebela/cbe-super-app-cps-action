package budget_category

import (
	"mime/multipart"
)

type CreateBudgetRequest struct {
	Name  string                `form:"name" example:"Monthly Groceries"`
	Color string                `form:"color" example:"#FF5733"`
	Icon  *multipart.FileHeader `form:"icon"`
	Type  string                `form:"type" example:"CB"`
}

type UpdateBudgetRequest struct {
	Name  string                `form:"name,omitempty" example:"Monthly Groceries"`
	Color string                `form:"color,omitempty" example:"#FF5733"`
	Icon  *multipart.FileHeader `form:"icon,omitempty"`
	Type  string                `form:"type,omitempty" example:"CB"`
}
