package customer

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/dto/customer"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/external_call/account_lookup"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"cbe-super-app-cps-action/internal/constants/types"

	customer_dto "cbe-super-app-cps-action/internal/constants/dto/customer"

	member "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type customerService struct {
	repo       storage.CustomerRepository
	redis      storage.RedisRepository
	bpsRepo    storage.BPSActionRepository
	cpsService service.CPSActionService
	cfg        *config.VaultConfig
	logger     utils.Logger
	smsService *lib.NotificationStore
	core       account_lookup.Account
}

func NewCustomerService(repo storage.CustomerRepository, bpsRepo storage.BPSActionRepository, cpsService service.CPSActionService, redis storage.RedisRepository, smsService *lib.NotificationStore, core account_lookup.Account, cfg *config.VaultConfig, logger utils.Logger) service.CustomerService {
	return &customerService{
		repo:       repo,
		bpsRepo:    bpsRepo,
		cpsService: cpsService,
		redis:      redis,
		cfg:        cfg,
		logger:     logger,
		smsService: smsService,
		core:       core,
	}
}

// GetCustomerActionLogByID implements [service.CustomerService].
func (s *customerService) GetCustomerActionLogByID(ctx context.Context, id string, filterParams types.Filter) (types.PaginatedResponse[[]customer_dto.CustomerActionLogResponse], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCustomerActionLog", "Customer", "GetCustomerActionLog")
	defer span.End()

	cus, err := s.bpsRepo.GetBPSActionByUserID(ctx, id, filterParams)
	if err != nil {
		span.AddEvent("Failed to fetch customer", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return types.PaginatedResponse[[]customer_dto.CustomerActionLogResponse]{}, err
	}
	response := MapBpsActionToCustomerLog(cus.Data)
	return types.PaginatedResponse[[]customer_dto.CustomerActionLogResponse]{
		Data: response,
		Meta: cus.Meta,
	}, nil
}

func (c *customerService) GetCustomersDetail(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*customer_dto.CustomerListResponse], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCustomersDetail", "Customer", "GetCustomersDetail")
	defer span.End()

	customers, err := c.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("Failed to fetch customers", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}

	return customers, nil
}

func (s *customerService) GetCustomerByID(ctx context.Context, id string) (customer.FindCustomerByIDResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCustomerByID", "Customer", "GetCustomerByID")
	defer span.End()

	cus, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to fetch customer", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return customer.FindCustomerByIDResponse{}, err
	}
	response := MapToDto(cus)
	return response, nil
}

func (s *customerService) GetLinkedAccount(ctx context.Context, id string) ([]model.LinkedAccount, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetLinkedAccount", "Customer", "GetLinkedAccount")
	defer span.End()

	accounts, err := s.repo.FetchLinkedAccount(ctx, id)
	if err != nil {
		span.AddEvent("Failed to fetch linked accounts", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_id", id),
		))
		return nil, err
	}
	return accounts, nil
}

func (s *customerService) GetBlockedCustomer(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*customer_dto.CustomerListResponse], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetBlockedCustomer", "Customer", "GetBlockedCustomer")
	defer span.End()

	if filterParams.Filters == nil {
		filterParams.Filters = make(map[string]interface{})
	}
	filterParams.Filters["is_blocked"] = true
	customers, err := s.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("Failed to fetch blocked customers", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	return customers, nil
}

// CreateEnableCustomerSession implements service.CustomerService.
func (c *customerService) CreateEnableCustomerSession(ctx context.Context, id string) (string, error) {
	log := local_util.LoggerFromCtx(ctx, c.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "CreateEnableCustomerSession", "Customer", "CreateEnableCustomerSession")
	defer span.End()

	log.Infof("[CustomerSvc][CreateEnableSession] id: %s", id)
	existing_otp, err := c.redis.Get(ctx, fmt.Sprintf("cps:action:otp:%s", id))
	if existing_otp != "" || err == nil {
		log.Errorf("[CustomerSvc][CreateEnableSession] OTP exists")
		span.AddEvent("OTP already exists", trace.WithAttributes(
			attribute.String("error", localization.ErrorOTPAlreadyExists.Code),
			attribute.String("id", id),
		))
		return "", errors.New(localization.ErrorOTPAlreadyExists.Code)
	}

	customer, err := c.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[CustomerSvc][CreateEnableSession] find err: %v", err)
		span.AddEvent("Failed to find customer", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return "", err
	}
	if customer.Enabled {
		log.Errorf("[CustomerSvc][CreateEnableSession] already enabled")
		span.AddEvent("Customer already enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorCustomerAlreadyEnabled.Code),
			attribute.String("id", id),
		))
		return "", errors.New(localization.ErrorCustomerAlreadyEnabled.Code)
	}

	otp := local_util.OTPGenerator(6)
	encryptedOTP, _, err := local_util.LocalEncryptPassword(otp, constants.OTP, constants.OTP, constants.OTP, c.cfg)
	if err != nil {
		log.Errorf("[CustomerSvc][CreateEnableSession] encrypt OTP err: %v", err)
		span.AddEvent("Failed to encrypt OTP", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return "", err
	}

	err = c.redis.Set(ctx, fmt.Sprintf("cps:action:otp:%s", id), encryptedOTP, constants.OtpExpirationTime)
	if err != nil {
		log.Errorf("[CustomerSvc][CreateEnableSession] redis set err: %v", err)
		span.AddEvent("Failed to store OTP in redis", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return "", err
	}

	go func() {

		if err := c.smsService.PublishMessage(context.Background(), types.SMSKafkaMessage{
			Recipient:   customer.PhoneNumber,
			MessageBody: fmt.Sprintf("Your Supper app verification OTP: %s ", otp),
		}); err != nil {
			log.Errorf("[CustomerSvc][CreateEnableSession] SMS err: %v", err)
		}
	}()

	log.Infof("[CustomerSvc][CreateEnableSession] OTP generated id: %s", id)
	return otp, nil
}

func (c *customerService) ApproveFaydaCustomer(ctx context.Context, id string, req customer.FaydaApproveRequest) error {
	log := local_util.LoggerFromCtx(ctx, c.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "ApproveFaydaCustomer", "Customer", "ApproveFaydaCustomer")
	defer span.End()

	log.Infof("[CustomerSvc][ApproveFayda] id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		log.Errorf("[CustomerSvc][ApproveFayda] incomplete user")
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.Incomplete),
			attribute.String("id", id),
		))
		return errors.New(constants.Incomplete)
	}

	customer, err := c.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[CustomerSvc][ApproveFayda] find err: %v", err)
		span.AddEvent("Failed to find customer", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	if customer.Enabled {
		log.Errorf("[CustomerSvc][ApproveFayda] already enabled")
		span.AddEvent("Customer already enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorCustomerAlreadyEnabled.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorCustomerAlreadyEnabled.Code)
	}

	// if customer.KYCLevel != uint8(constants.ONE) {
	// 	log.Errorf("[ApproveFaydaCustomer] user is not a Fayda registered customer")
	// 	span.AddEvent("User is not a Fayda registered customer", trace.WithAttributes(
	// 		attribute.String("error", localization.ErrorUserNotFaydaRegistered.Code),
	// 		attribute.String("id", id),
	// 	))
	// 	return errors.New(localization.ErrorUserNotFaydaRegistered.Code)
	// }

	updateData := *customer

	// updateData.FaydaRiskLevel = req.RiskLevel

	action := lib.CpsModelBuilder(id, makerData, customer, updateData, string(constants.RequestApproveFaydaCustomer), constants.UPDATE)

	if err := c.cpsService.CreateCPSAction(ctx, &action); err != nil {
		log.Errorf("[CustomerSvc][ApproveFayda] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	log.Infof("[CustomerSvc][ApproveFayda] request created id: %s", id)
	return nil
}

func (c *customerService) EnableCustomerByID(ctx context.Context, id string, user_otp string) error {
	log := local_util.LoggerFromCtx(ctx, c.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "EnableCustomerByID", "Customer", "EnableCustomerByID")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.IncompleteUserInfo),
			attribute.String("id", id),
		))
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	otp, err := c.redis.Get(ctx, fmt.Sprintf("cps:action:otp:%s", id))

	if otp == "" || err != nil {
		span.AddEvent("OTP not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorOTPNotFound.Code),
			attribute.String("id", id),
		))
		return fmt.Errorf("%s", localization.ErrorOTPNotFound.Code)
	}
	encryptedOTP, _, err := local_util.LocalEncryptPassword(user_otp, constants.OTP, constants.OTP, constants.OTP, c.cfg)
	if err != nil {
		span.AddEvent("Failed to encrypt OTP", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	if encryptedOTP != otp {
		log.Infof("[CustomerSvc][Enable] OTP mismatch")
		span.AddEvent("OTP mismatch", trace.WithAttributes(
			attribute.String("error", localization.ErrorOTPInvalid.Code),
			attribute.String("id", id),
		))
		return fmt.Errorf("%s", localization.ErrorOTPInvalid.Code)
	}

	customer, err := c.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find customer", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	new_customer := *customer
	new_customer.Enabled = true
	new_customer.IsBlocked = false

	action := lib.CpsModelBuilder(id, makerData, customer, new_customer, string(constants.RequestEnableDisableCustomer), constants.UPDATE)

	err = c.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		log.Errorf("[CustomerSvc][Enable] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	if err := c.redis.Delete(ctx, fmt.Sprintf("cps:action:otp:%s", id)); err != nil {
		log.Errorf("[CustomerSvc][Enable] redis delete err: %v", err)
	}
	log.Infof("[CustomerSvc][Enable] request created id: %s", id)
	return nil
}

func (c *customerService) DisableCustomerByID(ctx context.Context, id string, disable customer.CustomerDisableDTO) error {
	log := local_util.LoggerFromCtx(ctx, c.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "DisableCustomerByID", "Customer", "DisableCustomerByID")
	defer span.End()

	log.Infof("[CustomerSvc][Disable] id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		log.Errorf("[CustomerSvc][Disable] incomplete user")
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.IncompleteUserInfo),
			attribute.String("id", id),
		))
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	customer, err := c.repo.FindUserByUserCode(ctx, id)
	if err != nil {
		return err
	}

	if disable.Channel == "BOTH" && !customer.ISuperappEnabled && !customer.IsUSSDEnabled {
		return fmt.Errorf("Both channels are already disabled")
	} else if disable.Channel == "SUPPERAPP" && !customer.ISuperappEnabled {
		return fmt.Errorf("Supperapp channel is already disabled")
	} else if disable.Channel == "USSD" && !customer.IsUSSDEnabled {
		return fmt.Errorf("USSD channel is already disabled")
	}

	channel := disable.Channel
	if channel != "BOTH" && channel != "SUPPERAPP" && channel != "USSD" {
		log.Errorf("[CustomerSvc][Disable] invalid channel: %s", channel)
		span.AddEvent("Invalid channel", trace.WithAttributes(
			attribute.String("error", localization.ErrorInvalidInputParameter.Code),
			attribute.String("channel", channel),
		))
		return fmt.Errorf("%s", localization.ErrorInvalidInputParameter.Code)
	}

	new_customer := *customer
	new_customer.Enabled = false
	new_customer.IsBlocked = true
	new_customer.BlockedReason = disable.DisableReason
	new_customer.BlockedOn = member.BlockedOn(constants.CPS)

	currentAction := imodel.DisableCustomerCurrentAction{
		User:    new_customer,
		Channel: channel,
	}

	action := lib.CpsModelBuilder(id, makerData, customer, currentAction, string(constants.RequestEnableDisableCustomer), constants.UPDATE)

	err = c.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		log.Errorf("[CustomerSvc][Disable] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	log.Infof("[CustomerSvc][Disable] request created id: %s", id)
	return nil
}

func (c *customerService) BlockCustomerSession(ctx context.Context, id, BlockedReason string) error {
	log := local_util.LoggerFromCtx(ctx, c.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "BlockCustomerSession", "Customer", "BlockCustomerSession")
	defer span.End()

	log.Infof("[CustomerSvc][BlockCustomerSession] id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		log.Errorf("[CustomerSvc][BlockCustomerSession] incomplete user")
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.IncompleteUserInfo),
			attribute.String("id", id),
		))
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	customer, err := c.repo.FindUserByUserCode(ctx, id)
	if err != nil {
		return err
	}

	if customer.IsBlocked {
		log.Errorf("[CustomerSvc][BlockCustomerSession] already blocked")
		span.AddEvent("Customer already blocked", trace.WithAttributes(
			attribute.String("error", localization.ErrorCustomerAlreadyBlocked.Code),
			attribute.String("id", id),
		))
		return fmt.Errorf("%s", localization.ErrorCustomerAlreadyBlocked.Code)
	}

	new_customer := *customer
	new_customer.IsBlocked = true
	new_customer.BlockedReason = BlockedReason

	action := lib.CpsModelBuilder(new_customer.UserCode, makerData, customer, new_customer, string(constants.RequestBlockCustomer), constants.UPDATE)

	err = c.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		log.Errorf("[CustomerSvc][BlockCustomerSession] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	// entry := &imodel.CustomerBarUnBarReason{
	// 	UserID:    id,
	// 	Reason:    BlockedReason,
	// 	IsBarred:  true,
	// 	CreatedBy: makerData.UserCode,
	// 	CreatedAt: time.Now(),
	// }
	// if saveErr := c.repo.SaveBarUnBarReason(ctx, entry); saveErr != nil {
	// 	log.Errorf("[CustomerSvc][BlockCustomerSession] save reason err: %v", saveErr)
	// }

	log.Infof("[CustomerSvc][BlockCustomerSession] request created id: %s", id)
	return nil
}

func (c *customerService) UnBlockCustomerSession(ctx context.Context, id, reason string) error {
	log := local_util.LoggerFromCtx(ctx, c.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "BlockCustomerSession", "Customer", "BlockCustomerSession")
	defer span.End()

	log.Infof("[CustomerSvc][BlockCustomerSession] id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		log.Errorf("[CustomerSvc][BlockCustomerSession] incomplete user")
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.IncompleteUserInfo),
			attribute.String("id", id),
		))
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	customer, err := c.repo.FindUserByUserCode(ctx, id)
	if err != nil {
		return err
	}

	if !customer.IsBlocked {
		log.Errorf("[CustomerSvc][UnBlockCustomerSession] already unblocked")
		span.AddEvent("Customer already unblocked", trace.WithAttributes(
			attribute.String("error", localization.ErrorCustomerAlreadyUnBlocked.Code),
			attribute.String("id", id),
		))
		return fmt.Errorf("%s", localization.ErrorCustomerAlreadyUnBlocked.Code)
	}

	new_customer := *customer
	new_customer.IsBlocked = false
	new_customer.BlockedReason = reason

	action := lib.CpsModelBuilder(customer.UserCode, makerData, customer, new_customer, string(constants.RequestUnblockCustomer), constants.UPDATE)

	err = c.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		log.Errorf("[CustomerSvc][UnBlockCustomerSession] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	// entry := &imodel.CustomerBarUnBarReason{
	// 	UserID:    id,
	// 	Reason:    "",
	// 	IsBarred:  false,
	// 	CreatedBy: makerData.UserCode,
	// 	CreatedAt: time.Now(),
	// }
	// if saveErr := c.repo.SaveBarUnBarReason(ctx, entry); saveErr != nil {
	// 	log.Errorf("[CustomerSvc][UnBlockCustomerSession] save reason err: %v", saveErr)
	// }

	log.Infof("[CustomerSvc][UnBlockCustomerSession] request created id: %s", id)
	return nil
}

func (d *customerService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, d.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Customer", "Authorize")
	defer span.End()

	log.Infof("[CustomerSvc][Authorize] action: %s", cpsAction.RequestAction)
	actionData, err := local_util.JsonUnmarshal[member.User](cpsAction.CurrentAction)
	if err != nil {
		log.Errorf("[CustomerSvc][Authorize] unmarshal err: %v", err)
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		return nil, fmt.Errorf("failed to unmarshal CurrentAction to User: %v", err)
	}

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestEnableDisableCustomer):
		disableAction, unmarshalErr := local_util.JsonUnmarshal[imodel.DisableCustomerCurrentAction](cpsAction.CurrentAction)
		if unmarshalErr != nil || disableAction == nil || disableAction.Channel == "" {
			err := d.repo.EnableOrDisable(ctx, cpsAction.UniqueId, actionData.Enabled)
			if err != nil {
				log.Errorf("[CustomerSvc][Authorize] enable/disable err: %v", err)
				span.AddEvent("Customer enable/disable action failed", trace.WithAttributes(
					attribute.String("error", err.Error()),
					attribute.String("unique_id", cpsAction.UniqueId),
				))
				return nil, err
			}
		} else {
			if err := d.repo.DisableCustomerByChannel(ctx, cpsAction.UniqueId, disableAction.Channel); err != nil {
				log.Errorf("[CustomerSvc][Authorize] disable by channel err: %v", err)
				span.AddEvent("Customer disable by channel failed", trace.WithAttributes(
					attribute.String("error", err.Error()),
					attribute.String("unique_id", cpsAction.UniqueId),
					attribute.String("channel", disableAction.Channel),
				))
				return nil, err
			}
		}
		log.Infof("[CustomerSvc][Authorize] enable/disable done id: %s", cpsAction.UniqueId)
	case string(constants.RequestBlockCustomer):
		err := d.repo.BlockCustomerByUserCode(ctx, cpsAction.UniqueId)
		if err != nil {
			log.Errorf("[CustomerSvc][Authorize] block err: %v", err)
			span.AddEvent("Customer block action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		log.Infof("[CustomerSvc][Authorize] block done id: %s", cpsAction.UniqueId)
	case string(constants.RequestUnblockCustomer):
		err := d.repo.UNBlockCustomerByUserCode(ctx, cpsAction.UniqueId)
		if err != nil {
			log.Errorf("[CustomerSvc][Authorize] unblock err: %v", err)
			span.AddEvent("Customer unblock action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		log.Infof("[CustomerSvc][Authorize] unblock done id: %s", cpsAction.UniqueId)
	case string(constants.RequestApproveFaydaCustomer):
		err := d.repo.Update(ctx, cpsAction.UniqueId, *actionData)
		if err != nil {
			log.Errorf("[CustomerSvc][Authorize] fayda update err: %v", err)
			span.AddEvent("Failed to update Fayda customer", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		log.Infof("[CustomerSvc][Authorize] fayda approved id: %s", cpsAction.UniqueId)
	default:
		log.Errorf("[CustomerSvc][Authorize] unsupported: %s", cpsAction.RequestAction)
		span.AddEvent("Unsupported action", trace.WithAttributes(
			attribute.String("error", localization.MsgInvalidAction),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		return nil, fmt.Errorf("%s", localization.MsgInvalidAction)
	}
	log.Infof("[CustomerSvc][Authorize] done: %s", cpsAction.RequestAction)
	return cpsAction, nil
}

func (d *customerService) SearchCustomerByCIForAccountNumber(ctx context.Context, number string) (*customer_dto.CustomerListResponse, error) {
	log := local_util.LoggerFromCtx(ctx, d.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "SearchCustomerByCIForAccountNumber", "Customer", "SearchCustomerByCIForAccountNumber")
	defer span.End()
	log.Infof("[CustomerSvc][SearchByCI] value: %s", number)

	number = local_util.NormalizePhoneNumberOrReturnInput(number)
	customer, err := d.repo.SearchCustomerByCIForAccountNumber(ctx, number)
	if err != nil {
		span.AddEvent("Failed to fetch customer by CI", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("cif_or_account_number", number),
		))
		return nil, err
	}

	if customer.AccountNumber != "" {
		accountDetail, lookupErr := d.core.LookupAccountByAccountNumber(ctx, model.AccountLookUpRequest{
			AccountNumber: customer.AccountNumber,
		})
		if lookupErr != nil {
			log.Errorf("[CustomerSvc][SearchByCI] account lookup failed (restriction skipped): %v", lookupErr)
		} else if accountDetail != nil {
			customer.Restriction = accountDetail.Restriction
		}
	}

	return customer, nil
}

// SearchCustomerServiceLimitByCIF implements [service.CustomerService].
func (s *customerService) SearchCustomerServiceLimitByCIF(ctx context.Context, cif string, filterParam *types.Filter) (types.PaginatedResponse[[]customer_dto.CustomerServiceLimitResponses], error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	log.Infof("[CustomerSvc][SearchCustomerServiceLimitByCIF] cif: %s", cif)

	res, err := s.core.SearchCustomerServiceLimitByCIF(ctx, cif)
	if err != nil {
		log.Errorf("[CustomerSvc][SearchCustomerServiceLimitByCIF] core search err: %v", err)
		return types.PaginatedResponse[[]customer_dto.CustomerServiceLimitResponses]{}, err
	}

	// Sort by service code first, then channel, so merged output stays deterministic.
	sort.Slice(res.Detail, func(i, j int) bool {
		if res.Detail[i].ServiceCode == res.Detail[j].ServiceCode {
			return res.Detail[i].Channel < res.Detail[j].Channel
		}
		return res.Detail[i].ServiceCode < res.Detail[j].ServiceCode
	})

	data := make([]customer_dto.CustomerServiceLimitResponses, 0)
	indexByServiceCode := make(map[string]int)
	channel_name := ""

	res_cif := ""
	serviceCode := ""
	serviceName := ""
	limit := ""
	count := ""

	for _, detail := range res.Detail {
		log.Infof("[CustomerSvc][SearchCustomerServiceLimitByCIF] detail: %v", detail)
		if detail.CIF != "" {
			res_cif = detail.CIF
		}
		if detail.ServiceCode != "" {
			serviceCode = detail.ServiceCode
		}
		if detail.ServiceName != "" {
			serviceName = detail.ServiceName
		}

		if detail.Limit != "" {
			limit = detail.Limit
		}
		if detail.Count != "" {
			count = detail.Count
		}

		channelData := customer_dto.CustomerServiceLimitResponse{
			CIF:         res_cif,
			Channel:     channel_name,
			ServiceCode: serviceCode,
			ServiceName: serviceName,
			Limit:       limit,
			Count:       count,
		}

		if idx, exists := indexByServiceCode[serviceCode]; exists {
			data[idx].CustomerServiceLimitResponse = append(data[idx].CustomerServiceLimitResponse, channelData)
			continue
		}

		data = append(data, customer_dto.CustomerServiceLimitResponses{
			ServiceCode:                  serviceCode,
			ServiceName:                  serviceName,
			CustomerServiceLimitResponse: []customer_dto.CustomerServiceLimitResponse{channelData},
		})
		indexByServiceCode[serviceCode] = len(data) - 1
	}

	filteredData := data
	page := 1
	perPage := 10
	search := ""
	if filterParam != nil {
		page = filterParam.Page
		perPage = filterParam.PerPage
		search = filterParam.Search
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	search = strings.TrimSpace(strings.ToLower(search))
	if search != "" {
		filtered := make([]customer_dto.CustomerServiceLimitResponses, 0, len(data))
		for _, item := range data {
			if strings.Contains(strings.ToLower(item.ServiceCode), search) || strings.Contains(strings.ToLower(item.ServiceName), search) {
				filtered = append(filtered, item)
			}
		}
		filteredData = filtered
	}

	totalDocs := int64(len(filteredData))
	meta := local_util.BuildPaginationMeta(totalDocs, page, perPage)

	start := (page - 1) * perPage
	if start >= len(filteredData) {
		return types.PaginatedResponse[[]customer_dto.CustomerServiceLimitResponses]{
			Data: []customer_dto.CustomerServiceLimitResponses{},
			Meta: meta,
		}, nil
	}

	end := start + perPage
	if end > len(filteredData) {
		end = len(filteredData)
	}

	filteredData = filteredData[start:end]

	return types.PaginatedResponse[[]customer_dto.CustomerServiceLimitResponses]{
		Data: filteredData,
		Meta: meta,
	}, nil
}

func (d *customerService) GetCustomerDetailByID(ctx context.Context, id string) (*customer_dto.CustomerDetailResponse, error) {
	log := local_util.LoggerFromCtx(ctx, d.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "GetCustomerDetailByID", "Customer", "GetCustomerDetailByID")
	defer span.End()

	log.Infof("[CustomerSvc][GetDetail] id: %s", id)
	res, err := d.repo.FindCustomerDetailByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to fetch customer detail by id", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}

	// use external call to get missing account detail
	coreRes, err := d.core.CifSearch(ctx, res.PersonalInfo.CustomerNumber)
	if err != nil {
		log.Errorf("[CustomerSvc][GetCustomerDetailByID] CifSearch error: %v", err)
		span.AddEvent("Failed to search CIF", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("customer_number", res.PersonalInfo.CustomerNumber),
		))
		return nil, err
	}

	log.Infof("core result data---------------------: %v", coreRes)

	if len(coreRes) > 0 {
		cr := coreRes[0]
		if cr.BirthOfDate != "" {
			res.PersonalInfo.DateOfBirth = cr.BirthOfDate
		}
		if cr.Email != "" {
			res.PersonalInfo.Email = cr.Email
		}
		if cr.Gender != "" {
			res.PersonalInfo.Gender = cr.Gender
		}
		if cr.PhoneNumber != "" {
			res.PersonalInfo.PhoneNumber = cr.PhoneNumber
		}
	}

	if len(res.LinkedAccount) != 0 {
		for i, linkedAccount := range res.LinkedAccount {

			if len(coreRes) != 0 {
				for _, coreResBranch := range coreRes {
					if linkedAccount.AccountNumber == coreResBranch.AccountNumber {
						res.LinkedAccount[i].AccountBranchCode = coreResBranch.BranchCode
						res.LinkedAccount[i].AccountBranchName = coreResBranch.Branch
						break
					}
				}
			}
			// res.LinkedAccount[i].AccountBranchName = coreResBranch.Data.na
		}
	}
	return res, nil
}

func (c *customerService) GetCustomerBarUnBarReasons(ctx context.Context, userID string) ([]*imodel.CustomerBarUnBarReason, error) {
	log := local_util.LoggerFromCtx(ctx, c.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "GetCustomerBarUnBarReasons", "Customer", "GetCustomerBarUnBarReasons")
	defer span.End()

	log.Infof("[CustomerSvc][GetCustomerBarUnBarReasons] userID: %s", userID)
	reasons, err := c.repo.GetBarUnBarReasons(ctx, userID)
	if err != nil {
		log.Errorf("[CustomerSvc][GetCustomerBarUnBarReasons] err: %v", err)
		return nil, err
	}
	return reasons, nil
}
