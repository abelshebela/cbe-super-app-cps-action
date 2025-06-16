package bank

import "time"

type Bank struct {
	ID             string
	Name           string
	Logo           string
	Code           string
	BIC            string
	Enabled        bool
	CreatedAt      time.Time
	LastModifiedAt time.Time
}
