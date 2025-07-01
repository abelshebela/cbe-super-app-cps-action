package budget

import "net/http"

type BudgetPortHandler interface {
	CreateBudgetIcon(w http.ResponseWriter, r *http.Request)
	BudgetFetchIcons(w http.ResponseWriter, r *http.Request)
	BudgetUpdateIcon(w http.ResponseWriter, r *http.Request)
	BudgetCreateColor(w http.ResponseWriter, r *http.Request)
	BudgetFetchColors(w http.ResponseWriter, r *http.Request)
	BudgetUpdateColor(w http.ResponseWriter, r *http.Request)
	BudgetCheckerApproval(w http.ResponseWriter, r *http.Request)
}
