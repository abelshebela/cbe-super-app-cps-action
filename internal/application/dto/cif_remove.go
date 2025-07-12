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
	UserID    string `json:"user_id"`
	ActionID  string `json:"action_id"`
}

type CifRemoveCheckerRequest struct {
	ActionID     string `json:"action_id"`
	ServiceAction bool   `json:"service_action"`
}
