package account_block

import (
	"context"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"go.mongodb.org/mongo-driver/v2/bson"

	ctx_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	constant_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type AccountService struct {
	repo AccountBlockRepo
}

func NewAccountService(repo AccountBlockRepo) ApplicationServices {
	return &AccountService{repo: repo}
}

func (s *AccountService) GetBranch(ctx context.Context, branchCode string, filterParams *constant.Filter) (*model.Branch, error) {
	return s.repo.GetBranch(ctx, branchCode, filterParams)
}

func (s *AccountService) GetAllBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error) {
	return s.repo.GetAllBranches(ctx, region, district, filterParams)
}

func (s *AccountService) GetRegionByCode(ctx context.Context, regionCode string) (action.Region, error) {
	if regionCode == "" {
		return action.Region{}, fmt.Errorf("regionCode is required")
	}
	return s.repo.GetRegionByCode(ctx, regionCode)
}

func (s *AccountService) GetDistrictByCode(ctx context.Context, districtCode string) (action.District, error) {
	if districtCode == "" {
		return action.District{}, fmt.Errorf("districtCode is required")
	}
	return s.repo.GetDistrictByCode(ctx, districtCode)
}

func (s *AccountService) GetCityByCode(ctx context.Context, cityCode string) (action.City, error) {
	if cityCode == "" {
		return action.City{}, fmt.Errorf("cityCode is required")
	}
	return s.repo.GetCityByCode(ctx, cityCode)
}

func (s *AccountService) GetAllCities(ctx context.Context, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.City], error) {
	if filterParams == nil || filterParams.Page < 1 || filterParams.PerPage < 1 {
		return nil, fmt.Errorf("invalid pagination parameters")
	}
	return s.repo.GetAllCities(ctx, filterParams)
}

func (s *AccountService) GetAllDistricts(ctx context.Context, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.District], error) {
	if filterParams == nil || filterParams.Page < 1 || filterParams.PerPage < 1 {
		return nil, fmt.Errorf("invalid pagination parameters")
	}
	return s.repo.GetAllDistricts(ctx, filterParams)
}

func (s *AccountService) GetAllRegions(ctx context.Context, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Region], error) {
	if filterParams == nil || filterParams.Page < 1 || filterParams.PerPage < 1 {
		return nil, fmt.Errorf("invalid pagination parameters")
	}
	return s.repo.GetAllRegions(ctx, filterParams)
}

func generateCPSAction(ctx context.Context, codeType string, enabled bool, codes []string, actionType model.ActionType, requestActionType model.RequestAction) model.CPSAction {
	userPayload := ctx_utils.ExtractContext(ctx)

	var actions []interface{}
	for _, code := range codes {
		switch codeType {
		case "BRANCH":
			actions = append(actions, model.Branch{
				BranchCode: code,
				Enabled:    enabled,
			})
		case "REGION":
			actions = append(actions, model.Region{
				RegionCode: code,
				Enabled:    enabled,
			})
		case "DISTRICT":
			actions = append(actions, model.District{
				DistrictCode: code,
				Enabled:      enabled,
			})
		case "CITY":
			actions = append(actions, model.City{
				CityCode: code,
				Enabled:  enabled,
			})
		}
	}

	cpsAction := model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       utils.RandomGenerator(24),
		MakerID:          userPayload.UserID,
		MakerName:        userPayload.FullName,
		MakerPhoneNumber: userPayload.PhoneNumber,
		Department:       userPayload.Department,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(actionType),
		RequestAction:    string(requestActionType),
		PreviousAction:   nil,
		CurrentAction:    actions,
		CreatedAt:        time.Now(),
		MakerActionTime:  time.Now(),
	}

	return cpsAction
}

func (s *AccountService) EnableBranches(ctx context.Context, branchCodes []string, enabled bool) error {
	cpsAction := generateCPSAction(ctx, "BRANCH", enabled, branchCodes, model.ActionEnable, model.RequestEnableBranches)
	return s.repo.EnableOrDisable(ctx, "BRANCH", branchCodes, cpsAction, model.RequestEnableBranches)
}

func (s *AccountService) DisableBranches(ctx context.Context, branchCodes []string, enabled bool) error {
	cpsAction := generateCPSAction(ctx, "BRANCH", enabled, branchCodes, model.ActionDisable, model.RequestDisableBranches)
	return s.repo.EnableOrDisable(ctx, "BRANCH", branchCodes, cpsAction, model.RequestDisableBranches)
}

func (s *AccountService) EnableRegion(ctx context.Context, regionsCode []string, enabled bool) error {
	cpsAction := generateCPSAction(ctx, "REGION", enabled, regionsCode, model.ActionEnable, model.RequestEnableRegion)
	return s.repo.EnableOrDisable(ctx, "REGION", regionsCode, cpsAction, model.RequestEnableRegion)
}

func (s *AccountService) DisableRegion(ctx context.Context, regionsCode []string, enabled bool) error {
	cpsAction := generateCPSAction(ctx, "REGION", enabled, regionsCode, model.ActionDisable, model.RequestDisableRegion)
	return s.repo.EnableOrDisable(ctx, "REGION", regionsCode, cpsAction, model.RequestDisableRegion)
}

func (s *AccountService) EnableDistrict(ctx context.Context, districtsCode []string, enabled bool) error {
	cpsAction := generateCPSAction(ctx, "DISTRICT", enabled, districtsCode, model.ActionEnable, model.RequestEnableDistrict)
	return s.repo.EnableOrDisable(ctx, "DISTRICT", districtsCode, cpsAction, model.RequestEnableDistrict)
}

func (s *AccountService) DisableDistrict(ctx context.Context, districtsCode []string, enabled bool) error {
	cpsAction := generateCPSAction(ctx, "DISTRICT", enabled, districtsCode, model.ActionDisable, model.RequestDisableDistrict)
	return s.repo.EnableOrDisable(ctx, "DISTRICT", districtsCode, cpsAction, model.RequestDisableDistrict)
}

func (s *AccountService) EnableCity(ctx context.Context, citiesCode []string, enabled bool) error {
	cpsAction := generateCPSAction(ctx, "CITY", enabled, citiesCode, model.ActionEnable, model.RequestEnableCity)
	return s.repo.EnableOrDisable(ctx, "CITY", citiesCode, cpsAction, model.RequestEnableCity)
}

func (s *AccountService) DisableCity(ctx context.Context, citiesCode []string, enabled bool) error {
	cpsAction := generateCPSAction(ctx, "CITY", enabled, citiesCode, model.ActionDisable, model.RequestDisableCity)
	return s.repo.EnableOrDisable(ctx, "CITY", citiesCode, cpsAction, model.RequestDisableCity)
}

func (s AccountService) Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {

	switch action.RequestAction {
	case cps_const.RequestEnableBranches:
		return s.repo.AuthorizeEnableBranches(ctx, action)
	case cps_const.RequestDisableBranches:
		return s.repo.AuthorizeDisableBranches(ctx, action)
	case cps_const.RequestEnableRegion:
		return s.repo.AuthorizeEnableRegions(ctx, action)
	case cps_const.RequestDisableRegion:
		return s.repo.AuthorizeDisableRegions(ctx, action)
	case cps_const.RequestEnableDistrict:
		return s.repo.AuthorizeEnableDistricts(ctx, action)
	case cps_const.RequestDisableDistrict:
		return s.repo.AuthorizeDisableDistrict(ctx, action)
	case cps_const.RequestEnableCity:
		return s.repo.AuthorizeEnableCities(ctx, action)
	case cps_const.RequestDisableCity:
		return s.repo.AuthorizeDisableCities(ctx, action)
	default:
		return nil, fmt.Errorf("UNSUPPORTED_REQUEST_ACTION")

	}
}
