package cpsaction

type CPSActionCountResponse struct {
	Pending    int `json:"total_pending" bson:"pending"`
	Approved   int `json:"total_approved" bson:"approved"`
	Rejected   int `json:"total_rejected" bson:"rejected"`
	Canceled   int `json:"total_canceled" bson:"canceled"`
	Inprogress int `json:"total_inprogress_audit" bson:"inprogress"`
	Completed  int `json:"total_completed_audit" bson:"completed"`
}

type AutorizersLevelResponse struct {
	CheckerIndex int64 `json:"checker_index"`
	AuditorIndex int64 `json:"auditor_index"`
}
