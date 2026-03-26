package bps_action

type BPSActionCountResponse struct {
	Pending    int `json:"total_pending" bson:"pending"`
	Approved   int `json:"total_approved" bson:"approved"`
	Rejected   int `json:"total_rejected" bson:"rejected"`
	Inprogress int `json:"total_inprogress_audit" bson:"inprogress"`
	Completed  int `json:"total_completed_audit" bson:"completed"`
}
