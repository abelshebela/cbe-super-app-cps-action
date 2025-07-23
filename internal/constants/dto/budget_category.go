package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BudgetCategory struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name" bson:"name"`
	Logo           string             `json:"logo" bson:"logo"`
	Color          string             `json:"color" bson:"color"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time          `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt      *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

func NewBudgetCategory(req BudgetCategoryRequest) (*BudgetCategory, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	now := time.Now()
	return &BudgetCategory{
		ID:             req.ID,
		Name:           req.Name,
		Logo:           req.Logo,
		Color:          req.Color,
		CreatedAt:      now,
		LastModifiedAt: now,
	}, nil
}

type BudgetCategoryRequest struct {
	ID    primitive.ObjectID `json:"id,omitempty"`
	Name  string             `json:"name"`
	Logo  string             `json:"logo"`
	Color string             `json:"color"`
}

func (bcr *BudgetCategoryRequest) Validate() error {
	return validate.Struct(bcr)
}
