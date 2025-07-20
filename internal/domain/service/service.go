package service

import (
	"context"
	"errors"
	"time"

	"encoding/json"
	"fmt"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type ServiceDetails struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	ServiceID   string    `json:"service_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ServiceInterface interface {
	GetAllServiceDetails(ctx context.Context) ([]*Service, error)
	GetServiceDetailsByID(ctx context.Context, id string) (*Service, error)
	UpdateServiceDetailsRequest(ctx context.Context, id string, update *Service, makerID string) (string, error)
	UpdateServiceDetails(ctx context.Context, actionID string, approve bool, checkerID string, rejectionReason string) error
	UpdateServiceCap(ctx context.Context, id string, cap Cap) (Service, error)
	UpdateCapMinAmount(ctx context.Context, id string, minAmount uint64) error
	InitiateServiceFeeUpdate(ctx context.Context, cpsAction CPSAction) (UpdateServiceDetailsResponse, error)
	ApproveServiceFeeUpdate(ctx context.Context, action_code string) error
	RejectServiceFeeUpdate(ctx context.Context, action_code string, rejection_reason string) error
}

type ServiceStore struct {
	repository ServiceRepository
	actionRepo action.ActionRepository
	logger     utils.Logger
}

func NewServiceStore(repo ServiceRepository, actionRepo action.ActionRepository, logger utils.Logger) ServiceInterface {
	return &ServiceStore{
		repository: repo,
		actionRepo: actionRepo,
		logger:     logger,
	}
}

func (s *ServiceStore) handleActionApproval(ctx context.Context, cpsAction *action.CPSAction,
	approve bool,
	checkerID, rejectionReason string,
	updateFunc func() error,
) error {
	if cpsAction.ActionStatus != action.ActionPending {
		return errors.New(common_util.ActionNotPending)
	}

	now := time.Now()

	cpsAction.CheckerID = checkerID
	cpsAction.CheckerActionTime = &now
	cpsAction.LastModifiedAt = time.Now()

	if approve {
		cpsAction.ActionStatus = action.ActionApproved
		if err := updateFunc(); err != nil {
			return err
		}
	} else {
		cpsAction.ActionStatus = action.ActionRejected
		reason := "Rejected by checker"
		if rejectionReason != "" {
			reason = rejectionReason
		}
		cpsAction.RejectionReason = makeStringPointer(reason)
	}

	if err := s.actionRepo.UpdateCpsAction(ctx, *cpsAction); err != nil {
		s.logger.Errorf("failed to update CPS action: %v", err)
		return errors.New(common_util.FailedToUpdateAction)
	}

	return nil
}

func (s *ServiceStore) handleServiceFeeUpdate(ctx context.Context, cpsAction CPSAction, updateFunc func(context.Context, CPSAction) (UpdateServiceDetailsResponse, error)) (UpdateServiceDetailsResponse, error) {
	resp, err := updateFunc(ctx, cpsAction)
	if err != nil {
		var emptyResp UpdateServiceDetailsResponse
		return emptyResp, err
	}
	return resp, nil
}

func (s *ServiceStore) InitiateServiceFeeUpdate(ctx context.Context, cpsAction CPSAction) (UpdateServiceDetailsResponse, error) {
	return s.handleServiceFeeUpdate(ctx, cpsAction, s.repository.InitiateServiceFeeUpdate)
}

func (s *ServiceStore) ApproveServiceFeeUpdate(ctx context.Context, action_code string) error {
	return s.repository.ApproveServiceFeeUpdate(ctx, action_code)
}

func (s *ServiceStore) RejectServiceFeeUpdate(ctx context.Context, action_code string, rejection_reason string) error {
	return s.repository.RejectServiceFeeUpdate(ctx, action_code, rejection_reason)
}

func (s *ServiceStore) GetAllServiceDetails(ctx context.Context) ([]*Service, error) {
	return s.repository.GetAllServiceDetails(ctx)
}

func (s *ServiceStore) GetServiceDetailsByID(ctx context.Context, id string) (*Service, error) {
	service, err := s.repository.GetOneServiceDetail(ctx, id)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (s *ServiceStore) UpdateServiceDetailsRequest(ctx context.Context, id string, update *Service, makerID string) (string, error) {
	if err := validateService(*update); err != nil {
		return "", err
	}

	actionID := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})
	cpsAction := action.CPSAction{
		ActionCode:      actionID,
		MakerID:         makerID,
		MakerActionTime: time.Now(),
		UniqueId:        id,
		Department:      update.ServiceCode,
		ActionType:      action.ActionUpdate,
		RequestAction:   action.RequestUpdateServiceDetails,
		ActionStatus:    action.ActionPending,
		CurrentAction:   update,
		CreatedAt:       time.Now(),
		LastModifiedAt:  time.Now(),
	}

	createdAction, err := s.actionRepo.CreateCpsAction(ctx, cpsAction)
	if err != nil {
		s.logger.Errorf("failed to create CPS action: %v", err)
		return "", errors.New(common_util.FailedToCreateAction)
	}

	return createdAction.ActionCode, nil
}

func (s *ServiceStore) UpdateServiceDetails(ctx context.Context, actionID string, approve bool, checkerID string, rejectionReason string) error {
	cpsAction, err := s.actionRepo.FetchCpsActionById(ctx, actionID)
	if err != nil {
		s.logger.Errorf("failed to fetch CPS action: %v", err)
		return errors.New(common_util.FailedToFetchAction)
	}

	return s.handleActionApproval(ctx, &cpsAction, approve, checkerID, rejectionReason, func() error {

		serviceData, err := ConvertToService(cpsAction.CurrentAction)
		if err != nil {
			return err
		}
		if err := s.repository.UpdateOneServiceDetailRequest(ctx, cpsAction.UniqueId, *serviceData); err != nil {
			s.logger.Errorf("failed to update service details: %v", err)
			return errors.New(common_util.FailedToUpdateService)
		}
		return nil
	})
}

func (s *ServiceStore) UpdateServiceCap(ctx context.Context, id string, cap Cap) (Service, error) {
	svc, err := s.repository.GetOneServiceDetail(ctx, id)
	if err != nil {
		return Service{}, err
	}
	svc.Cap = cap
	return svc, s.repository.UpdateOneServiceDetailRequest(ctx, id, svc)
}

func validateService(service Service) error {
	return service.Validate()
}

func (s *ServiceStore) UpdateCapMinAmount(ctx context.Context, id string, minAmount uint64) error {
	_, err := s.repository.GetOneServiceDetail(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to fetch service: %v", err)
		return errors.New(common_util.ServiceNotFound)
	}

	if err := s.repository.UpdateCapMinAmount(ctx, id, minAmount); err != nil {
		s.logger.Errorf("failed to update cap min amount: %v", err)
		return errors.New(common_util.FailedToUpdateCapMin)
	}
	return nil
}

func makeStringPointer(s string) *string {
	return &s
}

// ConvertToService attempts to convert any interface{} to a *Service struct.
func ConvertToService(data interface{}) (*Service, error) {
	switch v := data.(type) {
	case *Service:
		return v, nil
	case Service:
		return &v, nil
	case map[string]interface{}:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal current action: %w", err)
		}
		var svc Service
		if err := json.Unmarshal(b, &svc); err != nil {
			return nil, fmt.Errorf("failed to unmarshal current action: %w", err)
		}
		return &svc, nil
	default:
		return nil, errors.New(common_util.InvalidActionData)
	}
}
