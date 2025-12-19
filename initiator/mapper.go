package initiator

import (
	"cbe-super-app-cps-action/internal/service"
)

// ServiceContainerToServiceLayer maps a ServiceContainer to a ServiceLayer
// Only maps the services that exist in both structures
func ServiceContainerToServiceLayer(container *service.ServiceContainer) *service.ServiceLayer {
	if container == nil {
		return nil
	}

	return &service.ServiceLayer{
		EventService: container.EventContainer,
		CPSAction:    container.CPSActionContainer,
		Feedback:     container.FeedbackContainer,
		Unlink:       container.UnlinkContainer,
		BpsUser:      container.BPSUserContainer,
		Advert:       container.AdContainer,
		PortalCard:   container.PortalCardContainer,
	}
}

// ServiceLayerToServiceContainer maps a ServiceLayer to a ServiceContainer
// Only maps the services that exist in ServiceLayer, other fields will be nil
func ServiceLayerToServiceContainer(layer *service.ServiceLayer) *service.ServiceContainer {
	if layer == nil {
		return nil
	}

	return &service.ServiceContainer{
		EventContainer:      layer.EventService,
		CPSActionContainer:  layer.CPSAction,
		FeedbackContainer:   layer.Feedback,
		UnlinkContainer:     layer.Unlink,
		BPSUserContainer:    layer.BpsUser,
		AdContainer:         layer.Advert,
		PortalCardContainer: layer.PortalCard,
		// All other fields will be nil/zero values
	}
}

// MergeServiceLayerIntoContainer merges ServiceLayer services into an existing ServiceContainer
// This preserves existing services in the container while updating with ServiceLayer services
func MergeServiceLayerIntoContainer(container *service.ServiceContainer, layer *service.ServiceLayer) {
	if container == nil || layer == nil {
		return
	}

	if layer.EventService != nil {
		container.EventContainer = layer.EventService
	}
	if layer.CPSAction != nil {
		container.CPSActionContainer = layer.CPSAction
	}
	if layer.Feedback != nil {
		container.FeedbackContainer = layer.Feedback
	}
	if layer.Unlink != nil {
		container.UnlinkContainer = layer.Unlink
	}
	if layer.BpsUser != nil {
		container.BPSUserContainer = layer.BpsUser
	}
	if layer.Advert != nil {
		container.AdContainer = layer.Advert
	}
	if layer.PortalCard != nil {
		container.PortalCardContainer = layer.PortalCard
	}
}

// ExtractServicesFromContainer extracts specific services from ServiceContainer
// Returns individual services that can be used to construct a ServiceLayer
func ExtractServicesFromContainer(container *service.ServiceContainer) (
	eventService service.EventService,
	cpsActionService service.CPSActionService,
	feedbackService service.FeedbackService,
	unlinkService service.UnlinkService,
	bpsUserService service.BPSUserService,
	advertService service.AdvertService,
	portalCardService service.PortalCardService,
) {
	if container == nil {
		return nil, nil, nil, nil, nil, nil, nil
	}

	return container.EventContainer,
		container.CPSActionContainer,
		container.FeedbackContainer,
		container.UnlinkContainer,
		container.BPSUserContainer,
		container.AdContainer,
		container.PortalCardContainer
}
