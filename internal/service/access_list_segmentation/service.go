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
	repo                  storage.AccessListSegmentationRepository
	accBlock              storage.AccountBlockRepository
	customerSeg           storage.CPSRolesRepository
	memberRepo            storage.CustomerRepository
	accessListServiceRepo storage.AppAccessListRepository
	cpsAction             service.CPSActionService
	logger                utils.Logger
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

		if action.SegmentType == "Block" {
			if err := a.repo.CreateBlockSegment(ctx, *action); err != nil {
				a.logger.Errorf("[Authorize] failed to create access list segmentation: %v", err)
				return nil, err
			}
		} else {
			err := a.repo.CreateAccountSegment(ctx, *action)
			if err != nil {
				a.logger.Errorf("[Authorize] failed to create access list segmentation: %v", err)
				return nil, err
			}
			// publish to kafka

		}
	case string(constants.RequestUpdateAccessListSegmentation):
		action, err := local_util.JsonUnmarshal[local_model.AccessListSegmentation](cpsAction.CurrentAction)
		if err != nil {
			a.logger.Errorf("[Authorize] failed to unmarshal current action: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}
		if action.ID.Hex() == "" {
			a.logger.Errorf("[Authorize] missing access list segmentation ID")
			return nil, errors.New(localization.ErrorAccessListSegmentationInvalidID.Code)
		}
		if err := a.repo.Update(ctx, action.ID.Hex(), *action); err != nil {
			a.logger.Errorf("[Authorize] failed to update access list segmentation: %v", err)
			return nil, err
		}
	case string(constants.RequestEnableDisableAccessListSegmentation):
		bulkDisable, err := local_util.JsonUnmarshal[access_list_segmentation_dto.BulkDisableAccessListSegmentationRequest](cpsAction.CurrentAction)
		if err != nil {
			a.logger.Errorf("[Authorize] failed to unmarshal current action: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		if err := a.repo.BulkDisable(ctx, *bulkDisable); err != nil {
			a.logger.Errorf("[Authorize] failed to enable or disable access list segmentation: %v", err)
			return nil, err
		}
		// publish to kafka
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

	found, err := a.accessListServiceRepo.FindByKeys(ctx, req.AccessListKeys)
	if err != nil {
		a.logger.Errorf("[Create] failed to find service by id: %v", err)
		return err
	}
	for _, key := range req.AccessListKeys {
		if _, ok := found[key]; !ok {
			a.logger.Errorf("[Create] access list key not found: %s", key)
			return fmt.Errorf("access list key not found: %s", key)
		}
		req.AccessListNames = append(req.AccessListNames, found[key])
	}

	if req.SegmentType == "Block" {
		if als, err := a.repo.FindBySegmentIDAndAccessListKeys(ctx, req.SegmentedID, req.AccessListKeys); err != nil && err.Error() != localization.ErrorAccessListSegmentationNotFound.Code {
			a.logger.Errorf("[Create] failed to find access list segmentation by id: %v", err)
			return errors.New(localization.ErrorUnexpectedError.Code)
		} else if als != nil {
			a.logger.Errorf("[Create] access list segmentation already exists with id: %s", req.SegmentedID)
			return fmt.Errorf("access list segmentation with id: %v and access list key: %v already exists", req.SegmentedID, als.AccessListKey)
		}

		if err := a.CheckALLIdsExist(ctx, req.Type, []string{req.SegmentedID}); err != nil {
			a.logger.Errorf("[Create] failed to check all IDs exist: %v", err)
			return err
		}
	} else {
		// add checks for
		// 1. if the passed segment code is valid
		// 2. if service id and segment code combination already exists
		seg, err := a.customerSeg.FindByCustomerSegmentation(ctx, req.SegmentCode)
		if err != nil || seg == nil {
			a.logger.Errorf("[Create] access list segmentation already exists with segmented id and service id: %v", err)
			return errors.New(localization.ErrorCustomerSegmentationCodeNotFound.Code)
		}

		if seg, err := a.repo.FindByAccountSegmentationAndAccessListKeys(ctx, req.SegmentCode, req.AccessListKeys); err != nil || seg != nil {
			a.logger.Errorf("[Create] access list segmentation already exists with segmentation code and service id: %v", err)
			return errors.New(localization.ErrorAccessListSegmentationNameAlreadyExists.Code)
		}
	}
	cpsAction := lib.CpsModelBuilder("", makerData, nil, req, string(constants.RequestCreateAccessListSegmentation), constants.CREATE)

	if err := a.cpsAction.CreateCPSAction(ctx, &cpsAction); err != nil {
		a.logger.Errorf("[Create] failed to create CPS action: %v", err)
		return err
	}
	return nil
}

// EnableDisableAccessListSegmentation implements service.AccessListSegmentationService.
func (a *AccessListSegmentationService) EnableDisableAccessListSegmentation(ctx context.Context, id string, enabled bool, keys []string) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplete := local_util.IsIncomplete(makerData); incomplete {
		a.logger.Errorf("[EnableDisable] incomplete user information")
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	accessListSegmentation, err := a.repo.FindAllBySegmentIDorSegmentCodeAndKeys(ctx, id, keys)
	if err != nil {
		a.logger.Errorf("[EnableDisable] failed to find access list segmentation by id: %v", err)
		return err
	}

	if accessListSegmentation != nil && len(keys) != len(accessListSegmentation) {
		a.logger.Errorf("[EnableDisable] some access list segmentation keys not found for id: %s", id)
		return errors.New(localization.ErrorAccessListSegmentationKeyNotFound.Code)
	}

	bulkDisable := access_list_segmentation_dto.BulkDisableAccessListSegmentationRequest{
		Keys: keys,
		ID:   id,
	}

	// updatedAccessListSegmentation := accessListSegmentation
	// for _, accessListSegmentation := range updatedAccessListSegmentation {
	// 	accessListSegmentation.Enabled = enabled

	// }

	cpsAction := lib.CpsModelBuilder(id, makerData, accessListSegmentation, bulkDisable, string(constants.RequestEnableDisableAccessListSegmentation), constants.UPDATE)

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
	if req.NewAccessListKey != "" {

		als, err := a.accessListServiceRepo.FindByKeys(ctx, []string{req.NewAccessListKey})
		if err != nil {
			a.logger.Errorf("[Create] failed to find service by id: %v", err)
			return err
		}
		if _, ok := als[req.NewAccessListKey]; !ok {
			a.logger.Errorf("[UpdateAccessListSegmentation] access list key not found: %s", req.NewAccessListKey)
			return fmt.Errorf("access list key not found: %s", req.NewAccessListKey)
		}
		updatedAccessListSegmentation.AccessListKey = req.NewAccessListKey
		updatedAccessListSegmentation.AccessListName = als[req.NewAccessListKey]

	}
	if req.Type != "" {
		updatedAccessListSegmentation.Type = req.Type
	}

	if req.SegmentCode != "" {
		seg, err := a.customerSeg.FindByCustomerSegmentation(ctx, req.SegmentCode)
		if err != nil || seg == nil {
			a.logger.Errorf("[Create] access list segmentation already exists with segmented id and service id: %v", err)
			return errors.New(localization.ErrorCustomerSegmentationCodeNotFound.Code)
		}
		// req.SegmentName = seg.CustomerSubSegments[0].CustomerGroup
		updatedAccessListSegmentation.SegmentationCode = req.SegmentCode
		// updatedAccessListSegmentation.SegmentationName = seg.CustomerSubSegment
	}

	if seg, err := a.repo.FindByAccountSegmentationAndAccessListKeys(ctx, req.NewSegmentedID, []string{req.NewAccessListKey}); err != nil || seg != nil {
		a.logger.Errorf("[UpdateAccessListSegmentation] access list segmentation already exists with segmented id and service id: %v", err)
		return errors.New(localization.ErrorAccessListSegmentationNameAlreadyExists.Code)
	}

	cpsAction := lib.CpsModelBuilder(req.ID, makerData, accessListSegmentation, updatedAccessListSegmentation, string(constants.RequestUpdateAccessListSegmentation), constants.UPDATE)

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

func (a *AccessListSegmentationService) GetAllAccessListSegmentationBySegmentIDorSegmentCode(ctx context.Context, segmentIdentifier string) ([]model.APPAccessList, []model.APPAccessList, error) {
	accessListSegmentation, err := a.repo.FindAllBySegmentIDorSegmentCode(ctx, segmentIdentifier)
	if err != nil {
		a.logger.Errorf("[GetAllAccessListSegmentationBySegmentIDorSegmentCode] failed to get access list segmentation: %v", err)
		return nil, nil, err
	}

	accessList := access_list_segmentation_core.FindNoneSegmentedAccessList(ctx, a.accessListServiceRepo, accessListSegmentation)

	realtions, err := a.repo.FindParentChildRelationship(ctx)
	if err != nil {
		a.logger.Errorf("[GetAllAccessListSegmentationBySegmentIDorSegmentCode] failed to get parent child relationship: %v", err)
		return nil, nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	accessList, segmentedAccessList := access_list_segmentation_core.MapParentChildRelationship(accessListSegmentation, realtions, &accessList)

	return accessList, segmentedAccessList, nil
}

func NewAccessListSegmentationService(repo storage.AccessListSegmentationRepository, cpsAction service.CPSActionService, accessListServiceRepo storage.AppAccessListRepository, accBlock storage.AccountBlockRepository, memberRepo storage.CustomerRepository, customerSeg storage.CPSRolesRepository, logger utils.Logger) service.AccessListSegmentationService {
	return &AccessListSegmentationService{
		repo:                  repo,
		cpsAction:             cpsAction,
		customerSeg:           customerSeg,
		memberRepo:            memberRepo,
		accessListServiceRepo: accessListServiceRepo,
		accBlock:              accBlock,
		logger:                logger,
	}
}
