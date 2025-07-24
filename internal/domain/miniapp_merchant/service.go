package miniappmerchant

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cps_constants "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type MiniAppMerchantService interface {
	CreateMiniAppMerchant(ctx context.Context, cpsAction *entities.CreateCPSAction) (*entities.CPSAction, error)
	UpdateMiniAppMerchant(ctx context.Context, cpsAction *entities.CreateCPSAction) (*entities.CPSAction, error)
	ListMiniAppMerchant(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*MiniAppMerchant], error)
	DetailMiniAppByID(ctx context.Context, id string) (*MiniAppMerchant, error)
	DeleteMiniAppMerchant(ctx context.Context, id string, req *entities.CreateCPSAction) (*entities.CPSAction, error)
	EnableOrDisableMerchant(ctx context.Context, id string, requestAction cps_constants.RequestAction, req *entities.CreateCPSAction) (*entities.CPSAction, error)
	Authorize(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
}
type MiniAppMerchantServiceImpl struct {
	repo   MiniAppMerchantRepository
	logger shared_utils.Logger
}

func NewMiniAppMerchantService(repo MiniAppMerchantRepository, logger shared_utils.Logger,
) MiniAppMerchantService {
	return &MiniAppMerchantServiceImpl{
		logger: logger,
		repo:   repo,
	}
}

func (s *MiniAppMerchantServiceImpl) getCurrentData(req *entities.CreateCPSAction) (*MiniAppMerchant, error) {

	// This works even if ActionData is map[string]interface{} or *MiniAppMerchant serialized from HTTP
	bytes, err := json.Marshal(req.ActionData)
	if err != nil {
		s.logger.Errorf("failed to marshal ActionData: %v", err)
		return nil, fmt.Errorf("MARSHAL_ERROR")
	}

	var merchant MiniAppMerchant
	if err := json.Unmarshal(bytes, &merchant); err != nil {
		s.logger.Errorf("failed to unmarshal ActionData into MiniAppMerchant: %v", err)
		return nil, fmt.Errorf("UNMARSHAL_ERROR")
	}

	return &merchant, nil
}

func (s *MiniAppMerchantServiceImpl) CreateMiniAppMerchant(ctx context.Context, req *entities.CreateCPSAction) (*entities.CPSAction, error) {
	data, err := s.getCurrentData(req)
	if err != nil {
		return nil, err
	}
	exist, err := s.repo.MiniAppMerchantInfoExists(ctx, CheckMiniAppMerchant{
		BankAccountNumber: data.BankAccountNumber,
		Email:             data.Email,
		PhoneNumber:       data.PhoneNumber,
	})

	if err != nil {
		return nil, err
	}

	if exist {
		return nil, fmt.Errorf(common_util.InformationAlreadyExistst)
	}

	now := time.Now()
	cps := &entities.CPSAction{
		ActionCode:       shared_utils.RandomGenerator(20),
		MakerID:          req.MakerUser.UserCode,
		MakerName:        req.MakerUser.FullName,
		MakerPhoneNumber: req.MakerUser.PhoneNumber,
		Department:       req.MakerUser.Department,
		CurrentAction: MiniAppMerchant{
			Code:              shared_utils.RandomGenerator(10), // merchantCode
			MerchantName:      data.MerchantName,
			MerchantType:      data.MerchantType,
			PhoneNumber:       data.PhoneNumber,
			Email:             data.Email,
			BankAccountNumber: data.BankAccountNumber,
			MiniAppIDs:        data.MiniAppIDs,
			Enabled:           false,
			IsDeleted:         false,
			CreatedAt:         now,
			LastModifiedAt:    now,
			KYC: KYC{
				Status: KYCStatusComplete,
				Representative: KYCInformation{
					Name:  data.MerchantName,
					Email: data.Email,
					Phone: data.PhoneNumber,
				},
			},
			Branches: []BranchInformation{},
		},
		RequestAction:   cps_constants.RequestCreateMiniAppMerchant,
		ActionStatus:    cps_constants.ActionPending,
		ActionType:      cps_constants.ActionCreate,
		MakerActionTime: now,
		CreatedAt:       now,
		LastModifiedAt:  now,
	}

	return cps, nil
}

func (s *MiniAppMerchantServiceImpl) UpdateMiniAppMerchant(ctx context.Context, req *entities.CreateCPSAction) (*entities.CPSAction, error) {
	data, err := s.getCurrentData(req)
	if err != nil {
		return nil, err
	}

	old, err := s.repo.DetailMiniAppByID(ctx, data.ID)
	if err != nil {
		return nil, err
	}

	exist, err := s.repo.MiniAppMerchantInfoExists(ctx, CheckMiniAppMerchant{
		BankAccountNumber: data.BankAccountNumber,
		Email:             data.Email,
		PhoneNumber:       data.PhoneNumber,
	})

	if err != nil {
		return nil, err
	}

	if exist {
		return nil, fmt.Errorf(common_util.InformationAlreadyExistst)
	}

	now := time.Now()
	cps := &entities.CPSAction{
		ActionCode:       shared_utils.RandomGenerator(20),
		MakerID:          req.MakerUser.UserCode,
		MakerName:        req.MakerUser.FullName,
		MakerPhoneNumber: req.MakerUser.PhoneNumber,
		Department:       req.MakerUser.Department,
		CurrentAction: MiniAppMerchant{
			ID:                old.ID,
			Code:              old.Code, // keep existing code
			MerchantName:      data.MerchantName,
			MerchantType:      data.MerchantType,
			PhoneNumber:       data.PhoneNumber,
			Email:             data.Email,
			BankAccountNumber: data.BankAccountNumber,
			MiniAppIDs:        data.MiniAppIDs,
			Enabled:           old.Enabled,
			IsDeleted:         old.IsDeleted,
			CreatedAt:         old.CreatedAt,
			LastModifiedAt:    now,
			KYC: KYC{
				Status: data.KYC.Status,
				Representative: KYCInformation{
					Name:  data.KYC.Representative.Name,
					Email: data.KYC.Representative.Email,
					Phone: data.KYC.Representative.Phone,
				},
			},
			Branches: data.Branches,
		},
		PreviousAction:  *old,
		ActionType:      cps_constants.ActionUpdate,
		ActionStatus:    cps_constants.ActionPending,
		RequestAction:   cps_constants.RequestUpdateMiniAppMerchant,
		MakerActionTime: now,
		CreatedAt:       now,
		LastModifiedAt:  now,
	}

	return cps, nil
}

func (s *MiniAppMerchantServiceImpl) DeleteMiniAppMerchant(ctx context.Context, id string, req *entities.CreateCPSAction) (*entities.CPSAction, error) {

	old, err := s.repo.DetailMiniAppByID(ctx, id)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	cps := &entities.CPSAction{
		ActionCode:       shared_utils.RandomGenerator(20),
		MakerID:          req.MakerUser.UserCode,
		MakerName:        req.MakerUser.FullName,
		MakerPhoneNumber: req.MakerUser.PhoneNumber,
		Department:       req.MakerUser.Department,
		CurrentAction: MiniAppMerchant{
			ID:             old.ID,
			IsDeleted:      true,
			LastModifiedAt: now,
			DeletedAt:      now,
		},
		PreviousAction:  *old,
		ActionType:      cps_constants.ActionDelete,
		ActionStatus:    cps_constants.ActionPending,
		RequestAction:   cps_constants.RequestDeleteMiniAppMerchant,
		MakerActionTime: now,
		CreatedAt:       now,
		LastModifiedAt:  now,
	}

	return cps, nil
}

func (s *MiniAppMerchantServiceImpl) EnableOrDisableMerchant(ctx context.Context, id string, requestAction cps_constants.RequestAction, req *entities.CreateCPSAction) (*entities.CPSAction, error) {
	existingMerchant, err := s.repo.DetailMiniAppByID(ctx, id)
	if err != nil {
		return nil, err
	}

	enable := requestAction == cps_constants.RequestEnableMiniAppMerchant
	if enable && existingMerchant.Enabled {
		return nil, fmt.Errorf(common_util.ErrAlreadyEnabled)
	}
	if !enable && !existingMerchant.Enabled {
		return nil, fmt.Errorf(common_util.ErrAlreadyDisabled)
	}

	now := time.Now()
	cpsAction := &entities.CPSAction{
		ActionCode:       shared_utils.RandomGenerator(20),
		MakerID:          req.MakerUser.UserCode,
		MakerName:        req.MakerUser.FullName,
		MakerPhoneNumber: req.MakerUser.PhoneNumber,
		Department:       req.MakerUser.Department,
		ActionStatus:     cps_constants.ActionPending,
		ActionType:       cps_constants.ActionUpdate,
		RequestAction:    requestAction,
		CurrentAction: MiniAppMerchant{
			ID:                existingMerchant.ID,
			Code:              existingMerchant.Code,
			MerchantName:      existingMerchant.MerchantName,
			MerchantType:      existingMerchant.MerchantType,
			PhoneNumber:       existingMerchant.PhoneNumber,
			Email:             existingMerchant.Email,
			BankAccountNumber: existingMerchant.BankAccountNumber,
			MiniAppIDs:        existingMerchant.MiniAppIDs,
			Enabled:           enable,
			IsDeleted:         existingMerchant.IsDeleted,
			CreatedAt:         existingMerchant.CreatedAt,
			LastModifiedAt:    now,
			DeletedAt:         existingMerchant.DeletedAt,
			KYC:               existingMerchant.KYC,
			Branches:          existingMerchant.Branches,
		},
		PreviousAction:  *existingMerchant,
		MakerActionTime: now,
		CreatedAt:       now,
		LastModifiedAt:  now,
	}

	return cpsAction, nil
}

func (s *MiniAppMerchantServiceImpl) Authorize(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {

	var currentAction MiniAppMerchant
	raw := cpsAction.CurrentAction
	tmpJSON, errr := json.Marshal(raw)
	if errr != nil {
		s.logger.Errorf("failed to marshal current action: %v", errr)
		return nil, fmt.Errorf(common_util.InvalidActionData)
	}

	errr = json.Unmarshal(tmpJSON, &currentAction)
	if errr != nil {
		s.logger.Errorf("failed to unmarshal current action to MiniAppMerchant: %v", errr)
		return nil, fmt.Errorf(common_util.InvalidActionData)
	}

	var (
		newMerchant *MiniAppMerchant
		err         error
	)

	switch cpsAction.RequestAction {
	case cps_constants.RequestCreateMiniAppMerchant:
		newMerchant, err = s.repo.CreateMiniAppMerchant(ctx, &currentAction)

	case cps_constants.RequestUpdateMiniAppMerchant:
		newMerchant, err = s.repo.UpdateMiniAppMerchant(ctx, &currentAction)

	case cps_constants.RequestDeleteMiniAppMerchant:
		newMerchant, err = s.repo.DeleteMiniAppMerchant(ctx, currentAction.ID)

	case cps_constants.RequestEnableMiniAppMerchant:
		newMerchant, err = s.repo.EnableMiniAppMerchant(ctx, currentAction.ID)

	case cps_constants.RequestDisableMiniAppMerchant:
		newMerchant, err = s.repo.DisableMiniAppMerchant(ctx, currentAction.ID)

	default:
		s.logger.Warnf("unsupported request action: %s", cpsAction.RequestAction)
		return nil, fmt.Errorf(common_util.ErrUnsupported)
	}

	if err != nil {
		return nil, err
	}

	cpsAction.CurrentAction = newMerchant
	return cpsAction, nil
}

func (s *MiniAppMerchantServiceImpl) ListMiniAppMerchant(ctx context.Context, filter *constant.Filter) (*common_util.PaginatedResponse[[]*MiniAppMerchant], error) {
	return s.repo.ListMiniAppMerchant(ctx, filter)
}

func (s *MiniAppMerchantServiceImpl) DetailMiniAppByID(ctx context.Context, id string) (*MiniAppMerchant, error) {
	return s.repo.DetailMiniAppByID(ctx, id)
}
