package model

type BudgetCategoryOracle struct {
	ID             string `bson:"id" json:"id,omitempty"`
	Name           string `bson:"name" json:"name"`
	AccountType    string `bson:"account_type" json:"account_type"` // Changed from Type to AccountType
	Color          string `bson:"color" json:"color"`
	Icon           string `bson:"icon" json:"icon"`
	IsEnabled      int    `bson:"is_enabled" json:"is_enabled"`
	IsDeleted      int    `bson:"is_deleted" json:"is_deleted"`
	CreatedAt      string `bson:"created_at" json:"created_at,omitempty"`             // Renamed from CreateAt
	LastModifiedAt string `bson:"last_modified_at" json:"last_modified_at,omitempty"` // Renamed from UpdateAt
	DeletedAt      string `bson:"deleted_at" json:"deleted_at,omitempty"`             // New field
}
