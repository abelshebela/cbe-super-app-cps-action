package wallet

import "time"

type Wallet struct {
	ID             string
	Name           string
	Code           string
	Avatat         string
	Enabled        bool
	IsDeleted      bool
	CreatedAt      time.Time
	LastModifiedAt time.Time
	DeletedAt      time.Time
}
