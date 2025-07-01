package dto

import "time"

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

type TicketType struct {
	TicketName      string `json:"ticket_name"`
	NumberOfTickets int    `json:"number_of_tickets"`
	TicketPrice     string `json:"ticket_price"`
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

type EventFetchRequest struct {
	RequestId string `json:"request_id"`
}
