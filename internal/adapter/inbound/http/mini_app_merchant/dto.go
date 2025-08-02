package miniappmerchant

import (
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)



type MiniAppMerchantDTO struct {
	Type                       string `json:"type"`
	MerchantName               string `json:"merchant_name"`
	MerchantRepresentativeName string `json:"merchant_representative_name"`
	PhoneNumber                string `json:"phone_number"`
	Email                      string `json:"email"`
	AccountNumber              string `json:"account_number"`
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

func (dto MiniAppMerchantDTO) IsEmpty() bool {
	return dto.Type == "" &&
		dto.MerchantName == "" &&
		dto.MerchantRepresentativeName == "" &&
		dto.PhoneNumber == "" &&
		dto.Email == "" &&
		dto.AccountNumber == ""
}

func (dto MiniAppMerchantDTO) Validate(isCreate bool) error {
	if !isCreate && dto.IsEmpty() {
		return nil
	}

	var rules []*validation.FieldRules

	if isCreate {
		rules = []*validation.FieldRules{
			validation.Field(&dto.Type, validation.Required.Error("type is required")),
			validation.Field(&dto.MerchantName, validation.Required.Error("merchant name is required")),
			validation.Field(&dto.MerchantRepresentativeName, validation.Required.Error("representative name is required")),
			validation.Field(&dto.PhoneNumber,
				validation.Required.Error("phone number is required"),
				is.Digit.Error("phone number must contain only digits"),
				validation.Length(9, 15).Error("phone number must be between 9 and 15 digits"),
				validation.Match(regexp.MustCompile(`^(?:\+251|251|0)9\d{8}$`)).Error("invalid phone number format"),
			),
			validation.Field(&dto.Email,
				validation.Required.Error("email is required"),
				is.Email.Error("email must be a valid email address"),
			),
			validation.Field(&dto.AccountNumber, validation.Required.Error("account number is required")),
		}
	} else {
		if dto.Type != "" {
			rules = append(rules, validation.Field(&dto.Type, validation.Required.Error("type is required")))
		}
		if dto.MerchantName != "" {
			rules = append(rules, validation.Field(&dto.MerchantName, validation.Required.Error("merchant name is required")))
		}
		if dto.MerchantRepresentativeName != "" {
			rules = append(rules, validation.Field(&dto.MerchantRepresentativeName, validation.Required.Error("representative name is required")))
		}
		if dto.PhoneNumber != "" {
			rules = append(rules, validation.Field(&dto.PhoneNumber,
				is.Digit.Error("phone number must contain only digits"),
				validation.Length(9, 15).Error("phone number must be between 9 and 15 digits"),
				validation.Match(regexp.MustCompile(`^(?:\+251|251|0)9\d{8}$`)).Error("invalid phone number format"),
			))
		}
		if dto.Email != "" {
			rules = append(rules, validation.Field(&dto.Email, is.Email.Error("email must be a valid email address")))
		}
		if dto.AccountNumber != "" {
			rules = append(rules, validation.Field(&dto.AccountNumber, validation.Required.Error("account number is required")))
		}
	}

	if len(rules) > 0 {
		if err := validation.ValidateStruct(&dto, rules...); err != nil {
			return err
		}
	}

	return nil
}