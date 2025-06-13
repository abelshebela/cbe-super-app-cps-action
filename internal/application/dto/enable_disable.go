package dto

type EnableDisableServiceMakerDtoRequest struct {
	ServiceId     string `json:"service_id"`
	ServiceAction bool   `json:"service_action"`
}

type EnableDisableServiceMakerDtoResponse struct {
	ActionId string `json:"action_id"`
}

type EnableDisableServiceCheckerDtoRequest struct {
	Action_Id     string `json:"service_id"`
	ServiceAction bool   `json:"service_action"`
}
