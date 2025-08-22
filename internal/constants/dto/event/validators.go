package eventdto

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (e EventRequest) IsEmpty() bool {
	return e.MerchantID == "" &&
		e.EventName == "" &&
		e.EventVenue == "" &&
		e.EventCity == "" &&
		e.StartDate.IsZero() &&
		e.DueDate.IsZero() &&
		e.CoverImage == nil &&
		e.EventDescription == "" &&
		e.TotalTicketCount == 0 &&
		(len(e.Tickets) == 0 || e.Tickets == nil)
}

func (e EventRequest) Validate(isCreate bool) error {
	if !isCreate && e.IsEmpty() {
		return nil
	}

	var rules []*validation.FieldRules

	if isCreate {
		rules = []*validation.FieldRules{
			validation.Field(&e.MerchantID,
				validation.Required.Error(localization.ErrorMerchantIDRequired.Code)),
			validation.Field(&e.EventName,
				validation.Required.Error(localization.ErrorEventNameRequired.Code)),
			validation.Field(&e.EventVenue,
				validation.Required.Error(localization.ErrorEventVenueRequired.Code)),
			validation.Field(&e.EventCity,
				validation.Required.Error(localization.ErrorEventCityRequired.Code)),
			validation.Field(&e.StartDate,
				validation.Required.Error(localization.ErrorStartDateRequired.Code),
				validation.By(validateStartDate(e))),
			validation.Field(&e.DueDate,
				validation.Required.Error(localization.ErrorDueDateRequired.Code),
				validation.By(validateDueDate(e))),
			validation.Field(&e.CoverImage,
				validation.Required.Error(localization.ErrorCoverImageRequired.Code),
				validation.By(validateCoverImage(e))),
			validation.Field(&e.EventDescription,
				validation.Required.Error(localization.ErrorEventDescriptionRequired.Code)),
			validation.Field(&e.TotalTicketCount,
				validation.Required.Error(localization.ErrorTotalTicketCountRequired.Code),
				validation.Min(uint(1)).Error(localization.ErrorInvalidTicketCount.Code)),
			validation.Field(&e.Tickets,
				validation.Required.Error(localization.ErrorTicketsRequired.Code),
				validation.By(validateTickets(e, isCreate))),
		}
	} else {
		if e.MerchantID != "" {
			rules = append(rules, validation.Field(&e.MerchantID, validation.Required.Error("merchant_id is required")))
		}
		if e.EventName != "" {
			rules = append(rules, validation.Field(&e.EventName, validation.Required.Error("event_name is required")))
		}
		if e.EventVenue != "" {
			rules = append(rules, validation.Field(&e.EventVenue, validation.Required.Error("event_venue is required")))
		}
		if e.EventCity != "" {
			rules = append(rules, validation.Field(&e.EventCity, validation.Required.Error("event_city is required")))
		}
		if !e.StartDate.IsZero() {
			rules = append(rules, validation.Field(&e.StartDate, validation.By(validateStartDate(e))))
		}
		if !e.DueDate.IsZero() {
			rules = append(rules, validation.Field(&e.DueDate, validation.By(validateDueDate(e))))
		}
		if e.CoverImage != nil {
			rules = append(rules, validation.Field(&e.CoverImage, validation.By(validateCoverImage(e))))
		}
		if e.EventDescription != "" {
			rules = append(rules, validation.Field(&e.EventDescription, validation.Required.Error("event_description is required")))
		}
		if e.TotalTicketCount != 0 {
			rules = append(rules, validation.Field(&e.TotalTicketCount,
				validation.Min(uint(1)).Error("total_ticket_count must be at least 1")))
		}
		if len(e.Tickets) > 0 {
			rules = append(rules, validation.Field(&e.Tickets, validation.By(validateTickets(e, isCreate))))
		}
	}

	// Run validations if any rules are accumulated
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
			return fmt.Errorf("start_date can't be before today")
		}
		return nil
	}
}

func validateDueDate(e EventRequest) validation.RuleFunc {
	return func(value interface{}) error {
		if !e.StartDate.IsZero() && !e.DueDate.IsZero() && e.StartDate.After(e.DueDate) {
			return fmt.Errorf("due_date can't be before start_date")
		}
		return nil
	}
}

func validateCoverImage(e EventRequest) validation.RuleFunc {
	return func(value interface{}) error {
		file, ok := value.(*multipart.FileHeader)
		if !ok || file == nil {
			return nil
		}
		if file.Size > (2 << 20) {
			return errors.New("FILE_TOO_LARGE")
		}
		if file.Header != nil {
			ct := file.Header.Get("Content-Type")
			if ct == "" || (ct != "image/jpeg" && ct != "image/png" && ct != "image/gif" && ct != "image/webp") {
				return fmt.Errorf("cover_image must be an image (jpeg, png, gif, webp)")
			}
		}
		return nil
	}
}

func validateTickets(e EventRequest, isCreate bool) validation.RuleFunc {
	return func(value interface{}) error {
		if (len(e.Tickets) == 0) && isCreate {
			return fmt.Errorf("at least one ticket is required")
		}
		var total uint
		for _, t := range e.Tickets {
			total += uint(t.NumberOfTicker)
		}
		if e.TotalTicketCount != 0 && total != e.TotalTicketCount {
			return fmt.Errorf("sum of all ticket counts must match total_ticket_count")
		}
		return nil
	}
}
