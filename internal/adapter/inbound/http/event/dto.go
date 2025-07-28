package event

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	evententity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type EventRequest struct {
	MerchantID       string                `json:"merchant_id"`
	EventName        string                `json:"event_name"`
	EventVenue       string                `json:"event_venue"`
	EventCity        string                `json:"event_city"`
	StartDate        time.Time             `json:"start_date"`
	DueDate          time.Time             `json:"due_date"`
	CoverImage       *multipart.FileHeader `json:"cover_image"`
	EventDescription string                `json:"event_description"`
	TotalTicketCount uint                  `json:"total_ticket_count"`
	Tickets          []evententity.Ticket  `json:"tickets"`
}

// validateRequiredString creates a required validation rule for a string field
func validateRequiredString(field *string, name string) *validation.FieldRules {
	return validation.Field(field, validation.Required.Error(fmt.Sprintf("%s is required", name)))
}

// validateStartDate ensures start_date is not before today
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

// validateDueDate ensures due_date is not before start_date
func validateDueDate(e EventRequest) validation.RuleFunc {
	return func(value interface{}) error {
		if !e.StartDate.IsZero() && !e.DueDate.IsZero() && e.StartDate.After(e.DueDate) {
			return fmt.Errorf("due_date can't be before start_date")
		}
		return nil
	}
}

// validateCoverImage ensures cover_image is a valid image file
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

// validateTickets ensures tickets are provided and match total_ticket_count
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

func (e EventRequest) Validate(isCreate bool) error {
	// For updates, return no error if the request is completely empty
	if !isCreate && e.IsEmpty() {
		return nil
	}

	var rules []*validation.FieldRules

	if isCreate {
		rules = []*validation.FieldRules{
			validation.Field(&e.MerchantID, validation.Required.Error("merchant_id is required")),
			validation.Field(&e.EventName, validation.Required.Error("event_name is required")),
			validation.Field(&e.EventVenue, validation.Required.Error("event_venue is required")),
			validation.Field(&e.EventCity, validation.Required.Error("event_city is required")),
			validation.Field(&e.StartDate, validation.Required.Error("start_date is required"), validation.By(validateStartDate(e))),
			validation.Field(&e.DueDate, validation.Required.Error("due_date is required"), validation.By(validateDueDate(e))),
			validation.Field(&e.CoverImage, validation.Required.Error("cover_image is required"), validation.By(validateCoverImage(e))),
			validation.Field(&e.EventDescription, validation.Required.Error("event_description is required")),
			validation.Field(&e.TotalTicketCount,
				validation.Required.Error("total_ticket_count is required"),
				validation.Min(uint(1)).Error("total_ticket_count must be at least 1"),
			),
			validation.Field(&e.Tickets, validation.Required.Error("tickets are required"), validation.By(validateTickets(e, isCreate))),
		}
	} else {
		// For update: only validate fields that are present (non-zero / non-empty)
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
		if e.Tickets != nil && len(e.Tickets) > 0 {
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

func ParseEventRequestFromMultipartForm(r *http.Request, isCreate bool) (EventRequest, error) {
	var req EventRequest

	// Parse cover image
	_, fileHeader, err := utils.ParseMultipartFormFile(r, "cover_image", 2<<20)
	if err != nil && (err.Error() != common_util.ErrMissingFile && !isCreate) {

		return req, err
	}
	req.CoverImage = fileHeader

	// Parse simple fields
	req.MerchantID = r.FormValue("merchant_id")
	req.EventName = r.FormValue("event_name")
	req.EventVenue = r.FormValue("event_venue")
	req.EventCity = r.FormValue("event_city")
	req.EventDescription = r.FormValue("event_description")

	// Parse start_date and due_date
	startDateStr := r.FormValue("start_date")
	dueDateStr := r.FormValue("due_date")
	if startDateStr != "" {
		t, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return req, err
		}
		req.StartDate = t
	}
	if dueDateStr != "" {
		t, err := time.Parse(time.RFC3339, dueDateStr)
		if err != nil {
			return req, err
		}
		req.DueDate = t
	}

	// Parse total_ticket_count
	totalTicketCountStr := r.FormValue("total_ticket_count")
	if totalTicketCountStr != "" {
		n, err := strconv.ParseUint(totalTicketCountStr, 10, 32)
		if err != nil {
			return req, err
		}
		req.TotalTicketCount = uint(n)
	}

	// Parse tickets
	tickets, err := parseTicketsFromForm(r.Form)
	if err != nil {
		return req, err
	}
	req.Tickets = tickets

	return req, nil
}

// EventResponse, TicketTypeResponse, TicketSaleResponse types unchanged
type EventResponse struct {
	EventID          string               `json:"event_id"`
	EventCode        string               `json:"event_code"`
	CoverImage       string               `json:"cover_image"`
	EventName        string               `json:"event_name"`
	EventDescription string               `json:"event_description"`
	TotalTicketCount int                  `json:"total_ticket_count"`
	EventStartDate   time.Time            `json:"event_startdate"`
	EventEndDate     time.Time            `json:"event_enddate"`
	EventStatus      string               `json:"event_status"`
	EventVenue       string               `json:"event_venue"`
	EventCity        string               `json:"event_city"`
	MerchantName     string               `json:"merchant_name"`
	MerchantPhone    string               `json:"merchant_phone"`
	TicketTypes      []TicketTypeResponse `json:"ticket_types"`
}

type TicketTypeResponse struct {
	TicketName  string `json:"ticket_name"`
	TicketPrice int    `json:"ticket_price"`
	TicketCount int    `json:"ticket_count"`
}
