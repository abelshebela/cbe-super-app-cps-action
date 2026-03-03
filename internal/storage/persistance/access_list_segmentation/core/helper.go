package core

import "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"

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
