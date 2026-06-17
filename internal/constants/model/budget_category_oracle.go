package model

type BudgetCategoryOracle struct {
	ID             string `bson:"id" json:"id,omitempty"`
	Name           string `bson:"name" json:"name"`
	AccountType    string `bson:"account_type" json:"account_type"`
	Color          string `bson:"color" json:"color"`
	Icon           string `bson:"icon" json:"icon"`
	IsEnabled      int    `bson:"is_enabled" json:"is_enabled"`
	IsDeleted      int    `bson:"is_deleted" json:"is_deleted"`
	CreatedAt      string `bson:"created_at" json:"created_at,omitempty"`
	LastModifiedAt string `bson:"last_modified_at" json:"last_modified_at,omitempty"`
	DeletedAt      string `bson:"deleted_at" json:"deleted_at,omitempty"`
}

// BudgetCategoryOracleCPSPayload is stored in the CPS action (MongoDB) with bool is_enabled/is_deleted
// so the checker sees true/false instead of 0/1.
type BudgetCategoryOracleCPSPayload struct {
	ID             string `bson:"id" json:"id,omitempty"`
	Name           string `bson:"name" json:"name"`
	AccountType    string `bson:"account_type" json:"account_type"`
	Color          string `bson:"color" json:"color"`
	Icon           string `bson:"icon" json:"icon"`
	IsEnabled      bool   `bson:"is_enabled" json:"is_enabled"`
	IsDeleted      bool   `bson:"is_deleted" json:"is_deleted"`
	CreatedAt      string `bson:"created_at" json:"created_at,omitempty"`
	LastModifiedAt string `bson:"last_modified_at" json:"last_modified_at,omitempty"`
}

func BudgetCategoryOracleToPayload(bc BudgetCategoryOracle) BudgetCategoryOracleCPSPayload {
	return BudgetCategoryOracleCPSPayload{
		ID:             bc.ID,
		Name:           bc.Name,
		AccountType:    bc.AccountType,
		Color:          bc.Color,
		Icon:           bc.Icon,
		IsEnabled:      bc.IsEnabled == 1,
		IsDeleted:      bc.IsDeleted == 1,
		CreatedAt:      bc.CreatedAt,
		LastModifiedAt: bc.LastModifiedAt,
	}
}

func BudgetCategoryPayloadToOracle(p BudgetCategoryOracleCPSPayload) BudgetCategoryOracle {
	isEnabled := 0
	if p.IsEnabled {
		isEnabled = 1
	}
	isDeleted := 0
	if p.IsDeleted {
		isDeleted = 1
	}
	return BudgetCategoryOracle{
		ID:             p.ID,
		Name:           p.Name,
		AccountType:    p.AccountType,
		Color:          p.Color,
		Icon:           p.Icon,
		IsEnabled:      isEnabled,
		IsDeleted:      isDeleted,
		CreatedAt:      p.CreatedAt,
		LastModifiedAt: p.LastModifiedAt,
	}
}
