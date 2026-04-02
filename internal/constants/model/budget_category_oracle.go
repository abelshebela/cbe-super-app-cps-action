package model

type BudgetCategoryOracle struct {
	ID        string `sqlx:"id" json:"id,omitempty"`
	Name      string `sqlx:"name" json:"name"`
	Color     string `sqlx:"color" json:"color"`
	Icon      string `sqlx:"icon" json:"icon"`
	Type      string `sqlx:"type" json:"type"`
	IsEnabled int    `sqlx:"is_enabled" json:"is_enabled"`
	IsDeleted int    `sqlx:"is_deleted" json:"is_deleted"`
	CreateAt  string `sqlx:"create_at" json:"create_at,omitempty"`
	UpdateAt  string `sqlx:"update_at" json:"update_at,omitempty"`
}
