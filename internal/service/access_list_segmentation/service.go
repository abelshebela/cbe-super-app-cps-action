package access_list_segmentation_service

import (
	"cbe-super-app-cps-action/internal/constants"
	access_list_segmentation_dto "cbe-super-app-cps-action/internal/constants/dto/access_list_segmentation"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	local_model "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	access_list_segmentation_core "cbe-super-app-cps-action/internal/service/access_list_segmentation/core"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"fmt"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"

	local_util "cbe-super-app-cps-action/pkgs/utils"
)

type AccessListSegmentationService struct {
	repo        storage.AccessListSegmentationRepository
	accBlock    storage.AccountBlockRepository
	memberRepo  storage.CustomerRepository
	serviceRepo storage.ServicesRepository
	cpsAction   service.CPSActionService
	logger      utils.Logger
}

// Authorize implements service.AccessListSegmentationService.
func (a *AccessListSegmentationService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	a.logger.Infof("[Authorize] authorizing access list segmentation CPS action with ID: %s", cpsAction.ID.Hex())

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateAccessListSegmentation):
		action, err := local_util.JsonUnmarshal[access_list_segmentation_dto.CreateAccessListSegmentationRequest](cpsAction.CurrentAction)
		if err != nil {
			a.logger.Errorf("[Authorize] failed to unmarshal current action: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		if err := a.repo.Create(ctx, *action); err != nil {
			a.logger.Errorf("[Authorize] failed to create access list segmentation: %v", err)
			return nil, err
		}
	case string(constants.RequestUpdateAccessListSegmentation):
		action, err := local_util.JsonUnmarshal[access_list_segmentation_dto.UpdateAccessListSegmentationRequest](cpsAction.CurrentAction)
		if err != nil {
			a.logger.Errorf("[Authorize] failed to unmarshal current action: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}
		if action.ID == "" {
			a.logger.Errorf("[Authorize] missing access list segmentation ID")
			return nil, errors.New(localization.ErrorAccessListSegmentationInvalidID.Code)
		}
		if err := a.repo.Update(ctx, action.ID, *action); err != nil {
			a.logger.Errorf("[Authorize] failed to update access list segmentation: %v", err)
			return nil, err
		}
	case string(constants.RequestEnableDisableAccessListSegmentation):
		action, err := local_util.JsonUnmarshal[local_model.AccessListSegmentation](cpsAction.CurrentAction)
		if err != nil {
			a.logger.Errorf("[Authorize] failed to unmarshal current action: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}
		if action.ID.IsZero() {
			a.logger.Errorf("[Authorize] missing access list segmentation ID")
			return nil, errors.New(localization.ErrorAccessListSegmentationInvalidID.Code)
		}
		if err := a.repo.EnableOrDisable(ctx, action.ID.Hex(), action.Enabled); err != nil {
			a.logger.Errorf("[Authorize] failed to enable or disable access list segmentation: %v", err)
			return nil, err
		}
	default:
		a.logger.Errorf("[Authorize] unknown request action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}
	return cpsAction, nil
}

// CreateAccessListSegmentation implements service.AccessListSegmentationService.
func (a *AccessListSegmentationService) CreateAccessListSegmentation(ctx context.Context, req access_list_segmentation_dto.CreateAccessListSegmentationRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		a.logger.Errorf("[Create] incomplete user information")
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	service, err := a.serviceRepo.FindByID(ctx, req.ServiceID)
	if err != nil {
		a.logger.Errorf("[Create] failed to find service by id: %v", err)
		return err
	}
	req.ServiceName = service.ServiceName

	if als, err := a.repo.FindByIDS(ctx, req.SegmentedID, req.Type); err != nil {
		a.logger.Errorf("[Create] failed to find access list segmentation by id: %v", err)
		return err
	} else if als != nil {
		a.logger.Errorf("[Create] access list segmentation already exists with id: %s", req.SegmentedID)
		return fmt.Errorf("access list segmentation with id: %v already exists", als.ID)
	}

	if err := a.CheckALLIdsExist(ctx, req.Type, req.SegmentedID); err != nil {
		a.logger.Errorf("[Create] failed to check all IDs exist: %v", err)
		return err
	}

	cpsAction := lib.CpsModelBuilder("", makerData, nil, req, string(constants.RequestCreateAccessListSegmentation), constants.CREATE)

	if err := a.cpsAction.CreateCPSAction(ctx, &cpsAction); err != nil {
		a.logger.Errorf("[Create] failed to create CPS action: %v", err)
		return err
	}
	return nil
}

// EnableDisableAccessListSegmentation implements service.AccessListSegmentationService.
func (a *AccessListSegmentationService) EnableDisableAccessListSegmentation(ctx context.Context, id string, enabled bool) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplete := local_util.IsIncomplete(makerData); incomplete {
		a.logger.Errorf("[EnableDisable] incomplete user information")
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	accessListSegmentation, err := a.repo.FindByID(ctx, id)
	if err != nil {
		a.logger.Errorf("[EnableDisable] failed to find access list segmentation by id: %v", err)
		return err
	}

	if accessListSegmentation != nil && accessListSegmentation.Enabled == enabled {
		a.logger.Errorf("[EnableDisable] access list segmentation already in the desired state: %v", enabled)
		if enabled {
			return errors.New(localization.ErrorAccessListSegmentationAlreadyEnabled.Code)
		} else {
			return errors.New(localization.ErrorAccessListSegmentationAlreadyDisabled.Code)
		}
	}

	updatedAccessListSegmentation := *accessListSegmentation
	updatedAccessListSegmentation.Enabled = enabled

	cpsAction := lib.CpsModelBuilder(id, makerData, accessListSegmentation, updatedAccessListSegmentation, string(constants.RequestEnableDisableAccessListSegmentation), constants.UPDATE)

	if err := a.cpsAction.CreateCPSAction(ctx, &cpsAction); err != nil {
		a.logger.Errorf("[EnableDisable] failed to create CPS action: %v", err)
		return err
	}
	return nil
}

// GetAccessListSegmentationByID implements service.AccessListSegmentationService.
func (a *AccessListSegmentationService) GetAccessListSegmentationByID(ctx context.Context, id string) (access_list_segmentation_dto.AccessListSegmentationResponse, error) {
	model, err := a.repo.FindByID(ctx, id)
	if err != nil {
		a.logger.Errorf("[GetAccessListSegmentationByID] failed to get access list segmentation by id: %v", err)
		return access_list_segmentation_dto.AccessListSegmentationResponse{}, err
	}
	resp := access_list_segmentation_core.MapModelToDTO(*model)

	return resp, nil
}

// GetAllAccessListSegmentation implements service.AccessListSegmentationService.
func (a *AccessListSegmentationService) GetAllAccessListSegmentation(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]access_list_segmentation_dto.AccessListSegmentationResponse], error) {
	paginatedModel, err := a.repo.FindAllWithPagination(ctx, filter)
	if err != nil {
		a.logger.Errorf("[GetAllAccessListSegmentation] failed to get access list segmentations: %v", err)
		return nil, err
	}
	dto := access_list_segmentation_core.ConvertPaginatedModelToDTO(*paginatedModel)
	return &dto, nil
}

// UpdateAccessListSegmentation implements service.AccessListSegmentationService.
func (a *AccessListSegmentationService) UpdateAccessListSegmentation(ctx context.Context, req access_list_segmentation_dto.UpdateAccessListSegmentationRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplete := local_util.IsIncomplete(makerData); incomplete {
		a.logger.Errorf("[UpdateAccessListSegmentation] incomplete user information")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	accessListSegmentation, err := a.repo.FindByID(ctx, req.ID)
	if err != nil {
		a.logger.Errorf("[UpdateAccessListSegmentation] failed to find access list segmentation by id: %v", err)
		return err
	}

	if accessListSegmentation == nil {
		a.logger.Errorf("[UpdateAccessListSegmentation] access list segmentation with with id:%s not found", req.ID)
		return errors.New(localization.ErrorAccessListSegmentationNotFound.Code)
	}

	updatedAccessListSegmentation := *accessListSegmentation
	if req.NewSegmentedID != "" {
		objID, err := bson.ObjectIDFromHex(req.NewSegmentedID)
		if err != nil {
			a.logger.Errorf("[UpdateAccessListSegmentation] invalid new segmented id: %v", err)
			return errors.New(localization.ErrorAccessListSegmentationNotFound.Code)
		}
		updatedAccessListSegmentation.SegmentedID = objID
	}
	if req.NewServiceID != "" {
		objID, err := bson.ObjectIDFromHex(req.NewServiceID)
		if err != nil {
			a.logger.Errorf("[UpdateAccessListSegmentation] invalid new service id: %v", err)
			return errors.New(localization.ErrorServiceIdRequired.Code)
		}
		updatedAccessListSegmentation.ServiceID = objID
		service, err := a.serviceRepo.FindByID(ctx, req.NewServiceID)
		if err != nil {
			a.logger.Errorf("[Create] failed to find service by id: %v", err)
			return err
		}
		req.NewServiceName = service.ServiceName

	}
	if req.Type != "" {
		updatedAccessListSegmentation.Type = req.Type
	}

	if seg, err := a.repo.FindBySegmentationAndServiceID(ctx, req.NewSegmentedID, req.NewServiceID); err != nil || seg != nil {
		a.logger.Errorf("[UpdateAccessListSegmentation] access list segmentation already exists with segmented id and service id: %v", err)
		return errors.New(localization.ErrorAccessListSegmentationNameAlreadyExists.Code)
	}

	cpsAction := lib.CpsModelBuilder(req.ID, makerData, accessListSegmentation, updatedAccessListSegmentation, string(constants.RequestEnableDisableAccessListSegmentation), constants.UPDATE)

	if err := a.cpsAction.CreateCPSAction(ctx, &cpsAction); err != nil {
		a.logger.Errorf("[UpdateAccessListSegmentation] failed to create CPS action: %v", err)
		return err
	}
	return nil
}

func (a *AccessListSegmentationService) CheckALLIdsExist(ctx context.Context, t string, ids []string) error {
	switch t {
	case "U":
		if r, err := a.memberRepo.FindCustomerByIDs(ctx, ids); err != nil {
			return err
		} else if len(r) != len(ids) {
			var modelMaps []map[string]interface{}
			for _, d := range r {
				modelMaps = append(modelMaps, map[string]interface{}{"_id": d.ID.Hex()})
			}
			missingids := access_list_segmentation_core.GetMissingIds(ids, modelMaps)
			return fmt.Errorf("some user ids do not exist:%v", missingids)
		}
	case "B":
		if r, err := a.accBlock.GetBranchesByIds(ctx, ids); err != nil {
			return err
		} else if len(r) != len(ids) {
			var modelMaps []map[string]interface{}
			for _, d := range r {
				modelMaps = append(modelMaps, map[string]interface{}{"_id": d.ID.Hex()})
			}
			missingids := access_list_segmentation_core.GetMissingIds(ids, modelMaps)
			return fmt.Errorf("some Branch ids do not exist:%v", missingids)
		}
	case "R":
		if r, err := a.accBlock.GetRegionsByIds(ctx, ids); err != nil {
			return err
		} else if len(r) != len(ids) {
			var modelMaps []map[string]interface{}
			for _, d := range r {
				modelMaps = append(modelMaps, map[string]interface{}{"_id": d.ID.Hex()})
			}
			missingids := access_list_segmentation_core.GetMissingIds(ids, modelMaps)
			return fmt.Errorf("some Region ids do not exist:%v", missingids)
		}
	case "D":
		if r, err := a.accBlock.GetDistrictsByIds(ctx, ids); err != nil {
			return err
		} else if len(r) != len(ids) {
			var modelMaps []map[string]interface{}
			for _, d := range r {
				modelMaps = append(modelMaps, map[string]interface{}{"_id": d.ID.Hex()})
			}
			missingids := access_list_segmentation_core.GetMissingIds(ids, modelMaps)
			return fmt.Errorf("some District ids do not exist:%v", missingids)
		}
	case "C":
		if r, err := a.accBlock.GetCitiesByIds(ctx, ids); err != nil {
			return err
		} else if len(r) != len(ids) {
			var modelMaps []map[string]interface{}
			for _, d := range r {
				modelMaps = append(modelMaps, map[string]interface{}{"_id": d.ID.Hex()})
			}
			missingids := access_list_segmentation_core.GetMissingIds(ids, modelMaps)
			return fmt.Errorf("some City ids do not exist:%v", missingids)
		}
	default:
		return errors.New("type must be one of: B,R, D, C, U")
	}
	return nil
}

func NewAccessListSegmentationService(repo storage.AccessListSegmentationRepository, cpsAction service.CPSActionService, serviceRepo storage.ServicesRepository, accBlock storage.AccountBlockRepository, memberRepo storage.CustomerRepository, logger utils.Logger) service.AccessListSegmentationService {
	return &AccessListSegmentationService{
		repo:        repo,
		cpsAction:   cpsAction,
		memberRepo:  memberRepo,
		serviceRepo: serviceRepo,
		accBlock:    accBlock,
		logger:      logger,
	}
}
