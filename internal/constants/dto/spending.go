package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Spending struct {
	ID                   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	SuperAppUserID       string             `json:"superapp_user_id" bson:"superapp_user_id"`
	BudgetID             primitive.ObjectID `json:"budget_id" bson:"budget_id"`
	SpendingAmount       float32            `json:"pending_amount" bson:"pending_amount"`
	TransactionReference string             `json:"transaction_reference" bson:"transaction_reference"`
	TransactionDate      time.Time          `json:"transaction_date" bson:"transaction_date"`
	CreatedAt            time.Time          `json:"created_at" bson:"created_at"`
	LastModifiedAt       time.Time          `json:"last_modified_at" bson:"last_modified_at"`
}

func NewSpending(req SpendingRequest) (*Spending, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Spending{
		ID:                   req.ID,
		SuperAppUserID:       req.SuperAppUserID,
		BudgetID:             req.BudgetID,
		SpendingAmount:       req.SpendingAmount,
		TransactionReference: req.TransactionReference,
		TransactionDate:      req.TransactionDate,
		CreatedAt:            now,
		LastModifiedAt:       now,
	}, nil
}

type SpendingRequest struct {
	ID                   primitive.ObjectID `json:"id,omitempty"`
	SuperAppUserID       string             `json:"superapp_user_id"`
	BudgetID             primitive.ObjectID `json:"budget_id"`
	SpendingAmount       float32            `json:"pending_amount"`
	TransactionReference string             `json:"transaction_reference"`
	TransactionDate      time.Time          `json:"transaction_date"`
}

func (sr *SpendingRequest) Validate() error {
	return validate.Struct(sr)
}
