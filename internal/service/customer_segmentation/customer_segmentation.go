package customersegmentation

import (
	"cbe-super-app-cps-action/internal/constants"
	cust_seg "cbe-super-app-cps-action/internal/constants/dto/customer_segmentation"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"fmt"
	"time"

	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service/customer_segmentation/core"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	imodel "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type customerSegmentationService struct {
	repo        storage.CustomerSegmentationRepository
	cpsRoleRepo storage.CPSRolesRepository
	cpsService  service.CPSActionService
	logger      utils.Logger
}

func NewCustomerSegmentation(repo storage.CustomerSegmentationRepository, cpsRoleRepo storage.CPSRolesRepository, cpsService service.CPSActionService, logger utils.Logger) *customerSegmentationService {
	return &customerSegmentationService{
		repo:        repo,
		cpsRoleRepo: cpsRoleRepo,
		cpsService:  cpsService,
		logger:      logger,
	}
}

func (s *customerSegmentationService) Create(ctx context.Context, req cust_seg.CreateCustomerSegmentationRequest) error {
	makerUser := local_util.ExtractUserFromContext(ctx)

	var newData imodel.CustomerSegmentation
	role, err := s.cpsRoleRepo.FindById(ctx, req.CustomerRole)
	if err != nil {
		if err.Error() == localization.ErrorResourceNotFound.Code {
			s.logger.Errorf("[CustSegSvc][Create] role not found: %s", req.CustomerRole)
			return fmt.Errorf("This role not found")
		} else {
			s.logger.Errorf("[CustSegSvc][Create] role fetch err: %s, %v", req.CustomerRole, err)
			return err
		}
	}

	if role != nil {
		newData.CustomerRole = imodel.CustomerRoleInfo{ID: role.ID, Name: role.Name}
		newData.CreatedAt = time.Now()
		newData.UpdatedAt = time.Now()
		newData.IsEnabled = true
		newData.IsDeleted = false
		for _, info := range req.CustomerSubSegments {
			seg := imodel.CustomerSubSegments{
				Name:            info.Name,
				CustomerGroup:   info.CustomerGroup,
				CustomerSegment: info.CustomerSegment,
			}
			newData.CustomerSubSegments = append(newData.CustomerSubSegments, seg)
		}
	}

	cpsActionData := lib.CpsModelBuilder("", makerUser, nil, core.MapCustomerSegmentationToMap(newData), string(constants.RequestCreateCustomerSegmentation), constants.CREATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		s.logger.Errorf("[CustSegSvc][Create] cps action err: %v", err)
		return err
	}
	s.logger.Infof("[CustSegSvc][Create] request created")

	return nil
}

func (s *customerSegmentationService) Update(ctx context.Context, id string, req cust_seg.UpdateCustomerSegmentationRequest) error {
	makerUser := local_util.ExtractUserFromContext(ctx)

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("[CustSegSvc][Update] find err: %v", err)
		return err
	}

	updated := *existing
	if req.CustomerSubSegments != nil {
		updated.CustomerSubSegments = make([]imodel.CustomerSubSegments, len(req.CustomerSubSegments))
		for i, sub := range req.CustomerSubSegments {
			updated.CustomerSubSegments[i] = imodel.CustomerSubSegments{
				Name:            sub.Name,
				CustomerGroup:   sub.CustomerGroup,
				CustomerSegment: sub.CustomerSegment,
			}
		}
	}

	// if len(req.CustomerSubSegments) > 0 && len(req.OldName) > 0 {
	// 	updates := make(map[string]imodel.CustomerSubSegments)

	// 	for i, sub := range req.CustomerSubSegments {
	// 		updates[req.OldName[i]] = imodel.CustomerSubSegments{
	// 			Name:            sub.Name,
	// 			CustomerGroup:   sub.CustomerGroup,
	// 			CustomerSegment: sub.CustomerSegment,
	// 		}
	// 	}

	// 	for i, existingSub := range updated.CustomerSubSegments {
	// 		if newSub, ok := updates[existingSub.Name]; ok {
	// 			updated.CustomerSubSegments[i] = newSub
	// 		}
	// 	}
	// }

	updated.UpdatedAt = time.Now()

	cpsActionData := lib.CpsModelBuilder(id, makerUser, core.MapCustomerSegmentationToMap(*existing), core.MapCustomerSegmentationToMap(updated), string(constants.RequestUpdateCustomerSegmentation), constants.UPDATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		s.logger.Errorf("[CustSegSvc][Update] cps action err: %v", err)
		return err
	}

	s.logger.Infof("[CustSegSvc][Update] request created")
	return nil
}

func (s *customerSegmentationService) isSubSegmentsEqual(existing, incoming []imodel.CustomerSubSegments) bool {
	if len(existing) != len(incoming) {
		return false
	}

	for i := range existing {
		if existing[i].CustomerGroup != incoming[i].CustomerGroup ||
			existing[i].CustomerSegment != incoming[i].CustomerSegment {
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
	makerUser := local_util.ExtractUserFromContext(ctx)

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("[CustSegSvc][EnableDisable] find err: %v", err)
		return err
	}

	if existing.IsEnabled == enable {
		s.logger.Warnf("[CustSegSvc][EnableDisable] already in state id: %s, enable: %v", id, enable)
		return errors.New("already in desired state")
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
		s.logger.Errorf("[CustSegSvc][EnableDisable] cps action err: %v", err)
		return err
	}

	return nil
}

func (s *customerSegmentationService) Delete(ctx context.Context, id string) error {
	makerUser := local_util.ExtractUserFromContext(ctx)

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("[CustSegSvc][Delete] find err: %v", err)
		return err
	}

	updated := *existing
	updated.IsDeleted = true
	cpsActionData := lib.CpsModelBuilder(id, makerUser, core.MapCustomerSegmentationToMap(*existing), core.MapCustomerSegmentationToMap(updated), string(constants.RequestDeleteCustomerSegmentation), constants.DELETE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		s.logger.Errorf("[CustSegSvc][Delete] cps action err: %v", err)
		return err
	}

	return nil
}

func (s *customerSegmentationService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	s.logger.Infof("[CustSegSvc][Authorize] action: %s", action.RequestAction)

	var (
		err error
		seg *imodel.CustomerSegmentation
	)

	seg, marshal_err := local_util.JsonUnmarshal[imodel.CustomerSegmentation](action.CurrentAction)
	if marshal_err != nil || seg == nil {
		s.logger.Errorf("[CustSegSvc][Authorize] unmarshal err: %v", marshal_err)
		return nil, marshal_err
	}

	switch action.RequestAction {
	case string(constants.RequestCreateCustomerSegmentation):
		err = s.repo.Create(ctx, seg)
		if err != nil {
			s.logger.Errorf("[CustSegSvc][Authorize] create err: %v", err)
			return nil, err
		}
		s.logger.Infof("[CustSegSvc][Authorize] created")
	case string(constants.RequestUpdateCustomerSegmentation):
		err = s.repo.Update(ctx, action.UniqueId, seg)
		if err != nil {
			s.logger.Errorf("[CustSegSvc][Authorize] update err: %v", err)
			return nil, err
		}
		s.logger.Infof("[CustSegSvc][Authorize] updated")
	case string(constants.RequestDeleteCustomerSegmentation):
		err = s.repo.Delete(ctx, action.UniqueId)
		if err != nil {
			s.logger.Errorf("[CustSegSvc][Authorize] delete err: %v", err)
			return nil, err
		}
		s.logger.Infof("[CustSegSvc][Authorize] deleted")
	case string(constants.RequestEnableCustomerSegmentation):
		err = s.repo.EnableOrDisable(ctx, action.UniqueId, true)
		if err != nil {
			s.logger.Errorf("[CustSegSvc][Authorize] enable err: %v", err)
			return nil, err
		}
		s.logger.Infof("[CustSegSvc][Authorize] enabled")
	case string(constants.RequestDisableCustomerSegmentation):
		err = s.repo.EnableOrDisable(ctx, action.UniqueId, false)
		if err != nil {
			s.logger.Errorf("[CustSegSvc][Authorize] disable err: %v", err)
			return nil, err
		}
		s.logger.Infof("[CustSegSvc][Authorize] disabled")
	default:
		s.logger.Errorf("[CustSegSvc][Authorize] unsupported action: %s", action.RequestAction)
		return nil, errors.New("unsupported action")
	}
	return action, nil
}
