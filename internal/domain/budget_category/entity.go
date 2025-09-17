package budget_category

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BudgetCategory struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name          string        `bson:"name,omitempty" json:"name,omitempty"`
	Icon          string        `bson:"icon,omitempty" json:"icon,omitempty"`
	Description   string        `bson:"description,omitempty" json:"description,omitempty"`
	CreatedAt     time.Time     `bson:"created_at,omitempty" json:"created_at,omitempty"`
	LastUpdatedAt time.Time     `bson:"last_updated_at,omitempty" json:"last_updated_at,omitempty"`
	IsDeleted     bool          `bson:"is_deleted" json:"is_deleted,omitempty"`
}
