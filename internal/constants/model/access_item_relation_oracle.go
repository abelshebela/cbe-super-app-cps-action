package model

// AccessItemRelationOracle maps ACCESS_ITEMS_RELATION rows in Oracle.
type AccessItemRelationOracle struct {
	ID string `json:"id,omitempty"`

	ParentKey string `json:"parent_key"`
	ChildKey  string `json:"child_key"`
}
