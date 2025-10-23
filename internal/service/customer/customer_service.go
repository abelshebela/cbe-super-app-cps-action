package customer

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/dto/customer"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/external_call"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"fmt"

	"cbe-super-app-cps-action/internal/constants/types"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type customerService struct {
	repo       storage.CustomerRepository
	redis      storage.RedisRepository
	cpsService service.CPSActionService
	cfg        *config.VaultConfig
	logger     utils.Logger
	smsService *external_call.SMSPersistence
}

func NewCustomerService(repo storage.CustomerRepository, cpsService service.CPSActionService, redis storage.RedisRepository, smsService *external_call.SMSPersistence, cfg *config.VaultConfig, logger utils.Logger) service.CustomerService {
	return &customerService{
		repo:       repo,
		cpsService: cpsService,
		redis:      redis,
		cfg:        cfg,
		logger:     logger,
		smsService: smsService,
	}
}

func (c *customerService) GetCustomersDetail(ctx context.Context, kyc_level int, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.User], error) {

	customers, err := c.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (s *customerService) GetCustomerByID(ctx context.Context, id string) (*model.User, error) {
	customer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *customerService) GetBlockedCustomer(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.User], error) {
	if filterParams.Filters == nil {
		filterParams.Filters = make(map[string]interface{})
	}
	filterParams.Filters["is_blocked"] = true
	customers, err := s.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		return nil, err
	}
	return customers, nil
}

// CreateEnableCustomerSession implements service.CustomerService.
func (c *customerService) CreateEnableCustomerSession(ctx context.Context, id string) (string, error) {
	existing_otp, err := c.redis.Get(ctx, fmt.Sprintf("cps:action:otp:%s", id))
	if existing_otp != "" || err == nil {
		return "", fmt.Errorf("%s", localization.ErrorOTPAlreadyExists.Code)
	}

	customer, err := c.repo.FindByID(ctx, id)
	if err != nil {
		return "", err
	}
	c.logger.Infof("Customer fetched for enabling session: %+v", customer.ID, customer.Enabled)
	if customer.Enabled {
		return "", fmt.Errorf("%s", localization.ErrorCustomerAlreadyEnabled.Code)
	}

	otp := local_util.OTPGenerator(6)
	encryptedOTP, _, err := local_util.LocalEncryptPassword(otp, constants.OTP, constants.OTP, constants.OTP, c.cfg)
	if err != nil {
		return "", err
	}

	err = c.redis.Set(ctx, fmt.Sprintf("cps:action:otp:%s", id), encryptedOTP, constants.OtpExpirationTime)
	if err != nil {
		return "", err
	}

	go func() {
		err = c.smsService.SendSMS(context.Background(), customer.PhoneNumber, fmt.Sprintf("Your Verification OTP is: %s", otp))
		if err != nil {
			c.logger.Errorf("failed to send OTP SMS to customer %s: %v", id, err)
		}
	}()

	c.logger.Infof("OTP for enabling customer with ID %s is %s", id, otp)
	return otp, nil
}

func (c *customerService) EnableCustomerByID(ctx context.Context, id string, user_otp string) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	otp, err := c.redis.Get(ctx, fmt.Sprintf("cps:action:otp:%s", id))

	if otp == "" || err != nil {
		return fmt.Errorf("%s", localization.ErrorOTPNotFound.Code)
	}
	encryptedOTP, _, err := local_util.LocalEncryptPassword(user_otp, constants.OTP, constants.OTP, constants.OTP, c.cfg)
	if err != nil {
		return err
	}

	if encryptedOTP != otp {
		c.logger.Infof("OTP mismatch: %s != %s", encryptedOTP, otp)
		return fmt.Errorf("%s", localization.ErrorOTPInvalid.Code)
	}

	customer, err := c.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	new_customer := *customer
	new_customer.Enabled = true

	action := lib.CpsModelBuilder(id, makerData, customer, new_customer, string(constants.RequestEnableDisableCustomer), constants.UPDATE)

	err = c.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		return err
	}

	if err := c.redis.Delete(ctx, fmt.Sprintf("cps:action:otp:%s", id)); err != nil {
		c.logger.Errorf("failed to delete OTP from redis for customer %s: %v", id, err)
	}
	return nil
}

func (c *customerService) DisableCustomerByID(ctx context.Context, id string, disable customer.CustomerDisableDTO) error {

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	customer, err := c.repo.FindByID(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		return fmt.Errorf("%s", code)
	} else if err != nil {
		return err
	}

	if !customer.Enabled {
		return fmt.Errorf("%s", localization.ErrorCustomerAlreadyDisabled.Code)
	}

	new_customer := *customer
	new_customer.Enabled = false
	new_customer.BlockedReason = disable.DisableReason

	if *disable.IsTemporary {
		new_customer.BlockedOn = constants.BPS
	} else {
		new_customer.BlockedOn = constants.CPS
	}

	action := lib.CpsModelBuilder(id, makerData, customer, new_customer, string(constants.RequestEnableDisableCustomer), constants.UPDATE)

	err = c.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		return err
	}

	return nil
}

func (d *customerService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {

	actionData, err := local_util.JsonUnmarshal[model.User](cpsAction.CurrentAction)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal CurrentAction to User: %v", err)
	}

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestEnableDisableCustomer):
		err := d.repo.EnableOrDisable(ctx, cpsAction.UniqueId, actionData.Enabled)
		if err != nil {
			d.logger.Errorf("Customer Enable Disable action  failed", "error", err)
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%s", localization.MsgInvalidAction)
	}
	return cpsAction, nil
}
