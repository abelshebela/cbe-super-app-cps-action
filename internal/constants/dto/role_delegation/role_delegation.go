package role_delegation_dto

type RoleDelegationRequest struct {
	UserID     string `json:"user_id"`
	JobTitleID string `json:"job_title"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
}
