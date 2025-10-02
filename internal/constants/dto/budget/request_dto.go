package budget

type BudgetCreateColor struct {
	Color string `json:"color" example:"#FF5733"`
}

type UpdateColorRequest struct {
	Color string `json:"color" example:"#33FF57"`
}
