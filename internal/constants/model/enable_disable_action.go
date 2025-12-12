package model

// EnableDisableAction represents an action to enable or disable an entity(s)
type EnableDisableAction struct {
	Codes  []string `json:"codes" bson:"codes"`
	Reason string   `json:"reason" bson:"reason"`
}
