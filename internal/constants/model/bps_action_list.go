package model

type BPSActionList struct {
	ActionName   string `bson:"action_name" json:"action_name"`
	ActionCode   string `bson:"action_code" json:"action_code"`
	Description  string `bson:"description" json:"description"`
	IsConfigured bool   `bson:"is_configured" json:"is_configured"`
}
