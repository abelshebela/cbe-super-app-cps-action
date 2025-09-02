package service_details

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"encoding/json"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ValidateTotalCapAgainstServices(ctx context.Context, serviceRepo storage.ServiceDetailsRepository, logger utils.Logger, newTotalCap uint64) error {
	logger.Infof("Validating new total cap %d against all service caps", newTotalCap)

	// Get all services with their caps
	projection := bson.M{
		"service_name": 1,
		"service_code": 1,
		"cap":          1,
	}

	services, err := serviceRepo.FindAllWithPagination(ctx, projection, types.Filter{Page: 1, PerPage: 1000})
	if err != nil {
		logger.Errorf("error fetching services for total cap validation: %v", err)
		return err
	}

	for _, service := range services.Data {
		if service.Cap.ISingleCap > newTotalCap ||
			service.Cap.IDailyCap > newTotalCap ||
			service.Cap.CorporateSingleCap > newTotalCap ||
			service.Cap.CorporateDailyCap > newTotalCap {

			logger.Warnf(
				"Service %s (ID: %s) has a cap exceeding new total cap %d",
				service.ServiceName, service.ID.Hex(), newTotalCap,
			)
			return errors.New(localization.ErrorTotalMaxTransferCannotBeLessExistTransfers.Code)
		}
	}

	logger.Infof("All service caps are within the new total cap %d", newTotalCap)
	return nil
}

func MapServiceDetailsForUpdate(existingService *model.ServiceDetails, newData interface{}) (*model.ServiceDetails, error) {
	updatedService := *existingService
	if err := BindAction(newData, &updatedService); err != nil {
		return nil, err
	}

	// Ensure the ID is preserved
	updatedService.ID = existingService.ID

	return &updatedService, nil
}

func MapHQForTotalCapUpdate(existingHQ *model.HQ, newTotalCap uint64) *model.HQ {
	updatedHQ := *existingHQ
	updatedHQ.TotalCap = newTotalCap
	return &updatedHQ
}

func BindAction(source any, target any) error {
	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}
