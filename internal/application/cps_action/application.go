package cpsaction

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type CPSActionApplication interface {
	ApproveCPSAction(ctx context.Context, action *entities.AuthorizeCPSAction) (*entities.CPSAction, error)
	RejectCPSAction(ctx context.Context, action *entities.AuthorizeCPSAction) (*entities.CPSAction, error)
	GetCPSActionsByDepartment(ctx context.Context, department string, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.CPSAction], error)
	GetCPSActionByID(ctx context.Context, id string) (*entities.CPSAction, error)
	GetCPSActionByActionCode(ctx context.Context, uniqueID string) (*entities.CPSAction, error)
}

type cpsActionApplication struct {
	service    service.CPSActionService
	services   application.Domain
	dispatcher Dispatcher
}

func NewCPSActionApplication(service service.CPSActionService, services application.Domain, dispatcher Dispatcher) CPSActionApplication {
	return &cpsActionApplication{
		service:    service,
		services:   services,
		dispatcher: dispatcher,
	}
}

func (a *cpsActionApplication) CPSActionExists(ctx context.Context, uniqueID string) (bool, error) {
	return a.service.CPSActionExists(ctx, uniqueID)
}

func PrettyPrintJSON(data interface{}) error {
	prettyJSON, err := json.MarshalIndent(data, "", "  ") // 2 spaces indent
	if err != nil {
		return err
	}

	fmt.Println(string(prettyJSON))
	return nil
}

func (a *cpsActionApplication) ApproveCPSAction(ctx context.Context, action *entities.AuthorizeCPSAction) (*entities.CPSAction, error) {
	cpsAction, err := a.service.ApproveCPSAction(ctx, action)
	if err != nil {
		return nil, err
	}

	PrettyPrintJSON(cpsAction)

	cpsRes, err := a.dispatcher.Authorize(ctx, cpsAction)
	if err != nil {
		return nil, err
	}

	return cpsRes, nil
}

func (a *cpsActionApplication) RejectCPSAction(ctx context.Context, action *entities.AuthorizeCPSAction) (*entities.CPSAction, error) {
	return a.service.RejectCPSAction(ctx, action)
}

func (a *cpsActionApplication) GetCPSActionsByDepartment(ctx context.Context, department string, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.CPSAction], error) {
	return a.service.GetCPSActionsByDepartment(ctx, department, filterParams)
}

func (a *cpsActionApplication) GetCPSActionByID(ctx context.Context, id string) (*entities.CPSAction, error) {
	return a.service.GetCPSActionByID(ctx, id)
}

func (a *cpsActionApplication) GetCPSActionByActionCode(ctx context.Context, uniqueID string) (*entities.CPSAction, error) {
	return a.service.GetCPSActionByActionCode(ctx, uniqueID)
}
