package budget_category

import "net/http"

type BudgetCategoryPortHandler interface {
	CreateBudgetCategory(w http.ResponseWriter, r *http.Request)
	BudgetCategoryFetch(w http.ResponseWriter, r *http.Request)
	BudgetCategoryUpdate(w http.ResponseWriter, r *http.Request)
	BudgetCategoryFetchById(w http.ResponseWriter, r *http.Request)
}
