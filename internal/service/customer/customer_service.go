package customer

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/dto/customer"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
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
	smsService *lib.NotificationStore
}

func NewCustomerService(repo storage.CustomerRepository, cpsService service.CPSActionService, redis storage.RedisRepository, smsService *lib.NotificationStore, cfg *config.VaultConfig, logger utils.Logger) service.CustomerService {
	return &customerService{
		repo:       repo,
		cpsService: cpsService,
		redis:      redis,
		cfg:        cfg,
		logger:     logger,
		smsService: smsService,
	}
}

func (c *customerService) GetCustomersDetail(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.User], error) {
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

func (s *customerService) GetLinkedAccount(ctx context.Context, customerNumber string) ([]*model.LinkedAccount, error) {
	return s.repo.FetchLinkedAccount(ctx, customerNumber)
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
	c.logger.Infof("[CreateEnableCustomerSession] creating enable customer session for id: %s", id)
	existing_otp, err := c.redis.Get(ctx, fmt.Sprintf("cps:action:otp:%s", id))
	if existing_otp != "" || err == nil {
		c.logger.Errorf("[CreateEnableCustomerSession] OTP already exists for customer")
		return "", errors.New(localization.ErrorOTPAlreadyExists.Code)
	}

	customer, err := c.repo.FindByID(ctx, id)
	if err != nil {
		c.logger.Errorf("[CreateEnableCustomerSession] failed to find customer: %v", err)
		return "", err
	}
	if customer.Enabled {
		c.logger.Errorf("[CreateEnableCustomerSession] customer already enabled")
		return "", errors.New(localization.ErrorCustomerAlreadyEnabled.Code)
	}

	otp := local_util.OTPGenerator(6)
	encryptedOTP, _, err := local_util.LocalEncryptPassword(otp, constants.OTP, constants.OTP, constants.OTP, c.cfg)
	if err != nil {
		c.logger.Errorf("[CreateEnableCustomerSession] failed to encrypt OTP: %v", err)
		return "", err
	}

	err = c.redis.Set(ctx, fmt.Sprintf("cps:action:otp:%s", id), encryptedOTP, constants.OtpExpirationTime)
	if err != nil {
		c.logger.Errorf("[CreateEnableCustomerSession] failed to store OTP in redis: %v", err)
		return "", err
	}

	go func() {

		if err := c.smsService.PublishMessage(context.Background(), types.SMSKafkaMessage{
			Recipient:   customer.PhoneNumber,
			MessageBody: fmt.Sprintf("Your Supper app verification OTP: %s ", otp),
		}); err != nil {
			c.logger.Errorf("[CreateEnableCustomerSession] failed to send OTP SMS: %v", err)
		}
	}()

	c.logger.Infof("[CreateEnableCustomerSession] OTP generated successfully for customer id: %s", id)
	return otp, nil
}

func (c *customerService) ApproveFaydaCustomer(ctx context.Context, id string, req customer.FaydaApproveRequest) error {
	c.logger.Infof("[ApproveFaydaCustomer] approving Fayda customer for id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		c.logger.Errorf("[ApproveFaydaCustomer] incomplete user data")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	customer, err := c.repo.FindByID(ctx, id)
	if err != nil {
		c.logger.Errorf("[ApproveFaydaCustomer] failed to find customer: %v", err)
		return err
	}
	if customer.Enabled {
		c.logger.Errorf("[ApproveFaydaCustomer] customer already enabled")
		return errors.New(localization.ErrorCustomerAlreadyEnabled.Code)
	}

	if customer.KYCLevel != uint8(constants.ONE) {
		c.logger.Errorf("[ApproveFaydaCustomer] user is not a Fayda registered customer")
		return errors.New(localization.ErrorUserNotFaydaRegistered.Code)
	}

	updateData := *customer

	updateData.FaydaRiskLevel = req.RiskLevel

	action := lib.CpsModelBuilder(id, makerData, customer, updateData, string(constants.RequestApproveFaydaCustomer), constants.UPDATE)

	if err := c.cpsService.CreateCPSAction(ctx, &action); err != nil {
		c.logger.Errorf("[ApproveFaydaCustomer] failed to create CPS action: %v", err)
		return err
	}
	c.logger.Infof("[ApproveFaydaCustomer] Fayda customer approval request created successfully for id: %s", id)
	return nil
}

func (c *customerService) EnableCustomerByID(ctx context.Context, id string, user_otp string) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
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
		return errors.New(localization.ErrorOTPInvalid.Code)
	}

	customer, err := c.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	new_customer := *customer
	new_customer.Enabled = true
	new_customer.IsBlocked = false

	action := lib.CpsModelBuilder(id, makerData, customer, new_customer, string(constants.RequestEnableDisableCustomer), constants.UPDATE)

	err = c.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		c.logger.Errorf("[EnableCustomerByID] failed to create CPS action: %v", err)
		return err
	}

	if err := c.redis.Delete(ctx, fmt.Sprintf("cps:action:otp:%s", id)); err != nil {
		c.logger.Errorf("[EnableCustomerByID] failed to delete OTP from redis: %v", err)
	}
	c.logger.Infof("[EnableCustomerByID] customer enable request created successfully for id: %s", id)
	return nil
}

func (c *customerService) DisableCustomerByID(ctx context.Context, id string, disable customer.CustomerDisableDTO) error {
	c.logger.Infof("[DisableCustomerByID] disabling customer for id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		c.logger.Errorf("[DisableCustomerByID] incomplete user data")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	customer, err := c.repo.FindByID(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		c.logger.Errorf("[DisableCustomerByID] customer not found")
		return errors.New(localization.ErrorResourceNotFound.Code)
	} else if err != nil {
		c.logger.Errorf("[DisableCustomerByID] failed to find customer: %v", err)
		return err
	}

	if !customer.Enabled {
		c.logger.Errorf("[DisableCustomerByID] customer already disabled")
		return errors.New(localization.ErrorCustomerAlreadyDisabled.Code)
	}

	new_customer := *customer
	new_customer.Enabled = false
	new_customer.IsBlocked = true
	new_customer.BlockedReason = disable.DisableReason

	if *disable.IsTemporary {
		new_customer.BlockedOn = constants.BPS
	} else {
		new_customer.BlockedOn = constants.CPS
	}

	action := lib.CpsModelBuilder(id, makerData, customer, new_customer, string(constants.RequestEnableDisableCustomer), constants.UPDATE)

	err = c.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		c.logger.Errorf("[DisableCustomerByID] failed to create CPS action: %v", err)
		return err
	}
	c.logger.Infof("[DisableCustomerByID] customer disable request created successfully for id: %s", id)
	return nil
}

func (d *customerService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	d.logger.Infof("[Authorize] authorizing customer action: %s", cpsAction.RequestAction)
	actionData, err := local_util.JsonUnmarshal[model.User](cpsAction.CurrentAction)
	if err != nil {
		d.logger.Errorf("[Authorize] failed to unmarshal CurrentAction to User: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestEnableDisableCustomer):
		err := d.repo.EnableOrDisable(ctx, cpsAction.UniqueId, actionData.Enabled)
		if err != nil {
			d.logger.Errorf("[Authorize] customer enable/disable action failed: %v", err)
			return nil, err
		}
		d.logger.Infof("[Authorize] customer enable/disable action completed successfully for id: %s", cpsAction.UniqueId)
	case string(constants.RequestApproveFaydaCustomer):
		err := d.repo.Update(ctx, cpsAction.UniqueId, *actionData)
		if err != nil {
			d.logger.Errorf("[Authorize] failed to update Fayda customer: %v", err)
			return nil, err
		}
		d.logger.Infof("[Authorize] Fayda customer approval completed successfully for id: %s", cpsAction.UniqueId)
	default:
		d.logger.Errorf("[Authorize] unsupported action: %s", cpsAction.RequestAction)
		return nil, fmt.Errorf("%s", localization.MsgInvalidAction)
	}
	d.logger.Infof("[Authorize] customer action authorized successfully: %s", cpsAction.RequestAction)
	return cpsAction, nil
}
