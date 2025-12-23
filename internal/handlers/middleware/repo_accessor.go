package middleware

import "cbe-super-app-cps-action/internal/storage"

// GetCPSActionApproveRepo exposes the configured CPSActionApproveIndexRepository so
// handlers can validate role-based checker_index for sequential approvals.
func GetCPSActionApproveRepo() storage.CPSActionApproveIndexRepository {
	return cpsApproveRepo
}
