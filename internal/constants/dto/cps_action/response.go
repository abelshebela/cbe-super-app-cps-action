package cpsaction

type CPSActionCountResponse struct {
	Pending  int `json:"total_pending"`
	Approved int `json:"total_approved"`
	Rejected int `json:"total_rejected"`
}
