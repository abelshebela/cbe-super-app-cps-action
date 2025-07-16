package dto

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type EventCreateRequest struct {
	MerchantName    string       `json:"merchant_name"`
	EventName       string       `json:"event_name"`
	Venue           string       `json:"venue"`
	City            string       `json:"city"`
	Date            time.Time    `json:"date"`
	Description     string       `json:"description"`
	TotalNoOfTicket int64        `json:"total_number_of_ticket"`
	TicketTypes     []TicketType `json:"ticket_types"`
}

func (e EventCreateRequest) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.MerchantName, validation.Required.Error("merchant_name is required")),
		validation.Field(&e.EventName, validation.Required.Error("event_name is required")),
		validation.Field(&e.Venue, validation.Required.Error("venue is required")),
		validation.Field(&e.City, validation.Required.Error("city is required")),
		validation.Field(&e.Date, validation.Required.Error("date is required")),
		validation.Field(&e.Description, validation.Required.Error("description is required")),
		validation.Field(&e.TotalNoOfTicket, validation.Required.Error("total_number_of_ticket is required")),
		validation.Field(&e.TicketTypes, validation.Required.Error("ticket_types is required")),
	)
}

type TicketType struct {
	TicketName      string `json:"ticket_name"`
	NumberOfTickets int    `json:"number_of_tickets"`
	TicketPrice     string `json:"ticket_price"`
}

func (t TicketType) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.TicketName, validation.Required.Error("ticket_name is required")),
		validation.Field(&t.NumberOfTickets, validation.Required.Error("number_of_tickets is required")),
		validation.Field(&t.TicketPrice, validation.Required.Error("ticket_price is required")),
	)
}

type EventResponse struct {
	Docs          []EventDTO `json:"docs"`
	TotalDocs     int        `json:"totalDocs"`
	Limit         int        `json:"limit"`
	TotalPages    int        `json:"totalPages"`
	Page          int        `json:"page"`
	PagingCounter int        `json:"pagingCounter"`
	HasPrevPage   bool       `json:"hasPrevPage"`
	HasNextPage   bool       `json:"hasNextPage"`
	PrevPage      *int       `json:"prevPage"`
	NextPage      *int       `json:"nextPage"`
}

// EventDTO represents an individual event
type EventDTO struct {
	EventID          string          `json:"event_id"`
	EventCode        string          `json:"event_code"`
	CoverImage       string          `json:"cover_image"`
	EventName        string          `json:"event_name"`
	EventDescription string          `json:"event_description"`
	TotalTicketCount int             `json:"total_ticket_count"`
	EventStartDate   time.Time       `json:"event_startdate"`
	EventEndDate     time.Time       `json:"event_enddate"`
	EventStatus      string          `json:"event_status"`
	EventVenue       string          `json:"event_venue"`
	EventCity        string          `json:"event_city"`
	MerchantName     string          `json:"merchant_name"`
	MerchantPhone    string          `json:"merchant_phone"`
	TicketTypes      []TicketTypeDTO `json:"ticket_types"`
	TicketSales      []TicketSaleDTO `json:"ticket_sales"`
}

// TicketTypeDTO represents a ticket type for an event
type TicketTypeDTO struct {
	TicketName  string `json:"ticket_name"`
	TicketPrice int    `json:"ticket_price"`
	TicketCount int    `json:"ticket_count"`
}

// TicketSaleDTO represents ticket sales information
type TicketSaleDTO struct {
	TicketName    string `json:"ticket_name"`
	TicketRevenue int    `json:"ticket_revenue"`
	TicketSold    int    `json:"ticket_sold"`
	TicketCount   int    `json:"ticket_count"`
}

type EventCheckerRequest struct {
	Request_Id string `json:"request_id"`
	Action     bool   `json:"action"`
}

func (e EventCheckerRequest) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Request_Id, validation.Required.Error("request_id is required")),
	)
}

type EventFetchRequest struct {
	RequestId string `json:"request_id"`
}

func (e EventFetchRequest) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.RequestId, validation.Required.Error("request_id is required")),
	)
}
