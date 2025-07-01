package service

import (
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"cbe-super-app-cps-action/internal/domain/action"
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

type ServiceRepository interface {
	GetServiceDetailsByID(ctx context.Context, id string) (ServiceDetails, error)
	UpdateServiceDetails(ctx context.Context, id string, update ServiceDetails) error
	FetchPendingActionsByUniqueID(ctx context.Context, uniqueID string) ([]action.CPSAction, error)
	UpdateCapMinAmount(ctx context.Context, id string, minAmount uint64) error
	ApproveServiceDetails(ctx context.Context, actionID string, approve bool, checkerID string, rejectionReason string) error
}

type ServiceInterface interface {
	GetAllServiceDetails(ctx context.Context) ([]*Service, error)
	GetServiceDetailsByID(ctx context.Context, id string) (*Service, error)
	UpdateServiceDetailsRequest(ctx context.Context, id string, update *Service, makerID string) (string, error)
	UpdateServiceDetails(ctx context.Context, actionID string, approve bool, checkerID string, rejectionReason string) error
	UpdateServiceCap(ctx context.Context, id string, cap Cap) (Service, error)

	UpdateCapMinAmount(ctx context.Context, id string, minAmount uint64) error
	ApproveServiceDetails(ctx context.Context, actionID string, approve bool, checkerID string, rejectionReason string) error

	InitiateServiceFeeUpdate(ctx context.Context, cpsAction CPSAction) (UpdateServiceDetailsResponse, error)
	ApproveServiceFeeUpdate(ctx context.Context, cpsAction CPSAction) (UpdateServiceDetailsResponse, error)
	RejectServiceFeeUpdate(ctx context.Context, cpsAction CPSAction) (UpdateServiceDetailsResponse, error)
}

type ServiceStore struct {
	repository Repository
	actionRepo action.Repository
	logger     utils.Logger
}

func NewServiceStore(repo Repository, actionRepo action.Repository, logger utils.Logger) ServiceInterface {
	return &ServiceStore{
		repository: repo,
		actionRepo: actionRepo,
		logger:     logger,
	}
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

	//actionID := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})
	actionID := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})
	cpsAction := action.CPSAction{
		ActionCode: actionID,
		Maker: action.User{
			UserID:      makerID,
			FullName:    "",
			PhoneNumber: "",
			Timestamp:   time.Now(),
		},
		Checker:        action.User{},
		UniqueId:       id,
		Department:     update.ServiceCode,
		ActionType:     action.ActionUpdate,
		RequestAction:  action.RequestUpdateServiceDetails,
		ActionStatus:   action.ActionPending,
		CurrentAction:  update,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
	}

	createdAction, err := s.actionRepo.CreateCpsAction(ctx, cpsAction)
	if err != nil {
		s.logger.Errorf("failed to create CPS action: %v", err)
		return "", errors.New("failed to create CPS action")
	}

	return createdAction.ActionCode, nil
}

func (s *ServiceStore) UpdateServiceDetails(ctx context.Context, actionID string, approve bool, checkerID string, rejectionReason string) error {
	cpsAction, err := s.actionRepo.FetchCpsActionById(ctx, actionID)
	if err != nil {
		s.logger.Errorf("failed to fetch CPS action: %v", err)
		return errors.New("failed to fetch CPS action")
	}

	if cpsAction.ActionStatus != action.ActionPending {
		return errors.New("action is not in pending status")
	}
	cpsAction.Checker = action.User{
		UserID:      checkerID,
		FullName:    "",
		PhoneNumber: "",
		Timestamp:   time.Now(),
	}
	cpsAction.LastModifiedAt = time.Now()
	if approve {
		cpsAction.ActionStatus = action.ActionApproved

		serviceData, ok := cpsAction.CurrentAction.(*Service)
		if !ok {
			return errors.New("invalid service data in action")
		}

		if err := s.repository.UpdateOneServiceDetailRequest(ctx, cpsAction.UniqueId, *serviceData); err != nil {
			s.logger.Errorf("failed to update service details: %v", err)
			return errors.New("failed to update service details")
		}
	} else {
		cpsAction.ActionStatus = action.ActionRejected
		cpsAction.LastModifiedAt = time.Now()
		if rejectionReason != "" {
			cpsAction.RejectionReason = makeStringPointer(rejectionReason)
		} else {
			cpsAction.RejectionReason = makeStringPointer("Rejected by checker")
		}
	}

	if err := s.actionRepo.UpdateCpsAction(ctx, cpsAction); err != nil {
		s.logger.Errorf("failed to update CPS action: %v", err)
		return errors.New("failed to update CPS action")
	}

	return nil
}

func (s *ServiceStore) UpdateServiceCap(ctx context.Context, id string, cap Cap) (Service, error) {
	svc, err := s.repository.GetOneServiceDetail(ctx, id)
	if err != nil {
		return Service{}, err
	}
	svc.Cap = cap
	return svc, s.repository.UpdateOneServiceDetailRequest(ctx, id, svc)
}
func (s *ServiceStore) InitiateServiceFeeUpdate(ctx context.Context, cpsAction CPSAction) (UpdateServiceDetailsResponse, error) {
	svc, err := s.repository.InitiateServiceFeeUpdate(ctx, cpsAction)
	if err != nil {
		var emptyResp UpdateServiceDetailsResponse
		return emptyResp, err
	}
	return svc, nil
}

func (s *ServiceStore) ApproveServiceFeeUpdate(ctx context.Context, cpsAction CPSAction) (UpdateServiceDetailsResponse, error) {
	svc, err := s.repository.ApproveServiceFeeUpdate(ctx, cpsAction)
	if err != nil {
		var emptyResp UpdateServiceDetailsResponse
		return emptyResp, err
	}
	return svc, nil
}
func (s *ServiceStore) RejectServiceFeeUpdate(ctx context.Context, cpsAction CPSAction) (UpdateServiceDetailsResponse, error) {
	svc, err := s.repository.RejectServiceFeeUpdate(ctx, cpsAction)
	if err != nil {
		var emptyResp UpdateServiceDetailsResponse
		return emptyResp, err
	}
	return svc, nil
}

func validateService(service Service) error {
	if service.ServiceCode == "" {
		return errors.New("service code cannot be empty")
	}
	if service.ServiceName == "" {
		return errors.New("service name cannot be empty")
	}
	if service.ServiceType == "" {
		return errors.New("service type cannot be empty")
	}
	return nil
}
func (s *ServiceStore) UpdateCapMinAmount(ctx context.Context, id string, minAmount uint64) error {
	_, err := s.repository.GetOneServiceDetail(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to fetch service: %v", err)
		return errors.New("service not found")
	}

	if err := s.repository.UpdateCapMinAmount(ctx, id, minAmount); err != nil {
		s.logger.Errorf("failed to update cap min amount: %v", err)
		return errors.New("failed to update cap min amount")
	}
	return nil
}

func (s *ServiceStore) ApproveServiceDetails(ctx context.Context, actionID string, approve bool, checkerID string, rejectionReason string) error {
	cpsAction, err := s.actionRepo.FetchCpsActionById(ctx, actionID)
	if err != nil {
		s.logger.Errorf("failed to fetch CPS action: %v", err)
		return errors.New("failed to fetch CPS action")
	}

	if cpsAction.ActionStatus != action.ActionPending {
		return errors.New("action is not in pending status")
	}
	cpsAction.Checker = action.User{
		UserID:      checkerID,
		FullName:    "",
		PhoneNumber: "",
		Timestamp:   time.Now(),
	}
	cpsAction.LastModifiedAt = time.Now()
	if approve {
		cpsAction.ActionStatus = action.ActionApproved

		serviceData, ok := cpsAction.CurrentAction.(*Service)
		if !ok {
			return errors.New("invalid service data in action")
		}

		if err := s.repository.UpdateOneServiceDetailRequest(ctx, cpsAction.UniqueId, *serviceData); err != nil {
			s.logger.Errorf("failed to update service details: %v", err)
			return errors.New("failed to update service details")
		}
	} else {
		cpsAction.ActionStatus = action.ActionRejected
		cpsAction.LastModifiedAt = time.Now()
		if rejectionReason != "" {
			cpsAction.RejectionReason = makeStringPointer(rejectionReason)
		} else {
			cpsAction.RejectionReason = makeStringPointer("Rejected by checker")
		}
	}

	if err := s.actionRepo.UpdateCpsAction(ctx, cpsAction); err != nil {
		s.logger.Errorf("failed to update CPS action: %v", err)
		return errors.New("failed to update CPS action")
	}

	return nil
}

func makeStringPointer(s string) *string {
	return &s
}
