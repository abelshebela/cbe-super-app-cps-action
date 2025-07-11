package account_validation

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ValidationRule struct {
	ID             bson.ObjectID
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
