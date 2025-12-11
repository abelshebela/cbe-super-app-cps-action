package model

// EnableDisableAction represents an action to enable or disable an entity(s)
type EnableDisableAction struct {
	ID      string `json:"id" bson:"id"`
	Name    string `json:"name" bson:"name"`
	Enabled bool   `json:"enabled" bson:"enabled"`
	Reason  string `json:"reason,omitempty" bson:"reason,omitempty"`
}
