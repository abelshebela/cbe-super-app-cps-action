package donation_oracle

import (
	"context"

	"cbe-super-app-cps-action/internal/constants"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
)

type publishingRepository struct {
	storage.DonationRepository
	producer kafka.ClientOrchestrationProducer
}

func NewPublishingRepository(
	delegate storage.DonationRepository,
	producer kafka.ClientOrchestrationProducer,
) storage.DonationRepository {
	return &publishingRepository{
		DonationRepository: delegate,
		producer:           producer,
	}
}

func (d *publishingRepository) Create(ctx context.Context, donation *imodel.DonationOracle) error {
	if err := d.DonationRepository.Create(ctx, donation); err != nil {
		return err
	}
	topic := string(constants.ClientOrchestrationDonationTopic)
	d.producer.PublishMessage(ctx, donation, topic, topic, "new donation created")
	return nil
}

func (d *publishingRepository) Update(ctx context.Context, id string, donation *imodel.DonationOracle) error {
	if err := d.DonationRepository.Update(ctx, id, donation); err != nil {
		return err
	}
	topic := string(constants.ClientOrchestrationDonationTopic)
	d.producer.PublishMessage(ctx, donation, topic, topic, "donation updated")
	return nil
}
