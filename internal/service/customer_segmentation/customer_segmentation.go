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

	var newData model.CustomerSegmentation
	role, err := s.cpsRoleRepo.FindById(ctx, req.CustomerRole)
	if err != nil {
		if err.Error() == localization.ErrorResourceNotFound.Code {
			s.logger.Errorf("CPS role {%s} cannot be found", req.CustomerRole)
			return fmt.Errorf("This role not found")
		} else {
			s.logger.Errorf("error while fetching CPS role {%s}: %v", req.CustomerRole, err)
			return err
		}
	}

	if role != nil {
		newData.CustomerRole = model.CustomerRoleInfo{ID: role.ID, Name: role.Name}
		newData.CreatedAt = time.Now()
		newData.UpdatedAt = time.Now()
		for _, info := range req.CustomerSubSegments {
			seg := model.CustomerSubSegments{
				Name:            info.Name,
				CustomerGroup:   info.CustomerGroup,
				CustomerSegment: info.CustomerSegment,
			}
			newData.CustomerSubSegments = append(newData.CustomerSubSegments, seg)
		}
	}

	cpsActionData := lib.CpsModelBuilder("", makerUser, nil, core.MapCustomerSegmentationToMap(newData), string(constants.RequestCreateCustomerSegmentation), constants.CREATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		s.logger.Errorf("[Create] failed to create CPS action %d: %v", err)
		return err
	}
	s.logger.Infof("[Create] customer segmentation creation request sent successfully")

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

	if len(req.CustomerSubSegments) > 0 && len(req.OldName) > 0 {
		updates := make(map[string]model.CustomerSubSegments)

		for i, sub := range req.CustomerSubSegments {
			updates[req.OldName[i]] = model.CustomerSubSegments{
				Name:            sub.Name,
				CustomerGroup:   sub.CustomerGroup,
				CustomerSegment: sub.CustomerSegment,
			}
		}

		for i, existingSub := range updated.CustomerSubSegments {
			if newSub, ok := updates[existingSub.Name]; ok {
				updated.CustomerSubSegments[i] = newSub
			}
		}
	}

	updated.UpdatedAt = time.Now()

	cpsActionData := lib.CpsModelBuilder(id, makerUser, core.MapCustomerSegmentationToMap(*existing), core.MapCustomerSegmentationToMap(updated), string(constants.RequestUpdateCustomerSegmentation), constants.UPDATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		s.logger.Errorf("[Update] failed to create CPS action: %v", err)
		return err
	}

	s.logger.Infof("[Update] customer segmentation update request created successfully")
	return nil
}

func (s *customerSegmentationService) isSubSegmentsEqual(existing, incoming []model.CustomerSubSegments) bool {
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
	cpsActionData := lib.CpsModelBuilder(id, makerUser, core.MapCustomerSegmentationToMap(*existing), core.MapCustomerSegmentationToMap(updated), string(constants.RequestDeleteCustomerSegmentation), constants.DELETE)

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

	seg, marshal_err := local_util.JsonUnmarshal[model.CustomerSegmentation](action.CurrentAction)
	if marshal_err != nil || seg == nil {
		s.logger.Errorf("[Authorize] failed to unmarshal current action: %v", marshal_err)
		return nil, marshal_err
	}

	switch action.RequestAction {
	case string(constants.RequestCreateCustomerSegmentation):
		err = s.repo.Create(ctx, seg)
		if err != nil {
			s.logger.Errorf("[Authorize] failed to create customer segmentation: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] customer segmentation created successfully")
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
