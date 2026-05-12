package donation_company_oracle

import (
	"context"

	"cbe-super-app-cps-action/internal/constants"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
)

type publishingRepository struct {
	storage.DonationCompanyRepository
	producer kafka.ClientOrchestrationProducer
}

func NewPublishingRepository(
	delegate storage.DonationCompanyRepository,
	producer kafka.ClientOrchestrationProducer,
) storage.DonationCompanyRepository {
	return &publishingRepository{
		DonationCompanyRepository: delegate,
		producer:                  producer,
	}
}

func (d *publishingRepository) Create(ctx context.Context, c *imodel.DonationCompanyOracle) error {
	if err := d.DonationCompanyRepository.Create(ctx, c); err != nil {
		return err
	}
	topic := string(constants.ClientOrchestrationDonationCompanyTopic)
	d.producer.PublishMessage(ctx, c, topic, topic, "new donation company created")
	return nil
}

func (d *publishingRepository) Update(ctx context.Context, id string, c *imodel.DonationCompanyOracle) error {
	if err := d.DonationCompanyRepository.Update(ctx, id, c); err != nil {
		return err
	}
	topic := string(constants.ClientOrchestrationDonationCompanyTopic)
	d.producer.PublishMessage(ctx, c, topic, topic, "donation company updated")
	return nil
}
