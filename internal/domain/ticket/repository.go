package ticket

import (
	"context"

	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
)

type Repository interface {
	CreateTicketAction(ctx context.Context, action domain.CPSAction) (domain.CPSAction, error)
	CreateTicket(ctx context.Context, ticket Ticket) error
	GetTicketByActionId(ctx context.Context, action_id string) (domain.CPSAction, error)
	UpdateTicketCpsAction(ctx context.Context, action domain.CPSAction) error
	FetchAllTickets(ctx context.Context, limit, offset int32) ([]Ticket, error)
}
