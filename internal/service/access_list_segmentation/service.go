package access_list_segmentation_service

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	access_list_segmentation_dto "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/access_list_segmentation"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	local_model "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/service"
	access_list_segmentation_core "github.com/abelshebela/cbe-super-app-cps-action/internal/service/access_list_segmentation/core"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"fmt"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"
)

type AccessListSegmentationService struct {
	repo                  storage.AccessListSegmentationRepositoryOracle
	accBlock              storage.AccountBlockRepository
	customerSeg           storage.CPSRolesRepository
	memberRepo            storage.CustomerRepository
	accessListServiceRepo storage.BulkServiceRepository
	cpsAction             service.CPSActionService
	logger                utils.Logger
}

// Authorize implements service.AccessListSegmentationService.
func (a *AccessListSegmentationService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)

	log.Infof("[AccessListSegSvc][Authorize] id: %s", cpsAction.ID.Hex())

	switch cpsAction.RequestAction {
	case string(constants.RequestDisableAccessListSegmentation), string(constants.RequestAccessListDisableCustomerSegmentation):
		action, err := local_util.JsonUnmarshal[access_list_segmentation_dto.CreateAccessListSegmentationRequest](cpsAction.CurrentAction)
		if err != nil {
			log.Errorf("[AccessListSegSvc][Authorize] unmarshal err: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		if strings.EqualFold(action.SegmentType, "block") {
			if err := a.repo.CreateBlockSegment(ctx, *action); err != nil {
				log.Errorf("[AccessListSegSvc][Authorize] create block err: %v", err)
				return nil, err
			}
		} else {
			err := a.repo.CreateAccountSegment(ctx, *action)
			if err != nil {
				log.Errorf("[AccessListSegSvc][Authorize] create acct err: %v", err)
				return nil, err
			}
			// publish to kafka

		}
	case string(constants.RequestUpdateAccessListSegmentation):
		action, err := local_util.JsonUnmarshal[local_model.AccessListSegmentation](cpsAction.CurrentAction)
		if err != nil {
			log.Errorf("[AccessListSegSvc][Authorize] unmarshal update err: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}
		if action.ID == "" {
			log.Errorf("[AccessListSegSvc][Authorize] missing ID")
			return nil, errors.New(localization.ErrorAccessListSegmentationInvalidID.Code)
		}
		if err := a.repo.Update(ctx, action.ID, *action); err != nil {
			log.Errorf("[AccessListSegSvc][Authorize] update err: %v", err)
			return nil, err
		}
	case string(constants.RequestEnableDisableAccessListSegmentation), string(constants.RequestEnableAccessListSegmentation), string(constants.RequestAccessListEnableCustomerSegmentation):
		bulkDisable, err := local_util.JsonUnmarshal[access_list_segmentation_dto.BulkDisableAccessListSegmentationRequest](cpsAction.CurrentAction)
		if err != nil {
			log.Errorf("[AccessListSegSvc][Authorize] unmarshal enable/disable err: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		if err := a.repo.BulkDisable(ctx, *bulkDisable); err != nil {
			log.Errorf("[AccessListSegSvc][Authorize] bulk disable err: %v", err)
			return nil, err
		}
		// publish to kafka
	default:
		log.Errorf("[AccessListSegSvc][Authorize] unknown: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}
	return cpsAction, nil
}

// CreateAccessListSegmentation implements service.AccessListSegmentationService.
func (a *AccessListSegmentationService) CreateAccessListSegmentation(ctx context.Context, req access_list_segmentation_dto.CreateAccessListSegmentationRequest) error {
	log := local_util.LoggerFromCtx(ctx, a.logger)

	log.Infof("[AccessListSegSvc][Create] segment id: %s, type: %s", req.SegmentationID, req.Type)
	var requestAction string
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		log.Errorf("[AccessListSegSvc][Create] incomplete user")
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	found, err := a.accessListServiceRepo.FindAllByKeys(ctx, req.AccessListKeys)
	if err != nil || found == nil {
		log.Errorf("[AccessListSegSvc][Create] find keys err: %v", err)
		return err
	}

	if strings.EqualFold(req.SegmentType, "block") {
		if als, err := a.repo.FindBySegmentIDAndAccessListKeys(ctx, req.SegmentationID, req.AccessListKeys); err != nil && err.Error() != localization.ErrorAccessListSegmentationNotFound.Code {
			log.Errorf("[AccessListSegSvc][Create] find seg err: %v", err)
			return errors.New(localization.ErrorUnexpectedError.Code)
		} else if als != nil {
			log.Errorf("[AccessListSegSvc][Create] already exists id: %s", req.SegmentationID)
			return fmt.Errorf("access list segmentation with id: %v and access list key: %v already exists", req.SegmentationID, als.AccessListKey)
		}

		if err := a.CheckALLIdsExist(ctx, req.Type, []string{req.SegmentationID}); err != nil {
			log.Errorf("[AccessListSegSvc][Create] check IDs err: %v", err)
			return err
		}
		requestAction = string(constants.RequestDisableAccessListSegmentation)
	} else {
		seg, err := a.customerSeg.FindByCustomerSegmentationByID(ctx, req.SegmentationID)
		if err != nil || seg == nil {
			log.Errorf("[AccessListSegSvc][Create] seg code not found: %v", err)
			return errors.New(localization.ErrorCustomerSegmentationCodeNotFound.Code)
		}

		existing, err := a.repo.FindByAccountSegmentationAndAccessListKeys(ctx, req.SegmentationID, req.AccessListKeys)
		if err != nil {
			log.Errorf("[AccessListSegSvc][Create] dup check err: %v", err)
			return errors.New(localization.ErrorUnexpectedError.Code)
		}
		if existing != nil {
			log.Errorf("[AccessListSegSvc][Create] seg+key exists id: %s", req.SegmentationID)
			return errors.New(localization.ErrorAccessListSegmentationNameAlreadyExists.Code)
		}
		requestAction = string(constants.RequestAccessListDisableCustomerSegmentation)
	}
	cpsAction := lib.CpsModelBuilder("", makerData, nil, req, requestAction, constants.CREATE)

	if err := a.cpsAction.CreateCPSAction(ctx, &cpsAction); err != nil {
		log.Errorf("[AccessListSegSvc][Create] cps action err: %v", err)
		return err
	}
	return nil
}

// EnableDisableAccessListSegmentation implements service.AccessListSegmentationService.
func (a *AccessListSegmentationService) EnableDisableAccessListSegmentation(ctx context.Context, id string, enabled bool, keys []string, segmentation_type string) error {
	log := local_util.LoggerFromCtx(ctx, a.logger)

	var requestAction string
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplete := local_util.IsIncomplete(makerData); incomplete {
		log.Errorf("[AccessListSegSvc][EnableDisable] incomplete user")
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}
	var accessListSegmentation []local_model.AccessListSegmentation
	var err error
	if segmentation_type == "block" {
		accessListSegmentation, err = a.repo.FindAllByBlockAndKeys(ctx, id, keys)

		requestAction = string(constants.RequestEnableAccessListSegmentation)

	} else if segmentation_type == "account" {
		accessListSegmentation, err = a.repo.FindAllByAccountAndKeys(ctx, id, keys)

		requestAction = string(constants.RequestAccessListEnableCustomerSegmentation)

	} else {
		log.Errorf("[AccessListSegSvc][EnableDisable] invalid segmentation type: %s", segmentation_type)
		return localization.ErrorUnexpectedError
	}

	if err != nil {
		log.Errorf("[AccessListSegSvc][EnableDisable] find err: %v", err)
		return err
	}

	if accessListSegmentation != nil && len(keys) != len(accessListSegmentation) {
		log.Errorf("[AccessListSegSvc][EnableDisable] keys not found id: %s", id)
		return errors.New(localization.ErrorAccessListSegmentationKeyNotFound.Code)
	}
	bulkDisable := access_list_segmentation_dto.BulkDisableAccessListSegmentationRequest{
		Keys:             keys,
		ID:               id,
		SegmentationType: segmentation_type,
	}
	bulkDisable.Enabled = enabled

	// updatedAccessListSegmentation := accessListSegmentation
	// for _, accessListSegmentation := range updatedAccessListSegmentation {
	// 	accessListSegmentation.Enabled = enabled

	// }

	cpsAction := lib.CpsModelBuilder(id, makerData, accessListSegmentation, bulkDisable, requestAction, constants.UPDATE)

	if err := a.cpsAction.CreateCPSAction(ctx, &cpsAction); err != nil {
		log.Errorf("[AccessListSegSvc][EnableDisable] cps action err: %v", err)
		return err
	}
	return nil
}

// GetAccessListSegmentationForAccountByID implements service.AccessListSegmentationService.
func (a *AccessListSegmentationService) GetAccessListSegmentationForAccountByID(ctx context.Context, id string) (access_list_segmentation_dto.AccessListSegmentationResponse, error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)

	model, err := a.repo.FindAccountSegmentByID(ctx, id)
	if err != nil {
		log.Errorf("[AccessListSegSvc][GetByID] err: %v", err)
		return access_list_segmentation_dto.AccessListSegmentationResponse{}, err
	}
	resp := access_list_segmentation_core.MapModelToDTO(*model)

	return resp, nil
}

// GetAccessListSegmentationForBlockByID implements service.AccessListSegmentationService.
func (a *AccessListSegmentationService) GetAccessListSegmentationForBlockByID(ctx context.Context, id string) (access_list_segmentation_dto.AccessListSegmentationResponse, error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)

	model, err := a.repo.FindBlockSegmentByID(ctx, id)
	if err != nil {
		log.Errorf("[AccessListSegSvc][GetByID] err: %v", err)
		return access_list_segmentation_dto.AccessListSegmentationResponse{}, err
	}
	resp := access_list_segmentation_core.MapModelToDTO(*model)

	return resp, nil
}

// GetAllAccessListSegmentation implements service.AccessListSegmentationService.
func (a *AccessListSegmentationService) GetAllAccessListSegmentation(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]access_list_segmentation_dto.AccessListSegmentationResponse], error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)

	paginatedModel, err := a.repo.FindAllWithPagination(ctx, filter)
	if err != nil {
		log.Errorf("[AccessListSegSvc][GetAll] err: %v", err)
		return nil, err
	}
	dto := access_list_segmentation_core.ConvertPaginatedModelToDTO(*paginatedModel)
	return &dto, nil
}

// UpdateAccessListSegmentation implements service.AccessListSegmentationService.
func (a *AccessListSegmentationService) UpdateAccessListSegmentation(ctx context.Context, req access_list_segmentation_dto.UpdateAccessListSegmentationRequest) error {
	log := local_util.LoggerFromCtx(ctx, a.logger)

	var requestAction string
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplete := local_util.IsIncomplete(makerData); incomplete {
		log.Errorf("[AccessListSegSvc][Update] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	accessListSegmentation, err := a.repo.FindAccountSegmentByID(ctx, req.ID)
	if err != nil {
		log.Errorf("[AccessListSegSvc][Update] find err: %v", err)
		return err
	}

	if accessListSegmentation == nil {
		log.Errorf("[AccessListSegSvc][Update] not found id: %s", req.ID)
		return errors.New(localization.ErrorAccessListSegmentationNotFound.Code)
	}

	updatedAccessListSegmentation := *accessListSegmentation
	if req.NewSegmentationID != "" {
		// objID, err := bson.ObjectIDFromHex(req.NewSegmentedID)
		// if err != nil {
		// 	log.Errorf("[AccessListSegSvc][Update] invalid seg id: %v", err)
		// 	return errors.New(localization.ErrorAccessListSegmentationNotFound.Code)
		// }
		updatedAccessListSegmentation.SegmentedID = req.NewSegmentationID
	}
	if req.NewAccessListKey != "" {

		als, err := a.accessListServiceRepo.FindByKeys(ctx, []string{req.NewAccessListKey})
		if err != nil {
			log.Errorf("[AccessListSegSvc][Update] find key err: %v", err)
			return err
		}
		if _, ok := als[req.NewAccessListKey]; !ok {
			log.Errorf("[AccessListSegSvc][Update] key not found: %s", req.NewAccessListKey)
			return fmt.Errorf("access list key not found: %s", req.NewAccessListKey)
		}
		updatedAccessListSegmentation.AccessListKey = req.NewAccessListKey
		updatedAccessListSegmentation.AccessListName = als[req.NewAccessListKey]

	}
	if req.Type != "" {
		updatedAccessListSegmentation.Type = req.Type
	}

	if req.NewSegmentationID != "" {
		seg, err := a.customerSeg.FindByCustomerSegmentationByID(ctx, req.NewSegmentationID)
		if err != nil || seg == nil {
			log.Errorf("[AccessListSegSvc][Update] seg code not found: %v", err)
			return errors.New(localization.ErrorCustomerSegmentationCodeNotFound.Code)
		}
		// req.SegmentName = seg.CustomerSubSegments[0].CustomerGroup
		updatedAccessListSegmentation.SegmentationCode = req.NewSegmentationID
		// updatedAccessListSegmentation.SegmentationName = seg.CustomerSubSegment
	}

	if seg, err := a.repo.FindByAccountSegmentationAndAccessListKeys(ctx, req.NewSegmentationID, []string{req.NewAccessListKey}); err != nil || seg != nil {
		log.Errorf("[AccessListSegSvc][Update] seg+key exists: %v", err)
		return errors.New(localization.ErrorAccessListSegmentationNameAlreadyExists.Code)
	}

	if req.Type == "block" {
		requestAction = string(constants.RequestUpdateAccessListSegmentation)
	} else if req.Type == "account" {
		requestAction = string(constants.RequestAccessListUpdateCustomerSegmentation)
	} else {
		log.Errorf("[AccessListSegSvc][Update] invalid type: %s", req.Type)
		return errors.New(localization.ErrorInvalidSegmentationType.Code)
	}

	cpsAction := lib.CpsModelBuilder(req.ID, makerData, accessListSegmentation, updatedAccessListSegmentation, requestAction, constants.UPDATE)

	if err := a.cpsAction.CreateCPSAction(ctx, &cpsAction); err != nil {
		log.Errorf("[AccessListSegSvc][Update] cps action err: %v", err)
		return err
	}
	return nil
}

func (a *AccessListSegmentationService) CheckALLIdsExist(ctx context.Context, t string, ids []string) error {
	log := local_util.LoggerFromCtx(ctx, a.logger)

	switch t {
	case "B":
		if r, err := a.accBlock.GetBranchesByIds(ctx, ids); err != nil {
			return err
		} else if len(r) != len(ids) {
			var modelMaps []map[string]interface{}
			for _, d := range r {
				modelMaps = append(modelMaps, map[string]interface{}{"_id": d.ID})
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
				modelMaps = append(modelMaps, map[string]interface{}{"_id": d.ID})
			}
			missingids := access_list_segmentation_core.GetMissingIds(ids, modelMaps)
			return fmt.Errorf("some Region ids do not exist:%v", missingids)
		}
	case "D":
		if a.accBlock == nil {
			log.Errorf("[AccessListSegSvc][CheckALLIdsExist] accBlock repo is nil")
			return errors.New(localization.ErrorUnexpectedError.Code)
		}
		if r, err := a.accBlock.GetDistrictsByIds(ctx, ids); err != nil {
			return err
		} else if len(r) != len(ids) {
			var modelMaps []map[string]interface{}
			for _, d := range r {
				modelMaps = append(modelMaps, map[string]interface{}{"_id": d.ID})
			}
			missingids := access_list_segmentation_core.GetMissingIds(ids, modelMaps)
			return fmt.Errorf("some District ids do not exist:%v", missingids)
		}
	// case "C":
	// 	if r, err := a.accBlock.GetCitiesByIds(ctx, ids); err != nil {
	// 		return err
	// 	} else if len(r) != len(ids) {
	// 		var modelMaps []map[string]interface{}
	// 		for _, d := range r {
	// 			modelMaps = append(modelMaps, map[string]interface{}{"_id": d.ID})
	// 		}
	// 		missingids := access_list_segmentation_core.GetMissingIds(ids, modelMaps)
	// 		return fmt.Errorf("some City ids do not exist:%v", missingids)
	// 	}
	default:
		return errors.New("type must be one of: B,R, D, C, U")
	}
	return nil
}

func (a *AccessListSegmentationService) GetAllAccessListSegmentationForBlock(ctx context.Context, segmentIdentifier string) ([]local_model.APPAccessList, []local_model.APPAccessList, []local_model.APPAccessList, error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)

	accessListSegmentation, err := a.repo.FindAllForBlock(ctx, segmentIdentifier)
	if err != nil {
		log.Errorf("[AccessListSegSvc][GetAllAccessListSegmentationForBlock] err: %v", err)
		return nil, nil, nil, err
	}
	accessListSegmentationFromParent, err := a.repo.FindAllForBlockParents(ctx, segmentIdentifier)
	if err != nil {
		log.Errorf("[AccessListSegSvc][GetAllAccessListSegmentationForBlock][FindAllForBlockParents] err: %v", err)
		return nil, nil, nil, err
	}

	accessList := access_list_segmentation_core.FindNoneSegmentedAccessList(ctx, a.accessListServiceRepo, accessListSegmentation, accessListSegmentationFromParent)

	realtions, err := a.repo.FindParentChildRelationship(ctx)
	if err != nil {
		log.Errorf("[AccessListSegSvc][GetAllAccessListSegmentationForBlock] parent-child err: %v", err)
		return nil, nil, nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	accessList = access_list_segmentation_core.MapParentChildRelationship(realtions, accessList)
	accessListSegmentation = access_list_segmentation_core.MapParentChildRelationship(realtions, accessListSegmentation)
	accessListSegmentationFromParent = access_list_segmentation_core.MapParentChildRelationship(realtions, accessListSegmentationFromParent)

	return accessList, accessListSegmentation, accessListSegmentationFromParent, nil
}

func (a *AccessListSegmentationService) GetAllAccessListSegmentationForAccount(ctx context.Context, segmentIdentifier string) ([]local_model.APPAccessList, []local_model.APPAccessList, error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)

	accessListSegmentation, err := a.repo.FindAllForAccount(ctx, segmentIdentifier)

	if err != nil {
		log.Errorf("[AccessListSegSvc][GetBySegID] err: %v", err)
		return nil, nil, err
	}

	accessList := access_list_segmentation_core.FindNoneSegmentedAccessList(ctx, a.accessListServiceRepo, accessListSegmentation, []local_model.APPAccessList{})

	realtions, err := a.repo.FindParentChildRelationship(ctx)
	if err != nil {
		log.Errorf("[AccessListSegSvc][GetBySegID] parent-child err: %v", err)
		return nil, nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	accessList = access_list_segmentation_core.MapParentChildRelationship(realtions, accessList)
	accessListSegmentation = access_list_segmentation_core.MapParentChildRelationship(realtions, accessListSegmentation)

	return accessList, accessListSegmentation, nil
}

func NewAccessListSegmentationService(repo storage.AccessListSegmentationRepositoryOracle, cpsAction service.CPSActionService, accessListServiceRepo storage.BulkServiceRepository, accBlock storage.AccountBlockRepository, memberRepo storage.CustomerRepository, customerSeg storage.CPSRolesRepository, logger utils.Logger) service.AccessListSegmentationService {
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
