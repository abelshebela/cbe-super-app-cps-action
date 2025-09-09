package dto

type EnableDisableServiceMakerDtoRequest struct {
	ServiceId     string `json:"service_id" bson:"service_id"`
	ServiceAction bool   `json:"service_action" bson:"service_action"`
}

type EnableDisableServiceMakerDtoResponse struct {
	ActionId string `json:"action_id" bson:"action_id"`
}

type EnableDisableServiceCheckerDtoRequest struct {
	Action_Id     string `json:"action_code" bson:"action_code"`
	ServiceAction bool   `json:"action" bson:"action"`
	RejectReason  string `json:"rejected_reason" bson:"rejected_reason"`
}
