package core

import (
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
)

func ChangeTopicName(cfg *config.VaultConfig) string {
	var topic string
	switch cfg.GoEnv {
	case "qa":
		topic = "customer_segmentation_qa"
	case "uat":
		topic = "customer_segmentation_uat"
	case "dev":
		topic = "customer_segmentation_dev"
	case "production":
		topic = "customer_segmentation_production"
	case "staging":
		topic = "customer_segmentation_staging"
	}
	return topic
}

func GetAllBranches(ctx context.Context, segmentedID string, repo storage.AccountBlockRepository) []string {
	// Implementation for getting all branches
	branches, err := repo.GetAllBranches(ctx, segmentedID)
	if err != nil {
		// Handle error appropriately
		return []string{}
	}
	// Extract branch IDs from the retrieved branches
	var branchIDs []string
	for _, branch := range branches {
		branchIDs = append(branchIDs, branch.Code)
	}
	return branchIDs
}
