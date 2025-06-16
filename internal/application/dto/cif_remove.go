package dto

type CifRemoveMakerRequest struct {
	UserId string `json:"user_id"`
	Cif    string `json:"cif"`
}

type CifSearchRequest struct {
	Cif string `json:"cif"`
}

type CifSearchResposnse struct {
	Cif string `json:"cif"`
}

type CifRemoveMakerResponse struct {
	UserId    string `json:"user_id"`
	Action_Id string `json:"action_id"`
}

type CifRemoveCheckerRequest struct {
	Action_Id     string `json:"service_id"`
	ServiceAction bool   `json:"service_action"`
}
