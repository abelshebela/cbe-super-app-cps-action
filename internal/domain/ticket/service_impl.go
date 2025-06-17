package ticket

import "context"

type TicketStore struct {
	repository Repository
}
type TicketService interface {
	CreateTicketAction(ctx context.Context, ticket Ticket, makerId string) (string, error)
	CheckTicket(ctx context.Context, actionId string, action bool, checkerId string) error
	FetchAllTickets(ctx context.Context, limit, offset int32) ([]Ticket, error)
}

func NewService(repository Repository) TicketService {
	return &TicketStore{
		repository: repository,
	}
}
