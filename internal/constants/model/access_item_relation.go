package model

type AccessItemRelation struct {
	ParentKey string `json:"parent_key" bson:"parent_key"`
	ChildKey  string `json:"key" bson:"key"`
}
