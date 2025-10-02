package eventdto

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"mime/multipart"
	"time"
)

type EventRequest struct {
	AccountNumber       string    `json:"account_number"`
	MerchantID          string    `json:"merchant_id"`
	MercahntName        string    `json:"merchant_name"`
	MerchantPhoneNumber string    `json:"merchant_phone_number"`
	MerchantEmail       string    `json:"merchant_email"`
	EventName           string    `json:"event_name"`
	EventVenue          string    `json:"event_venue"`
	EventCity           string    `json:"event_city"`
	StartDate           time.Time `json:"start_date"`
	DueDate             time.Time `json:"due_date"`
	CoverImage       *multipart.FileHeader `json:"cover_image" swaggertype:"string" format:"binary"`
	EventDescription string                `json:"event_description"`
	TotalTicketCount uint                  `json:"total_ticket_count"`
	Tickets          []types.Ticket        `json:"tickets" swaggertype:"array,object"`
}
