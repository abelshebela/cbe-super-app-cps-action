package walletDto

import "time"

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
	Enabled          bool                 `json:"enabled"`
	MerchantName     string               `json:"merchant_name"`
	MerchantPhone    string               `json:"merchant_phone"`
	TicketTypes      []TicketTypeResponse `json:"ticket_types"`
}

type TicketTypeResponse struct {
	TicketName  string `json:"ticket_name"`
	TicketPrice int    `json:"ticket_price"`
	TicketCount int    `json:"ticket_count"`
}
