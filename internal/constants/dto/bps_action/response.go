package bps_action

type BPSActionCountResponse struct {
	Pending    int `json:"total_pending"`
	Approved   int `json:"total_approved"`
	Rejected   int `json:"total_rejected"`
	Inprogress int `json:"total_inprogress_audit"`
	Completed  int `json:"total_completed_audit"`
}
