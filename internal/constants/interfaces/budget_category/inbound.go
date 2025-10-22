package budget_category

import "net/http"

type BudgetCategoryPortHandler interface {
	CreateBudgetCategory(w http.ResponseWriter, r *http.Request)
	UpdateBudgetCategory(w http.ResponseWriter, r *http.Request)
	GetBudgetCategoryByID(w http.ResponseWriter, r *http.Request)
	GetAllBudgetCategories(w http.ResponseWriter, r *http.Request)
	DeleteBudgetCategory(w http.ResponseWriter, r *http.Request)
	EnableBudgetCategory(w http.ResponseWriter, r *http.Request)
	DisableBudgetCategory(w http.ResponseWriter, r *http.Request)
}