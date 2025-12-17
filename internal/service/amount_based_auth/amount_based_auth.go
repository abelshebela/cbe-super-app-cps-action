package amount_based_auth

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	cpsaction "cbe-super-app-cps-action/internal/service/cps_action"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

	shared_constant "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/constants"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	core "cbe-super-app-cps-action/internal/service/amount_based_auth/core"

	amountauthdto "cbe-super-app-cps-action/internal/constants/dto/amount_based_auth"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type amountBasedAuthService struct {
	Repository storage.AmountBasedAuthRepository
	cpsService service.CPSActionService
	logger     utils.Logger
	cfg        *config.VaultConfig
}

func NewAmountBasedAuthService(repository storage.AmountBasedAuthRepository, cpsService service.CPSActionService, cfg *config.VaultConfig, logger utils.Logger) service.AmountBasedAuthService {
	return &amountBasedAuthService{
		Repository: repository,
		cpsService: cpsService,
		logger:     logger,
		cfg:        cfg,
	}
}

// Authorize handles persistence for amount-based auth actions
func (s *amountBasedAuthService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Amount Based Auth", "Authorize")
	defer span.End()

	s.logger.Infof("[Authorize] authorizing amount-based auth action: %s", action.RequestAction)

	// Unmarshal the current action data
	// Try direct type assertion first

	currentAction, err := local_util.JsonUnmarshal[map[string]interface{}](action.CurrentAction)
	if err != nil {
		span.AddEvent("Failed to extract currentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", action.UniqueId),
		))
		s.logger.Errorf("Failed to extract currentAction from action.CurrentAction: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	s.logger.Infof("[Authorize] processing action: %s", action.RequestAction)

	result := (*currentAction)

	switch result["method"] {
	case "OPEN":
		s.logger.Infof("[Authorize] processing OPEN tier update")
		data, err := local_util.JsonUnmarshal[map[string]interface{}](result["data"])
		if err != nil {
			s.logger.Errorf("Error occurred when extracting data tier: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		openTierMap, ok := (*data)["open"].(map[string]interface{})
		if !ok {
			s.logger.Errorf("Error occurred when casting open tier to map: %v", (*data)["open"])
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		openTier, err := local_util.JsonUnmarshal[model.AuthTier](openTierMap)
		if err != nil {
			s.logger.Errorf("Error occurred when converting open tier map to struct: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		pinTierMap, ok := (*data)["pin"].(map[string]interface{})
		if !ok {
			s.logger.Errorf("Error occurred when casting pin tier to map: %v", (*data)["pin"])
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		pinTier, err := local_util.JsonUnmarshal[model.AuthTier](pinTierMap)
		if err != nil {
			s.logger.Errorf("Failed to unmarshal PIN tier data: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		if err := s.Repository.Update(ctx, (*data)["open_id"].(string), openTier); err != nil {
			s.logger.Errorf("Failed to update OPEN tier data: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		if err := s.Repository.Update(ctx, (*data)["pin_id"].(string), pinTier); err != nil {
			s.logger.Errorf("Failed to update PIN tier data: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		s.logger.Infof("[Authorize] OPEN tier update completed successfully")
	case "PIN":
		s.logger.Infof("[Authorize] processing PIN tier update")
		data, err := local_util.JsonUnmarshal[map[string]interface{}](result["data"])
		if err != nil {
			s.logger.Errorf("Error occcure when extracting data tier error: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		openTierMap, ok := (*data)["open"].(map[string]interface{})
		if !ok {
			s.logger.Errorf("Error occurred when casting open tier to map: %v", (*data)["open"])
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		openTier, err := local_util.JsonUnmarshal[model.AuthTier](openTierMap)
		if err != nil {
			s.logger.Errorf("Error occcure when extracting open tier error: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		pinTierMap, ok := (*data)["pin"].(map[string]interface{})
		if !ok {
			s.logger.Errorf("Error occurred when casting pin tier to map: %v", (*data)["pin"])
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		pinTier, err := local_util.JsonUnmarshal[model.AuthTier](pinTierMap)
		if err != nil {
			s.logger.Errorf("Failed to unmarshal PIN tier data: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		otpPinTierMap, ok := (*data)["otp_pin"].(map[string]interface{})
		if !ok {
			s.logger.Errorf("Error occurred when casting OTP_PIN tier to map: %v", (*data)["otp_pin"])
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		otpPinTier, err := local_util.JsonUnmarshal[model.AuthTier](otpPinTierMap)
		if err != nil {
			s.logger.Errorf("Failed to unmarshal OTP_PIN tier data: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		if err := s.Repository.Update(ctx, (*data)["open_id"].(string), openTier); err != nil {
			s.logger.Errorf("Failed to update OPEN tier data: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		if err := s.Repository.Update(ctx, (*data)["pin_id"].(string), pinTier); err != nil {
			s.logger.Errorf("Failed to update PIN tier data: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		if err := s.Repository.Update(ctx, (*data)["otp_pin_id"].(string), otpPinTier); err != nil {
			s.logger.Errorf("Failed to update OTP_PIN tier data: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}
	case "OTP_PIN":
		s.logger.Infof("[Authorize] processing OTP_PIN tier update")
		data, err := local_util.JsonUnmarshal[map[string]interface{}](result["data"])
		if err != nil {
			s.logger.Errorf("Error occcure when extracting data tier error: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		pinTierMap, ok := (*data)["pin"].(map[string]interface{})
		if !ok {
			s.logger.Errorf("Error occurred when casting pin tier to map: %v", (*data)["pin"])
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		pinTier, err := local_util.JsonUnmarshal[model.AuthTier](pinTierMap)
		if err != nil {
			s.logger.Errorf("Error occcure when extracting pin tier error: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		otpPinTierMap, ok := (*data)["otp_pin"].(map[string]interface{})
		if !ok {
			s.logger.Errorf("Error occurred when casting OTP_PIN tier to map: %v", (*data)["otp_pin"])
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		otpPinTier, err := local_util.JsonUnmarshal[model.AuthTier](otpPinTierMap)
		if err != nil {
			s.logger.Errorf("Failed to unmarshal OTP_PIN tier data: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		if err := s.Repository.Update(ctx, (*data)["otp_pin_id"].(string), otpPinTier); err != nil {
			s.logger.Errorf("Failed to update OTP_PIN tier data: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		if err := s.Repository.Update(ctx, (*data)["pin_id"].(string), pinTier); err != nil {
			s.logger.Errorf("Failed to update PIN tier data: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		if err := s.Repository.Update(ctx, (*data)["otp_pin_id"].(string), otpPinTier); err != nil {
			s.logger.Errorf("Failed to update OTP_PIN tier data: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}
	default:
		s.logger.Errorf("Unsupported method: %v", result["method"])
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	s.logger.Infof("[Authorize] authorization completed successfully for action: %s", action.RequestAction)
	return action, nil
}

// FindAllWithPagination retrieves all amount-based auth tiers with pagination
func (s *amountBasedAuthService) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.AuthTier], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "Amount Based Auth", "FindAllWithPagination")
	defer span.End()

	result, err := s.Repository.FindAllWithPagination(ctx, filterParam)
	if err != nil {
		span.AddEvent("[FindAllWithPagination] failed to fetch amount-based auth tiers", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		s.logger.Errorf("[FindAllWithPagination] failed to fetch amount-based auth tiers: %v", err)
		return nil, err
	}
	return result, nil
}

// UpdateAmountBasedAuth updates any tier type and applies appropriate cascading logic
func (s *amountBasedAuthService) UpdateAmountBasedAuth(ctx context.Context, id string, method shared_constant.Method, request amountauthdto.UpdateAmountBasedAuthRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateAmountBasedAuth", "Amount Based Auth", "UpdateAmountBasedAuth")
	defer span.End()

	// Validate the request based on method
	if !request.Validate(method) {
		span.AddEvent("Invalid amount values", trace.WithAttributes(
			attribute.String("method", string(method)),
			attribute.Int64("min_amount", int64(request.MinAmount)),
			attribute.Int64("max_amount", int64(request.MaxAmount)),
		))
		s.logger.Errorf("Invalid amount values for method %s: MinAmount=%d, MaxAmount=%d", method, request.MinAmount, request.MaxAmount)
		return errors.New(localization.ErrorInvalidAmounts.Code)
	}

	existingTier, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find existing tier", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		s.logger.Errorf(" Failed to find existing tier by ID %s: %v", id, err)
		return err
	}
	if existingTier.Method != method {
		span.AddEvent("Mismatched method", trace.WithAttributes(
			attribute.String("expected", string(method)),
			attribute.String("got", string(existingTier.Method)),
			attribute.String("id", id),
		))
		s.logger.Errorf("Mismatched method for tier ID %s: expected %s, got %s", id, existingTier.Method, method)
		return errors.New(localization.ErrorInvalidMethod.Code)
	}

	notModifiedTier := *existingTier
	now := time.Now()

	switch method {
	case shared_constant.OPEN:
		// For OPEN: only MaxAmount is updated, preserve MinAmount
		existingTier.MaxAmount = request.MaxAmount
		// Fetch PIN tier to cascade min change
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

		// Persist updates: update OPEN first, then PIN
		existingTier.LastModified = now

		pinTier.LastModified = now

		data := map[string]interface{}{
			"method": "OPEN",
			"data": map[string]interface{}{
				"open_id": existingTier.ID.Hex(),
				"pin_id":  pinTier.ID.Hex(),
				"open":    existingTier,
				"pin":     pinTier,
			},
		}
		// Create cps action model
		cpsActionData := lib.CpsModelBuilder(id, local_util.ExtractUserFromContext(ctx), notModifiedTier, data, string(cpsaction.RequestUpdateAmountBasedAuth), constants.UPDATE)

		if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
			s.logger.Errorf("Failed to create CPS action: %v", err)
			return err
		}

		s.logger.Infof("OTP_PIN tier update request completed successfully")

		return nil

	case shared_constant.PIN:
		// For PIN: both MinAmount and MaxAmount can be updated
		existingTier.MinAmount = request.MinAmount
		existingTier.MaxAmount = request.MaxAmount

		s.logger.Infof("Updating PIN tier - existingTier: ID=%s, MinAmount=%d, MaxAmount=%d",
			existingTier.ID.Hex(), existingTier.MinAmount, existingTier.MaxAmount)

		// Fetch OPEN and OTP_PIN tiers for validation constraints
		s.logger.Infof("Fetching OPEN tiers...")
		openTiers, err := s.Repository.FindAll(ctx, bson.M{"method": constants.OPEN, "is_deleted": false}, bson.M{})
		if err != nil {
			s.logger.Errorf("Failed to fetch OPEN tiers: %v", err)
			return err
		}
		s.logger.Infof("Found %d OPEN tiers", len(openTiers))

		if len(openTiers) == 0 {
			s.logger.Errorf("No OPEN tiers found")
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		if openTiers[0] == nil {
			s.logger.Errorf("First OPEN tier is nil")
			return errors.New(localization.ErrorUnexpectedError.Code)
		}

		s.logger.Infof("Fetching OTP_PIN tiers...")
		otpPinTiers, err := s.Repository.FindAll(ctx, bson.M{"method": constants.OTPANDPIN, "is_deleted": false}, bson.M{})
		if err != nil {
			s.logger.Errorf("Failed to fetch OTP_PIN tiers: %v", err)
			return err
		}
		s.logger.Infof("Found %d OTP_PIN tiers", len(otpPinTiers))

		if len(otpPinTiers) == 0 {
			s.logger.Errorf("No OTP_PIN tiers found")
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		if otpPinTiers[0] == nil {
			s.logger.Errorf("First OTP_PIN tier is nil")
			return errors.New(localization.ErrorUnexpectedError.Code)
		}

		openTier := openTiers[0]
		otpPinTier := otpPinTiers[0]

		s.logger.Infof("OPEN tier: ID=%s, MinAmount=%d, MaxAmount=%d",
			openTier.ID.Hex(), openTier.MinAmount, openTier.MaxAmount)
		s.logger.Infof("OTP_PIN tier: ID=%s, MinAmount=%d, MaxAmount=%d",
			otpPinTier.ID.Hex(), otpPinTier.MinAmount, otpPinTier.MaxAmount)

		// Validate the user's PIN values against OPEN and OTP_PIN constraints
		s.logger.Infof("Applying PIN update validation...")
		if err := core.ApplyPinUpdate(existingTier, openTier, otpPinTier); err != nil {
			s.logger.Errorf("PIN update validation failed: %v", err)
			return err
		}

		// Persist all modified tiers: PIN, OPEN, and OTP_PIN
		s.logger.Infof("Updating existing tier...")

		s.logger.Infof("Updating OPEN tier...")
		openTier.LastModified = now

		s.logger.Infof("Updating OTP_PIN tier...")
		otpPinTier.LastModified = now

		data := map[string]interface{}{
			"method": "PIN",
			"data": map[string]interface{}{
				"open_id":    openTier.ID.Hex(),
				"pin_id":     existingTier.ID.Hex(),
				"otp_pin_id": otpPinTier.ID.Hex(),
				"open":       openTier,
				"pin":        existingTier,
				"otp_pin":    otpPinTier,
			},
		}

		// Create cps action model
		cpsActionData := lib.CpsModelBuilder(id, local_util.ExtractUserFromContext(ctx), notModifiedTier, data, string(cpsaction.RequestUpdateAmountBasedAuth), constants.UPDATE)

		if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
			s.logger.Errorf("Failed to create CPS action: %v", err)
			return err
		}
		s.logger.Infof("[UpdateAmountBasedAuth] PIN tier update request created successfully")
		return nil

	case shared_constant.OTPANDPIN:
		// For OTP_PIN: only MinAmount is updated, preserve MaxAmount
		existingTier.MinAmount = request.MinAmount

		// Fetch PIN tier to cascade max change
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

		// Persist updates: update OTP_PIN first, then PIN
		existingTier.LastModified = now

		data := map[string]interface{}{
			"method": "OTP_PIN",
			"data": map[string]interface{}{
				"pin_id":     pinTier.ID.Hex(),
				"otp_pin_id": existingTier.ID.Hex(),
				"pin":        pinTier,
				"otp_pin":    existingTier,
			},
		}
		// Create cps action model
		cpsActionData := lib.CpsModelBuilder(id, local_util.ExtractUserFromContext(ctx), notModifiedTier, data, string(cpsaction.RequestUpdateAmountBasedAuth), constants.UPDATE)

		if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
			s.logger.Errorf("Failed to create CPS action: %v", err)
			return err
		}

		s.logger.Infof("OTP_PIN tier update request completed successfully")

		return nil
	}

	return errors.New(localization.ErrorInvalidMethod.Code)
}
