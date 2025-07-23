package dto

import (
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

type Budget struct {
	ID                                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	SuperAppUserID                    string             `json:"superapp_user_id" bson:"superapp_user_id"`
	BudgetCategory                    primitive.ObjectID `json:"budget_category" bson:"budget_category"`
	BudgetType                        string             `json:"budget_type" bson:"budget_type"`
	StartingDate                      time.Time          `json:"starting_date" bson:"starting_dat"`
	EndingDate                        time.Time          `json:"ending_date" bson:"ending_date"`
	OverspendNotification             bool               `json:"overspend_notification" bson:"overspend_notification"`
	LimitedBudgetExceededNotifiaction bool               `json:"limited_budget_exceeded_notification" bson:"limited_budget_exceeded_notification"`
	Spending                          float64            `json:"speding" bson:"spending"`
	CreatedAt                         time.Time          `json:"created_at" bson:"created_at"`
	LastModifiedAt                    time.Time          `json:"lasT_modified_at" bson:"last_modified_at"`
	DeletedAt                         *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

func NewBudget(req BudgetRequest) (*Budget, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Budget{
		ID:                                req.ID,
		SuperAppUserID:                    req.SuperAppUserID,
		BudgetCategory:                    req.BudgetCategory,
		BudgetType:                        req.BudgetType,
		StartingDate:                      req.StartingDate,
		EndingDate:                        req.EndingDate,
		OverspendNotification:             req.OverspendNotification,
		LimitedBudgetExceededNotifiaction: req.LimitedBudgetExceededNotifiaction,
		Spending:                          req.Spending,
		CreatedAt:                         now,
		LastModifiedAt:                    now,
	}, nil
}

type BudgetRequest struct {
	ID                                primitive.ObjectID `json:"id,omitempty"`
	SuperAppUserID                    string             `json:"superapp_user_id"`
	BudgetCategory                    primitive.ObjectID `json:"budget_category"`
	BudgetType                        string             `json:"budget_type"`
	StartingDate                      time.Time          `json:"starting_date"`
	EndingDate                        time.Time          `json:"ending_date"`
	OverspendNotification             bool               `json:"overspend_notification"`
	LimitedBudgetExceededNotifiaction bool               `json:"limited_budget_exceeded_notification"`
	Spending                          float64            `json:"speding"`
}

func (br *BudgetRequest) Validate() error {
	return validate.Struct(br)
}
