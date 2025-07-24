package miniappmerchant

import (
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateMiniAppMerchantDTO struct {
	Type                       string `json:"type" validate:"required"`
	MerchantName               string `json:"merchant_name" validate:"required"`
	MerchantRepresentativeName string `json:"merchant_representative_name" validate:"required"`
	PhoneNumber                string `json:"phone_number" validate:"required"`
	Email                      string `json:"email" validate:"required,email"`
	AccountNumber              string `json:"account_number" validate:"required"`
}

type UpdateMiniAppMerchantDTO struct {
	Type                       *string `json:"type,omitempty"`
	MerchantName               *string `json:"merchant_name,omitempty"`
	MerchantRepresentativeName *string `json:"merchant_representative_name,omitempty"`
	PhoneNumber                *string `json:"phone_number,omitempty"`
	Email                      *string `json:"email,omitempty"`
	AccountNumber              *string `json:"account_number,omitempty"`
	MiniAppID                  *string `json:"mini_app_id,omitempty"`
	Enabled                    *bool   `json:"enabled,omitempty"`
}

type MiniAppMerchantResponseDTO struct {
	ID            string    `json:"id"`
	Code          string    `json:"code"`
	Type          string    `json:"type"`
	MerchantName  string    `json:"merchant_name"`
	KYC           KYCDTO    `json:"kyc"`
	AccountNumber string    `json:"account_number"`
	MiniAppIDs    []string  `json:"mini_app_ids"`
	Enabled       bool      `json:"enabled"`
	IsDeleted     bool      `json:"is_deleted"`
	CreatedAt     time.Time `json:"created_at"`
	LastModified  time.Time `json:"last_modified"`
}

type KYCDTO struct {
	Status         string            `json:"status"`
	Representative RepresentativeDTO `json:"representative"`
}

type RepresentativeDTO struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func (dto CreateMiniAppMerchantDTO) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.Type, validation.Required.Error("type is required")),
		validation.Field(&dto.MerchantName, validation.Required.Error("merchant name is required")),
		validation.Field(&dto.MerchantRepresentativeName, validation.Required.Error("representative name is required")),
		validation.Field(&dto.PhoneNumber,
			validation.Required.Error("phone number is required"),
			is.Digit.Error("phone number must contain only digits"),
			validation.Length(9, 15).Error("phone number must be between 9 and 15 digits"),
			validation.Match(regexp.MustCompile(`^(?:\+251|251|0)9\d{8}$`)).Error("invalid  phone number format"),
		),
		validation.Field(&dto.Email,
			validation.Required.Error("email is required"),
			is.Email.Error("email must be a valid email address"),
		),
		validation.Field(&dto.AccountNumber, validation.Required.Error("account number is required")),
	)
}

func (dto UpdateMiniAppMerchantDTO) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.PhoneNumber,
			validation.When(dto.PhoneNumber != nil,
				is.Digit.Error("phone number must contain only digits"),
				validation.Length(9, 15).Error("phone number must be between 9 and 15 digits"),
				validation.Match(regexp.MustCompile(`^(?:\+251|251|0)9\d{8}$`)).Error("invalid phone number format"),
			),
		),
		validation.Field(&dto.Email,
			validation.When(dto.Email != nil,
				is.Email.Error("email must be a valid email address"),
			),
		),
	)
}
