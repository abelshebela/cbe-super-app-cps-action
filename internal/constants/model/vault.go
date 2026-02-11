package model

// type VaultGroupCategory struct {
// 	ID         string     `json:"id" bson:"id"`
// 	Name       string     `json:"name" bson:"name"`
// 	CoverImage string     `json:"cover_image" bson:"cover_image"`
// 	IsActive   bool       `json:"is_active" bson:"is_active"`
// 	IsDeleted  bool       `json:"is_deleted" bson:"is_deleted"`
// 	CreatedAt  time.Time  `json:"created_at" bson:"created_at"`
// 	UpdatedAt  time.Time  `json:"updated_at" bson:"updated_at"`
// 	DeletedAt  *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
// 	CreatedBy  string     `json:"created_by" bson:"created_by"`
// 	UpdatedBy  string     `json:"updated_by" bson:"updated_by"`
// }

// type UpdateVaultGroupCategory struct {
// 	Name       string    `json:"name" bson:"name"`
// 	CoverImage *string   `json:"cover_image" bson:"cover_image"`
// 	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
// 	UpdatedBy  *string   `json:"updated_by" bson:"updated_by"`
// }

type Tiers struct {
	Name     string `json:"name" bson:"name"`
	Interest string `json:"interest" bson:"interest"`
	Min      string `json:"min" bson:"min"`
	Max      string `json:"max" bson:"max"`
	Status   string `json:"status" bson:"status"`
}

type Category struct {
	ID          string `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Interest    string `json:"interest" bson:"interest"`
	Tiers       Tiers  `json:"tiers" bson:"tiers"`
	Deadlock    string `json:"deadlock" bson:"deadlock"`
	Status      string `json:"status" bson:"status"`
	IsEnabled   string `json:"is_enabled" bson:"is_enabled"`
	CreatedAt   string `json:"created_at" bson:"created_at"`
	UpdatedAt   string `json:"updated_at" bson:"updated_at"`
}

type VaultTransaction struct {
	ID string `json:"id" bson:"id"`
}

type VaultWithdrawal struct {
	ID                    string `json:"id" bson:"id"`
	LockedVaultID         string `json:"locked_vault_id" bson:"locked_vault_id"`
	WithdrawalAmount      string `json:"withdrawal_amount" bson:"withdrawal_amount"`
	WithdrawerName        string `json:"withdrawer_name" bson:"withdrawer_name"`
	WithdrawerPhoneNumber string `json:"withdrawer_phone_number" bson:"withdrawer_phone_number"`
}
