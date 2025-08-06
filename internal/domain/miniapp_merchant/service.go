package miniappmerchant

import (
	"context"
	"fmt"
	"time"

	cps_constants "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/entity"
	user_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/service"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type MiniAppMerchantService interface {
	CreateMiniAppMerchant(ctx context.Context, data *MiniAppMerchantRequest) (*MiniAppMerchant, error)
	UpdateMiniAppMerchant(ctx context.Context, id string, data *MiniAppMerchantRequest) (*MiniAppMerchant, *MiniAppMerchant, error)
	ListMiniAppMerchant(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*MiniAppMerchant], error)
	DetailMiniAppByID(ctx context.Context, id string) (*MiniAppMerchant, error)
	DeleteMiniAppMerchant(ctx context.Context, id string) (*MiniAppMerchant, *MiniAppMerchant, error)
	EnableOrDisableMerchant(ctx context.Context, id string, enable bool) (*MiniAppMerchant, *MiniAppMerchant, error)
	Authorize(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	AddMiniApp(ctx context.Context, merchantID string, miniApp MiniApps) error
	UpdateMiniAppEnabledState(ctx context.Context, merchantID string, miniAppID string, enabled bool) error
	SoftDeleteMiniApp(ctx context.Context, merchantID string, miniAppID string) error
}
type MiniAppMerchantServiceImpl struct {
	repo        MiniAppMerchantRepository
	logger      shared_utils.Logger
	userService user_service.CustomerService
}

func NewMiniAppMerchantService(repo MiniAppMerchantRepository, userService user_service.CustomerService, logger shared_utils.Logger) MiniAppMerchantService {
	return &MiniAppMerchantServiceImpl{
		logger:      logger,
		repo:        repo,
		userService: userService,
	}
}

func (s *MiniAppMerchantServiceImpl) CreateMiniAppMerchant(ctx context.Context, data *MiniAppMerchantRequest) (*MiniAppMerchant, error) {
	exist, err := s.repo.MiniAppMerchantInfoExists(ctx, CheckMiniAppMerchant{
		BankAccountNumber: data.AccountNumber,
		Email:             data.Email,
		PhoneNumber:       data.PhoneNumber,
	}, nil)

	if err != nil {
		return nil, err
	}

	if exist {
		return nil, fmt.Errorf(common_util.InformationAlreadyExistst)
	}

	now := time.Now()
	res := &MiniAppMerchant{
		Code:              shared_utils.RandomGenerator(10),
		MerchantName:      data.MerchantName,
		MerchantType:      data.Type,
		PhoneNumber:       data.PhoneNumber,
		Email:             data.Email,
		BankAccountNumber: data.AccountNumber,
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
		MiniApps: []MiniApps{},
	}

	return res, nil
}

func (s *MiniAppMerchantServiceImpl) UpdateMiniAppMerchant(ctx context.Context, id string, data *MiniAppMerchantRequest) (*MiniAppMerchant, *MiniAppMerchant, error) {

	old, err := s.repo.DetailMiniAppByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	exist, err := s.repo.MiniAppMerchantInfoExists(ctx,
		CheckMiniAppMerchant{
			BankAccountNumber: data.AccountNumber,
			Email:             data.Email,
			PhoneNumber:       data.PhoneNumber,
		}, &MiniAppMerchantExistOptions{
			ExcludeID: id,
		})

	if err != nil {
		return nil, nil, err
	}

	if exist {
		return nil, nil, fmt.Errorf(common_util.InformationAlreadyExistst)
	}

	now := time.Now()
	cps := &MiniAppMerchant{
		ID:                old.ID,
		Code:              old.Code,
		MerchantName:      nonEmptyString(data.MerchantName, old.MerchantName),
		MerchantType:      nonEmptyString(data.Type, old.MerchantType),
		PhoneNumber:       nonEmptyString(data.PhoneNumber, old.PhoneNumber),
		Email:             nonEmptyString(data.Email, old.Email),
		BankAccountNumber: nonEmptyString(data.AccountNumber, old.BankAccountNumber),
		Enabled:           old.Enabled,
		IsDeleted:         old.IsDeleted,
		CreatedAt:         old.CreatedAt,
		LastModifiedAt:    now,
		KYC: KYC{
			Status: old.KYC.Status,
			Representative: KYCInformation{
				Name:  nonEmptyString(data.MerchantRepresentativeName, old.KYC.Representative.Name),
				Email: nonEmptyString(data.Email, old.KYC.Representative.Email),
				Phone: nonEmptyString(data.PhoneNumber, old.KYC.Representative.Phone),
			},
		},
		Branches: old.Branches,
	}

	return cps, old, nil
}

func (s *MiniAppMerchantServiceImpl) DeleteMiniAppMerchant(ctx context.Context, id string) (*MiniAppMerchant, *MiniAppMerchant, error) {

	old, err := s.repo.DetailMiniAppByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	now := time.Now()
	cur := *old

	cur.DeletedAt = now
	cur.LastModifiedAt = now
	cur.IsDeleted = true
	return &cur, old, nil
}

func (s *MiniAppMerchantServiceImpl) EnableOrDisableMerchant(ctx context.Context, id string, enabled bool) (*MiniAppMerchant, *MiniAppMerchant, error) {
	existingMerchant, err := s.repo.DetailMiniAppByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	if existingMerchant.Enabled && enabled {
		return nil, existingMerchant, fmt.Errorf(common_util.ErrAlreadyEnabled)
	}
	if !existingMerchant.Enabled && !enabled {
		return nil, existingMerchant, fmt.Errorf(common_util.ErrAlreadyDisabled)
	}

	now := time.Now()
	cur := *existingMerchant
	cur.Enabled = enabled
	cur.LastModifiedAt = now

	return &cur, existingMerchant, nil
}

func (s *MiniAppMerchantServiceImpl) Authorize(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {

	var currentAction MiniAppMerchant

	bindErr := common_util.BindAction(cpsAction.CurrentAction, &currentAction)
	if bindErr != nil {
		s.logger.Errorf("failed to bind current action to MiniApp: %v", bindErr)
		return nil, fmt.Errorf(common_util.InvalidActionData)
	}

	var (
		newMerchant *MiniAppMerchant
		err         error
	)

	switch cpsAction.RequestAction {
	case cps_constants.RequestCreateMiniAppMerchant:
		txCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		err = s.repo.RunInTransaction(txCtx, func(ctx context.Context) error {
			newMerchant, err = s.repo.CreateMiniAppMerchant(ctx, &currentAction)

			if err != nil {
				return err
			}

			_, err = s.userService.CreateUserMiniAppMerchant(ctx, &entity.User{
				FullName:    newMerchant.KYC.Representative.Name,
				PhoneNumber: newMerchant.KYC.Representative.Phone,
				Email:       newMerchant.KYC.Representative.Email,
			})

			if err != nil && err.Error() != common_util.AuthUserAlreadyExists {
				s.logger.Errorf("error occured on create minin app merchant %s", err.Error())
				return err
			}

			return nil
		})

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

func (s *MiniAppMerchantServiceImpl) AddMiniApp(ctx context.Context, merchantID string, miniApp MiniApps) error {
	return s.repo.AddMiniApp(ctx, merchantID, miniApp)
}
func (s *MiniAppMerchantServiceImpl) UpdateMiniAppEnabledState(ctx context.Context, merchantID string, miniAppID string, enabled bool) error {
	return s.repo.UpdateMiniAppEnabledState(ctx, merchantID, miniAppID, enabled)
}
func (s *MiniAppMerchantServiceImpl) SoftDeleteMiniApp(ctx context.Context, merchantID string, miniAppID string) error {
	return s.repo.SoftDeleteMiniApp(ctx, merchantID, miniAppID)
}
