package customersegmentation

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	cust_seg "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/customer_segmentation"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/service"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/service/customer_segmentation/core"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type customerSegmentationService struct {
	repo             storage.CustomerSegmentationRepository
	cpsRoleRepo      storage.CPSRolesRepository
	cpsService       service.CPSActionService
	segmentationRepo storage.AccessListSegmentationRepository
	logger           utils.Logger
}

func NewCustomerSegmentation(repo storage.CustomerSegmentationRepository, cpsRoleRepo storage.CPSRolesRepository, cpsService service.CPSActionService, segmentationRepo storage.AccessListSegmentationRepository, logger utils.Logger) *customerSegmentationService {
	return &customerSegmentationService{
		repo:             repo,
		cpsRoleRepo:      cpsRoleRepo,
		cpsService:       cpsService,
		segmentationRepo: segmentationRepo,
		logger:           logger,
	}
}

func (s *customerSegmentationService) Create(ctx context.Context, req cust_seg.CreateCustomerSegmentationRequest) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	var newData imodel.CustomerSegmentation
	role, err := s.cpsRoleRepo.FindById(ctx, req.CustomerRole)
	if err != nil {
		if err.Error() == localization.ErrorResourceNotFound.Code {
			log.Errorf("[CustSegSvc][Create] role not found: %s", req.CustomerRole)
			return fmt.Errorf("This role not found")
		} else {
			log.Errorf("[CustSegSvc][Create] role fetch err: %s, %v", req.CustomerRole, err)
			return err
		}
	}

	exist, err := s.repo.CheckIfCustomerSubSegmentExists(ctx, role.ID)
	if err != nil {
		log.Errorf("[CustSegSvc][Create] sub-segment check err: %s, %v", role.ID, err)
		return err
	}
	if exist {
		log.Errorf("[CustSegSvc][Create] sub-segment already exists: %s", role.ID)
		return fmt.Errorf("Customer segmentation already exists")
	}

	if role != nil {
		newData.CustomerRole = imodel.CustomerRoleInfo{ID: role.ID, Name: role.Name, Label: role.Lable}
		newData.CreatedAt = time.Now()
		newData.UpdatedAt = time.Now()
		newData.IsEnabled = true
		newData.IsDeleted = false
		newData.Customer = req.Customer
		newData.SyncSegmentsFromCustomer()
	}

	cpsActionData := lib.CpsModelBuilder("", makerUser, nil, core.MapCustomerSegmentationToMap(newData), string(constants.RequestCreateCustomerSegmentation), constants.CREATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[CustSegSvc][Create] cps action err: %v", err)
		return err
	}
	log.Infof("[CustSegSvc][Create] request created")

	return nil
}

func (s *customerSegmentationService) Update(ctx context.Context, id string, req cust_seg.UpdateCustomerSegmentationRequest) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[CustSegSvc][Update] find err: %v", err)
		return err
	}

	updated := *existing
	updated.Customer = req.Customer
	updated.UpdatedAt = time.Now()

	cpsActionData := lib.CpsModelBuilder(id, makerUser, core.MapCustomerSegmentationToMap(*existing), core.MapCustomerSegmentationToMap(updated), string(constants.RequestUpdateCustomerSegmentation), constants.UPDATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[CustSegSvc][Update] cps action err: %v", err)
		return err
	}

	log.Infof("[CustSegSvc][Update] request created")
	return nil
}

func (s *customerSegmentationService) isSubSegmentsEqual(existing, incoming []imodel.CustomerEntry) bool {
	if len(existing) != len(incoming) {
		return false
	}

	for i := range existing {
		if existing[i].Group.CustGroup != incoming[i].Group.CustGroup ||
			existing[i].Segment.CustSegmentName != incoming[i].Segment.CustSegmentName {
			return false
		}
	}

	return true
}

func (s *customerSegmentationService) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]imodel.CustomerSegmentation], error) {
	return s.repo.FindAllWithPagination(ctx, *filterParam)
}

func (s *customerSegmentationService) FindById(ctx context.Context, id string) (*imodel.CustomerSegmentation, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *customerSegmentationService) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[CustSegSvc][EnableDisable] find err: %v", err)
		return err
	}

	if existing.IsEnabled && enable {
		log.Warnf("[CustSegSvc][EnableDisable] already in state id: %s, enabled: %v", id, enable)
		return errors.New("database key already enabled")
	}

	if !existing.IsEnabled && !enable {
		log.Warnf("[CustSegSvc][EnableDisable] already in state id: %s, disabled: %v", id, enable)
		return errors.New("database key already disabled")
	}

	updated := *existing
	updated.IsEnabled = enable
	updated.UpdatedAt = time.Now()

	var action constants.RequestAction
	if enable {
		action = constants.RequestEnableCustomerSegmentation
	} else {
		action = constants.RequestDisableCustomerSegmentation
	}

	cpsActionData := lib.CpsModelBuilder(id, makerUser, core.MapCustomerSegmentationToMap(*existing), core.MapCustomerSegmentationToMap(updated), string(action), constants.UPDATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[CustSegSvc][EnableDisable] cps action err: %v", err)
		return err
	}

	return nil
}

func (s *customerSegmentationService) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[CustSegSvc][Delete] find err: %v", err)
		return err
	}

	updated := *existing
	updated.IsDeleted = true
	updated.UpdatedAt = time.Now()
	cpsActionData := lib.CpsModelBuilder(id, makerUser, core.MapCustomerSegmentationToMap(*existing), core.MapCustomerSegmentationToMap(updated), string(constants.RequestDeleteCustomerSegmentation), constants.DELETE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[CustSegSvc][Delete] cps action err: %v", err)
		return err
	}

	return nil
}

func (s *customerSegmentationService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	log.Infof("[CustSegSvc][Authorize] action: %s", action.RequestAction)

	var (
		err error
		seg *imodel.CustomerSegmentation
	)

	switch action.RequestAction {
	case string(constants.RequestDeleteCustomerSegmentation):
		err = s.repo.Delete(ctx, action.UniqueId)
		if err != nil {
			log.Errorf("[CustSegSvc][Authorize] delete err: %v", err)
			return nil, err
		}
		log.Infof("[CustSegSvc][Authorize] deleted")
		return action, nil
	}

	seg, marshal_err := local_util.JsonUnmarshal[imodel.CustomerSegmentation](action.CurrentAction)
	if marshal_err != nil || seg == nil {
		log.Errorf("[CustSegSvc][Authorize] unmarshal err: %v", marshal_err)
		return nil, marshal_err
	}
	switch action.RequestAction {
	case string(constants.RequestCreateCustomerSegmentation):
		err = s.repo.Create(ctx, seg)
		if err != nil {
			log.Errorf("[CustSegSvc][Authorize] create err: %v", err)
			return nil, err
		}
		log.Infof("[CustSegSvc][Authorize] created")
	case string(constants.RequestUpdateCustomerSegmentation):
		err = s.repo.Update(ctx, action.UniqueId, seg)
		if err != nil {
			log.Errorf("[CustSegSvc][Authorize] update err: %v", err)
			return nil, err
		}
		log.Infof("[CustSegSvc][Authorize] updated")
	case string(constants.RequestEnableCustomerSegmentation):
		err = s.repo.EnableOrDisable(ctx, action.UniqueId, true)
		if err != nil {
			log.Errorf("[CustSegSvc][Authorize] enable err: %v", err)
			return nil, err
		}
		log.Infof("[CustSegSvc][Authorize] enabled")
	case string(constants.RequestDisableCustomerSegmentation):
		err = s.repo.EnableOrDisable(ctx, action.UniqueId, false)
		if err != nil {
			log.Errorf("[CustSegSvc][Authorize] disable err: %v", err)
			return nil, err
		}
		log.Infof("[CustSegSvc][Authorize] disabled")
	default:
		log.Errorf("[CustSegSvc][Authorize] unsupported action: %s", action.RequestAction)
		return nil, errors.New("unsupported action")
	}
	return action, nil
}
