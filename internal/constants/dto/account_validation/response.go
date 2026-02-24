package accountvalidation

// GetAccountValidationResponse represents the response for getting account validation
type GetAccountValidationResponse struct {
	Validation ValidationRuleDTO `json:"validation"`
}
