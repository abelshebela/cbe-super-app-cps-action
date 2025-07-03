package account_validation

import "time"

type ValidationRule struct {
	ID             string
	EntityType     string
	ValidationFor  string
	Identifier     string
	MinLength      uint8
	MaxLength      uint8
	Enabled        bool
	IsDeleted      bool
	CreatedAt      time.Time
	LastModifiedAt time.Time
	ServiceID      string
}
