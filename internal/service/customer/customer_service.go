package customer

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/dto/customer"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"

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
}

func NewCustomerService(repo storage.CustomerRepository, bpsRepo storage.BPSActionRepository, cpsService service.CPSActionService, redis storage.RedisRepository, smsService *lib.NotificationStore, cfg *config.VaultConfig, logger utils.Logger) service.CustomerService {
	return &customerService{
		repo:       repo,
		bpsRepo:    bpsRepo,
		cpsService: cpsService,
		redis:      redis,
		cfg:        cfg,
		logger:     logger,
		smsService: smsService,
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
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateEnableCustomerSession", "Customer", "CreateEnableCustomerSession")
	defer span.End()

	c.logger.Infof("[CustomerSvc][CreateEnableSession] id: %s", id)
	existing_otp, err := c.redis.Get(ctx, fmt.Sprintf("cps:action:otp:%s", id))
	if existing_otp != "" || err == nil {
		c.logger.Errorf("[CustomerSvc][CreateEnableSession] OTP exists")
		span.AddEvent("OTP already exists", trace.WithAttributes(
			attribute.String("error", localization.ErrorOTPAlreadyExists.Code),
			attribute.String("id", id),
		))
		return "", errors.New(localization.ErrorOTPAlreadyExists.Code)
	}

	customer, err := c.repo.FindByID(ctx, id)
	if err != nil {
		c.logger.Errorf("[CustomerSvc][CreateEnableSession] find err: %v", err)
		span.AddEvent("Failed to find customer", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return "", err
	}
	if customer.Enabled {
		c.logger.Errorf("[CustomerSvc][CreateEnableSession] already enabled")
		span.AddEvent("Customer already enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorCustomerAlreadyEnabled.Code),
			attribute.String("id", id),
		))
		return "", errors.New(localization.ErrorCustomerAlreadyEnabled.Code)
	}

	otp := local_util.OTPGenerator(6)
	encryptedOTP, _, err := local_util.LocalEncryptPassword(otp, constants.OTP, constants.OTP, constants.OTP, c.cfg)
	if err != nil {
		c.logger.Errorf("[CustomerSvc][CreateEnableSession] encrypt OTP err: %v", err)
		span.AddEvent("Failed to encrypt OTP", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return "", err
	}

	err = c.redis.Set(ctx, fmt.Sprintf("cps:action:otp:%s", id), encryptedOTP, constants.OtpExpirationTime)
	if err != nil {
		c.logger.Errorf("[CustomerSvc][CreateEnableSession] redis set err: %v", err)
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
			c.logger.Errorf("[CustomerSvc][CreateEnableSession] SMS err: %v", err)
		}
	}()

	c.logger.Infof("[CustomerSvc][CreateEnableSession] OTP generated id: %s", id)
	return otp, nil
}

func (c *customerService) ApproveFaydaCustomer(ctx context.Context, id string, req customer.FaydaApproveRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "ApproveFaydaCustomer", "Customer", "ApproveFaydaCustomer")
	defer span.End()

	c.logger.Infof("[CustomerSvc][ApproveFayda] id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		c.logger.Errorf("[CustomerSvc][ApproveFayda] incomplete user")
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.Incomplete),
			attribute.String("id", id),
		))
		return errors.New(constants.Incomplete)
	}

	customer, err := c.repo.FindByID(ctx, id)
	if err != nil {
		c.logger.Errorf("[CustomerSvc][ApproveFayda] find err: %v", err)
		span.AddEvent("Failed to find customer", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	if customer.Enabled {
		c.logger.Errorf("[CustomerSvc][ApproveFayda] already enabled")
		span.AddEvent("Customer already enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorCustomerAlreadyEnabled.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorCustomerAlreadyEnabled.Code)
	}

	// if customer.KYCLevel != uint8(constants.ONE) {
	// 	c.logger.Errorf("[ApproveFaydaCustomer] user is not a Fayda registered customer")
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
		c.logger.Errorf("[CustomerSvc][ApproveFayda] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	c.logger.Infof("[CustomerSvc][ApproveFayda] request created id: %s", id)
	return nil
}

func (c *customerService) EnableCustomerByID(ctx context.Context, id string, user_otp string) error {
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
		c.logger.Infof("[CustomerSvc][Enable] OTP mismatch")
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
		c.logger.Errorf("[CustomerSvc][Enable] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	if err := c.redis.Delete(ctx, fmt.Sprintf("cps:action:otp:%s", id)); err != nil {
		c.logger.Errorf("[CustomerSvc][Enable] redis delete err: %v", err)
	}
	c.logger.Infof("[CustomerSvc][Enable] request created id: %s", id)
	return nil
}

func (c *customerService) DisableCustomerByID(ctx context.Context, id string, disable customer.CustomerDisableDTO) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DisableCustomerByID", "Customer", "DisableCustomerByID")
	defer span.End()

	c.logger.Infof("[CustomerSvc][Disable] id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		c.logger.Errorf("[CustomerSvc][Disable] incomplete user")
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.IncompleteUserInfo),
			attribute.String("id", id),
		))
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	customer, err := c.repo.FindByID(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		c.logger.Errorf("[CustomerSvc][Disable] not found")
		span.AddEvent("Customer not found", trace.WithAttributes(
			attribute.String("error", code),
			attribute.String("id", id),
		))
		return fmt.Errorf("%s", code)
	} else if err != nil {
		c.logger.Errorf("[CustomerSvc][Disable] find err: %v", err)
		span.AddEvent("Failed to find customer", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	if !customer.Enabled {
		c.logger.Errorf("[CustomerSvc][Disable] already disabled")
		span.AddEvent("Customer already disabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorCustomerAlreadyDisabled.Code),
			attribute.String("id", id),
		))
		return fmt.Errorf("%s", localization.ErrorCustomerAlreadyDisabled.Code)
	}

	new_customer := *customer
	new_customer.Enabled = false
	new_customer.IsBlocked = true
	new_customer.BlockedReason = disable.DisableReason

	if *disable.IsTemporary {
		new_customer.BlockedOn = member.BlockedOn(constants.BPS)
	} else {
		new_customer.BlockedOn = member.BlockedOn(constants.CPS)
	}

	action := lib.CpsModelBuilder(id, makerData, customer, new_customer, string(constants.RequestEnableDisableCustomer), constants.UPDATE)

	err = c.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		c.logger.Errorf("[CustomerSvc][Disable] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	c.logger.Infof("[CustomerSvc][Disable] request created id: %s", id)
	return nil
}

func (d *customerService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Customer", "Authorize")
	defer span.End()

	d.logger.Infof("[CustomerSvc][Authorize] action: %s", cpsAction.RequestAction)
	actionData, err := local_util.JsonUnmarshal[member.User](cpsAction.CurrentAction)
	if err != nil {
		d.logger.Errorf("[CustomerSvc][Authorize] unmarshal err: %v", err)
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		return nil, fmt.Errorf("failed to unmarshal CurrentAction to User: %v", err)
	}

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestEnableDisableCustomer):
		err := d.repo.EnableOrDisable(ctx, cpsAction.UniqueId, actionData.Enabled)
		if err != nil {
			d.logger.Errorf("[CustomerSvc][Authorize] enable/disable err: %v", err)
			span.AddEvent("Customer enable/disable action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		d.logger.Infof("[CustomerSvc][Authorize] enable/disable done id: %s", cpsAction.UniqueId)
	case string(constants.RequestApproveFaydaCustomer):
		err := d.repo.Update(ctx, cpsAction.UniqueId, *actionData)
		if err != nil {
			d.logger.Errorf("[CustomerSvc][Authorize] fayda update err: %v", err)
			span.AddEvent("Failed to update Fayda customer", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		d.logger.Infof("[CustomerSvc][Authorize] fayda approved id: %s", cpsAction.UniqueId)
	default:
		d.logger.Errorf("[CustomerSvc][Authorize] unsupported: %s", cpsAction.RequestAction)
		span.AddEvent("Unsupported action", trace.WithAttributes(
			attribute.String("error", localization.MsgInvalidAction),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		return nil, fmt.Errorf("%s", localization.MsgInvalidAction)
	}
	d.logger.Infof("[CustomerSvc][Authorize] done: %s", cpsAction.RequestAction)
	return cpsAction, nil
}

func (d *customerService) SearchCustomerByCIForAccountNumber(ctx context.Context, number string) (*customer_dto.CustomerListResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "SearchCustomerByCIForAccountNumber", "Customer", "SearchCustomerByCIForAccountNumber")
	defer span.End()
	d.logger.Infof("[CustomerSvc][SearchByCI] value: %s", number)

	number = local_util.NormalizePhoneNumberOrReturnInput(number)
	customer, err := d.repo.SearchCustomerByCIForAccountNumber(ctx, number)
	if err != nil {
		span.AddEvent("Failed to fetch customer by CI", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("cif_or_account_number", number),
		))
		return nil, err
	}
	return customer, nil
}

func (d *customerService) GetCustomerDetailByID(ctx context.Context, id string) (*customer_dto.CustomerDetailResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCustomerDetailByID", "Customer", "GetCustomerDetailByID")
	defer span.End()
	d.logger.Infof("[CustomerSvc][GetDetail] id: %s", id)
	res, err := d.repo.FindCustomerDetailByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to fetch customer detail by id", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}
	return res, nil
}
