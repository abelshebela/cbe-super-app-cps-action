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
	local_util "cbe-super-app-cps-action/pkgs/utils"

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

func (s *customerSegmentationService) CreateBulk(ctx context.Context, reqs []cust_seg.CreateCustomerSegmentationRequest) error {
	makerUser := local_util.ExtractUserFromContext(ctx)

	var newData []model.CustomerSegmentation
	var notFoundRoles []string
	for _, req := range reqs {
		role, err := s.cpsRoleRepo.FindById(ctx, req.CustomerRole)
		if err != nil {
			if err.Error() == localization.ErrorResourceNotFound.Code {
				s.logger.Errorf("CPS role {%s} cannot be found", req.CustomerRole)
				notFoundRoles = append(notFoundRoles, req.CustomerRole)
				continue
			} else {
				s.logger.Errorf("error while fetching CPS role {%s}: %v", req.CustomerRole, err)
				return err
			}
		}

		if role != nil {
			seg := model.CustomerSegmentation{
				CustomerRole:       req.CustomerRole,
				CustomerSegment:    req.CustomerSegment,
				CustomerSubSegment: req.CustomerSubSegment,
				CustomerGroup:      req.CustomerGroup,
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			}
			newData = append(newData, seg)
		}
	}

	if len(notFoundRoles) > 0 {
		return fmt.Errorf("these roles cannot be found: %s", notFoundRoles)

	}

	cpsActionData := lib.CpsModelBuilder("", makerUser, nil, newData, string(constants.RequestCreateCustomerSegmentation), constants.CREATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		s.logger.Errorf("[CreateBulk] failed to create CPS action %d: %v", err)
		return err
	}
	s.logger.Infof("[CreateBulk] customer segmentation creation request sent successfully")

	return nil
}

func (s *customerSegmentationService) Update(ctx context.Context, id string, req cust_seg.UpdateCustomerSegmentationRequest) error {
	makerUser := local_util.ExtractUserFromContext(ctx)

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("[Update] failed to find existing customer segmentation: %v", err)
		return err
	}

	updated := *existing

	if req.CustomerRole != nil {
		updated.CustomerRole = *req.CustomerRole
	}
	if req.CustomerSegment != nil {
		updated.CustomerSegment = *req.CustomerSegment
	}
	if req.CustomerSubSegment != nil {
		updated.CustomerSubSegment = *req.CustomerSubSegment
	}
	if req.CustomerGroup != nil {
		updated.CustomerGroup = *req.CustomerGroup
	}
	updated.UpdatedAt = time.Now()

	cpsActionData := lib.CpsModelBuilder(id, makerUser, existing, updated, string(constants.RequestUpdateCustomerSegmentation), constants.UPDATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		s.logger.Errorf("[Update] failed to create CPS action: %v", err)
		return err
	}

	s.logger.Infof("[Update] customer segmentation update request created successfully")
	return nil
}

func (s *customerSegmentationService) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]*model.CustomerSegmentation], error) {
	return s.repo.FindAllWithPagination(ctx, *filterParam)
}

func (s *customerSegmentationService) FindById(ctx context.Context, id string) (*model.CustomerSegmentation, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *customerSegmentationService) Delete(ctx context.Context, id string) error {
	makerUser := local_util.ExtractUserFromContext(ctx)

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("[Delete] failed to find existing customer segmentation: %v", err)
		return err
	}

	updated := *existing
	updated.IsDeleted = true
	cpsActionData := lib.CpsModelBuilder(id, makerUser, existing, updated, string(constants.RequestDeleteCustomerSegmentation), constants.DELETE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		s.logger.Errorf("[Delete] failed to create CPS action: %v", err)
		return err
	}

	return nil
}

func (s *customerSegmentationService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	s.logger.Infof("[Authorize] authorizing customer segmentation action: %s", action.RequestAction)

	var (
		err error
		seg *model.CustomerSegmentation
	)
	if action.ActionType != constants.CREATE {
		seg, marshal_err := local_util.JsonUnmarshal[model.CustomerSegmentation](action.CurrentAction)
		if marshal_err != nil || seg == nil {
			s.logger.Errorf("[Authorize] failed to unmarshal current action: %v", marshal_err)
			return nil, marshal_err
		}
	}

	switch action.RequestAction {
	case string(constants.RequestCreateCustomerSegmentation):
		cusSegs, marshal_err := local_util.JsonUnmarshal[[]model.CustomerSegmentation](action.CurrentAction)
		if marshal_err != nil || cusSegs == nil {
			s.logger.Errorf("[Authorize] failed to unmarshal current action: %v", marshal_err)
			return nil, marshal_err
		}

		for _, action := range *cusSegs {
			err = s.repo.Create(ctx, &action)
			if err != nil {
				s.logger.Errorf("[Authorize] failed to create customer segmentation: %v", err)
				return nil, err
			}
			s.logger.Infof("[Authorize] customer segmentation created successfully")
		}
	case string(constants.RequestUpdateCustomerSegmentation):
		err = s.repo.Update(ctx, action.UniqueId, seg)
		if err != nil {
			s.logger.Errorf("[Authorize] failed to update customer segmentation: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] customer segmentation updated successfully")
	case string(constants.RequestDeleteCustomerSegmentation):
		err = s.repo.Delete(ctx, action.UniqueId)
		if err != nil {
			s.logger.Errorf("[Authorize] failed to delete customer segmentation: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] customer segmentation deleted successfully")
	default:
		s.logger.Errorf("[Authorize] unsupported action: %s", action.RequestAction)
		return nil, errors.New("unsupported action")
	}
	return action, nil
}
