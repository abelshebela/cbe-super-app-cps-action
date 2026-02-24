<<<<<<< HEAD
package amount_based_auth

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	local_model "cbe-super-app-cps-action/internal/constants/model"
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

	s.logger.Infof("[AmountAuthSvc][Authorize] action: %s", action.RequestAction)

	switch action.RequestAction {
	case string(constants.RequestCreateAmountBasedAuth):
		return s.authorizeCreate(ctx, action, span)
	case string(constants.RequestResetAmountBasedAuth):
		return s.authorizeReset(ctx, action, span)
	case string(constants.RequestUpdateAmountBasedAuth):
		return s.authorizeUpdate(ctx, action, span)
	default:
		s.logger.Errorf("[AmountAuthSvc][Authorize] unsupported request action: %s", action.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}
}

// authorizeCreate persists new currency tiers when a CREATE action is approved
func (s *amountBasedAuthService) authorizeCreate(ctx context.Context, action *model.CPSAction, span trace.Span) (*model.CPSAction, error) {
	s.logger.Infof("[AmountAuthSvc][authorizeCreate] processing CREATE")

	currentAction, err := local_util.JsonUnmarshal[map[string]interface{}](action.CurrentAction)
	if err != nil {
		span.RecordError(err)
		s.logger.Errorf("[AmountAuthSvc][authorizeCreate] unmarshal currentAction err: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	tiersRaw, ok := (*currentAction)["tiers"]
	if !ok {
		s.logger.Errorf("[AmountAuthSvc][authorizeCreate] missing tiers in action data")
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	tiers, err := local_util.JsonUnmarshal[[]local_model.AuthTier](tiersRaw)
	if err != nil {
		s.logger.Errorf("[AmountAuthSvc][authorizeCreate] unmarshal tiers err: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	if err := s.Repository.CreateMany(ctx, *tiers); err != nil {
		s.logger.Errorf("[AmountAuthSvc][authorizeCreate] create tiers err: %v", err)
		return nil, err
	}

	s.logger.Infof("[AmountAuthSvc][authorizeCreate] CREATE done")
	return action, nil
}

// authorizeReset deletes old currency tiers and creates new ones when a RESET action is approved
func (s *amountBasedAuthService) authorizeReset(ctx context.Context, action *model.CPSAction, span trace.Span) (*model.CPSAction, error) {
	s.logger.Infof("[AmountAuthSvc][authorizeReset] processing RESET")

	currentAction, err := local_util.JsonUnmarshal[map[string]interface{}](action.CurrentAction)
	if err != nil {
		span.RecordError(err)
		s.logger.Errorf("[AmountAuthSvc][authorizeReset] unmarshal currentAction err: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	currencyRaw, ok := (*currentAction)["currency"]
	if !ok {
		s.logger.Errorf("[AmountAuthSvc][authorizeReset] missing currency in action data")
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}
	currency, ok := currencyRaw.(string)
	if !ok {
		s.logger.Errorf("[AmountAuthSvc][authorizeReset] invalid currency type")
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	// Delete existing tiers for this currency
	if err := s.Repository.DeleteByCurrency(ctx, currency); err != nil {
		s.logger.Errorf("[AmountAuthSvc][authorizeReset] delete tiers err: %v", err)
		return nil, err
	}

	tiersRaw, ok := (*currentAction)["tiers"]
	if !ok {
		s.logger.Errorf("[AmountAuthSvc][authorizeReset] missing tiers in action data")
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	tiers, err := local_util.JsonUnmarshal[[]local_model.AuthTier](tiersRaw)
	if err != nil {
		s.logger.Errorf("[AmountAuthSvc][authorizeReset] unmarshal tiers err: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	if err := s.Repository.CreateMany(ctx, *tiers); err != nil {
		s.logger.Errorf("[AmountAuthSvc][authorizeReset] create tiers err: %v", err)
		return nil, err
	}

	s.logger.Infof("[AmountAuthSvc][authorizeReset] RESET done")
	return action, nil
}

// authorizeUpdate handles the legacy per-method update flow
func (s *amountBasedAuthService) authorizeUpdate(ctx context.Context, action *model.CPSAction, span trace.Span) (*model.CPSAction, error) {
	s.logger.Infof("[AmountAuthSvc][authorizeUpdate] processing UPDATE")

	currentAction, err := local_util.JsonUnmarshal[map[string]interface{}](action.CurrentAction)
	if err != nil {
		span.AddEvent("Failed to extract currentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", action.UniqueId),
		))
		s.logger.Errorf("[AmountAuthSvc][authorizeUpdate] extract currentAction err: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	result := *currentAction

	switch result["method"] {
	case "OPEN":
		return s.authorizeUpdateOpen(ctx, action, result, span)
	case "PIN":
		return s.authorizeUpdatePin(ctx, action, result, span)
	case "OTP_PIN":
		return s.authorizeUpdateOtpPin(ctx, action, result, span)
	default:
		s.logger.Errorf("[AmountAuthSvc][authorizeUpdate] unsupported method: %v", result["method"])
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}
}

func (s *amountBasedAuthService) authorizeUpdateOpen(ctx context.Context, action *model.CPSAction, result map[string]interface{}, span trace.Span) (*model.CPSAction, error) {
	s.logger.Infof("[AmountAuthSvc][Authorize] processing OPEN tier")
	data, err := local_util.JsonUnmarshal[map[string]interface{}](result["data"])
	if err != nil {
		s.logger.Errorf("[AmountAuthSvc][Authorize] extract data tier err: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	openTierMap, ok := (*data)["open"].(map[string]interface{})
	if !ok {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	openTier, err := local_util.JsonUnmarshal[local_model.AuthTier](openTierMap)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	pinTierMap, ok := (*data)["pin"].(map[string]interface{})
	if !ok {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	pinTier, err := local_util.JsonUnmarshal[local_model.AuthTier](pinTierMap)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	if err := s.Repository.Update(ctx, (*data)["open_id"].(string), openTier); err != nil {
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}
	if err := s.Repository.Update(ctx, (*data)["pin_id"].(string), pinTier); err != nil {
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	s.logger.Infof("[AmountAuthSvc][Authorize] OPEN tier update done")
	return action, nil
}

func (s *amountBasedAuthService) authorizeUpdatePin(ctx context.Context, action *model.CPSAction, result map[string]interface{}, span trace.Span) (*model.CPSAction, error) {
	s.logger.Infof("[AmountAuthSvc][Authorize] processing PIN tier")
	data, err := local_util.JsonUnmarshal[map[string]interface{}](result["data"])
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	openTierMap, ok := (*data)["open"].(map[string]interface{})
	if !ok {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	openTier, err := local_util.JsonUnmarshal[local_model.AuthTier](openTierMap)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	pinTierMap, ok := (*data)["pin"].(map[string]interface{})
	if !ok {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	pinTier, err := local_util.JsonUnmarshal[local_model.AuthTier](pinTierMap)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	otpPinTierMap, ok := (*data)["otp_pin"].(map[string]interface{})
	if !ok {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	otpPinTier, err := local_util.JsonUnmarshal[local_model.AuthTier](otpPinTierMap)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	if err := s.Repository.Update(ctx, (*data)["open_id"].(string), openTier); err != nil {
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}
	if err := s.Repository.Update(ctx, (*data)["pin_id"].(string), pinTier); err != nil {
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}
	if err := s.Repository.Update(ctx, (*data)["otp_pin_id"].(string), otpPinTier); err != nil {
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	s.logger.Infof("[AmountAuthSvc][Authorize] PIN tier update done")
	return action, nil
}

func (s *amountBasedAuthService) authorizeUpdateOtpPin(ctx context.Context, action *model.CPSAction, result map[string]interface{}, span trace.Span) (*model.CPSAction, error) {
	s.logger.Infof("[AmountAuthSvc][Authorize] processing OTP_PIN tier")
	data, err := local_util.JsonUnmarshal[map[string]interface{}](result["data"])
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	pinTierMap, ok := (*data)["pin"].(map[string]interface{})
	if !ok {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	pinTier, err := local_util.JsonUnmarshal[local_model.AuthTier](pinTierMap)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	otpPinTierMap, ok := (*data)["otp_pin"].(map[string]interface{})
	if !ok {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	otpPinTier, err := local_util.JsonUnmarshal[local_model.AuthTier](otpPinTierMap)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	if err := s.Repository.Update(ctx, (*data)["pin_id"].(string), pinTier); err != nil {
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}
	if err := s.Repository.Update(ctx, (*data)["otp_pin_id"].(string), otpPinTier); err != nil {
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	s.logger.Infof("[AmountAuthSvc][Authorize] OTP_PIN tier update done")
	return action, nil
}

// FindAllWithPagination retrieves all amount-based auth tiers grouped by currency
func (s *amountBasedAuthService) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]amountauthdto.CurrencyGroup], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "Amount Based Auth", "FindAllWithPagination")
	defer span.End()

	result, err := s.Repository.FindAllWithPagination(ctx, filterParam)
	if err != nil {
		span.AddEvent("[FindAllWithPagination] failed to fetch amount-based auth tiers", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		s.logger.Errorf("[AmountAuthSvc][FindAll] fetch err: %v", err)
		return nil, err
	}

	// Group tiers by currency
	grouped := groupTiersByCurrency(result.Data)

	return &types.PaginatedResponse[[]amountauthdto.CurrencyGroup]{
		Data: grouped,
		Meta: result.Meta,
	}, nil
}

func groupTiersByCurrency(tiers []local_model.AuthTier) []amountauthdto.CurrencyGroup {
	orderMap := map[constants.CurrencyType]int{}
	groupMap := map[constants.CurrencyType]*amountauthdto.CurrencyGroup{}

	for _, t := range tiers {
		if _, exists := groupMap[t.Currency]; !exists {
			orderMap[t.Currency] = len(orderMap)
			groupMap[t.Currency] = &amountauthdto.CurrencyGroup{
				Currency: t.Currency,
				Tiers:    []amountauthdto.TierResponse{},
			}
		}
		groupMap[t.Currency].Tiers = append(groupMap[t.Currency].Tiers, amountauthdto.TierResponse{
			ID:        t.ID.Hex(),
			Method:    t.Method,
			MinAmount: t.MinAmount,
			MaxAmount: t.MaxAmount,
			Enabled:   t.Enabled,
		})
	}

	groups := make([]amountauthdto.CurrencyGroup, len(groupMap))
	for currency, group := range groupMap {
		groups[orderMap[currency]] = *group
	}
	return groups
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
		s.logger.Errorf("[AmountAuthSvc][Update] invalid amounts method=%s min=%d max=%d", method, request.MinAmount, request.MaxAmount)
		return errors.New(localization.ErrorInvalidAmounts.Code)
	}

	existingTier, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find existing tier", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		s.logger.Errorf("[AmountAuthSvc][Update] find tier id=%s err: %v", id, err)
		return err
	}
	localMethod := constants.Method(method)
	if existingTier.Method != localMethod {
		span.AddEvent("Mismatched method", trace.WithAttributes(
			attribute.String("expected", string(method)),
			attribute.String("got", string(existingTier.Method)),
			attribute.String("id", id),
		))
		s.logger.Errorf("[AmountAuthSvc][Update] method mismatch id=%s expected=%s got=%s", id, existingTier.Method, method)
		return errors.New(localization.ErrorInvalidMethod.Code)
	}

	notModifiedTier := *existingTier
	now := time.Now()
	currency := string(existingTier.Currency)

	switch method {
	case shared_constant.OPEN:
		existingTier.MaxAmount = request.MaxAmount
		pinTiers, err := s.Repository.FindAll(ctx, bson.M{"method": constants.PIN, "currency": currency, "is_deleted": false}, bson.M{})
		if err != nil {
			return err
		}
		if len(pinTiers) == 0 {
			return errors.New(localization.ErrorFileNotFound.Code)
		}

		if err := core.ApplyOpenUpdate(existingTier, &pinTiers[0]); err != nil {
			return err
		}

		existingTier.LastModified = now
		pinTiers[0].LastModified = now

		data := map[string]interface{}{
			"method": "OPEN",
			"data": map[string]interface{}{
				"open_id": existingTier.ID.Hex(),
				"pin_id":  pinTiers[0].ID.Hex(),
				"open":    existingTier,
				"pin":     &pinTiers[0],
			},
		}
		cpsActionData := lib.CpsModelBuilder(id, local_util.ExtractUserFromContext(ctx), notModifiedTier, data, string(cpsaction.RequestUpdateAmountBasedAuth), constants.UPDATE)

		if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
			s.logger.Errorf("[AmountAuthSvc][Update] cps action err: %v", err)
			return err
		}

		s.logger.Infof("[AmountAuthSvc][Update] OPEN tier request done")
		return nil

	case shared_constant.PIN:
		existingTier.MinAmount = request.MinAmount
		existingTier.MaxAmount = request.MaxAmount

		openTiers, err := s.Repository.FindAll(ctx, bson.M{"method": constants.OPEN, "currency": currency, "is_deleted": false}, bson.M{})
		if err != nil {
			return err
		}
		if len(openTiers) == 0 {
			return errors.New(localization.ErrorFileNotFound.Code)
		}

		otpPinTiers, err := s.Repository.FindAll(ctx, bson.M{"method": constants.OTPANDPIN, "currency": currency, "is_deleted": false}, bson.M{})
		if err != nil {
			return err
		}
		if len(otpPinTiers) == 0 {
			return errors.New(localization.ErrorFileNotFound.Code)
		}

		if err := core.ApplyPinUpdate(existingTier, &openTiers[0], &otpPinTiers[0]); err != nil {
			return err
		}

		openTiers[0].LastModified = now
		otpPinTiers[0].LastModified = now

		data := map[string]interface{}{
			"method": "PIN",
			"data": map[string]interface{}{
				"open_id":    openTiers[0].ID.Hex(),
				"pin_id":     existingTier.ID.Hex(),
				"otp_pin_id": otpPinTiers[0].ID.Hex(),
				"open":       &openTiers[0],
				"pin":        existingTier,
				"otp_pin":    &otpPinTiers[0],
			},
		}

		cpsActionData := lib.CpsModelBuilder(id, local_util.ExtractUserFromContext(ctx), notModifiedTier, data, string(cpsaction.RequestUpdateAmountBasedAuth), constants.UPDATE)

		if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
			s.logger.Errorf("[AmountAuthSvc][Update] cps action err: %v", err)
			return err
		}
		s.logger.Infof("[AmountAuthSvc][Update] PIN tier request done")
		return nil

	case shared_constant.OTPANDPIN:
		existingTier.MinAmount = request.MinAmount

		pinTiers, err := s.Repository.FindAll(ctx, bson.M{"method": constants.PIN, "currency": currency, "is_deleted": false}, bson.M{})
		if err != nil {
			return err
		}
		if len(pinTiers) == 0 {
			return errors.New(localization.ErrorFileNotFound.Code)
		}

		if err := core.ApplyOtpPinUpdate(existingTier, &pinTiers[0]); err != nil {
			return err
		}

		existingTier.LastModified = now

		data := map[string]interface{}{
			"method": "OTP_PIN",
			"data": map[string]interface{}{
				"pin_id":     pinTiers[0].ID.Hex(),
				"otp_pin_id": existingTier.ID.Hex(),
				"pin":        &pinTiers[0],
				"otp_pin":    existingTier,
			},
		}
		cpsActionData := lib.CpsModelBuilder(id, local_util.ExtractUserFromContext(ctx), notModifiedTier, data, string(cpsaction.RequestUpdateAmountBasedAuth), constants.UPDATE)

		if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
			s.logger.Errorf("[AmountAuthSvc][Update] cps action err: %v", err)
			return err
		}

		s.logger.Infof("[AmountAuthSvc][Update] OTP_PIN tier request done")
		return nil
	}

	return errors.New(localization.ErrorInvalidMethod.Code)
}

// AddCurrency creates tier configuration for a new currency
func (s *amountBasedAuthService) AddCurrency(ctx context.Context, request amountauthdto.AddCurrencyRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "AddCurrency", "Amount Based Auth", "AddCurrency")
	defer span.End()

	// Check duplicate currency
	exists, err := s.Repository.CurrencyExists(ctx, string(request.Currency))
	if err != nil {
		s.logger.Errorf("[AmountAuthSvc][AddCurrency] check currency err: %v", err)
		return err
	}
	if exists {
		s.logger.Warnf("[AmountAuthSvc][AddCurrency] currency already exists: %s", request.Currency)
		return errors.New(localization.ErrorCurrencyAlreadyExists.Code)
	}

	// Validate tier cascade
	if err := core.ValidateTierCascade(request.Tiers); err != nil {
		s.logger.Errorf("[AmountAuthSvc][AddCurrency] tier cascade validation err: %v", err)
		return err
	}

	// Build tier models
	tiers := core.BuildTiersFromRequest(request.Currency, request.Tiers)

	// CPS action data
	data := map[string]interface{}{
		"currency": string(request.Currency),
		"tiers":    tiers,
	}

	cpsActionData := lib.CpsModelBuilder(
		string(request.Currency),
		local_util.ExtractUserFromContext(ctx),
		nil,
		data,
		string(constants.RequestCreateAmountBasedAuth),
		constants.CREATE,
	)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		s.logger.Errorf("[AmountAuthSvc][AddCurrency] cps action err: %v", err)
		return err
	}

	s.logger.Infof("[AmountAuthSvc][AddCurrency] currency %s add request sent", request.Currency)
	return nil
}

// ResetConfig resets tier configuration for an existing currency
func (s *amountBasedAuthService) ResetConfig(ctx context.Context, currency constants.CurrencyType, request amountauthdto.ResetConfigRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "ResetConfig", "Amount Based Auth", "ResetConfig")
	defer span.End()

	// Check currency exists
	exists, err := s.Repository.CurrencyExists(ctx, string(currency))
	if err != nil {
		s.logger.Errorf("[AmountAuthSvc][ResetConfig] check currency err: %v", err)
		return err
	}
	if !exists {
		s.logger.Warnf("[AmountAuthSvc][ResetConfig] currency not found: %s", currency)
		return errors.New(localization.ErrorCurrencyNotFound.Code)
	}

	// Validate tier cascade
	if err := core.ValidateTierCascade(request.Tiers); err != nil {
		s.logger.Errorf("[AmountAuthSvc][ResetConfig] tier cascade validation err: %v", err)
		return err
	}

	// Get existing tiers as previous action
	existingTiers, err := s.Repository.FindByCurrency(ctx, string(currency))
	if err != nil {
		s.logger.Errorf("[AmountAuthSvc][ResetConfig] fetch existing tiers err: %v", err)
		return err
	}

	// Build new tier models
	tiers := core.BuildTiersFromRequest(currency, request.Tiers)

	// CPS action data
	data := map[string]interface{}{
		"currency": string(currency),
		"tiers":    tiers,
	}

	cpsActionData := lib.CpsModelBuilder(
		string(currency),
		local_util.ExtractUserFromContext(ctx),
		existingTiers,
		data,
		string(constants.RequestResetAmountBasedAuth),
		constants.UPDATE,
	)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		s.logger.Errorf("[AmountAuthSvc][ResetConfig] cps action err: %v", err)
		return err
	}

	s.logger.Infof("[AmountAuthSvc][ResetConfig] currency %s reset request sent", currency)
	return nil
}
=======
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

	s.logger.Infof("[AmountAuthSvc][Authorize] action: %s", action.RequestAction)

	// Unmarshal the current action data
	// Try direct type assertion first

	currentAction, err := local_util.JsonUnmarshal[map[string]interface{}](action.CurrentAction)
	if err != nil {
		span.AddEvent("Failed to extract currentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", action.UniqueId),
		))
		s.logger.Errorf("[AmountAuthSvc][Authorize] extract currentAction err: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	s.logger.Infof("[AmountAuthSvc][Authorize] processing: %s", action.RequestAction)

	result := *currentAction

	switch result["method"] {
	case "OPEN":
		s.logger.Infof("[AmountAuthSvc][Authorize] processing OPEN tier")
		data, err := local_util.JsonUnmarshal[map[string]interface{}](result["data"])
		if err != nil {
			s.logger.Errorf("[AmountAuthSvc][Authorize] extract data tier err: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		openTierMap, ok := (*data)["open"].(map[string]interface{})
		if !ok {
			s.logger.Errorf("[AmountAuthSvc][Authorize] cast open tier err: %v", (*data)["open"])
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		openTier, err := local_util.JsonUnmarshal[model.AuthTier](openTierMap)
		if err != nil {
			s.logger.Errorf("[AmountAuthSvc][Authorize] unmarshal open tier err: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		pinTierMap, ok := (*data)["pin"].(map[string]interface{})
		if !ok {
			s.logger.Errorf("[AmountAuthSvc][Authorize] cast pin tier err: %v", (*data)["pin"])
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		pinTier, err := local_util.JsonUnmarshal[model.AuthTier](pinTierMap)
		if err != nil {
			s.logger.Errorf("[AmountAuthSvc][Authorize] unmarshal pin tier err: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		if err := s.Repository.Update(ctx, (*data)["open_id"].(string), openTier); err != nil {
			s.logger.Errorf("[AmountAuthSvc][Authorize] update OPEN tier err: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		if err := s.Repository.Update(ctx, (*data)["pin_id"].(string), pinTier); err != nil {
			s.logger.Errorf("[AmountAuthSvc][Authorize] update PIN tier err: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		s.logger.Infof("[AmountAuthSvc][Authorize] OPEN tier update done")
	case "PIN":
		s.logger.Infof("[AmountAuthSvc][Authorize] processing PIN tier")
		data, err := local_util.JsonUnmarshal[map[string]interface{}](result["data"])
		if err != nil {
			s.logger.Errorf("[AmountAuthSvc][Authorize] extract data tier err: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		openTierMap, ok := (*data)["open"].(map[string]interface{})
		if !ok {
			s.logger.Errorf("[AmountAuthSvc][Authorize] cast open tier err: %v", (*data)["open"])
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		openTier, err := local_util.JsonUnmarshal[model.AuthTier](openTierMap)
		if err != nil {
			s.logger.Errorf("[AmountAuthSvc][Authorize] unmarshal open tier err: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		pinTierMap, ok := (*data)["pin"].(map[string]interface{})
		if !ok {
			s.logger.Errorf("[AmountAuthSvc][Authorize] cast pin tier err: %v", (*data)["pin"])
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		pinTier, err := local_util.JsonUnmarshal[model.AuthTier](pinTierMap)
		if err != nil {
			s.logger.Errorf("[AmountAuthSvc][Authorize] unmarshal pin tier err: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		otpPinTierMap, ok := (*data)["otp_pin"].(map[string]interface{})
		if !ok {
			s.logger.Errorf("[AmountAuthSvc][Authorize] cast otp_pin tier err: %v", (*data)["otp_pin"])
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		otpPinTier, err := local_util.JsonUnmarshal[model.AuthTier](otpPinTierMap)
		if err != nil {
			s.logger.Errorf("[AmountAuthSvc][Authorize] unmarshal otp_pin tier err: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		if err := s.Repository.Update(ctx, (*data)["open_id"].(string), openTier); err != nil {
			s.logger.Errorf("[AmountAuthSvc][Authorize] update OPEN tier err: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		if err := s.Repository.Update(ctx, (*data)["pin_id"].(string), pinTier); err != nil {
			s.logger.Errorf("[AmountAuthSvc][Authorize] update PIN tier err: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		if err := s.Repository.Update(ctx, (*data)["otp_pin_id"].(string), otpPinTier); err != nil {
			s.logger.Errorf("[AmountAuthSvc][Authorize] update OTP_PIN tier err: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}
	case "OTP_PIN":
		s.logger.Infof("[AmountAuthSvc][Authorize] processing OTP_PIN tier")
		data, err := local_util.JsonUnmarshal[map[string]interface{}](result["data"])
		if err != nil {
			s.logger.Errorf("[AmountAuthSvc][Authorize] extract data tier err: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		pinTierMap, ok := (*data)["pin"].(map[string]interface{})
		if !ok {
			s.logger.Errorf("[AmountAuthSvc][Authorize] cast pin tier err: %v", (*data)["pin"])
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		pinTier, err := local_util.JsonUnmarshal[model.AuthTier](pinTierMap)
		if err != nil {
			s.logger.Errorf("[AmountAuthSvc][Authorize] unmarshal pin tier err: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		otpPinTierMap, ok := (*data)["otp_pin"].(map[string]interface{})
		if !ok {
			s.logger.Errorf("[AmountAuthSvc][Authorize] cast otp_pin tier err: %v", (*data)["otp_pin"])
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		otpPinTier, err := local_util.JsonUnmarshal[model.AuthTier](otpPinTierMap)
		if err != nil {
			s.logger.Errorf("[AmountAuthSvc][Authorize] unmarshal otp_pin tier err: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		if err := s.Repository.Update(ctx, (*data)["otp_pin_id"].(string), otpPinTier); err != nil {
			s.logger.Errorf("[AmountAuthSvc][Authorize] update OTP_PIN tier err: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		if err := s.Repository.Update(ctx, (*data)["pin_id"].(string), pinTier); err != nil {
			s.logger.Errorf("[AmountAuthSvc][Authorize] update PIN tier err: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		if err := s.Repository.Update(ctx, (*data)["otp_pin_id"].(string), otpPinTier); err != nil {
			s.logger.Errorf("[AmountAuthSvc][Authorize] update OTP_PIN tier err: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}
	default:
		s.logger.Errorf("[AmountAuthSvc][Authorize] unsupported method: %v", result["method"])
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	s.logger.Infof("[AmountAuthSvc][Authorize] done: %s", action.RequestAction)
	return action, nil
}

// FindAllWithPagination retrieves all amount-based auth tiers with pagination
func (s *amountBasedAuthService) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.AuthTier], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "Amount Based Auth", "FindAllWithPagination")
	defer span.End()

	result, err := s.Repository.FindAllWithPagination(ctx, filterParam)
	if err != nil {
		span.AddEvent("[FindAllWithPagination] failed to fetch amount-based auth tiers", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		s.logger.Errorf("[AmountAuthSvc][FindAll] fetch err: %v", err)
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
		s.logger.Errorf("[AmountAuthSvc][Update] invalid amounts method=%s min=%d max=%d", method, request.MinAmount, request.MaxAmount)
		return errors.New(localization.ErrorInvalidAmounts.Code)
	}

	existingTier, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find existing tier", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		s.logger.Errorf("[AmountAuthSvc][Update] find tier id=%s err: %v", id, err)
		return err
	}
	if existingTier.Method != method {
		span.AddEvent("Mismatched method", trace.WithAttributes(
			attribute.String("expected", string(method)),
			attribute.String("got", string(existingTier.Method)),
			attribute.String("id", id),
		))
		s.logger.Errorf("[AmountAuthSvc][Update] method mismatch id=%s expected=%s got=%s", id, existingTier.Method, method)
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

		if err := core.ApplyOpenUpdate(existingTier, &pinTiers[0]); err != nil {
			return err
		}

		// Persist updates: update OPEN first, then PIN
		existingTier.LastModified = now

		pinTiers[0].LastModified = now

		data := map[string]interface{}{
			"method": "OPEN",
			"data": map[string]interface{}{
				"open_id": existingTier.ID.Hex(),
				"pin_id":  pinTiers[0].ID.Hex(),
				"open":    existingTier,
				"pin":     &pinTiers[0],
			},
		}
		// Create cps action model
		cpsActionData := lib.CpsModelBuilder(id, local_util.ExtractUserFromContext(ctx), notModifiedTier, data, string(cpsaction.RequestUpdateAmountBasedAuth), constants.UPDATE)

		if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
			s.logger.Errorf("[AmountAuthSvc][Update] cps action err: %v", err)
			return err
		}

		s.logger.Infof("[AmountAuthSvc][Update] OPEN tier request done")

		return nil

	case shared_constant.PIN:
		// For PIN: both MinAmount and MaxAmount can be updated
		existingTier.MinAmount = request.MinAmount
		existingTier.MaxAmount = request.MaxAmount

		s.logger.Infof("[AmountAuthSvc][Update] PIN tier id=%s min=%d max=%d",
			existingTier.ID.Hex(), existingTier.MinAmount, existingTier.MaxAmount)

		// Fetch OPEN and OTP_PIN tiers for validation constraints
		s.logger.Infof("[AmountAuthSvc][Update] fetching OPEN tiers")
		openTiers, err := s.Repository.FindAll(ctx, bson.M{"method": constants.OPEN, "is_deleted": false}, bson.M{})
		if err != nil {
			s.logger.Errorf("[AmountAuthSvc][Update] fetch OPEN tiers err: %v", err)
			return err
		}
		s.logger.Infof("[AmountAuthSvc][Update] found %d OPEN tiers", len(openTiers))

		if len(openTiers) == 0 {
			s.logger.Errorf("[AmountAuthSvc][Update] no OPEN tiers found")
			return errors.New(localization.ErrorFileNotFound.Code)
		}

		s.logger.Infof("[AmountAuthSvc][Update] fetching OTP_PIN tiers")
		otpPinTiers, err := s.Repository.FindAll(ctx, bson.M{"method": constants.OTPANDPIN, "is_deleted": false}, bson.M{})
		if err != nil {
			s.logger.Errorf("[AmountAuthSvc][Update] fetch OTP_PIN tiers err: %v", err)
			return err
		}
		s.logger.Infof("[AmountAuthSvc][Update] found %d OTP_PIN tiers", len(otpPinTiers))

		if len(otpPinTiers) == 0 {
			s.logger.Errorf("[AmountAuthSvc][Update] no OTP_PIN tiers found")
			return errors.New(localization.ErrorFileNotFound.Code)
		}

		s.logger.Infof("[AmountAuthSvc][Update] validating PIN update")
		if err := core.ApplyPinUpdate(existingTier, &openTiers[0], &otpPinTiers[0]); err != nil {
			s.logger.Errorf("[AmountAuthSvc][Update] PIN validation err: %v", err)
			return err
		}

		// Persist all modified tiers: PIN, OPEN, and OTP_PIN
		s.logger.Infof("[AmountAuthSvc][Update] updating existing tier")

		s.logger.Infof("[AmountAuthSvc][Update] updating OPEN tier")
		openTiers[0].LastModified = now

		s.logger.Infof("[AmountAuthSvc][Update] updating OTP_PIN tier")
		otpPinTiers[0].LastModified = now

		data := map[string]interface{}{
			"method": "PIN",
			"data": map[string]interface{}{
				"open_id":    openTiers[0].ID.Hex(),
				"pin_id":     existingTier.ID.Hex(),
				"otp_pin_id": otpPinTiers[0].ID.Hex(),
				"open":       &openTiers[0],
				"pin":        existingTier,
				"otp_pin":    &otpPinTiers[0],
			},
		}

		// Create cps action model
		cpsActionData := lib.CpsModelBuilder(id, local_util.ExtractUserFromContext(ctx), notModifiedTier, data, string(cpsaction.RequestUpdateAmountBasedAuth), constants.UPDATE)

		if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
			s.logger.Errorf("[AmountAuthSvc][Update] cps action err: %v", err)
			return err
		}
		s.logger.Infof("[AmountAuthSvc][Update] PIN tier request done")
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

		if err := core.ApplyOtpPinUpdate(existingTier, &pinTiers[0]); err != nil {
			return err
		}

		// Persist updates: update OTP_PIN first, then PIN
		existingTier.LastModified = now

		data := map[string]interface{}{
			"method": "OTP_PIN",
			"data": map[string]interface{}{
				"pin_id":     pinTiers[0].ID.Hex(),
				"otp_pin_id": existingTier.ID.Hex(),
				"pin":        &pinTiers[0],
				"otp_pin":    existingTier,
			},
		}
		// Create cps action model
		cpsActionData := lib.CpsModelBuilder(id, local_util.ExtractUserFromContext(ctx), notModifiedTier, data, string(cpsaction.RequestUpdateAmountBasedAuth), constants.UPDATE)

		if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
			s.logger.Errorf("[AmountAuthSvc][Update] cps action err: %v", err)
			return err
		}

		s.logger.Infof("[AmountAuthSvc][Update] OTP_PIN tier request done")

		return nil
	}

	return errors.New(localization.ErrorInvalidMethod.Code)
}
>>>>>>> 3640c5b5dde77249222ec7dd9f90a5770ab0dbc4
