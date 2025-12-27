package model

type CPSActionList struct {
	ActionName     string `bson:"action_name" json:"action_name"`
	ActionCode     string `bson:"action_code" json:"action_code"`
	Description    string `bson:"description" json:"description"`
	PortalCardName string `bson:"portal_card_name" json:"portal_card_name"`
	IsConfigured   bool   `bson:"is_configured" json:"is_configured"`
}
