package middleware

import (
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
)

// GetCPSActionApproveRepo exposes the configured CPSActionApproveIndexRepository so
// handlers can validate role-based checker_index for sequential approvals.
func GetCPSActionApproveRepo() storage.CPSActionApproveIndexRepository {
	return cpsApproveRepo
}

func GetBPSActionApproveRepo() storage.BPSActionApproveIndexRepository {
	return bpsApproveRepo
}

func GetClientOrchestrationProducer() *kafka.ClientOrchestrationProducer {
	return clientOrchestrationProducer
}
