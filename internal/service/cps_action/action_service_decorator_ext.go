package cpsaction

import (
	"context"
)

// ReverseCPSAction delegates to the base implementation to satisfy service.CPSActionService
func (s *cpsActionServiceWithRoles) ReverseCPSAction(ctx context.Context, actionCode string) error {
	return s.base.ReverseCPSAction(ctx, actionCode)
}
