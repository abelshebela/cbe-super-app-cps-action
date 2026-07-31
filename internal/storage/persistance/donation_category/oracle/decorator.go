package donation_category_oracle

import (
	"context"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/kafka"
)

// publishingRepository wraps a bare repo and publishes Kafka events on writes.
type publishingRepository struct {
	storage.DonationCategoryRepository
	producer kafka.ClientOrchestrationProducer
}

func NewPublishingRepository(
	delegate storage.DonationCategoryRepository,
	producer kafka.ClientOrchestrationProducer,
) storage.DonationCategoryRepository {
	return &publishingRepository{
		DonationCategoryRepository: delegate,
		producer:                   producer,
	}
}

func (d *publishingRepository) Create(ctx context.Context, c *imodel.DonationCategoryOracle) error {
	if err := d.DonationCategoryRepository.Create(ctx, c); err != nil {
		return err
	}
	topic := string(constants.ClientOrchestrationDonationCategoryTopic)
	d.producer.PublishMessage(ctx, c, topic, topic, "new donation category created")
	return nil
}

func (d *publishingRepository) Update(ctx context.Context, id string, c *imodel.DonationCategoryOracle) error {
	if err := d.DonationCategoryRepository.Update(ctx, id, c); err != nil {
		return err
	}
	topic := string(constants.ClientOrchestrationDonationCategoryTopic)
	d.producer.PublishMessage(ctx, c, topic, topic, "donation category updated")
	return nil
}

func (d *publishingRepository) EnableDisable(ctx context.Context, id string, enable bool) error {
	if err := d.DonationCategoryRepository.EnableDisable(ctx, id, enable); err != nil {
		return err
	}
	topic := string(constants.ClientOrchestrationDonationCategoryTopic)
	payload := map[string]any{"id": id, "enabled": enable}
	d.producer.PublishMessage(ctx, payload, topic, topic, "donation category enable/disable updated")
	return nil
}
