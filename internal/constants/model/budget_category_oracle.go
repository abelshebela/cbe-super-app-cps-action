package model

type BudgetCategoryOracle struct {
	ID        string `bson:"id" json:"id,omitempty"`
	Name      string `bson:"name" json:"name"`
	Color     string `bson:"color" json:"color"`
	Icon      string `bson:"icon" json:"icon"`
	Type      string `bson:"type" json:"type"`
	IsEnabled int    `bson:"is_enabled" json:"is_enabled"`
	IsDeleted int    `bson:"is_deleted" json:"is_deleted"`
	CreateAt  string `bson:"create_at" json:"create_at,omitempty"`
	UpdateAt  string `bson:"update_at" json:"update_at,omitempty"`
}
