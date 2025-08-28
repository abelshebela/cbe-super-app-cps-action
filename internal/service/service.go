package service

import (
	"context"
	"mime/multipart"

	bank_dto "cbe-super-app-cps-action/internal/constants/dto/bank"
	eventdto "cbe-super-app-cps-action/internal/constants/dto/event"
	fbdto "cbe-super-app-cps-action/internal/constants/dto/feedback"
	walletDto "cbe-super-app-cps-action/internal/constants/dto/wallet"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
)

type CPSActionService interface {
	ApproveCPSAction(ctx context.Context, action *model.CPSAction) error
	CreateCPSAction(ctx context.Context, action *model.CPSAction) error
	RejectCPSAction(ctx context.Context, action_code string, action *model.CPSAction) error
	GetCPSActionsByDepartment(ctx context.Context, department string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error)
	GetCPSActionByID(ctx context.Context, id, department string) (*model.CPSAction, error)
	GetCPSActionByUniqueID(ctx context.Context, id, department string) (*model.CPSAction, error)
	GetCPSActionByActionCode(ctx context.Context, uniqueID, department string) (*model.CPSAction, error)
}

type BranchService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type BudgetService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type BulkService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type CPSUserService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type CustomerService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type DepartmentService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type DonationService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type EventService interface {
	CreateEvent(ctx context.Context, event eventdto.EventRequest) error
	UpdateEvent(ctx context.Context, id string, event eventdto.EventRequest) error
	DeleteEvent(ctx context.Context, id string) error
	EnableDisableEvent(ctx context.Context, id string, enable bool) error
	FetchEventByID(ctx context.Context, id string) (*model.Event, error)
	FetchEvent(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Event], error)
	Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error)
}

type FaydaAccountService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type FeedbackService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	CreateFeedback(ctx context.Context, req fbdto.FeedbackRequest, userID string) (*model.Feedback, error)
	GetFeedbackByID(ctx context.Context, id string) (*model.Feedback, error)
	GetFeedbacks(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Feedback], error)
}

type HQService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type MiniAppService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type MiniAppMerchantService interface {
	DetailMiniAppByID(ctx context.Context, id string) (*model.MiniAppMerchant, error)
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type NotificationService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type PasswordRuleService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type PermissionService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type PortalCardService interface {
	GetAll(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]*model.Card], error)
	ValidatePortalCard(ctx context.Context, names []string) (bool, error)
}

type ProductCodeService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type ServiceService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type UnlinkService interface {
	GetUserByAccount(ctx context.Context, accNumber string) (*model.User, error)
	GetAllArchivedUser(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.ArchivedUser], error)
	UnlinkUserCif(ctx context.Context, userCode string) error
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type WalletService interface {
	CreateWallet(ctx context.Context, req walletDto.WalletRequest) error
	UpdateWallet(ctx context.Context, id string, req walletDto.WalletRequest) error
	DeleteWallet(ctx context.Context, id string) error
	EnableOrDisableWallet(ctx context.Context, id string, enable bool) error
	GetWallet(ctx context.Context, id string) (*model.Wallet, error)
	GetAllWallet(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]*model.Wallet], error)
	Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error)
}

type AccountBlockService interface {
	GetBranchByCode(ctx context.Context, branchCode string) (*model.Branch, error)
	GetAllBranches(ctx context.Context, filter *types.Filter) (*types.PaginatedResponse[[]*model.Branch], error)
	GetRegionByCode(ctx context.Context, regionCode string) (*model.Region, error)
	GetAllRegions(ctx context.Context, filter *types.Filter) (*types.PaginatedResponse[[]*model.Region], error)
	GetDistrictByCode(ctx context.Context, districtCode string) (*model.District, error)
	GetAllDistricts(ctx context.Context, filter *types.Filter) (*types.PaginatedResponse[[]*model.District], error)
	GetCityByCode(ctx context.Context, cityCode string) (*model.City, error)
	GetAllCities(ctx context.Context, filter *types.Filter) (*types.PaginatedResponse[[]*model.City], error)
	EnableOrDisableBranches(ctx context.Context, branchCodes []string, enabled bool) error
	EnableOrDisableRegions(ctx context.Context, regionsCode []string, enabled bool) error
	EnableOrDisableDistricts(ctx context.Context, districtsCode []string, enabled bool) error
	EnableOrDisableCities(ctx context.Context, citiesCode []string, enabled bool) error
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type AccountValidationService interface {
	Update(ctx context.Context, id string, rule *model.ValidationRule) error
	FindById(ctx context.Context, id string) (*model.ValidationRule, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.ValidationRule], error)
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type ActionService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type AdvertService interface {
	CreateAdvert(ctx context.Context, ad *model.Advert, bannerImage *multipart.FileHeader) error
	FetchAdverts(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Advert], error)
	FetchAdvertByID(ctx context.Context, id string) (*model.Advert, error)
	UpdateAdvert(ctx context.Context, id string, ad *model.Advert, bannerImage *multipart.FileHeader) error
	DeleteAdvert(ctx context.Context, id string) error
	EnableDisableAdvert(ctx context.Context, id string, enable bool) error
	Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error)
}

type AmountBasedAuthService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type AvatarService interface {
	CreateAvatar(ctx context.Context, avatar *model.Avatar) error
	UpdateAvatar(ctx context.Context, id string, avatar *model.Avatar) error
	DeleteAvatar(ctx context.Context, id string) error
	FetchAllAvatar(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]*model.Avatar], error)
	FetchAvatarById(ctx context.Context, id string) (*model.Avatar, error)
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type BankService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	GetAllBank(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Bank], error)

	GetOneBank(ctx context.Context, id string) (*model.Bank, error)

	CreateOneBank(ctx context.Context, req bank_dto.CreateBankRequest) error
	UpdateOneBank(ctx context.Context, id string, req bank_dto.UpdateBankRequest) error
	DeleteOneBank(ctx context.Context, id string) error

	EnableOrDisableBank(ctx context.Context, id string, enableDisable bool) error
	UpdateLogo(ctx context.Context, id string, logo bank_dto.UpdateLogo) error
}

type BudgetCategoryService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type BPSUserService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) error
	FetchUserByUserCode(ctx context.Context, userCode string) (*model.BPSUser, error)
	GetAllBPSUsers(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.BPSUser], error)
	UpdateBpsUser(ctx context.Context, userCode string, status bool) error
}

type AccountSearchService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type KeyGeneratorService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type ServiceContainer struct {
	AccountBlockContainer    AccountBlockService
	AccountContainer         AccountValidationService
	ActionContainer          ActionService
	AdContainer              AdvertService
	AmountBasedAuthContainer AmountBasedAuthService
	AvatarDomian             AvatarService
	BankContainer            BankService
	BPSUserContainer         BPSUserService
	BudgetCategoryContainer  BudgetCategoryService
	BudgetContainer          BudgetService
	CPSActionContainer       CPSActionService
	CPSUserContainer         CPSUserService
	CustomerContainer        CustomerService
	DepartmentContainer      DepartmentService
	EventContainer           EventService // fully not ready
	FaydaContainer           FaydaAccountService
	FeedbackContainer        FeedbackService
	HQContainer              HQService
	MiniAppContainer         MiniAppService
	PasswordRuleContainer    PasswordRuleService
	PermissionContainer      PermissionService
	PortalCardContainer      PortalCardService
	UnlinkContainer          UnlinkService
	WalletContainer          WalletService
	MiniAppMerchantContainer MiniAppMerchantService
	AccountLookup            AccountSearchService
	BulkServiceContainer     BulkService
	ServiceCheckContainer    ServiceService
	KeyGenService            KeyGeneratorService
	NotificationService      NotificationService
	ProductCodeService       ProductCodeService
	DonationContainer        DonationService
}
