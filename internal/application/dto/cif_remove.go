package dto

type CifRemoveMakerRequest struct {
	UserID string   `json:"user_id"`
	Cif    []string `json:"cif"`
}

type CifSearchRequest struct {
	Cif string `json:"cif"`
}

type CifSearchResposnse struct {
	Cif string `json:"cif"`
}

type CifRemoveMakerResponse struct {
	UserID     string `json:"user_id"`
	ActionCode string `json:"action_code"`
}

type CifRemoveCheckerRequest struct {
	ActionID      string `json:"action_id"`
	ServiceAction bool   `json:"service_action"`
	RejectReason  string `json:"rejection_reason"`
}
