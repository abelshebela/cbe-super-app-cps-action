package department_dto

type CreateDepartmentRequest struct {
	Department  string   `json:"department" bson:"department" example:"IT Department"`
	PortalCards []string `json:"portal_cards" bson:"portal_cards" example:"[\"card1\", \"card2\"]"`
}

type UpdateDepartmentRequest struct {
	Department  string   `json:"department,omitempty" bson:"department,omitempty" example:"IT Department Updated"`
	PortalCards []string `json:"portal_cards,omitempty" bson:"portal_cards,omitempty" example:"[\"card1\", \"card2\"]"`
}
