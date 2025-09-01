package amount_based_auth

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	cpsaction "cbe-super-app-cps-action/internal/service/cps_action"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	// added for cascading logic
	core "cbe-super-app-cps-action/internal/service/amount_based_auth/core"

	amountauthdto "cbe-super-app-cps-action/internal/constants/dto/amount_based_auth"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type amountBasedAuthService struct {
	Repository storage.AmountBasedAuthRepository 
	cpsService service.CPSActionService
	logger utils.Logger
	minioClient config.MinioClientInterface
	bucketName string
	cfg *config.VaultConfig
}
	
func NewAmountBasedAuthService(repository storage.AmountBasedAuthRepository, cpsService service.CPSActionService, minioClient config.MinioClientInterface, bucketName string, cfg *config.VaultConfig, logger utils.Logger) service.AmountBasedAuthService {
	return &amountBasedAuthService{
		Repository: repository,
		cpsService: cpsService,
		logger: logger,
		minioClient: minioClient,
		bucketName: bucketName,
		cfg: cfg,
	}
}

// Authorize handles persistence for amount-based auth actions
func (s *amountBasedAuthService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	s.logger.Infof("Authorizing amount-based auth action, action: %s", action.RequestAction)

	var authTier *model.AuthTier
	if err := local_util.BindAction(action.CurrentAction, &authTier); err != nil {
		s.logger.Errorf("Failed to bind current action to auth tier: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	var err error	
	switch action.RequestAction {
	case string(cpsaction.RequestUpdateAmountBasedAuth):
		authTier.LastModified = time.Now()
		err = s.Repository.Update(ctx, authTier.ID.Hex(), authTier)
		
	default:
		s.logger.Errorf("Unsupported action requested, action: %s", action.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	if err != nil {
		s.logger.Errorf("Failed to process amount-based auth action, action: %s, error: %v", action.RequestAction, err)
		return nil, err
	}

	action.CurrentAction = authTier
	s.logger.Infof("Authorization completed for action, action: %s, id: %s", action.RequestAction, authTier.ID)
	return action, nil
}

// FindAllWithPagination retrieves all amount-based auth tiers with pagination
func (s *amountBasedAuthService) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.AuthTier], error) {
	return s.Repository.FindAllWithPagination(ctx, filterParam)
}

// UpdateOpenTier updates OPEN tier's max and cascades changes to PIN min
func (s *amountBasedAuthService) UpdateOpenTier(ctx context.Context, id string, request amountauthdto.UpdateOpenTierRequest) error {
	if request.MaxAmount == 0 {
		return errors.New(localization.ErrorInvalidAmounts.Code)
	}

	existingTier, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingTier.Method != constants.OPEN {
		return errors.New(localization.ErrorInvalidMethod.Code)
	}

	existingTier.MaxAmount = request.MaxAmount

	pinTiers, err := s.Repository.FindAll(ctx, bson.M{"method": constants.PIN, "is_deleted": false}, bson.M{})
	if err != nil {
		return err
	}
	if len(pinTiers) == 0 {
		return errors.New(localization.ErrorFileNotFound.Code)
	}
	pinTier := pinTiers[0]

	if err := core.ApplyOpenUpdate(existingTier, pinTier); err != nil {
		return err
	}

	now := time.Now()
	existingTier.LastModified = now
	if err := s.Repository.Update(ctx, id, existingTier); err != nil {
		return err
	}
	pinTier.LastModified = now
	if err := s.Repository.Update(ctx, pinTier.ID.Hex(), pinTier); err != nil {
		return err
	}
	return nil
}

// UpdatePinTier updates PIN tier's min and max and validates against OPEN and OTP_PIN constraints
func (s *amountBasedAuthService) UpdatePinTier(ctx context.Context, id string, request amountauthdto.UpdatePinTierRequest) error {
	if request.MinAmount == 0 || request.MaxAmount == 0 {
		return errors.New(localization.ErrorInvalidAmounts.Code)
	}

	existingTier, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingTier.Method != constants.PIN {
		return errors.New(localization.ErrorInvalidMethod.Code)
	}

	// Set the user's requested values
	existingTier.MinAmount = request.MinAmount
	existingTier.MaxAmount = request.MaxAmount

	// Fetch OPEN and OTP_PIN tiers for validation constraints
	openTiers, err := s.Repository.FindAll(ctx, bson.M{"method": constants.OPEN, "is_deleted": false}, bson.M{})
	if err != nil {
		return err
	}
	if len(openTiers) == 0 {
		return errors.New(localization.ErrorFileNotFound.Code)
	}
	otpPinTiers, err := s.Repository.FindAll(ctx, bson.M{"method": constants.OTPANDPIN, "is_deleted": false}, bson.M{})
	if err != nil {
		return err
	}
	if len(otpPinTiers) == 0 {
		return errors.New(localization.ErrorFileNotFound.Code)
	}

	openTier := openTiers[0]
	otpPinTier := otpPinTiers[0]

	// Validate the user's PIN values against OPEN and OTP_PIN constraints
	if err := core.ApplyPinUpdate(existingTier, openTier, otpPinTier); err != nil {
		return err
	}

	existingTier.LastModified = time.Now()
	return s.Repository.Update(ctx, id, existingTier)
}

// UpdateOtpPinTier updates OTP_PIN tier's min and cascades to PIN's max
func (s *amountBasedAuthService) UpdateOtpPinTier(ctx context.Context, id string, request amountauthdto.UpdateOtpPinTierRequest) error {
	if request.MinAmount == 0 {
		return errors.New(localization.ErrorInvalidAmounts.Code)
	}

	existingTier, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingTier.Method != constants.OTPANDPIN {
		return errors.New(localization.ErrorInvalidMethod.Code)
	}

	existingTier.MinAmount = request.MinAmount

	pinTiers, err := s.Repository.FindAll(ctx, bson.M{"method": constants.PIN, "is_deleted": false}, bson.M{})
	if err != nil {
		return err
	}
	if len(pinTiers) == 0 {
		return errors.New(localization.ErrorFileNotFound.Code)
	}
	pinTier := pinTiers[0]

	if err := core.ApplyOtpPinUpdate(existingTier, pinTier); err != nil {
		return err
	}

	now := time.Now()
	existingTier.LastModified = now
	if err := s.Repository.Update(ctx, id, existingTier); err != nil {
		return err
	}
	pinTier.LastModified = now
	if err := s.Repository.Update(ctx, pinTier.ID.Hex(), pinTier); err != nil {
		return err
	}
	return nil
}




