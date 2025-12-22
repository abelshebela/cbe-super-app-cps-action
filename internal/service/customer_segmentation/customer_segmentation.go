package customersegmentation

import (
	"cbe-super-app-cps-action/internal/constants"
	cust_seg "cbe-super-app-cps-action/internal/constants/dto/customer_segmentation"
	"cbe-super-app-cps-action/internal/constants/lib"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"time"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type customerSegmentationService struct {
	repo       storage.CustomerSegmentationRepository
	cpsService service.CPSActionService
	logger     utils.Logger
}

func NewCustomerSegmentation(repo storage.CustomerSegmentationRepository, cpsService service.CPSActionService, logger utils.Logger) *customerSegmentationService {
	return &customerSegmentationService{repo: repo, cpsService: cpsService, logger: logger}
}

func (s *customerSegmentationService) Create(ctx context.Context, req cust_seg.CreateCustomerSegmentationRequest) error {
	makerUser := local_util.ExtractUserFromContext(ctx)

	seg := &imodel.CustomerSegmentation{
		CustomerRole:       req.CustomerRole,
		CustomerSegment:    req.CustomerSegment,
		CustomerSubSegment: req.CustomerSubSegment,
		CustomerGroup:      req.CustomerGroup,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	cpsActionData := lib.CpsModelBuilder("", makerUser, nil, seg, string(constants.RequestCreateCustomerSegmentation), constants.CREATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		s.logger.Errorf("[Create] failed to create CPS action: %v", err)
		return err
	}

	s.logger.Infof("[Create] customer segmentation creation request created successfully")
	return nil
}

func (s *customerSegmentationService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	s.logger.Infof("[Authorize] authorizing customer segmentation action: %s", action.RequestAction)
	var err error
	seg, marshal_err := local_util.JsonUnmarshal[imodel.CustomerSegmentation](action.CurrentAction)
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
	default:
		s.logger.Errorf("[Authorize] unsupported action: %s", action.RequestAction)
		return nil, errors.New("unsupported action")
	}
	return action, nil
}
