package factory

import (
	"cbe-super-app-cps-action/internal/constants"
	eventdto "cbe-super-app-cps-action/internal/constants/dto/event"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// Account Block Service
type accountBlockService struct {
	repo   storage.AccountBlockRepository
	logger utils.Logger
}

func (a *accountBlockService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	a.logger.Infof("Account Block service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Account Validation Service
type accountValidationService struct {
	repo   storage.UserRepository
	logger utils.Logger
}

func (a *accountValidationService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	a.logger.Infof("Account Validation service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Account Search Service
type accountSearchService struct {
	repo   storage.AccountAPIPort
	logger utils.Logger
}

func (a *accountSearchService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	a.logger.Infof("Account Search service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Action Service
type actionService struct {
	logger utils.Logger
}

func (a *actionService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	a.logger.Infof("Action service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Advert Service
type advertService struct {
	repo   storage.AdvertRepository
	logger utils.Logger
}

func (a *advertService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	a.logger.Infof("Advert service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Amount Based Auth Service
type amountBasedAuthService struct {
	repo   storage.AmountBasedAuthRepository
	logger utils.Logger
}

func (a *amountBasedAuthService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	a.logger.Infof("Amount Based Auth service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Avatar Service
type avatarService struct {
	repo   storage.AvatarRepository
	logger utils.Logger
}

func (a *avatarService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	a.logger.Infof("Avatar service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Bank Service
type bankService struct {
	repo   storage.BankRepository
	logger utils.Logger
}

func (b *bankService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	b.logger.Infof("Bank service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// BPS User Service
type bpsUserService struct {
	repo   storage.BPSUserRepository
	logger utils.Logger
}

func (b *bpsUserService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	b.logger.Infof("BPS User service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Budget Category Service
type budgetCategoryService struct {
	repo   storage.AmountBasedAuthRepository
	logger utils.Logger
}

func (b *budgetCategoryService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	b.logger.Infof("Budget Category service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Budget Service
type budgetService struct {
	repo   storage.AmountBasedAuthRepository
	logger utils.Logger
}

func (b *budgetService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	b.logger.Infof("Budget service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Bulk Service
type bulkService struct {
	repo   storage.AccountBlockRepository
	logger utils.Logger
}

func (b *bulkService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	b.logger.Infof("Bulk service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// CPS Action Service
type cpsActionService struct {
	logger utils.Logger
}

func (c *cpsActionService) ApproveCPSAction(ctx context.Context, action *model.CPSAction) error {
	c.logger.Infof("CPS Action service approving action: %s", action.ActionCode)
	return nil
}

func (c *cpsActionService) RejectCPSAction(ctx context.Context, action_code string, action *model.CPSAction) error {
	c.logger.Infof("CPS Action service rejecting action: %s", action_code)
	return nil
}

func (c *cpsActionService) GetCPSActionsByDepartment(ctx context.Context, department string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	c.logger.Infof("CPS Action service getting actions for department: %s", department)
	return nil, nil
}

func (c *cpsActionService) GetCPSActionByID(ctx context.Context, id, department string) (*model.CPSAction, error) {
	c.logger.Infof("CPS Action service getting action by ID: %s", id)
	return nil, nil
}
func (s *cpsActionService) CPSActionExists(ctx context.Context, user model.CheckCPSAction) (bool, error) {
	return true, nil
}
func (c *cpsActionService) GetCPSActionByActionCode(ctx context.Context, uniqueID, department string) (*model.CPSAction, error) {
	c.logger.Infof("CPS Action service getting action by code: %s", uniqueID)
	return nil, nil
}
func (c *cpsActionService) GetCPSActionByUniqueID(ctx context.Context, uniqueID, department string) (*model.CPSAction, error) {
	c.logger.Infof("CPS Action service getting action by code: %s", uniqueID)
	return nil, nil
}
func (ca *cpsActionService) CreateCPSAction(ctx context.Context, cpsAction *model.CPSAction) error {

	return nil

}

// CPS User Service
type cpsUserService struct {
	repo   storage.CpsUserRepository
	logger utils.Logger
}

func (c *cpsUserService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	c.logger.Infof("CPS User service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Customer Service
type customerService struct {
	repo   storage.UserRepository
	logger utils.Logger
}

func (c *customerService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	c.logger.Infof("Customer service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Department Service
type departmentService struct {
	repo   storage.AppAccessListRepository
	logger utils.Logger
}

func (d *departmentService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	d.logger.Infof("Department service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Donation Service
type donationService struct {
	repo   storage.DonationRepository
	logger utils.Logger
}

func (d *donationService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	d.logger.Infof("Donation service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Event Service
type eventService struct {
	repo   storage.EventRepository
	logger utils.Logger
}

func (e *eventService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	e.logger.Infof("Event service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}
func (e *eventService) CreateEvent(ctx context.Context, event eventdto.EventRequest) error {
	return errors.New("CreateEvent: not implemented")
}

func (e *eventService) UpdateEvent(ctx context.Context, id string, event eventdto.EventRequest) error {
	return errors.New("UpdateEvent: not implemented")
}

func (e *eventService) DeleteEvent(ctx context.Context, id string) error {
	return errors.New("DeleteEvent: not implemented")
}

func (e *eventService) EnableDisableEvent(ctx context.Context, id string, enable bool) error {
	return errors.New("EnableDisableEvent: not implemented")
}

func (e *eventService) FetchEventByID(ctx context.Context, id string) (*model.Event, error) {
	return nil, errors.New("FetchEventByID: not implemented")
}

func (e *eventService) FetchEvent(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Event], error) {
	return nil, errors.New("FetchEvent: not implemented")
}

// Fayda Service
type faydaService struct {
	repo   storage.LinkedAccountRepository
	logger utils.Logger
}

func (f *faydaService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	f.logger.Infof("Fayda service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Feedback Service
type feedbackService struct {
	repo   storage.FeedbackRepository
	logger utils.Logger
}

func (f *feedbackService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	f.logger.Infof("Feedback service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// HQ Service
type hqService struct {
	repo   storage.HQRepository
	logger utils.Logger
}

func (h *hqService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	h.logger.Infof("HQ service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Key Generator Service
type keyGeneratorService struct {
	logger utils.Logger
}

func (k *keyGeneratorService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	k.logger.Infof("Key Generator service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Mini App Service
type miniAppService struct {
	repo   storage.MiniAppRepository
	logger utils.Logger
}

func (m *miniAppService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	m.logger.Infof("Mini App service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Mini App Merchant Service
type miniAppMerchantService struct {
	repo   storage.MiniAppMerchantRepository
	logger utils.Logger
}

func (m *miniAppMerchantService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	m.logger.Infof("Mini App Merchant service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}
func (m *miniAppMerchantService) DetailMiniAppByID(ctx context.Context, id string) (*model.MiniAppMerchant, error) {

	m.logger.Infof("Mini App Merchant service authorizing action: %s", constants.ActionCode)

	return m.repo.DetailMiniAppByID(ctx, id)
}

// Notification Service
type notificationService struct {
	repo   storage.NotificationRepository
	logger utils.Logger
}

func (n *notificationService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	n.logger.Infof("Notification service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Password Rule Service
type passwordRuleService struct {
	repo   storage.PasswordRuleRepository
	logger utils.Logger
}

func (p *passwordRuleService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	p.logger.Infof("Password Rule service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Permission Service
type permissionService struct {
	logger utils.Logger
}

func (p *permissionService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	p.logger.Infof("Permission service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Portal Card Service
type portalCardService struct {
	repo   storage.PortalCardRepository
	logger utils.Logger
}

func (p *portalCardService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	p.logger.Infof("Portal Card service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Product Code Service
type productCodeService struct {
	logger utils.Logger
}

func (p *productCodeService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	p.logger.Infof("Product Code service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Service Service
type serviceService struct {
	repo   storage.ServiceDetailsRepository
	logger utils.Logger
}

func (s *serviceService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	s.logger.Infof("Service service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Unlink Service
type unlinkService struct {
	logger utils.Logger
}

func (u *unlinkService) Authorize(ctx context.Context, cpsAction *model.CPSAction) error {
	u.logger.Infof("Unlink service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return nil
}

func (u *unlinkService) GetUserByAccount(ctx context.Context, accNumber string) (*model.ArchivedUser, error) {
	return nil, errors.New("GetUserByAccount: not implemented")
}

func (u *unlinkService) GetAllArchivedUser(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[*model.ArchivedUser], error) {
	return nil, errors.New("GetAllArchivedUser: not implemented")
}

func (u *unlinkService) UnlinkUserCif(ctx context.Context, userCode string) error {
	return errors.New("UnlinkUserCif: not implemented")
}

// Wallet Service
type walletService struct {
	repo   storage.WalletRepository
	logger utils.Logger
}

func (w *walletService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	w.logger.Infof("Wallet service authorizing action: %s", cpsAction.ActionCode)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}
