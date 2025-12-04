package budget_category

type BudgetCategoryResponse struct {
	ID        string `json:"id" example:"507f1f77bcf86cd799439011"`
	Name      string `json:"name" example:"Monthly Groceries"`
	Color     string `json:"color" example:"#FF5733"`
	Icon      string `json:"icon" example:"shopping-cart"`
	Enabled   bool   `json:"enabled" example:"true"`
	CreatedAt string `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}
