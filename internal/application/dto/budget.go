package dto

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

var (
	ErrActionApproved = errors.New("This action already approved")
	ErrActionRejected = errors.New("This action already rejected")
)

type CreateBudgetCategoryRequest struct {
	Name        string `json:"name" form:"name"`
	Icon        string `json:"icon" form:"icon"`
	BucketName  string `json:"bucket_name,omitempty" form:"bucket_name,omitempty"`
	ObjectName  string `json:"object_name,omitempty" form:"object_name,omitempty"`
	Description string `json:"description,omitempty" form:"description,omitempty"`
}

type UpdateBudgetCategoryRequest struct {
	ID          string `json:"id,omitempty" bson:"id,omitempty"`
	ActionCode  string `json:"action_code,omitempty" bson:"action_code,omitempty"`
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Icon        string `json:"icon,omitempty" bson:"icon,omitempty"`
	BucketName  string `json:"bucket_name,omitempty" bson:"bucket_name,omitempty"`
	ObjectName  string `json:"object_name,omitempty" bson:"object_name,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}

type GetBudgetCategoryRequest struct {
	ID string `json:"id"`
}

type DeleteBudgetCategoryRequest struct {
	ID string `json:"id"`
}

type GetAllBudgetCategoryRequest struct {
	Search string `json:"search,omitempty"`
	Page   int64  `json:"page,omitempty"`
	Limit  int64  `json:"limit,omitempty"`
}

type ApproveBudgetCategoryRequest struct {
	ActionID string `json:"action_id"`
	Approve  bool   `json:"approve"`
	Reason   string `json:"reason"`
}

func (r CreateBudgetCategoryRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name,
			validation.Required.Error("name is required"),
			validation.Length(2, 100).Error("name must be between 2-100 characters"),
		),
		validation.Field(&r.Description,
			validation.Length(0, 500).Error("description cannot exceed 500 characters"),
		),
	)
}

func (r UpdateBudgetCategoryRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID,
			validation.Required.Error("ID is required"),
		),
		validation.Field(&r.Name,
			validation.When(r.Name != "", validation.Length(2, 100).Error("name must be between 2-100 characters")),
		),
		validation.Field(&r.Icon,
			validation.When(r.Icon != "", validation.By(validateIcon)),
		),
		validation.Field(&r.Description,
			validation.Length(0, 500).Error("description cannot exceed 500 characters"),
		),
	)
}

func (r GetBudgetCategoryRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID,
			validation.Required.Error("ID is required"),
		),
	)
}

func (r DeleteBudgetCategoryRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID,
			validation.Required.Error("ID is required"),
			is.MongoID.Error("invalid ID format"),
		),
	)
}

func (r GetAllBudgetCategoryRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Page,
			validation.Min(1).Error("page must be at least 1"),
		),
		validation.Field(&r.Limit,
			validation.Min(1).Error("limit must be at least 1"),
			validation.Max(100).Error("limit cannot exceed 100"),
		),
		validation.Field(&r.Search,
			validation.Length(0, 100).Error("search query cannot exceed 100 characters"),
		),
	)
}

func validateIcon(value interface{}) error {
	s, _ := value.(string)
	// Check if it's a URL
	if err := is.URL.Validate(s); err == nil {
		return nil
	}
	// Check if it's base64 encoded image
	if err := is.Base64.Validate(s); err == nil {
		return nil
	}
	return validation.NewError("validation_invalid_icon", "icon must be a valid URL or base64 encoded image")
}

func (r ApproveBudgetCategoryRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ActionID,
			validation.Required.Error("action_id is required"),
			is.MongoID.Error("invalid action_id format"),
		),
		validation.Field(&r.Approve,
			validation.Required.Error("approve is required"),
		),
		validation.Field(&r.Reason,
			validation.Required.Error("reason is required"),
		),
	)
}
