package ussd_merchant_dto

import (
	"errors"
	"strings"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/pkgs/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const settlementDirectTransfer = "Direct Transfer"

func (r CreateUssdMerchantRequest) Validate() error {
	errs := validation.Errors{}

	// Required fields on create
	if strings.TrimSpace(r.SettlementMethod) == "" {
		errs["settlement_method"] = validation.NewError("settlement_method", "settlement_method is required")
	}
	if strings.TrimSpace(r.Name) == "" {
		errs["name"] = validation.NewError("name", "name is required")
	}
	if strings.TrimSpace(r.PhoneNumber) == "" {
		errs["phone_number"] = validation.NewError("phone_number", "phone_number is required")
	}
	if strings.TrimSpace(r.Email) == "" {
		return validation.NewError("email", "email is required")
	}

	if strings.TrimSpace(r.Service) == "" {
		errs["service"] = validation.NewError("service", "service is required")
	}

	if r.Email != "" {
		if err := validation.Validate(r.Email, validation.By(func(value interface{}) error {
			s, _ := value.(string)
			if !strings.Contains(s, "@") || !strings.Contains(s, ".") {
				return validation.NewError("email", "invalid email format")
			}
			return nil
		})); err != nil {
			errs["email"] = err
		}
	}

	if err := validation.Validate(
		strings.ToUpper(strings.TrimSpace(r.SettlementMethod)),
		validation.In(
			string(constants.SettlementMethodDirect),
			string(constants.SettlementMethodGL),
			string(constants.SettlementMethodMultiAccount),
		),
	); err != nil {
		errs["settlement_method"] = validation.NewError(
			"settlement_method",
			"invalid settlement_method; must be one of DIRECT, GL, MULTI_ACCOUNT",
		)
	}

	// Logo required on create and must be a valid image (<= 15MB)
	if r.Logo == nil {
		errs["logo"] = validation.NewError("logo", "logo is required")
	} else {
		if !utils.IsValidImage(r.Logo) {
			errs["logo"] = errors.New("invalid logo file type")
		} else if r.Logo.Size > (15 << 20) {
			errs["logo"] = validation.NewError("logo", "file too large")
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (r UpdateUssdMerchantRequest) Validate() error {
	errs := validation.Errors{}

	// When settlement method or account number is provided, enforce the conditional
	if r.SettlementMethod != "" || r.AccountNumber != "" {

		if err := validation.Validate(
			strings.ToUpper(strings.TrimSpace(r.SettlementMethod)),
			validation.In(
				string(constants.SettlementMethodDirect),
				string(constants.SettlementMethodGL),
				string(constants.SettlementMethodMultiAccount),
			),
		); err != nil {
			errs["settlement_method"] = validation.NewError(
				"settlement_method",
				"invalid settlement_method; must be one of DIRECT, GL, MULTI_ACCOUNT",
			)
		}

	}

	// Optional email validation when provided
	if r.Email != "" {
		if err := validation.Validate(r.Email, validation.By(func(value interface{}) error {
			s, _ := value.(string)
			if !strings.Contains(s, "@") || !strings.Contains(s, ".") {
				return validation.NewError("email", "invalid email format")
			}
			return nil
		})); err != nil {
			errs["email"] = err
		}
	}

	// Logo optional on update but validate if provided
	if r.Logo != nil {
		if !utils.IsValidImage(r.Logo) {
			errs["logo"] = errors.New("invalid logo file type")
		} else if r.Logo.Size > (15 << 20) {
			errs["logo"] = validation.NewError("logo", "file too large")
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}
