package eventdto

import (
	"errors"
	"fmt"
	"mime/multipart"
	"regexp"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (e EventRequest) IsEmpty() bool {
	return strings.TrimSpace(e.MerchantID) == "" &&
		strings.TrimSpace(e.EventName) == "" &&
		strings.TrimSpace(e.EventVenue) == "" &&
		strings.TrimSpace(e.EventCity) == "" &&
		e.StartDate.IsZero() &&
		e.DueDate.IsZero() &&
		e.CoverImage == nil &&
		strings.TrimSpace(e.EventDescription) == "" &&
		e.TotalTicketCount == 0 &&
		(len(e.Tickets) == 0 || e.Tickets == nil)
}

func (e EventRequest) Validate(isCreate bool) error {
	if !isCreate && e.IsEmpty() {
		return nil
	}

	var rules []*validation.FieldRules

	validateString := func(fieldName, value string, required bool) validation.RuleFunc {
		return func(_ interface{}) error {
			if required && strings.TrimSpace(value) == "" {
				return fmt.Errorf("%s is required", fieldName)
			}
			if matched, _ := regexp.MatchString(`[<>$%]`, value); matched {
				return fmt.Errorf("%s contains invalid characters", fieldName)
			}
			return nil
		}
	}

	if isCreate {
		rules = []*validation.FieldRules{
			validation.Field(&e.MerchantID, validation.By(validateString("merchant_id", e.MerchantID, true))),
			validation.Field(&e.EventName, validation.By(validateString("event_name", e.EventName, true))),
			validation.Field(&e.EventVenue, validation.By(validateString("event_venue", e.EventVenue, true))),
			validation.Field(&e.EventCity, validation.By(validateString("event_city", e.EventCity, true))),
			validation.Field(&e.EventDescription, validation.By(validateString("event_description", e.EventDescription, true))),
			validation.Field(&e.StartDate, validation.Required.Error(localization.ErrorStartDateRequired.Code), validation.By(validateStartDate(e))),
			validation.Field(&e.DueDate, validation.Required.Error(localization.ErrorDueDateRequired.Code), validation.By(validateDueDate(e))),
			validation.Field(&e.CoverImage, validation.Required.Error(localization.ErrorCoverImageRequired.Code), validation.By(validateImage)),
			validation.Field(&e.TotalTicketCount, validation.Required.Error(localization.ErrorTotalTicketCountRequired.Code), validation.Min(uint(1)).Error(localization.ErrorInvalidTicketCount.Code)),
			validation.Field(&e.Tickets, validation.Required.Error(localization.ErrorTicketsRequired.Code), validation.By(validateTickets(e, isCreate))),
		}
	} else {
		if strings.TrimSpace(e.MerchantID) != "" {
			rules = append(rules, validation.Field(&e.MerchantID, validation.By(validateString("merchant_id", e.MerchantID, true))))
		}
		if strings.TrimSpace(e.EventName) != "" {
			rules = append(rules, validation.Field(&e.EventName, validation.By(validateString("event_name", e.EventName, true))))
		}
		if strings.TrimSpace(e.EventVenue) != "" {
			rules = append(rules, validation.Field(&e.EventVenue, validation.By(validateString("event_venue", e.EventVenue, true))))
		}
		if strings.TrimSpace(e.EventCity) != "" {
			rules = append(rules, validation.Field(&e.EventCity, validation.By(validateString("event_city", e.EventCity, true))))
		}
		if !e.StartDate.IsZero() {
			rules = append(rules, validation.Field(&e.StartDate, validation.By(validateStartDate(e))))
		}
		if !e.DueDate.IsZero() {
			rules = append(rules, validation.Field(&e.DueDate, validation.By(validateDueDate(e))))
		}
		if e.CoverImage != nil {
			rules = append(rules, validation.Field(&e.CoverImage, validation.By(validateImage)))
		}
		if strings.TrimSpace(e.EventDescription) != "" {
			rules = append(rules, validation.Field(&e.EventDescription, validation.By(validateString("event_description", e.EventDescription, true))))
		}
		if e.TotalTicketCount != 0 {
			rules = append(rules, validation.Field(&e.TotalTicketCount, validation.Min(uint(1)).Error("total_ticket_count must be >=1")))
		}
		if len(e.Tickets) > 0 {
			rules = append(rules, validation.Field(&e.Tickets, validation.By(validateTickets(e, isCreate))))
		}
	}

	if len(rules) > 0 {
		if err := validation.ValidateStruct(&e, rules...); err != nil {
			return err
		}
	}
	return nil
}

func validateStartDate(e EventRequest) validation.RuleFunc {
	return func(value interface{}) error {
		t, ok := value.(time.Time)
		if !ok || t.IsZero() {
			return nil
		}
		if t.Before(time.Now().Truncate(24 * time.Hour)) {
			return errors.New("start_date cannot be before today")
		}
		return nil
	}
}

func validateDueDate(e EventRequest) validation.RuleFunc {
	return func(value interface{}) error {
		if !e.StartDate.IsZero() && !e.DueDate.IsZero() && e.StartDate.After(e.DueDate) {
			return errors.New("due_date cannot be before start_date")
		}
		return nil
	}
}

func validateImage(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok || file == nil {
		return localization.ErrorMissingOrInvalidImage
	}
	if !utils.IsValidImage(file) {
		return localization.ErrorMissingOrInvalidImage
	}

	if file.Size > (15 << 20) {
		return validation.NewError("logo", localization.MsgFileTooLarge)
	}

	return nil
}

func validateTickets(e EventRequest, isCreate bool) validation.RuleFunc {
	return func(value interface{}) error {
		if len(e.Tickets) == 0 && isCreate {
			return errors.New("at least one ticket is required")
		}
		var total uint
		for _, t := range e.Tickets {
			total += uint(t.NumberOfTicker)
		}
		if e.TotalTicketCount != 0 && total != e.TotalTicketCount {
			return errors.New("sum of ticket numbers must match total_ticket_count")
		}
		return nil
	}
}
