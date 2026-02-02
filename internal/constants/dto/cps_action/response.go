package cpsaction

type CPSActionCountResponse struct {
	Pending    int `json:"total_pending"`
	Approved   int `json:"total_approved"`
	Rejected   int `json:"total_rejected"`
	Canceled   int `json:"total_canceled"`
	Inprogress int `json:"total_inprogress_audit"`
	Completed  int `json:"total_completed_audit"`
}

type AutorizersLevelResponse struct {
	CheckerIndex int64 `json:"checker_index"`
	AuditorIndex int64 `json:"auditor_index"`
}
