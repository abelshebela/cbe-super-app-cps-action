package model

type VaultAmountTier struct {
	ID              string  `json:"id"`
	VaultCategoryID string  `json:"vault_category_id"`
	MinAmount       float64 `json:"min_amount"`
	MaxAmount       float64 `json:"max_amount"`
	Interest        float64 `json:"interest"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	DeletedAt       *string `json:"deleted_at"`
	IsDeleted       bool    `json:"is_deleted"`
}
