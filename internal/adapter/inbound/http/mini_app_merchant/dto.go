package miniappmerchant

import (
	"regexp"
	"strings"
	"time"

	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp_merchant"
	utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
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
	ID            string            `json:"id"`
	Code          string            `json:"code"`
	Type          string            `json:"type"`
	MerchantName  string            `json:"merchant_name"`
	KYC           KYCDTO            `json:"kyc"`
	AccountNumber string            `json:"account_number"`
	MiniApps      []domain.MiniApps `json:"mini_apps"`
	Enabled       bool              `json:"enabled"`
	IsDeleted     bool              `json:"is_deleted"`
	CreatedAt     time.Time         `json:"created_at"`
	LastModified  time.Time         `json:"last_modified"`
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
		strings.TrimSpace(dto.MerchantName) == "" &&
		strings.TrimSpace(dto.MerchantRepresentativeName) == "" &&
		strings.TrimSpace(dto.PhoneNumber) == "" &&
		strings.TrimSpace(dto.Email) == "" &&
		strings.TrimSpace(dto.AccountNumber) == ""
}

func (dto MiniAppMerchantDTO) Validate(isCreate bool) error {
	if !isCreate && dto.IsEmpty() {
		return nil
	}

	var rules []*validation.FieldRules

	if isCreate {
		rules = []*validation.FieldRules{
			validation.Field(&dto.Type,
				validation.Required.Error("type is required"),
				validation.By(utils.NoSpecialChars),
				validation.By(utils.TrimWhiteSpace),
			),
			validation.Field(&dto.MerchantName,
				validation.Required.Error("merchant name is required"),
				validation.By(utils.TrimWhiteSpace),
				validation.By(utils.NoSpecialChars),
			),
			validation.Field(&dto.MerchantRepresentativeName,
				validation.Required.Error("representative name is required"),
				validation.By(utils.TrimWhiteSpace),
				validation.By(utils.NoSpecialChars),
			),
			validation.Field(&dto.PhoneNumber,
				validation.Required.Error("phone number is required"),
				is.Digit.Error("phone number must contain only digits"),
				validation.Length(9, 15).Error("phone number must be between 9 and 15 digits"),
				validation.Match(regexp.MustCompile(`^(?:\+251|251|0)9\d{8}$`)).Error("invalid phone number format"),
				validation.By(utils.TrimWhiteSpace),
				validation.By(utils.NoSpecialChars),
			),
			validation.Field(&dto.Email,
				validation.Required.Error("email is required"),
				is.Email.Error("email must be a valid email address"),
				validation.By(utils.TrimWhiteSpace),
			),
			validation.Field(&dto.AccountNumber,
				validation.Required.Error("account number is required"),
				validation.By(utils.TrimWhiteSpace),
				validation.By(utils.NoSpecialChars),
			),
		}
	} else {
		if strings.TrimSpace(dto.Type) != "" {
			rules = append(rules, validation.Field(&dto.Type,
				validation.Required.Error("type is required"),
				validation.By(utils.NoSpecialChars)),
			)
		}
		if strings.TrimSpace(dto.MerchantName) != "" {
			rules = append(rules, validation.Field(&dto.MerchantName, validation.By(utils.NoSpecialChars)))
		}
		if strings.TrimSpace(dto.MerchantRepresentativeName) != "" {
			rules = append(rules, validation.Field(&dto.MerchantRepresentativeName, validation.By(utils.NoSpecialChars)))
		}
		if strings.TrimSpace(dto.PhoneNumber) != "" {
			rules = append(rules, validation.Field(&dto.PhoneNumber,
				is.Digit.Error("phone number must contain only digits"),
				validation.Length(9, 15).Error("phone number must be between 9 and 15 digits"),
				validation.Match(regexp.MustCompile(`^(?:\+251|251|0)9\d{8}$`)).Error("invalid phone number format"),
				validation.By(utils.NoSpecialChars),
			))
		}
		if strings.TrimSpace(dto.Email) != "" {
			rules = append(rules, validation.Field(&dto.Email, is.Email.Error("email must be a valid email address"), validation.By(utils.NoSpecialChars)))
		}
		if strings.TrimSpace(dto.AccountNumber) != "" {
			rules = append(rules, validation.Field(&dto.AccountNumber, validation.By(utils.NoSpecialChars)))
		}
	}

	if len(rules) > 0 {
		if err := validation.ValidateStruct(&dto, rules...); err != nil {
			return err
		}
	}

	return nil
}
