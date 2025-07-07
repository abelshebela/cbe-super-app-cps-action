package budget_category

import "net/http"

type BudgetCategoryInbound interface {
	CreateBudgetCategory(w http.ResponseWriter, r *http.Request)
	UpdateBudgetCategory(w http.ResponseWriter, r *http.Request)
	DeleteBudgetCategory(w http.ResponseWriter, r *http.Request)
	GetBudgetCategory(w http.ResponseWriter, r *http.Request)
	GetAllBudgetCategory(w http.ResponseWriter, r *http.Request)
	ApproveBudgetCategoryActionHTTP(w http.ResponseWriter, r *http.Request)
}
