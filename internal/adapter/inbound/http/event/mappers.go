package event

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"mime/multipart"

	eventdto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
	evententity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
)

func parseTicketsFromForm(form url.Values) ([]evententity.Ticket, error) {
	var tickets []evententity.Ticket
	for i := 0; ; i++ {
		name := form.Get(fmt.Sprintf("ticket_name[%d]", i))
		if name == "" {
			break
		}

		category := form.Get(fmt.Sprintf("ticket_category[%d]", i))
		tType := form.Get(fmt.Sprintf("ticket_type[%d]", i))
		priceStr := form.Get(fmt.Sprintf("ticket_price[%d]", i))
		numStr := form.Get(fmt.Sprintf("ticket_number_of_ticker[%d]", i))

		price, err := strconv.ParseUint(priceStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid price at index %d", i)
		}
		number, err := strconv.ParseUint(numStr, 10, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid number of ticker at index %d", i)
		}

		tickets = append(tickets, evententity.Ticket{
			Name:           name,
			Category:       category,
			Type:           tType,
			Price:          price,
			NumberOfTicker: uint8(number),
		})
	}
	return tickets, nil
}

// ToDomainEventRequest maps HTTP EventRequest DTO to domain EventRequest DTO
func ToDomainEventRequest(req struct {
	MerchantID       string
	EventName        string
	EventVenue       string
	EventCity        string
	StartDate        time.Time
	DueDate          time.Time
	CoverImage       *multipart.FileHeader
	EventDescription string
	TotalTicketCount uint
	Tickets          []evententity.Ticket
}) eventdto.EventRequest {
	return eventdto.EventRequest{
		MerchantID:       req.MerchantID,
		EventName:        req.EventName,
		EventVenue:       req.EventVenue,
		EventCity:        req.EventCity,
		StartDate:        req.StartDate,
		DueDate:          req.DueDate,
		CoverImage:       req.CoverImage,
		EventDescription: req.EventDescription,
		TotalTicketCount: req.TotalTicketCount,
		Tickets:          req.Tickets,
	}
}

// ToEventResponse maps domain Event entity to HTTP EventResponse
func ToEventResponse(e *evententity.Event) EventResponse {
	return EventResponse{
		EventID:          e.ID,
		EventCode:        e.EventCode,
		CoverImage:       e.EventInformation.Cover,
		EventName:        e.EventName,
		EventDescription: e.EventInformation.Description,
		TotalTicketCount: int(e.TicketInformation.TotalNumberOfTicket),
		EventStartDate:   e.CreatedAt,
		EventEndDate:     e.EventInformation.DueDate,
		EventStatus:      string(e.Status),
		EventVenue:       e.EventVenue,
		Enabled:          e.Enabled,
		EventCity:        e.EventCity,
		MerchantName:     e.MerchantInformation.MercahntName,
		MerchantPhone:    e.MerchantInformation.MerchantPhoneNumber,
		TicketTypes:      toTicketTypeResponses(e.Ticket),
	}
}

func toTicketTypeResponses(tickets []evententity.Ticket) []TicketTypeResponse {
	var res []TicketTypeResponse
	for _, t := range tickets {
		res = append(res, TicketTypeResponse{
			TicketName:  t.Name,
			TicketPrice: int(t.Price),
			TicketCount: int(t.NumberOfTicker),
		})
	}
	return res
}
