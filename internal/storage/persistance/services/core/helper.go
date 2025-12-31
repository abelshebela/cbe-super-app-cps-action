package core

import (
	"time"

	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/imodel"
	imodel "cbe-super-app-cps-action/internal/constants/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func MapToServiceUpdate(service imodel.Service, existing imodel.Service) bson.M {
	update := bson.M{
		"last_modified_at": time.Now(),
	}

	if service.ServiceCode != "" {
		update["service_code"] = service.ServiceCode
	}

	if service.ServiceKey != "" {
		update["service_key"] = service.ServiceKey
	}

	if service.ServiceName != "" {
		update["service_name"] = service.ServiceName
	}

	if service.ProductGlAccount != "" {
		update["product_gl_account"] = service.ProductGlAccount
	}

	if len(service.Tiers) > 0 {
		update["tiers"] = service.Tiers
	}

	for i, _ := range service.ServiceList {
		existing := existing.ServiceList[i]
		service.ServiceList[i].IsEnabled = existing.IsEnabled
	}
	update["service_list"] = service.ServiceList

	if len(service.ServiceList) > 0 {
		update["service_list"] = service.ServiceList
	}

	if service.Cap != (imodel.Cap{}) {
		update["cap"] = service.Cap
	}

	update["enabled"] = existing.Enabled
	update["is_deleted"] = existing.IsDeleted

	return update
}
